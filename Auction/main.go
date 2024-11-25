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

	auction "github.com/olinesk/Disys7sem/Auction/proto"
)



// Initialize the logger to output logs to the console
func init(){
    // Set up log
	f, err := os.OpenFile("golang-demo.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)

	if err != nil {
		log.Fatalf("Could not open file: %v", err)
	}
	defer f.Close()
	log.SetOutput(f)
}

// FrontEnd struct represents the frontend server that interacts with multiple replication servers
type FrontEnd struct {
    replicationClient []auction.AuctionClient // List of clients connected to replication servers
    auction.UnimplementedAuctionServer // Embedding to implement the AuctionServer interface
}

func main(){
    
}