package main

import (
	"context"
	"fmt"
	"net"
	"os"
	"strconv"
	"sync"
	"time"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/durationpb"

	auction "github.com/olinesk/Disys7sem/Auction/proto"
)

type ReplicationManager struct {
	auction.UnimplementedAuctionServer
	endTime 	time.Time
	highestBid 	*auction.Bid
	bidLock 	sync.Mutex
}

func main() {
	arg1, _ := strconv.ParseInt(os.Args[1], 10, 32)
	ownPort := 5010 + (arg1 * 10)

	replicationManager := &ReplicationManager{highestBid: &auction.Bid{BidderId: -1, Amount: 0}}

	//Setup of API/FrontEnd serverside

	listener, err := net.Listen("tcp", fmt.Sprintf(":%v", ownPort))
	if err != nil {
		log.Fatalf("Error while creating server: %v \n", err)
	}

	log.Printf("Starting server at port: %v", ownPort)

	grpcServer := grpc.NewServer()
	auction.RegisterAuctionServer(grpcServer, replicationManager)
	grpcServer.Serve(listener)
}

func (replicationManager *ReplicationManager) MakeBid(ctx context.Context, bid *auction.Bid) (*auction.Ack, error) {
	if replicationManager.endTime.IsZero() {
		auctionEndTime := bid.TimeStamp.AsTime().Add(time.Minute)
		replicationManager.endTime = auctionEndTime
		log.Printf("Setting end-time as: %v \n", auctionEndTime)
	}

	if replicationManager.endTime.Before(time.Now()) {
		log.Printf("Bidder %d is trying to bid after the auction ended! \n", bid.BidderId)
		return &auction.Ack{}, nil
	}

	if replicationManager.setHighestBid(bid) {
		log.Printf("Bidder %d has the highest bid of: %d monnays \n", bid.BidderId, bid.Amount)
	} else {
		log.Printf("The bid %d from bidder %d was too low. \n", bid.Amount, bid.BidderId)
	}

	return &auction.Ack{}, nil
}

func (replicationManager *ReplicationManager) GetStatus (ctx context.Context, req *auction.StatusReq) (*auction.Status, error) {
	log.Printf("Bidder: %d is requesting a status of the auction. \n", req.BidderId)

	highestBid := replicationManager.getHighestBid()
	timeLeft := time.Until(replicationManager.endTime)

	status := &auction.Status{
		TimeLeft: durationpb.New(timeLeft),
		HighestBid: highestBid.Amount,
		BidderId: highestBid.BidderId,
	}

	return status, nil
}

func (replicationManager *ReplicationManager) setHighestBid (bid *auction.Bid) bool {
	replicationManager.bidLock.Lock()
	defer replicationManager.bidLock.Unlock()
	if bid.Amount < replicationManager.highestBid.Amount {
		return false
	}
	replicationManager.highestBid = bid
	return true
}

func (replicationManager *ReplicationManager) getHighestBid() *auction.Bid {
	replicationManager.bidLock.Lock()
	defer replicationManager.bidLock.Unlock()
	return replicationManager.highestBid
}