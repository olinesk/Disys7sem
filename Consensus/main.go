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
	"golang.org/x/tools/go/analysis/passes/defers"
	"google.golang.org/grpc"
	"google.golang.org/grpc/balancer/grpclb/state"
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

var mutex = &sync.Mutex{}

func init() {
	grpcLog = log.NewLogger(os.Stdout, os.Stdout, os.Stdout)
}

func main() {

	// Port
	arg1, _ := strconv.ParseInt(os.Args[1], 32)
	myPort := int32(arg1) + 5001

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Nodes
	n := &Node {
		id: myPort,
		state: StoppedEating,
		time: 0,
		clients: make(map[int32]proto.CakeServiceClient),
		ctx: ctx,
	}

	// Making a listener tcp on myPort
	list, err := net.Listen("tcp", fmt.Sprintf(":%v", myPort))
	if err != nil {
		grpcLog.Errorf("I'm sorry, I failed to listen to port %v :-(", err)
	}
	
	// Make a grpc server
	grpcServer := grpc.NewServer()

	// Register server
	proto.RegisterCakeServiceServer(grpcServer, n)

	go func ()  {
		if err := grpcServer.Serve(list); err != nil {
			grpcLog.Errorf("Failed to serve %v, I'm sorry I let you down :-(", err)
		}
	}()

	for i := 0; i < node_amount; i++ {
		port := int32(5001) + int32(i)
		
		if port == myPort {
			continue
		}

		var conn *grpc.ClientConn
		grpcLog.Infof("Really trying to dial %v\n", port)
		conn, err := grpc.Dial(fmt.Sprintf(":%v", port), grpc.WithInsecure(), grpc.WithBlock())

		if err != nil {
			grpcLog.Errorf("I'm so sorry man, could not connect to: %s", err)
		}

		defer conn.Close()
		c := proto.NewCakeServiceClient(conn)

		n.clients[port] = c
	}

	for {
		if n.state == StoppedEating {
			sleepTime(randomNumberGenerator(2, 7))
			grpcLog.Infof("Try to take crit func dude: %v : %v\n", n.id, n.time)
			n.requestToEat()
		}		
	}
}

type Server struct {
	proto.UnimplementedReceiveServer
	id int32
	timeStamp int32
}


func (s *Server) EatCake (ctx context.Context, Message *proto.Request) (*proto.Reply, error) {
	//some code here

    ack :=  // make an instance of your return type
    return (ack, nil)
}