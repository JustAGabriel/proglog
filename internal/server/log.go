package server

import (
	"fmt"
	"sync"
)

var ErrOffsetNotFound = fmt.Errorf("offset not found")

type Record struct {
	Value  []byte `json:"value"`
	Offset uint   `json:"offset"`
}

type Log struct {
	mtx     sync.Mutex
	records []Record
}

func NewLog() *Log {
	return &Log{}
}

func (l *Log) Append(record Record) (uint, error) {
	l.mtx.Lock()
	defer l.mtx.Unlock()
	record.Offset = uint(len(l.records))
	l.records = append(l.records, record)
	return record.Offset, nil
}

func (l *Log) Read(offset uint) (Record, error) {
	l.mtx.Lock()
	defer l.mtx.Unlock()
	if offset >= uint(len(l.records)) {
		return Record{}, ErrOffsetNotFound
	}
	return l.records[offset], nil
}
