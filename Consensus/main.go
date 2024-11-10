package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"sync"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	proto "github.com/olinesk/Disys7sem/proto"
)

func main() {

	// Make a grpc server
	grpcServer := grpc.NewServer()

	// Make an instance of your struct
	server := &Server {

	}

	// Register your server struct
	proto.RegisterReceiveServer(grpcServer, server)

	// Start serving
	grpcServer.Serve(list)
}

type Server struct {
	proto.UnimplementedReceiveServer
	id int32
	timeStamp int32


}

type Connection struct {
	server proto.ReceiveClient
	serverConnection *proto.ClientConn
}

func (s *Server) GetRequest (cnx context.Context, Message *proto.Request) (*proto.Reply, error) {
	//some code here
    ...
    ack :=  // make an instance of your return type
    return (ack, nil)
}