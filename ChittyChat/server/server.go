package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"sync"

	"google.golang.org/grpc"

	proto "ChittyChat/proto"
)

type Connection struct {
	stream proto.ChittyChat_JoinServer
	name   string
	active bool
	error  chan error
}
type Server struct {
	proto.UnimplementedChittyChatServer // Necessary
	name                                string
	port                                int
	users                               map[string]*Connection
}

var (
	port                     = flag.Int("port", 0, "server port number")
	serverLamportTime int64 = 0
	users                    = make(map[string]*Connection)
)

func main() {

	//set up log
	f, err := os.OpenFile("golang-demo.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)

	if err != nil {
		log.Fatalf("Could not open file: %v", err)
	}
	defer f.Close()
	log.SetOutput(f)

	// Parses the flags and gets port from commandline
	flag.Parse()

	//Creates server struct
	server := &Server{
		name:  "serverName",
		port:  *port,
		users: users,
	}

	// Starting the server
	go startServer(server)

	//Keeping the server running until quit
	for {

	}
}

func startServer(server *Server) {

	// Create a new grpc server
	grpcServer := grpc.NewServer()

	// Make the server listen at the given port (convert int port to string)
	listener, err := net.Listen("tcp", ":"+strconv.Itoa(server.port))

	if err != nil {
		log.Fatalf("Could not create the server %v", err)
	}
	log.Printf("Started server at port: %d\n", server.port)

	// Register the grpc server and serve its listener
	proto.RegisterChittyChatServer(grpcServer, server)
	serveError := grpcServer.Serve(listener)
	if serveError != nil {
		log.Fatalf("Could not serve listener")
	}
}

func (s *Server) Join(in *proto.Connect, stream proto.ChittyChat_JoinServer) error {

	con := &Connection{
		stream: stream,
		name:   in.User.Name,
		active: true,
		error:  make(chan error),
	}

	log.Printf("Participant " + in.User.Name + " joined Chitty-Chat at Lamport time: " + "%v", serverLamportTime)

	s.users[in.User.Name] = con

	//Participants in chat get notified when new user joins
	userJoinedChat := &proto.ChatMessage{
		UserName:  in.User.Name,
		Content:   in.User.Name + " joined the chat",
		TimeStamp: in.User.Timestamp,
	}

	if serverLamportTime < int64(userJoinedChat.TimeStamp) {
		serverLamportTime = int64(userJoinedChat.TimeStamp)
	}

	s.Publish(con.stream.Context(), userJoinedChat)

	return <-con.error
}

func (s *Server) Publish(context context.Context, in *proto.ChatMessage) (*proto.Close, error) {

	//Making sure goroutines gets to finish
	wait := sync.WaitGroup{}

	//Used to know if goroutines are finished
	done := make(chan int)
	if serverLamportTime < int64(in.TimeStamp) {
		serverLamportTime = int64(in.TimeStamp)
	}

	log.Printf(in.UserName + " sent message: \"" + in.Content + "\", at Lamport time: " + "%v", serverLamportTime)

	serverLamportTime++

	for _, con := range s.users {
		wait.Add(1)

		fmt.Printf("Server Time: %v " + "\n", serverLamportTime)
		go func(content *proto.ChatMessage, con *Connection) {
			serverLamportTime++
			defer wait.Done()

			if con.active {
				msgToBeSent := &proto.ChatMessage{
					UserName:  in.UserName,
					Content:   in.UserName + ": " + content.Content,
					TimeStamp: serverLamportTime,
				}

				msgToBeSent.Content += "\n"

				err := con.stream.Send(msgToBeSent)

				if err != nil {
					log.Printf("Participant " + content.UserName + " left Chitty-Chat at Lamport time: " + "%v", serverLamportTime)
					con.active = false
					con.error <- err
				}
			}
		}(in, con)
	}

	go func() {
		wait.Wait()
		close(done)
	}()

	<-done
	return &proto.Close{}, nil
}

func (s *Server) Leave(in *proto.Connect, stream proto.ChittyChat_LeaveServer) error {

	for name := range s.users {
		if name == in.User.Name {
			delete(s.users, name)
			log.Printf("Participant " + in.User.Name + " left Chitty-Chat at Lamport time: " + "%v", serverLamportTime)
		}
	}
	return nil
}
