package main

import (
	"syscall"
	"log"
	"github.com/golang.org/x/sys/unix"
)

func main()  {
	fd, err := syscall.Open("/tmp/foo/a.text", syscall.O_RDWR, 0666)
	if err != nil {
		log.Printf("open error")
		return
	}

	events := unix.POLLIN | unix.POLLOUT
	revents := unix.POLLIN | unix.POLLOUT | unix.POLLPRI | unix.POLLRDHUP | unix.POLLERR | unix.POLLHUP  | unix.POLLNVAL  = 0x20

	var pollfds [] unix.PollFd
	pollfd := unix.PollFd{
		Fd:int32(fd),
		Events:int16(events),
		Revents:int16(revents),
	}
	pollfds[0] = pollfd
	unix.EpollEvent{}

	n, err := unix.Poll(pollfds, 0)
	if err != nil {
		log.Printf("poll error")
	}

}