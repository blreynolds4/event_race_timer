package events

import (
	"fmt"
	"time"

	"github.com/go-redis/redis/v7"
)

type StartEvent struct {
	Source    string
	StartTime time.Time
}

type FinishEvent struct {
	Source     string
	Bib        string
	FinishTime time.Time
}

type EventTarget interface {
	SendStart(se StartEvent) error
	SendFinish(fe FinishEvent) error
}

type redisStreamEventTarget struct {
	client *redis.Client
	stream string
}

func NewRedisStreamEventTarget(c *redis.Client, name string) EventTarget {
	return &redisStreamEventTarget{
		client: c,
		stream: name,
	}
}

func (rset *redisStreamEventTarget) SendStart(se StartEvent) error {
	addArgs := redis.XAddArgs{
		Stream: rset.stream,
		Values: map[string]interface{}{
			"event_type": "start",
			"start_time": se.StartTime.UnixMilli(),
			"source":     se.Source,
		},
	}
	result := rset.client.XAdd(&addArgs)
	if result.Err() != nil {
		return result.Err()
	}

	fmt.Println("ok -", result.Val())
	return nil
}

func (rset *redisStreamEventTarget) SendFinish(fe FinishEvent) error {
	addArgs := redis.XAddArgs{
		Stream: rset.stream,
		Values: map[string]interface{}{
			"event_type":  "finish",
			"bib":         fe.Bib,
			"finish_time": fe.FinishTime.UnixMilli(),
			"source":      fe.Source,
		},
	}
	result := rset.client.XAdd(&addArgs)
	if result.Err() != nil {
		return result.Err()
	}

	fmt.Println("ok -", result.Val())
	return nil
}
