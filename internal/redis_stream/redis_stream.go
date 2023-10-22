package redis_stream

import (
	"blreynolds4/event-race-timer/internal/stream"
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
		resultMsg.Data, err = rs.decodeMessageData(redisMsg.Values[dataKey])
		if err != nil {
			return false, err
		}

		rs.lastMsgId = redisMsg.ID

		return true, nil
	}

	return false, nil
}

func (rs *RedisEventStream) decodeMessageData(data any) ([]byte, error) {
	switch data.(type) {
	case []byte:
		return data.([]byte), nil
	case string:
		stringData := data.(string)
		return []byte(stringData), nil
	default:
		return nil, fmt.Errorf("unknown msg data type")
	}

}

func (rs *RedisEventStream) RangeQueryMin() string {
	return "-"
}

func (rs *RedisEventStream) ExclusiveQueryStart(id string) string {
	return "(" + id
}

func (rs *RedisEventStream) RangeQueryMax() string {
	return "+"
}

func (rs *RedisEventStream) GetMessageRange(ctx context.Context, startId, endId string, resultMessages []stream.Message) (int, error) {
	if len(resultMessages) == 0 {
		return 0, fmt.Errorf("can't get message range with empty buffer")
	}

	data, err := rs.client.XRangeN(ctx, rs.stream, startId, endId, int64(len(resultMessages))).Result()
	if err != nil && err != redis.Nil {
		return 0, err
	}

	// put result messages into the result slice.  Only return up to the length of what was returned
	resultCount := len(data)
	for i := 0; i < resultCount; i++ {
		decodedData, err := rs.decodeMessageData(data[i].Values[dataKey])
		if err != nil {
			return 0, err
		}

		resultMessages[i] = stream.Message{ID: data[i].ID, Data: decodedData}
	}

	return resultCount, nil
}
