package main

import (
	"fmt"
	"github.com/163yun/mountcap"
	"github.com/golang/glog"
	"time"
)

func main() {
	pollChanged := make(chan bool, 128)
	tick := time.NewTicker(1 * time.Second)
	defer tick.Stop()
	quitCh := make(chan error, 1)
	go mountcap.PollMount(pollChanged, quitCh)

	count := 1

	for {
		select {
		case ch := <-pollChanged:
			glog.Infoln("Got a change:", ch)
			if ch {
				glog.Infoln("mount changed!")
			}
		case <-tick.C:
			glog.Infoln("heartbeat!")
			if count == 10 {
				quitCh <- fmt.Errorf("ready to exit")
				glog.Infoln("called exit")
				return
			}
		}
		count++
	}
}
