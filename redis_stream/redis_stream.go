package redis_stream

import (
	"blreynolds4/event-race-timer/stream"
	"context"
	"time"

	redis "github.com/redis/go-redis/v9"
)

type redisEventStream struct {
	client    *redis.Client
	stream    string
	lastMsgId string
}

func NewRedisStreamReader(c *redis.Client, name string) stream.Reader {
	result := newRedisStream(c, name)
	return &result
}

func NewRedisStreamWriter(c *redis.Client, name string) stream.Writer {
	result := newRedisStream(c, name)
	return &result
}

func newRedisStream(c *redis.Client, name string) redisEventStream {
	return redisEventStream{
		client:    c,
		stream:    name,
		lastMsgId: "0",
	}
}

func (rs *redisEventStream) SendMessage(ctx context.Context, sm stream.Message) error {
	addArgs := redis.XAddArgs{
		Stream: rs.stream,
		ID:     sm.ID,
		Values: sm.Values,
	}

	result := rs.client.XAdd(ctx, &addArgs)
	if result.Err() != nil {
		return result.Err()
	}

	return nil
}

func (rs *redisEventStream) GetMessage(ctx context.Context, timeout time.Duration) (stream.Message, error) {
	result := stream.Message{} // isValid() is false until successful read

	data, err := rs.client.XRead(ctx, &redis.XReadArgs{
		Streams: []string{rs.stream, rs.lastMsgId},
		//count is number of entries we want to read from redis
		Count: 1,
		//we use the block command to make sure if no entry is found we wait
		//timeout duration, 0 is forever
		Block: timeout,
	}).Result()
	if err != nil && err != redis.Nil {
		return stream.Message{}, err
	}

	if err != redis.Nil && len(data[0].Messages) > 0 {
		msg := data[0].Messages[0]
		result.ID = msg.ID
		result.Values = msg.Values

		rs.lastMsgId = msg.ID

		return result, nil
	}

	return result, nil
}

func (rs *redisEventStream) GetMessageRange(ctx context.Context, startId, endId string) ([]stream.Message, error) {
	data, err := rs.client.XRange(ctx, rs.stream, startId, endId).Result()
	if err != nil && err != redis.Nil {
		return nil, err
	}

	// convert the data to RaceEvents and return them
	result := make([]stream.Message, 0)
	for _, msg := range data {
		result = append(result, stream.Message{ID: msg.ID, Values: msg.Values})
	}

	return result, nil
}
