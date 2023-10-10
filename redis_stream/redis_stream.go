package redis_stream

import (
	"blreynolds4/event-race-timer/stream"
	"context"
	"fmt"
	"time"

	redis "github.com/redis/go-redis/v9"
)

// dataKey is used to store message payloads in AddXArg Values map
const dataKey = "data"

type RedisEventStream struct {
	client    *redis.Client
	stream    string
	lastMsgId string
}

func NewRedisEventStream(c *redis.Client, name string) *RedisEventStream {
	return &RedisEventStream{
		client:    c,
		stream:    name,
		lastMsgId: "0",
	}
}

func (rs *RedisEventStream) SendMessage(ctx context.Context, sm stream.Message) error {
	addArgs := redis.XAddArgs{
		Stream: rs.stream,
		ID:     sm.ID,
		Values: map[string]interface{}{
			dataKey: sm.Data,
		},
	}

	result := rs.client.XAdd(ctx, &addArgs)
	if result.Err() != nil {
		return result.Err()
	}

	return nil
}

func (rs *RedisEventStream) GetMessage(ctx context.Context, timeout time.Duration, resultMsg *stream.Message) (bool, error) {
	data, err := rs.client.XRead(ctx, &redis.XReadArgs{
		Streams: []string{rs.stream, rs.lastMsgId},
		//count is number of entries we want to read from redis
		Count: 1,
		//we use the block argument to make sure if no entry is found we wait
		//timeout duration, 0 is forever
		Block: timeout,
	}).Result()

	if err != nil && err != redis.Nil {
		return false, err
	}

	if err != redis.Nil && len(data[0].Messages) > 0 {
		redisMsg := data[0].Messages[0]
		resultMsg.ID = redisMsg.ID
		if _, ok := redisMsg.Values[dataKey].([]byte); !ok {
			return false, fmt.Errorf("unknown msg data type")
		}
		resultMsg.Data = redisMsg.Values[dataKey].([]byte)

		rs.lastMsgId = redisMsg.ID

		return true, nil
	}

	return false, nil
}

func (rs *RedisEventStream) GetMessageRange(ctx context.Context, startId, endId string, resultMessages []stream.Message) (int, error) {
	if len(resultMessages) == 0 {
		return 0, fmt.Errorf("can't get message range with empty buffer")
	}
	data, err := rs.client.XRangeN(ctx, rs.stream, startId, endId, int64(len(resultMessages))).Result()
	if err != nil && err != redis.Nil {
		return 0, err
	}

	// put result messages into the result slice.  Only return up to the capacity of the slice
	bufferLength := len(resultMessages)
	for i := 0; i < bufferLength; i++ {
		rawMsgData := data[i].Values[dataKey]
		if _, ok := rawMsgData.([]byte); !ok {
			return 0, fmt.Errorf("unknown msg data type in range")
		}
		resultMessages[i] = stream.Message{ID: data[i].ID, Data: rawMsgData.([]byte)}
	}

	return len(resultMessages), nil
}
