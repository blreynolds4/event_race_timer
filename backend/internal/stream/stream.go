package stream

import (
	"context"
	"time"
)

// Message represents the data sent on a stream.
// It needs an id and dictionary
type Message struct {
	ID   string
	Data []byte
}

type Writer interface {
	SendMessage(ctx context.Context, sm Message) error
}

type Reader interface {
	GetMessage(ctx context.Context, timeout time.Duration, msg *Message) (bool, error)
	GetMessageRange(ctx context.Context, startId, endId string, msgs []Message) (int, error)
	RangeQueryMin() string
	ExclusiveQueryStart(string) string
	RangeQueryMax() string
}

type ReaderWriter interface {
	Reader
	Writer
}
