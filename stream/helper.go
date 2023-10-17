package stream

import (
	"context"
	"time"
)

type MockStream struct {
	Send   func(ctx context.Context, sm Message) error
	Get    func(ctx context.Context, timeout time.Duration, msg *Message) (bool, error)
	Range  func(ctx context.Context, startId, endId string, msgs []Message) (int, error)
	Events []Message
}

func (mes *MockStream) SendMessage(ctx context.Context, sm Message) error {
	if mes.Send != nil {
		return mes.Send(ctx, sm)
	}

	mes.Events = append(mes.Events, sm)
	return nil
}

func (mes *MockStream) GetMessage(ctx context.Context, timeout time.Duration, msg *Message) (bool, error) {
	if mes.Get != nil {
		return mes.Get(ctx, timeout, msg)
	}
	if len(mes.Events) > 0 {
		*msg = mes.Events[0]
		mes.Events = mes.Events[1:]
		return true, nil
	}

	return false, nil
}

func (mes *MockStream) RangeQueryMin() string {
	return "-"
}

func (mes *MockStream) ExclusiveQueryStart(string) string {
	return "("
}

func (mes *MockStream) RangeQueryMax() string {
	return "+"
}

func (mes *MockStream) GetMessageRange(ctx context.Context, startId, endId string, msgs []Message) (int, error) {
	if mes.Range != nil {
		return mes.Range(ctx, startId, endId, msgs)
	}
	if len(mes.Events) > 0 {
		msgs[0] = mes.Events[0]
		mes.Events = mes.Events[1:]
		// returning 1 because it should be count returned, not buffer size
		// mock only returns 1 at a time
		return 1, nil
	}

	return 0, nil
}
