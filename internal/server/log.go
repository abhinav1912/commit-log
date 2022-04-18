package server

import (
	"fmt"
	"sync"
)

type Log struct {
	mutex   sync.Mutex
	records []Record
}

type Record struct {
	Value  []byte `json:"value"`
	Offset uint64 `json:"offset"`
}

var ErrorOffsetNotFound = fmt.Errorf("Offset not found.")

func NewLog() *Log {
	return &Log{}
}

func (commit *Log) Append(record Record) (uint64, error) {
	commit.mutex.Lock()
	defer commit.mutex.Unlock()
	record.Offset = uint64(len(commit.records))
	commit.records = append(commit.records, record)
	return record.Offset, nil
}

func (commit *Log) Read(offset uint64) (Record, error) {
	commit.mutex.Lock()
	defer commit.mutex.Unlock()
	if offset >= uint64(len(commit.records)) {
		return Record{}, ErrorOffsetNotFound
	}
	return commit.records[offset], nil
}
