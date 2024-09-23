package main

import (
	"fmt"
	"sync"
	"time"
)


func sendingSyn(SYN int, status chan bool) {
	fmt.Printf("Client is sending us this: %d\n", SYN)
	receiveSyn(SYN, status)
}

func receiveSyn(SYN int, status chan bool) {
	fmt.Printf("I received this from client: %d\n", SYN)
	status <- true 
}

func sendingSynAck(SYN int, status2 chan bool){
	SYN++
	fmt.Printf("Server is sending us this: %d\n", SYN)
	receiveSynAck(SYN, status2)
}

func receiveSynAck(SYN int, status2 chan bool){
	fmt.Printf("I received this from server: %d\n", SYN)
	status2 <- true
}

func transferData(SYN int, data chan string){
	var toBeSend string = "This is data, very nice yes"
	fmt.Printf("We have this SYN: %d\n", SYN)
	fmt.Printf("We will send this data: %s\n", toBeSend)
}

func receiveData(SYN int, data chan string, toBeSend ){
	data <- toBeSend

}
	

func start() {
	client := make(chan bool)
	server := make(chan bool)

	SYN := 100
	go sendingSyn(SYN, client)
	status := <- client
	if status {
		go sendingSynAck(SYN, server)
		status2 := <- server
		if status2 {
			data := make(chan string)
			go transferData(SYN+1, data)
			status3 := <- data 
			println(status3)
		}
	}
}
