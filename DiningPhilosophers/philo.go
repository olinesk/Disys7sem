package main

import (
	"fmt"
	"time"
)

var philoDone = make(chan bool)    // Channel to signal when a philosopher is done eating 3 times
var allPhiloDone = make(chan bool) // Channel to signal when all philosophers are done

/*
func randomTime(max int) {
	time.Sleep(time.Millisecond * time.Duration(rand.Intn(max*10)))
}*/

func fork(leftUp, leftDown, rightUp, rightDown chan string) {
	for {
		select {
		case <-leftDown:
			leftUp <- "free fork!"
			<-leftDown

		case <-rightDown:
			rightUp <- "free fork!"
			<-rightDown
		}
	}
}

func philo(id int, leftUp, leftDown, rightUp, rightDown chan string) {
	var philoHowFullo int
	// Run loop until each philo has eaten 3 times
	for philoHowFullo < 3 {

		fmt.Printf("Philo %d thinko\n", id)
		time.Sleep(time.Duration(1+id) * time.Second)

		// Deadlock avoidance: Philosopher 4 picks the right fork first, others pick the left first
		if id == 4 {
			rightDown <- "fork free?"
			<-rightUp

			leftDown <- "fork free?"
			<-leftUp
		} else {
			leftDown <- "fork free?"
			<-leftUp

			rightDown <- "fork free?"
			<-rightUp
		}

		philoHowFullo++
		fmt.Printf("Philo %d eato\n", id)
		time.Sleep(2 * time.Second)

		// Release forks after eating
		leftDown <- "Im done, delicious!"
		rightDown <- "Im done, delicious!"
	}

	fmt.Printf("Philo %d is fullo\n", id)

	philoDone <- true

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

func countPhiloDone() {
	var count int
	for count < 5 {
		<-philoDone
		count++
	}
	allPhiloDone <- true
}

func main() {
	// Creating channels for each philo, 10 because each philo has 2 forks.
	/*channels := make([]chan string, 5) // 2 channels per philo
	for i := 0; i < 5; i++ {
		channels[i] = make(chan string)
	}*/
	// Create channels for each fork (two channels per fork)
	forks := make([][2]chan string, 5)
	for i := range forks {
		forks[i][0] = make(chan string) // left input
		forks[i][1] = make(chan string) // right output
	}

	// Start fork goroutines
	for i := 0; i < 5; i++ {
		go fork(forks[i][0], forks[i][1], forks[(i+1)%5][0], forks[(i+1)%5][1])
	}

	// Start philosopher goroutines
	for i := 0; i < 5; i++ {
		go philo(i, forks[i][0], forks[i][1], forks[(i+1)%5][0], forks[(i+1)%5][1])
	}

	go countPhiloDone()

	<-allPhiloDone

	//time.Sleep(10 * time.Second)

	// Let the main function wait for philosophers to finish eating
	time.Sleep(15 * time.Second)

	fmt.Println("All philo fullo:)")

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
