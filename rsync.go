package main

import (
	"bufio"
	"crypto/md5"
	"fmt"
	"os"
)

const (
	BlockSize 	= 1024
	M 			= 1 << 16
)

type blockChecksum struct {
	index uint64
	weakChecksum uint32
	md5Checksum []byte
}

// 计算一个数据块的弱滚动校验和
func weakRollingChecksum(block []byte) uint32{
	var a, b uint32 = 0, 1
	for  i := range block{
		a += uint32(block[i])
		b += (uint32(len(block)) + 1) * uint32(block[i])
	}
	return (a % M) + (1 << 16 * (b % M))
}

//计算一个数据块的MD5校验和
func strongChecksum(block []byte) []byte{
	h := md5.New()
	h.Write(block)
	return h.Sum(nil)
}


func calculateBlockNumbers(pathname string) uint64 {
	fileInfo, err := os.Stat(pathname)
	if err != nil {
		panic(err)
	}
	size := fileInfo.Size()
	number := size / BlockSize
	if size % BlockSize != 0 {
		number += 1
	}
	return uint64(number)
}

//计算一个文件的弱滚动校验和以及MD5校验和
func calculateFileChecksum(pathname string) []blockChecksum{
	file, err := os.Open(pathname)
	if err != nil {
		panic(err)
	}
	checksumListLen := calculateBlockNumbers(pathname)
	var fileChecksumList = [checksumListLen]blockChecksum {}
	fileReader := bufio.NewReader(file)
	buf := make([]byte, BlockSize)
	var i uint64 = 0
	for {
		n, _ := fileReader.Read(buf)
		if n == 0 {
			return nil
		}
		fmt.Println(n)
		weakSum := weakRollingChecksum(buf)
		md5Sum := strongChecksum(buf)
		fileChecksumList[i] = blockChecksum{i, weakSum, md5Sum,}
	}
	return fileChecksumList
}

func main()  {
	filename := "/app/go_app/truman/truman"
	calculateFileChecksum(filename)
}
