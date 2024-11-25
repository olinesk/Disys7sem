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

var grpcLog glog.LoggerV2

// Initialize the logger to output logs to the console
func init(){
    grpcLog = glog.NewLoggerV2(os.Stdout, os.Stdout, os.Stout)
}

// FrontEnd struct represents the frontend server that interacts with multiple replication servers
type FrontEnd struct {
    replicationClient []auction.AuctionClient // List of clients connected to replication servers
    auction.UnimplementedAuctionServer // Embedding to implement the AuctionServer interface
}

func main(){
    
}