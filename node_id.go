package it

import (
	"sync/atomic"
	"time"
)

var nextSerial int64

func NewNodeID() int64 {
	// 34-bits of timestamp, 29-bits of serial
	return time.Now().Unix()<<29 + atomic.AddInt64(&nextSerial, 1)%(1<<29)
}
