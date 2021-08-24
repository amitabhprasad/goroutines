package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

func fanInExample() {
	eve := make(chan int)
	odd := make(chan int)
	fanIn := make(chan int)
	go processEvents(eve, odd)
	go collate(eve, odd, fanIn)

	for v := range fanIn {
		fmt.Println(v)
	}
}

func processEvents(o, e chan<- int) {
	for i := 0; i < 10; i++ {
		if i%2 == 0 {
			e <- i
		} else {
			o <- i
		}
	}
	close(e)
	close(o)
}

func collate(o, e <-chan int, fanIn chan<- int) {
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		for v := range o {
			fanIn <- v
		}
		wg.Done()
	}()

	go func() {
		for v := range e {
			fanIn <- v
		}
		wg.Done()
	}()
	wg.Wait()
	close(fanIn)
}

func fanInExampleRob() {
	c := fanIn(boring("Joe"), boring("Ann"))
	for i := 0; i < 10; i++ {
		fmt.Println(<-c)
	}
	fmt.Println("You guys are all boring...Let me try with range clause")
	fmt.Println("Iterating using range clause...this will go on for ever")
	for v := range c {
		fmt.Println(v)
	}
}

func boring(s string) <-chan string {
	c := make(chan string)
	go func() {
		for i := 0; ; i++ {
			c <- fmt.Sprintf("%s %d", s, i)
			time.Sleep(time.Duration(rand.Intn(1e3)) * time.Millisecond)
		}
	}()
	return c
}

func fanIn(input1, input2 <-chan string) <-chan string {
	c := make(chan string)
	go func() {
		for {
			c <- <-input1
		}
	}()
	go func() {
		for {
			c <- <-input2
		}
	}()
	return c
}
