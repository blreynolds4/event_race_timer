package redis_stream

import (
	"blreynolds4/event-race-timer/stream"
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/go-redis/redismock/v9"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
)

func TestConstructors(t *testing.T) {
	db, _ := redismock.NewClientMock()
	actualR := NewRedisStreamReader(db, t.Name())
	_, ok := actualR.(*redisEventStream)
	assert.True(t, ok)

	actualW := NewRedisStreamWriter(db, t.Name())
	_, ok = actualW.(*redisEventStream)
	assert.True(t, ok)
}

func TestSendMessage(t *testing.T) {
	db, mock := redismock.NewClientMock()

	msg := stream.Message{
		Values: map[string]interface{}{
			"key": "value",
		},
	}

	expectedArgs := &redis.XAddArgs{
		Stream: "stream",
		Values: msg.Values,
	}

	mock.ExpectXAdd(expectedArgs).SetVal("newId")

	rs := newRedisStream(db, "stream")
	err := rs.SendMessage(context.TODO(), msg)
	assert.NoError(t, err)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Error(err)
	}
}

func TestSendMessageFails(t *testing.T) {
	db, mock := redismock.NewClientMock()

	msg := stream.Message{
		Values: map[string]interface{}{
			"key": "value",
		},
	}

	expectedArgs := &redis.XAddArgs{
		Stream: "stream",
		Values: msg.Values,
	}

	expErr := fmt.Errorf("FAIL")
	mock.ExpectXAdd(expectedArgs).SetErr(expErr)

	rs := newRedisStream(db, "stream")
	err := rs.SendMessage(context.TODO(), msg)
	assert.Equal(t, expErr, err)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Error(err)
	}
}

func TestGetMessage(t *testing.T) {
	db, mock := redismock.NewClientMock()

	msg := stream.Message{
		ID: "msgId",
		Values: map[string]interface{}{
			"key": "value",
		},
	}

	expTimeout := time.Second

	streamName := "stream"
	expectedArgs := &redis.XReadArgs{
		Streams: []string{streamName, "0"},
		Count:   1,
		Block:   expTimeout,
	}

	expStream := []redis.XStream{
		{
			Stream: streamName,
			Messages: []redis.XMessage{
				{
					ID:     msg.ID,
					Values: msg.Values,
				},
			},
		}}
	mock.ExpectXRead(expectedArgs).SetVal(expStream)

	rs := newRedisStream(db, "stream")
	actualMsg, err := rs.GetMessage(context.TODO(), time.Second)
	assert.NoError(t, err)
	assert.Equal(t, msg, actualMsg)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Error(err)
	}
}

func TestGetMessageNoMessages(t *testing.T) {
	db, mock := redismock.NewClientMock()

	expTimeout := time.Second

	streamName := "stream"
	expectedArgs := &redis.XReadArgs{
		Streams: []string{streamName, "0"},
		Count:   1,
		Block:   expTimeout,
	}

	expStream := []redis.XStream{
		{
			Stream:   streamName,
			Messages: []redis.XMessage{},
		}}
	mock.ExpectXRead(expectedArgs).SetVal(expStream)

	rs := newRedisStream(db, "stream")
	actualMsg, err := rs.GetMessage(context.TODO(), time.Second)
	assert.NoError(t, err)
	assert.False(t, actualMsg.IsValid())

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Error(err)
	}
}

func TestGetMessageReadFail(t *testing.T) {
	db, mock := redismock.NewClientMock()

	expTimeout := time.Second

	streamName := "stream"
	expectedArgs := &redis.XReadArgs{
		Streams: []string{streamName, "0"},
		Count:   1,
		Block:   expTimeout,
	}

	expectedErr := fmt.Errorf("FAIL")
	mock.ExpectXRead(expectedArgs).SetErr(expectedErr)

	rs := newRedisStream(db, "stream")
	_, err := rs.GetMessage(context.TODO(), time.Second)
	assert.Equal(t, expectedErr, err)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Error(err)
	}
}

func TestGetMessageRange(t *testing.T) {
	db, mock := redismock.NewClientMock()

	streamName := "stream"
	expMsgs := []stream.Message{
		{
			ID: "msgId",
			Values: map[string]interface{}{
				"key": "value",
			},
		},
	}
	rawMsgs := []redis.XMessage{
		{
			ID:     expMsgs[0].ID,
			Values: expMsgs[0].Values,
		},
	}
	mock.ExpectXRange(streamName, "start", "end").SetVal(rawMsgs)

	rs := newRedisStream(db, streamName)
	actualMsg, err := rs.GetMessageRange(context.TODO(), "start", "end")
	assert.NoError(t, err)
	assert.Equal(t, expMsgs, actualMsg)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Error(err)
	}
}

func TestGetMessageRangeError(t *testing.T) {
	db, mock := redismock.NewClientMock()

	expErr := fmt.Errorf("Fail")
	streamName := "stream"
	mock.ExpectXRange(streamName, "start", "end").SetErr(expErr)

	rs := newRedisStream(db, streamName)
	_, err := rs.GetMessageRange(context.TODO(), "start", "end")
	assert.Equal(t, expErr, err)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Error(err)
	}
}
