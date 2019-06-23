package main

import (
	"errors"
	"net/rpc"
)

type RsyncServer int

func (r *RsyncClient) FileHashTable(fileName *string, hashTable *[]*entryPoint) error {
	checksumList := calculateFileChecksum(*fileName)
	hashTable = buildHashTable(checksumList)
	return nil
}
