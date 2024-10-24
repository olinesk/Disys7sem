package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"sync"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	proto "ChittyChat/proto"
)

var (
	userPort = flag.Int("uPort", 0, "user port number")
	serverPort = flag.Int("sPort", 0, "server port number (same port number as the one used for server)")
	userName = flag.String("name", "Clara", "name of user")
)

var JoinClient proto.ChittyChatClient
var wait *sync.WaitGroup

func init() {
	wait = &sync.WaitGroup{}
}

func main() {
	flag.Parse();

	//call correct server
	con, err := grpc.Dial(":"+strconv.Itoa(*serverPort), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Could not connect to port: %d", *serverPort)
	}

	JoinClient = proto.NewChittyChatClient(con)

	//create a new client
	client := &proto.User{
		Name: *userName,
		Timestamp: 0,
	}

	connectToServer(client)

	done := make(chan int)

	wait.Add(1)
	go func() {
		defer wait.Done()

		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			client.Timestamp++

			input := scanner.Text()
			if input == "exit" {
				JoinClient.Leave(context.Background(), &proto.Connect{
						User: client,
						Active: false,
				})

				LeaveMsg := &proto.ChatMessage{
					UserName: client.Name,
					Content: *userName + " left the chat",
					TimeStamp: client.Timestamp,
				}
				_, err := JoinClient.Publish(context.Background(), LeaveMsg)
				if err != nil {
					log.Println(err.Error())
				}
				<-done
			} else {
				message := &proto.ChatMessage{
					UserName: client.Name,
					Content: input,
					TimeStamp: client.Timestamp,
				}

				_, err := JoinClient.Publish(context.Background(), message)
				if err != nil {
					log.Println(err.Error())
				}
			}
		}
	}()

	go func() {
		wait.Wait()
		close(done)
	}()

	<-done

}

func connectToServer(user *proto.User) error {
	var streamError error

	user.Timestamp++

	stream, err := JoinClient.Join(context.Background(), &proto.Connect{
		User: user,
		Active: true,
	})

	if err != nil {
		return fmt.Errorf("failed connection: %v", err)
	}
	wait.Add(1)
	go func(str proto.ChittyChat_JoinClient) {
		defer wait.Done()

		for {
			message, err := str.Recv()

			if message.TimeStamp > user.Timestamp {
				user.Timestamp = message.TimeStamp
			}

			user.Timestamp++

			if err != nil {
				streamError = fmt.Errorf("error message: %v", err)
			}

			fmt.Printf("Lamport time: " + strconv.FormatInt(user.Timestamp, 10) + "\n" + message.Content)
		}
	}(stream)

	return streamError
}

