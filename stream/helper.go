package stream

import (
	"context"
	"fmt"
	"time"
)

type MockStream struct {
	Send   func(ctx context.Context, sm Message) error
	Get    func(ctx context.Context, timeout time.Duration) (Message, error)
	Range  func(ctx context.Context, startId, endId string) ([]Message, error)
	Events []Message
}

func (mes *MockStream) SendMessage(ctx context.Context, sm Message) error {
	if mes.Send != nil {
		return mes.Send(ctx, sm)
	}

	mes.Events = append(mes.Events, sm)
	return nil
}

func (mes *MockStream) GetMessage(ctx context.Context, timeout time.Duration) (Message, error) {
	if mes.Get != nil {
		return mes.Get(ctx, timeout)
	}
	fmt.Println("mes.Events", mes.Events)
	if len(mes.Events) > 0 {
		result := mes.Events[0]
		mes.Events = mes.Events[1:]
		fmt.Println("result", result)
		return result, nil
	}

	fmt.Println("default return")
	return Message{}, nil
}

func (mes *MockStream) GetMessageRange(ctx context.Context, startId, endId string) ([]Message, error) {
	if mes.Get != nil {
		return mes.Range(ctx, startId, endId)
	}
	if len(mes.Events) > 0 {
		result := mes.Events[0]
		mes.Events = mes.Events[1:]
		fmt.Println("result", result)
		return []Message{result}, nil
	}

	fmt.Println("default return")
	return []Message{}, nil
}
