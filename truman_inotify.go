package main

import (
	"log"
	"github.com/golang.org/x/sys/unix"
	"strings"
	"unsafe"
)

const agnosticEvents = unix.IN_MOVED_TO | unix.IN_MOVED_FROM |
	unix.IN_CREATE | unix.IN_ATTRIB | unix.IN_MODIFY |
	unix.IN_MOVE_SELF | unix.IN_DELETE | unix.IN_DELETE_SELF

var flags uint32 = agnosticEvents

func main(){
	fd, error := unix.InotifyInit1(unix.IN_CLOEXEC)
	if fd == -1 {
		log.Fatal(error)
	}

	wd, error := unix.InotifyAddWatch(fd, "/tmp/foo", flags)
	if wd == -1 {
		log.Fatal(error)
	}

	var buf [unix.SizeofInotifyEvent * 4096 ]byte
	n, error := unix.Read(fd, buf[:])

	var offset uint32
	// We don't know how many events we just read into the buffer
	// While the offset points to at least one whole event...
	for offset <= uint32(n-unix.SizeofInotifyEvent) {
		// Point "raw" to the event in the buffer
		raw := (*unix.InotifyEvent)(unsafe.Pointer(&buf[offset]))

		nameLen := uint32(raw.Len)


		bytes := (*[unix.PathMax]byte)(unsafe.Pointer(&buf[offset+unix.SizeofInotifyEvent]))
		name := "/" + strings.TrimRight(string(bytes[0:nameLen]), "\000")


	log.Printf("haha")
	log.Printf(string(name))
	log.Printf("lala")


}
