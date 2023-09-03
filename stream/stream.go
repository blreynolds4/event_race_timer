package stream

import (
	"context"
	"time"
)

// Message represents the data sent on a stream.
// It needs an id and dictionary
type Message struct {
	ID     string
	Values map[string]interface{}
}

func (msg Message) IsValid() bool {
	// a message is valid if it has and ID and non nil Values
	return (len(msg.ID) > 0) && (msg.Values != nil)
}

type Writer interface {
	SendMessage(ctx context.Context, sm Message) error
}

type Reader interface {
	GetMessage(ctx context.Context, timeout time.Duration) (Message, error)
	GetMessageRange(ctx context.Context, startId, endId string) ([]Message, error)
}
