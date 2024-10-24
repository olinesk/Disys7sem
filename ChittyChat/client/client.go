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
	clientPort = flag.Int("cPort", 0, "client port number")
	serverPort = flag.Int("sPort", 0, "server port number (must match port used)")
	clientName = flag.String("name", "Clara", "name of client")
)

var JoinClient proto.ChittyChatClient
var wait *sync.WaitGroup

func init() {
	wait = &sync.WaitGroup{}
}

func main() {
	flag.Parse();

	//call correct server
	con, err := grpc.Dial("localhost:"+strconv.Itoa(*serverPort), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Could not connect to port %d", *serverPort)
	}

	JoinClient = proto.NewChittyChatClient(con)

	//create a new client
	client := &proto.User{
		Name: *clientName,
		Timestamp: 0,
	}

	connectToServer(client)


}

func connectToServer() {
	
}

