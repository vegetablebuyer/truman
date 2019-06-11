package main

import (
	"flag"
	"fmt"
	"os"
)

func init(){
	flag.Parse()
}

var srcFile = flag.String("srcFile", "", "the source file")
var dstFile = flag.String("dstFile", "", "the destination file")

func main()  {
	if *srcFile == "" {
		fmt.Println("do not specify a source file")
		os.Exit(1)
	}
	if *dstFile == "" {
		fmt.Println("do not specify a destination file")
		os.Exit(1)
	}
	checksumList := calculateFileChecksum(*dstFile)
	hashTable := buildHashTable(checksumList)
	calculateDiffer(*srcFile, *hashTable)
}
