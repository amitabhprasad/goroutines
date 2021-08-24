package main

import (
	"time"
	"sync"
	"context"
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"runtime"
)

const workerPoolSize=4

func main(){
	//producerConsumer()
	//fanInExample()
	//fanOutExample()
	fanInExampleRob()
}
func producerConsumer(){
	var wg sync.WaitGroup
	wg.Add(workerPoolSize)
	fmt.Printf("Starting producer consumer....with num of routines %v\n ",runtime.NumGoroutine())

	ctx, cancel := context.WithCancel(context.Background())
	consumer := Consumer{
		ingestChan : make(chan string),
		jobsChan: make(chan string),
	}
	producer := Producer{
		callback : consumer.callbackFunc,
		producerID: "P101",
	}
	producer2 := Producer{
		callback : consumer.callbackFunc,
		producerID: "P201",
	}
	
	go producer.start(100)
	fmt.Println("Num of go routines with producer 1 ",runtime.NumGoroutine())
	go producer2.start(100)
	fmt.Println("Num of go routines with producer 2 ",runtime.NumGoroutine())
	go consumer.startConsumer(ctx)
	fmt.Println("Num of go routines with consumerStarting ",runtime.NumGoroutine())

	for i := 0; i < workerPoolSize; i++ {
		go consumer.workerFunc(&wg, i)
	}

	// Handle sigterm and await termChan signal
	termChan := make(chan os.Signal)
	signal.Notify(termChan, syscall.SIGINT, syscall.SIGTERM)

	<-termChan         // Blocks here until interrupted

	// Handle shutdown
	fmt.Println("*********************************\nShutdown signal received\n*********************************")
	fmt.Println("Num of go routines ",runtime.NumGoroutine())
	cancel()
	wg.Wait()
	// go func() {
	// 	time.Sleep(5 * time.Second)
	// 	fmt.Println("Enough time close this now ")
	// 	cancel()
	// }()
	fmt.Println("Num of go routines at the end ",runtime.NumGoroutine())
	fmt.Println("All workers done, shutting down!")
}

type Producer struct{
	callback func(event string)
	producerID string
}
 func (p Producer) start(eventIndex int){
	fmt.Println("Starting producer ....",runtime.NumGoroutine())
	 for {
		 event := fmt.Sprintf("%v#%v",p.producerID,eventIndex)
		 p.callback(event)
		 fmt.Println("Produced event ",event)
		 eventIndex++
		 time.Sleep(time.Millisecond * 100)
	 }
 }
type Consumer struct{
	ingestChan chan string
	jobsChan chan string
}

func (c Consumer) callbackFunc(event string){
	c.ingestChan <- event
}

func (c Consumer) startConsumer(ctx context.Context){
	fmt.Println("Starting  consumer.... num of go routines now ",runtime.NumGoroutine())
	for{
		select{
		case <-ctx.Done():
			fmt.Println("Consumer received cancellation signal, closing jobsChan!")
			close(c.jobsChan)
			fmt.Println("Consumer closed jobsChan")
			return
		case job := <-c.ingestChan:
			c.jobsChan<- job
		}
	}
}

func (c Consumer) workerFunc(wg *sync.WaitGroup,index int){
	defer wg.Done()
	fmt.Printf("Worker %d starting number of routines ( %v )\n", index, runtime.NumGoroutine())
	for eventIndex := range c.jobsChan{
		// simulate work  taking between 1-3 seconds
		fmt.Printf("Worker %d started job %v\n", index, eventIndex)
		time.Sleep(time.Millisecond * time.Duration(1000+rand.Intn(2000)))
		fmt.Printf("Worker %d finished processing job %v\n", index, eventIndex)
	}
	fmt.Printf("Worker %d interrupted\n", index)
}