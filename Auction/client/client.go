package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strconv"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"
	"google.golang.org/grpc/credentials/insecure"

	auction "github.com/olinesk/Disys7sem/Auction/proto"
)

type Bidder struct {
	client auction.AuctionClient
	id int32
}

func main() {
	id, _ := strconv.ParseInt(os.Args[1], 10, 32)

	frontEndPort := 8010 + (id * 10)
	log.Printf("Trying to dial %d", frontEndPort)
	conn, err := grpc.Dial(fmt.Sprintf(":%v", frontEndPort), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Could not connect to port: %d", frontEndPort)
	}
	log.Printf("Dialed %d", frontEndPort)

	bidder := &Bidder{id: int32(id), client: auction.NewAuctionClient(conn)}

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		amount, err := strconv.ParseInt(scanner.Text(), 10, 32)
		
		if err != nil {
			continue
		}

		if amount == 0 {
			req := auction.StatusReq{BidderId: bidder.id}
			status, err := bidder.client.GetStatus(context.Background(), &req)
			if err != nil {
				log.Fatalf(err.Error())
				continue
			}

			timeLeft := status.TimeLeft.AsDuration()

			if timeLeft <= 0 {
				log.Println("Auction is finished.")
			} else {
				log.Printf("Time left of auction: %v \n", timeLeft)
			}

			log.Printf("Bidder: %d has the highest bid of: %d moneys", status.BidderId, status.HighestBid)
			continue
		}

		bid := auction.Bid {
			BidderId: bidder.id,
			Amount: int32(amount),
			TimeStamp: timestamppb.Now(),
		}

		_, bidErr := bidder.client.MakeBid(context.Background(), &bid)
		if bidErr != nil {
			log.Fatalf(err.Error())
			continue
		}

		log.Printf("Bid: %d - sent! \n", amount)
	}
} 


