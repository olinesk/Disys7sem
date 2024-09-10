package main

import (
	"fmt"
	"math/rand"
	"time"
)

func randomTime(max int) {
	time.Sleep(time.Millisecond * time.Duration(rand.Intn(max*10)))
}

func fork(leftUp, leftDown, rightUp, rightDown chan string){
	for {
		select {
		case <- leftDown:
			leftUp <- "free fork!"
			<- leftDown
		
		case <- rightDown:
			rightUp <- "free fork!"
			<- rightDown
		}
	}
}

func philo(id int, leftUp, leftDown, rightUp, rightDown chan string) {
	var philoHowFullo int
	// Run loop until each philo has eaten 3 times
	for philoHowFullo < 3 {
		fmt.Printf("Philo %d thinko\n", id)
		randomTime(2)
		
		// if-statement avoiding deadlocks
		if id == 4 {
			rightDown <- "fork free?"
			<- rightUp

			leftDown <- "fork free?"
			<- leftUp
		} else {
			leftDown <- "fork free?"
			<- leftUp

			rightDown <- "fork free?"
			<- rightUp
		}

		philoHowFullo++
		fmt.Printf("Philo %d eato\n", id)
		randomTime(3)

		leftDown <- "Im done, delicious!"
		rightDown <- "Im done, delicious!"
	}
	
	fmt.Printf("Philo %d is fullo\n", id)
	
	
	/*
	for i := 0; i < 10; i++ {

		leftDown <- "fork free?"
		<- leftUp 

		rightDown <- "fork free?"
		<- rightUp

		fmt.Println("Im eating with free forks, yummyyy")

		leftDown <- "Im done, delicious!"
		rightDown <- "Im done, delicious!"
	}

	fmt.Println("This philo is fullo")*/
}

func main() {


}






/*
var p1, p2, p3, p4, p5 string 
var f1, f2, f3, f4, f5 bool

func philo1(){
	if f1 == true && f5 == true {
		fmt.Println("eating")
	} else {fmt.Println("thinking")}
}

func philo2(){
	if f1 == true && f2 == true {
		fmt.Println("eating")
	} else {fmt.Println("thinking")}
}

func philo3(){
	if f2 == true && f3 == true {
		fmt.Println("eating")
	} else {fmt.Println("thinking")}
}

func philo4(){
	if f3 == true && f4 == true {
		fmt.Println("eating")
	} else {fmt.Println("thinking")}
}

func philo5(){
	if f4 == true && f5 == true {
		fmt.Println("eating")
	} else {fmt.Println("thinking")}
}
*/




