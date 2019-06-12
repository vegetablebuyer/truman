package main

import (
	"crypto/md5"
	"fmt"
	"io/ioutil"
)

const (
	BlockSize    = 12
	M            = 1 << 16
	HashTableLen = 1 << 16
)

type entryPoint struct {
	weakChecksum uint32
	md5Checksum []byte
	next *entryPoint
}

// The checksum of a data block
type blockChecksum struct {
	index uint64
	weakChecksum uint32
	md5Checksum []byte
}

// The hash function
func hashFunc(weakChecksum uint32)uint16{
	return uint16(weakChecksum % HashTableLen)
}

// Calculate the hash table of the checksum list
func buildHashTable(checksumList *[]blockChecksum)*[]*entryPoint{
	hashTable := make([]*entryPoint, HashTableLen)
	for _, block := range *checksumList{
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

// Calculate the weak rolling checksum
func weakRollingChecksum(block []byte) (uint32, uint32, uint32){
	var a, b uint32 = 0, 1
	for i := range block{
		a += uint32(block[i])
		b += (uint32(len(block)-1) - uint32(i) + 1) * uint32(block[i])
	}
	return (a % M) + (1 << 16 * (b % M)), a % M, b % M
}

// Calculate the md5sum of a data block
func strongChecksum(block []byte) []byte{
	h := md5.New()
	h.Write(block)
	return h.Sum(nil)
}

// Return the block numbers of a file
// The result will be different due to the constant BlockSize
func calculateBlockNumbers(fileBytes []byte) uint64 {
	size := len(fileBytes)
	number := size / BlockSize
	if size % BlockSize != 0 {
		number += 1
	}
	return uint64(number)
}

// Compare two int numbers and return the smaller
func min(a int, b int) int {
	if a < b {
		return a
	}
	return b
}

// Divide the file pathname into blocks
// Calculate the weak rolling checksum and md5sum of every block
// Store the results into slice fileChecksumList
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

// Tell the differences between the source file and the destination file
// The result will be different due to the constant BlockSize
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