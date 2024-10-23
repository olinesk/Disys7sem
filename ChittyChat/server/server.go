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
	id string
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
	fmt.Println(".:server is starting:.")

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
		id: in.User.Id,
		name: in.User.Name,
		active: true,
		error: make(chan error),
	}

	s.users[in.User.Name] = con

	//Participants in chat get notified when new user joins
	userJoinedChat := &proto.ChatMessage{
		UserName: in.User.Name,
		Content: in.User.Name + " joined the chat ",
		TimeStamp: in.User.Utime,
		UserID: in.User.Id,
	}
}