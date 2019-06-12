package main

import (
	"crypto/md5"
	"fmt"
	"io/ioutil"
)

const (
	BlockSize 	 = 12
	M         	 = 1 << 16
	HashTableLen = 1 << 16
)

type entryPoint struct {
	weakChecksum uint32
	md5Checksum []byte
	next *entryPoint
}

type blockChecksum struct {
	index uint64
	weakChecksum uint32
	md5Checksum []byte
}

// 哈希函数
func hashFunc(weakChecksum uint32)uint16{
	return uint16(weakChecksum % HashTableLen)
}

// 以文件块的滚动校验和为哈希key计算哈希表
func buildHashTable(checksumList *[]blockChecksum)*[]*entryPoint{
	hashTable := make([]*entryPoint, HashTableLen)
	for _, block := range *checksumList{
		// 哈希函数为简单的取摸运算
		hashValue := hashFunc(block.weakChecksum)
		entry := hashTable[hashValue]
		if entry == nil{
			hashTable[hashValue] = &entryPoint{
				block.weakChecksum,
				block.md5Checksum,
				nil}
		} else {
			for {
				if entry.next == nil{
					entry.next = &entryPoint{
						block.weakChecksum,
						block.md5Checksum,
						nil}
					break
				} else {
					entry = entry.next
				}
			}
		}
	}
	return &hashTable
}

// 计算一个数据块的弱滚动校验和
func weakRollingChecksum(block []byte) (uint32, uint32, uint32){
	var a, b uint32 = 0, 1
	for i := range block{
		a += uint32(block[i])
		b += (uint32(len(block)-1) - uint32(i) + 1) * uint32(block[i])
	}
	return (a % M) + (1 << 16 * (b % M)), a % M, b % M
}

//计算一个数据块的MD5校验和
func strongChecksum(block []byte) []byte{
	h := md5.New()
	h.Write(block)
	return h.Sum(nil)
}

func calculateBlockNumbers(fileBytes []byte) uint64 {
	size := len(fileBytes)
	number := size / BlockSize
	if size % BlockSize != 0 {
		number += 1
	}
	return uint64(number)
}

//比较两个int数字的大小
func min(a int, b int) int {
	if a < b {
		return a
	}
	return b
}

//计算一个文件的弱滚动校验和以及MD5校验和
func calculateFileChecksum(pathname string) *[]blockChecksum{
	fileBytes, err := ioutil.ReadFile(pathname)
	if err != nil {
		panic(err)
	}
	checksumListLen := calculateBlockNumbers(fileBytes)
	var fileChecksumList = make([]blockChecksum, checksumListLen)
	for i := range fileChecksumList{
		startByte := i * BlockSize
		endByte := min((i+1)*BlockSize, len(fileBytes))
		buf := fileBytes[startByte:endByte]
		weakSum, _, _ := weakRollingChecksum(buf)
		md5Sum := strongChecksum(buf)
		fileChecksumList[i] = blockChecksum{uint64(i), weakSum, md5Sum,}
	}
	return &fileChecksumList
}

//计算源文件与目标文件不同的需要传输的文件块
func calculateDiffer(srcFile string, HashTable []*entryPoint)  {
	fileBytes, err := ioutil.ReadFile(srcFile)
	fmt.Println(len(fileBytes))
	if err != nil{
		panic(err)
	}
	var startByte, previousMatch int
	var weakSum, a, b uint32
	var isRolling bool
	Loop:
	for startByte < len(fileBytes){
		endByte := min(startByte + BlockSize, len(fileBytes))
		buf := fileBytes[startByte:endByte]
		if isRolling {
			a = (a - uint32(fileBytes[startByte-1]) + uint32(fileBytes[endByte-1])) % M
			b = (b - uint32(endByte-startByte)*uint32(fileBytes[startByte-1]) + a) % M
			weakSum = a + M*b
		}else {
			weakSum, a, b = weakRollingChecksum(buf)
		}
		hashValue := hashFunc(weakSum)
		if HashTable[hashValue] != nil {
			md5Sum := strongChecksum(buf)
			entry := HashTable[hashValue]
			for entry != nil{
				if entry.weakChecksum == weakSum && string(entry.md5Checksum) == string(md5Sum) {
					if isRolling {
						fmt.Println("differ")
						fmt.Println(string(fileBytes[previousMatch:startByte]))
						isRolling = false
					}
					startByte += BlockSize
					previousMatch = endByte
					goto Loop
				}
				entry = entry.next
			}
			isRolling = true
			startByte ++
		}else{
			isRolling = true
			startByte ++
		}
	}
	if isRolling {
		fmt.Println("differ")
		fmt.Println(string(fileBytes[previousMatch:startByte]))
		isRolling = false
	}
}