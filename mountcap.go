package mountcap

/*
#include<stdio.h>
#include<fcntl.h>
#include<poll.h>
#include<stdlib.h>
#include<unistd.h>

int catchMountFileChanged() {
        int mfd = open("/proc/self/mountinfo", O_RDONLY, 0);
        struct pollfd pfd;
        int rv;
        int changes = 0;
        pfd.fd = mfd;
        pfd.events = POLLERR | POLLPRI;
        pfd.revents = 0;
        while ((rv = poll(&pfd, 1, -1)) >= 0) {
                if (pfd.revents & POLLERR) {
                        fprintf(stdout, "Mount points changed. %d.\n", changes++);
                        close(mfd);
                        return 1;
                }
                pfd.revents = 0;
        }
        close(mfd);
        return 0;
}
*/
import "C"

import (
	"github.com/golang/glog"
	"github.com/docker/docker/pkg/mount"
)

const (
	DEFAULT_CHANNEL_LEN = 128
)

type MountInfo struct {
	Infos []*mount.Info
}

func InitCapturer(size int, quit chan error) chan MountInfo {
	if size == 0 {
		size = DEFAULT_CHANNEL_LEN
	}
	eventQueue := make(chan MountInfo, DEFAULT_CHANNEL_LEN)
	go ForeverPoll(eventQueue, quit)
	return eventQueue
}


func ForeverPoll(queue chan MountInfo, quit chan error) {
	var changed int
	for {
		select {
		case <- quit:
			glog.Infof("Existing forever poll for mount changed.")
			return
		default:
			var changedC C.int = C.catchMountFileChanged()
			changed = int(changedC)
			if changed == 1 {
				glog.V(6).Infof("catch a Mount event, will check /proc/self/mountinfo and NewFsInfo again.")
				mounts, err := mount.GetMounts()
				if err != nil {
					glog.Warningf("mountinfo changed but failed to get mounts on node. %s ", err.Error())
				}
				queue <- MountInfo{mounts}
			}
		}
	}
}
