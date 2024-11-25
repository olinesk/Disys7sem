package main

import (
	"context"
	"fmt"
	"net"
	"os"
	"strconv"
	"time"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	auction "github.com/olinesk/Disys7sem/Auction/proto"
)

// FrontEnd struct represents the frontend server that interacts with multiple replication servers
type FrontEnd struct {
    replicationClients []auction.AuctionClient // List of clients connected to replication servers
    auction.UnimplementedAuctionServer // Embedding to implement the AuctionServer interface
}

func main() {
    // Set up log
	f, err := os.OpenFile("golang-demo.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)

	if err != nil {
		log.Fatalf("Could not open file: %v", err)
	}
	defer f.Close()
	log.SetOutput(f)

	//Setup of frontEnd clientside
	frontEnd := &FrontEnd{}

	numOfRepServs := 3
	for i := 0; i < numOfRepServs; i++ {
		repServPort := 5010 + (10 * i)

		conn, err := grpc.Dial(fmt.Sprintf(":%v", repServPort), grpc.WithTimeout(3*time.Second), grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
		if err != nil {
			log.Fatalf("Could not listen to port: %v, error: %v \n", repServPort, err)
			continue
		}
		defer conn.Close()

		replicationClient := auction.NewAuctionClient(conn)

		frontEnd.replicationClients = append(frontEnd.replicationClients, replicationClient)
	}

	//Setup of frontEnd serverside
	arg1, _ := strconv.ParseInt(os.Args[1], 10, 32)
	ownPort := 8010 + (arg1 * 10)
	
	listener, err := net.Listen("tcp", fmt.Sprintf(":%v", ownPort))
	if err != nil {
		log.Fatalf("Error creating the server %v \n", err)
	}

	log.Printf("Starting the server on port: %d \n", ownPort)

	grpcServer := grpc.NewServer()
	auction.RegisterAuctionServer(grpcServer, frontEnd)
	grpcServer.Serve(listener)
}

func (frontEnd *FrontEnd) MakeBid(ctx context.Context, bid *auction.Bid) (*auction.Ack, error) {
	ackChannel := make(chan *auction.Ack)

	log.Printf("Bid: %d from bidder: %d \n", bid.Amount, bid.BidderId)

	for index, replicationClient := range frontEnd.replicationClients {
		go func (client auction.AuctionClient, index int) {
			ack, err := client.MakeBid(context.Background(), bid)
			if err != nil {
				frontEnd.removeReplicationManager(index)
				return
			}
			ackChannel <- ack
		} (replicationClient, index)
	}

	return <-ackChannel, nil
}

func (frontEnd *FrontEnd) GetStatus(ctx context.Context, req *auction.StatusReq) (*auction.Status, error) {
	statusChannel := make(chan *auction.Status)

	log.Println("Status requested.")

	for index, replicationClient := range frontEnd.replicationClients {
		go func (client auction.AuctionClient, index int) {
			status, err := client.GetStatus(context.Background(), req)
			if err != nil {
				frontEnd.removeReplicationManager(index)
				return
			}
			statusChannel <- status
		} (replicationClient, index) 
	}
	
	log.Println("Returning requested status.")
	return <-statusChannel, nil
}

func (frontEnd *FrontEnd) removeReplicationManager(index int) {
	log.Printf("A Replication Manager at index: %d is not responding \n", index)
	frontEnd.replicationClients = append(frontEnd.replicationClients[:index], frontEnd.replicationClients[index+1:]...)
}