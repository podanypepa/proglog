package server

import (
	"fmt"
	"sync"
)

// Record type
type Record struct {
	Value  []byte `json:"value"`
	Offset uint64 `josn:"offset"`
}

// Log type
type Log struct {
	m       sync.Mutex
	records []Record
}

// NewLog constructor
func NewLog() *Log {
	return &Log{}
}

func (l *Log) Append(r Record) (offset uint64, err error) {
	l.m.Lock()
	defer l.m.Unlock()

	offset = uint64(len(l.records))
	l.records = append(l.records, r)

	return offset, nil
}

func (l *Log) Read(offset uint64) (r Record, err error) {
	l.m.Lock()
	defer l.m.Unlock()

	if offset > uint64(len(l.records)) {
		return Record{}, fmt.Errorf("not found offset")
	}

	return l.records[offset], nil
}
