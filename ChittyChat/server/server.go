package main

import (
	"context"
	"flag"
	"log"
	"fmt"
	"strconv"
	"net"
	"os"
	"sync"

	"google.golang.org/grpc"

	proto "ChittyChat/proto"
)

type Connection struct {
	stream proto.ChittyChat_JoinServer
	name string
	active bool
	error chan error
}
type Server struct {
	proto.UnimplementedChittyChatServer // Necessary
	name                             string
	port                             int
	users 							 map[string]*Connection
}

var (
port = flag.Int("port", 0, "server port number")
	serverLamportClock int64 = 0;
	users = make(map[string]*Connection)
)


func main() {

	//set up log
	f, err := os.OpenFile("golang-demo.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	
	if err != nil {
		log.Fatalf("Could not open: %v", err)
	}
	defer f.Close()
	log.SetOutput(f)

	// Parses the flags and gets port from commandline
	flag.Parse()

	//Creates server struct
	server := &Server {
		name: "serverName",
		port: *port,
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

	con := &Connection {
		stream: stream,
		name: in.User.Name,
		active: true,
		error: make(chan error),
	}

	s.users[in.User.Name] = con

	//Participants in chat get notified when new user joins
	userJoinedChat := &proto.ChatMessage{
		UserName: in.User.Name,
		Content: in.User.Name + " joined the chat ",
	}

	if serverLamportClock < int64(userJoinedChat.TimeStamp) {
		serverLamportClock = int64(userJoinedChat.TimeStamp)
	}

	s.Publish(con.stream.Context(), userJoinedChat)

	log.Printf("User " + in.User.Name + " joined the chat at " + "%v", serverLamportClock);

	return <-con.error
}

func (s *Server) Publish(context context.Context, in *proto.ChatMessage)(*proto.Close, error) {
	
	//Making sure goroutines gets to finish
	wait := sync.WaitGroup{}; 
	done := make(chan int) 
	if serverLamportClock < int64(in.TimeStamp) {
		serverLamportClock = int64(in.TimeStamp)
	}

	serverLamportClock++

	for _, con := range s.users {
		wait.Add(1)

		fmt.Printf("Server Time: %v", serverLamportClock)
		go func (content *proto.ChatMessage, con *Connection) {
			serverLamportClock++
			defer wait.Done()

			if con.active {
				msgToBeSent := &proto.ChatMessage{
				UserName:	in.UserName,
				Content:	in.UserName + ": " + content.Content,
				TimeStamp:	serverLamportClock,
				}

				msgToBeSent.Content += "\n"

				err := con.stream.Send(msgToBeSent)

				if err != nil {
					log.Printf("User: " + content.UserName + " left the chat at " + "%v", serverLamportClock);
					con.active = false
					con.error <- err
				}
			}
		} (in, con)
	}

	log.Printf(in.UserName + " sent message " + in.Content + " at: " + "%v", serverLamportClock);

	go func() {
		wait.Wait()
		close(done)
	}()
	
	<- done
	return &proto.Close{}, nil
}

func (s *Server) Leave(in *proto.Connect, stream proto.ChittyChat_LeaveServer) error {

	for name := range s.users {
		if name == in.User.Name {
			delete(s.users, name)
			log.Printf("User: " + in.User.Name + " left the chat at " + "%v", serverLamportClock);
		}
	}
	return nil
}



