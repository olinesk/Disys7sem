package main

import (
	"fmt"
	"time"
	"math/rand"
)


func sendingSyn(SYN int, status chan bool) {
	fmt.Printf("Client is sending server this: %d\n", SYN)
	go receiveSyn(SYN, status)
}

func receiveSyn(SYN int, status chan bool) {
	fmt.Printf("Server received this from client: %d\n", SYN)
	status <- true 
}

func sendingSynAck(SYN int, status2 chan bool){
	SYN++
	fmt.Printf("Server is sending client this: %d\n", SYN)
	go receiveSynAck(SYN, status2)
}

func receiveSynAck(SYN int, status2 chan bool){
	fmt.Printf("Client received this from server: %d\n", SYN)
	status2 <- true
}

func transferData(SYN int, data chan string){
	var toBeSend string = String(9)
	fmt.Printf("Client is sending this SYN to server: %d\n", SYN)
	fmt.Printf("Client is sending this data to server: %s\n", toBeSend)
	go receiveData(SYN, data, toBeSend)
}

func receiveData(SYN int, data chan string, toBeSend string){
	data <- toBeSend
	fmt.Printf("Server received this SYN from client: %d\n", SYN)
	fmt.Printf("Server received this data from client: %s\n", toBeSend)
}

// funtions for random string generation
// source: https://www.calhoun.io/creating-random-strings-in-go/
const charset = "abcdefghijklmnopqrstuvwxyz" +
	"ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

var seededRand *rand.Rand = rand.New(
	rand.NewSource(time.Now().UnixNano()))

func StringWithCharset(length int, charset string) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

func String(length int) string {
	return StringWithCharset(length, charset)
}


func start(i int) {
	client := make(chan bool)
	server := make(chan bool)
	SYN := rand.Intn(100)

	go sendingSyn(SYN, client)
	status := <- client

	if status {
		go sendingSynAck(SYN, server)
		status2 := <- server
		if status2 {
			randomLoss := rand.Intn(15)
			if  randomLoss == 1 {
				fmt.Printf("packet loss on SYN: %d\n", SYN) //15% packet loss simulated
			} else {
				data := make(chan string)
				go transferData(SYN+1, data)
				status3 := <- data 
				println(status3)
			}
		}
	} else {
		fmt.Println("Error connection status")
	}

	if i == total {
		done <- true
	}
}

var done chan bool
var total int

func main() {
	done = make(chan bool)
	total = 10

	for i := 0; i <= total; i++ {
		go start(i)
		time.Sleep(time.Second)
	}
	<- done
}
