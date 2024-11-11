package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"net"
	"os"
	"strconv"
	"sync"
	"time"

	"golang.org/x/net/context/ctxhttp"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	proto "github.com/olinesk/Disys7sem/proto"
)

var grpcLog log.Logger
var node_amount = 3

type State int

const (
	StoppedEating State = iota
	WantsToEat
	Eating
)

type Node struct {
	proto.UnimplementedCakeServiceServer
	id 				int32
	state 			State
	time 			int32
	timeLock		sync.Mutex
	requestTime		int32
	eatWG			sync.WaitGroup
	clients 		map[int32]proto.CakeServiceClient
	ctx				context.Context	
}

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

func (s *Server) EatCake (cnx context.Context, Message *proto.Request) (*proto.Reply, error) {
	//some code here

    ack :=  // make an instance of your return type
    return (ack, nil)
}