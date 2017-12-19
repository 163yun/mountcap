package main

import (
	"github.com/golang/glog"
	"github.com/163yun/mountcap"
	"time"
	"fmt"
)

func main() {
	quit := make(chan error)
	mountQueue := mountcap.InitCapturer(10, quit)
	tickquit := time.NewTicker(1 * time.Second)
	defer tickquit.Stop()

	go func(){
		time.Sleep(10 * time.Second)
		quit <- fmt.Errorf("time to end")
		glog.Infoln("call stop")
	}()

	for {
		select {
		case <-tickquit.C:
			glog.Infof("heartbeat!")
		case q := <-quit:
			glog.Infoln(q)
			goto ForEnd
		case mountInfo := <-mountQueue:
			glog.Infof("catch mount changed! New total mount info :%v", mountInfo)
		}
	}
	ForEnd:
	glog.Infof("end.byebye")
}