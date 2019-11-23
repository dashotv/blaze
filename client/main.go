package main

import (
	"fmt"
	"os"
	"time"

	"github.com/nats-io/nats.go"

	"github.com/dashotv/mercury"
)

func main() {
	m, err := mercury.New("blaze", nats.DefaultURL)
	if err != nil {
		panic(err)
	}

	fmt.Println("starting receiver...")
	channel := make(chan *mercury.Message, 5)
	m.Receiver("blaze", channel)

	for {
		select {
		case r := <-channel:
			fmt.Printf("received: %#v\n", r)
		case <-time.After(30 * time.Second):
			fmt.Println("timeout")
			os.Exit(0)
		}
	}
}
