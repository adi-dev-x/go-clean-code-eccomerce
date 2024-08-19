package main

import (
	"fmt"
	"time"
)

// Function that simulates a goroutine sending messages
func worker(id int, ch chan<- string) {
	for i := 0; i < 3; i++ {
		time.Sleep(time.Millisecond * 500) // Simulate work
		ch <- fmt.Sprintf("Message %d from worker %d", i, id)
	}
}

func main() {
	// Buffered channel with a capacity of 5
	// bufferedCh := make(chan string, 9)
	bufferedCh := make(chan string)

	// Launch worker goroutines
	go worker(1, bufferedCh)
	go worker(2, bufferedCh)
	go worker(3, bufferedCh)

	// Receive messages from the buffered channel
	fmt.Println("Receiving from buffered channel:")
	for i := 0; i < 9; i++ { // Expecting 9 messages in total
		msg := <-bufferedCh
		fmt.Println("Received:", msg)
	}

	// Close the channel
	close(bufferedCh)
}
