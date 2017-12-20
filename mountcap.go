package mountcap

import (
	"fmt"
	syscall "golang.org/x/sys/unix"
	"os"
	"unsafe"
)

func PollMount(changed chan bool) {
	f, _ := os.Open("/proc/self/mountinfo")
	defer f.Close()
	fdOri := f.Fd()
	var fd *int32 = (*int32)(unsafe.Pointer(&fdOri))
	pollFd := syscall.PollFd{
		Fd:      *fd,
		Events:  syscall.POLLERR | syscall.POLLPRI,
		Revents: 0,
	}
	for {
		ret, _ := syscall.Poll([]syscall.PollFd{pollFd}, -1)
		if ret >= 0 {
			fmt.Printf("pollfd: %+v, get and: %v", pollFd, pollFd.Revents&syscall.POLLERR)
			//if (pollFd.Revents & syscall.POLLERR) == 1 {
			changed <- true
			//}
		}
		pollFd.Revents = 0
	}
}
