# goroutines
# Do not communicate by sharing memory; instead, share memory by communicating
## Contain working examples of some of the typical multi-threaded uses cases

## Channels:
- Channel another type in Go
- Better ways of writing concurrent code in go
- Concurrent is design pattern that allows code to run in parallel and take advantage of CPU's
- Blocks untill send and receive happens at the same time
- Works only with multiple Go routines with single go routines, program will be blocked
```
func main() {
	c := make(chan int)
	c <- 42
	fmt.Println(<-c)
}
```
but this will work 
```
func main() {
	c := make(chan int)
	go func(){c <- 42}()
	fmt.Println(<-c)
}
```
- Buffered channel is another way to use channels
- Channels allows to pass values between go-routines
- when channel is closed `ok` will be set to false

### Directional channel
- send only channle
    - can only write to channel
    - chan<-
- receive only channle
    - can only read from channel
    - <-chan
- One can read data from channel untill `channel` is closed

### Concurrent programming design patterns
- Channel blocks and there should be writer to the channel and corresponding reader as well
    - Range clause iterate over the channel until it's closed
- Using Select clause to read from channels 
- Fan in
    - Break a piece of work into multiple Go routines and once done use a different `channel`
    to colate the result
- Fan out
    - Run code in parallel, break large chunk of code and break them to work one at a time

### Context
- tool to be used with the concurrent design patterns, 
    - if a process launches set of go-routines
        then all the go-routines launched through the process should gets cancelled.
    - Helps to pass around variable related to request
    - In Go Servers, each incoming request is handled in its own goroutines. Application can start its own go routines, might have to pass some incoming data, cancellation signals etc. This is done
    by using context.
    - context contains request scoped values
