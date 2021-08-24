package main

import (
	"fmt"
	"math/rand"
	"runtime"
	"sync"
	"time"
)

func fanOutExample() {
	fmt.Println("Number of active go routines ", runtime.NumGoroutine())
	c1, c2 := make(chan int), make(chan int)
	go populateEvents(c1)
	//go fanInOut(c1, c2)

	go fanInOutWithThrottle(c1, c2)

	//go populateEventsForEver(c1)
	//go fanInOutWithThrottleInfiniteLoop(c1, c2)

	for v := range c2 {
		fmt.Println("processed: ", v)
	}
	fmt.Println("Exiting")
}

func populateEvents(c chan<- int) {
	for i := 0; i < 21; i++ {
		c <- i
	}
	close(c)
}

func populateEventsForEver(c chan<- int) {
	count := 0
	for {
		c <- count
		count++
	}
}

func fanInOut(c1 <-chan int, c2 chan<- int) {
	var wg sync.WaitGroup
	for v := range c1 {
		wg.Add(1)
		go func(i int) {
			c2 <- encodeVideo(i)
			wg.Done()
		}(v)
	}
	wg.Wait()
	close(c2)
}

func fanInOutWithThrottle(c1 <-chan int, c2 chan<- int) {
	var wg sync.WaitGroup
	const MaxGoRoutines = 5
	wg.Add(MaxGoRoutines)
	for i := 0; i < MaxGoRoutines; i++ {
		go func() {
			for v := range c1 {
				//for v := 0; v < 10; v++ {
				func(i int) {
					c2 <- encodeVideo(i)
				}(v)
			}
			wg.Done()
		}()
	}
	wg.Wait()
	close(c2)
}

func fanInOutWithThrottleInfiniteLoop(c1 <-chan int, c2 chan<- int) {
	//var wg sync.WaitGroup
	const MaxGoRoutines = 5
	//wg.Add(MaxGoRoutines)
	for i := 0; i <= MaxGoRoutines; i++ {
		go func() {
			for v := range c1 {
				func(i int) {
					c2 <- encodeVideo(i)
				}(v)
			}
			//wg.Done()
		}()
	}
	//wg.Wait()
	//close(c2)
}

func encodeVideo(i int) int {
	fmt.Println("Number of active go routines ", runtime.NumGoroutine())
	time.Sleep(time.Millisecond * time.Duration(rand.Intn(100)))
	return i + rand.Intn(1000)
}
