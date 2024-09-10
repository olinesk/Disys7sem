package main

import (
	"fmt"
	"time"
)

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

func main() {


}





