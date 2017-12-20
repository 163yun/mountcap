package main

import (
	"fmt"
	"github.com/163yun/mountcap"
	"time"
)

func main() {
	pollChanged := make(chan bool, 128)
	tick := time.NewTicker(1 * time.Second)
	defer tick.Stop()
	go mountcap.PollMount(pollChanged)

	for {
		select {
		case ch := <-pollChanged:
			fmt.Println("Got a change:", ch)
			if ch {
				fmt.Println("mount changed!")
			}
		case <-tick.C:
			fmt.Println("heartbeat!")
		}
	}
}
