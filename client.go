package main

import (
	"errors"
	"net/rpc"
)

type RsyncClient int

func (r *RsyncClient) FileToSync(fileName *string, hashTable *[]*entryPoint) error {
	checksumList := calculateFileChecksum(*fileName)
	hashTable = buildHashTable(checksumList)
	return nil
}
