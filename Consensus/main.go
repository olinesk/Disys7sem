package main

import (
	"context"
	"fmt"
	"log"
	"math/rand/v2"
	"net"
	"os"
	"strconv"
	"sync"
	"time"


	"google.golang.org/grpc"

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

	// Set up log
	f, err := os.OpenFile("golang-demo.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)

	if err != nil {
		log.Fatalf("Could not open file: %v", err)
	}
	defer f.Close()
	log.SetOutput(f)

	// Port
	arg1, _ := strconv.ParseInt(os.Args[1], 10, 32)
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
		grpcLog.Fatalf("I'm sorry, I failed to listen to port %v :-(", err)
	}
	
	// Make a grpc server
	grpcServer := grpc.NewServer()

	// Register server
	proto.RegisterCakeServiceServer(grpcServer, n)

	go func ()  {
		if err := grpcServer.Serve(list); err != nil {
			grpcLog.Fatalf("Failed to serve %v, I'm sorry I let you down :-(", err)
		}
	}()

	for i := 0; i < node_amount; i++ {
		port := int32(5001) + int32(i)
		
		if port == myPort {
			continue
		}

		var conn *grpc.ClientConn
		grpcLog.Printf("Really trying to dial %v\n", port)
		conn, err := grpc.Dial(fmt.Sprintf(":%v", port), grpc.WithInsecure(), grpc.WithBlock())

		if err != nil {
			grpcLog.Fatalf("I'm so sorry man, could not connect to: %s", err)
		}

		defer conn.Close()
		c := proto.NewCakeServiceClient(conn)

		n.clients[port] = c
	}

	for {
		if n.state == StoppedEating {
			sleepTime(rand.IntN(10))
			grpcLog.Printf("Try to take crit func dude: %v : %v\n", n.id, n.time)
			n.requestToEat()
		}		
	}
}

func (n *Node) EatCake(ctx context.Context, req *proto.Request) (*proto.Reply, error) {
	
	if n.time < req.Time {
		n.time = req.Time
	}
	n.time++

	if n.state == Eating || (n.state == WantsToEat && n.requestIsBefore(req)) {
		n.eatWG.Wait()
	}

	rep := &proto.Reply{Id: n.id, TimeStamp: n.time}
	return rep, nil
}

func (n *Node) requestToEat() {
	n.updateClock(n.time)
	n.requestTime = n.time
	n.eatWG.Add(1)
	n.state = WantsToEat

	req := proto.Request{
		Time: n.requestTime,
		Id: n.id,
	}

	var wG sync.WaitGroup

	for _, client := range n.clients {
		wG.Add(1)
		go func (c proto.CakeServiceClient)  {
			defer wG.Done()

			reply, _ := c.EatCake(n.ctx, &req)

			n.updateClock(reply.TimeStamp)
			grpcLog.Printf("Reply from %d at time: %v", reply.Id, reply.TimeStamp)
		}(client)
	}

	wG.Wait()
	n.state = Eating
	n.Eat()
}

func (n *Node) Eat() {

	defer n.noMoreCakeForMe()

	grpcLog.Printf("Eating... at time: %v", n.time)

	sleepTime(rand.IntN(10))

	grpcLog.Printf("I'm done eating momma, at time: %v", n.time)
}

func (n *Node) noMoreCakeForMe() {
	n.state = StoppedEating
	n.eatWG.Done()
}

func (n *Node) updateClock(time int32) {
	n.timeLock.Lock()
	defer n.timeLock.Unlock()

	if n.time < time {
		n.time = time
	}

	n.time++
}

func (n *Node) requestIsBefore(req *proto.Request) bool {
	
	if req.Time < n.time || (req.Time == n.time && req.Id < n.time) {
		return true
	}
	return false
}

func sleepTime(n int) {
	time.Sleep(time.Duration(n) * time.Second)
}
