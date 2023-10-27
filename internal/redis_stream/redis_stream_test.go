package redis_stream

import (
	"blreynolds4/event-race-timer/internal/stream"
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/go-redis/redismock/v9"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
)

func TestSendMessage(t *testing.T) {
	db, mock := redismock.NewClientMock()

	msg := stream.Message{
		Data: []byte("hello"),
	}

	expectedArgs := &redis.XAddArgs{
		Stream: "stream",
		Values: map[string]interface{}{
			dataKey: msg.Data,
		},
	}

	mock.ExpectXAdd(expectedArgs).SetVal("newId")

	rs := NewRedisEventStream(db, "stream")
	err := rs.SendMessage(context.TODO(), msg)
	assert.NoError(t, err)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Error(err)
	}
}

func TestSendMessageFails(t *testing.T) {
	db, mock := redismock.NewClientMock()

	msg := stream.Message{
		Data: []byte("hello"),
	}

	expectedArgs := &redis.XAddArgs{
		Stream: "stream",
		Values: map[string]interface{}{
			dataKey: msg.Data,
		},
	}

	expErr := fmt.Errorf("FAIL")
	mock.ExpectXAdd(expectedArgs).SetErr(expErr)

	rs := NewRedisEventStream(db, "stream")
	err := rs.SendMessage(context.TODO(), msg)
	assert.Equal(t, expErr, err)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Error(err)
	}
}

func TestGetMessage(t *testing.T) {
	db, mock := redismock.NewClientMock()

	msg := stream.Message{
		ID:   "msgId",
		Data: []byte("hello"),
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
					ID: msg.ID,
					Values: map[string]interface{}{
						dataKey: msg.Data,
					},
				},
			},
		}}
	mock.ExpectXRead(expectedArgs).SetVal(expStream)

	rs := NewRedisEventStream(db, "stream")

	var actualMsg stream.Message
	gotMsg, err := rs.GetMessage(context.TODO(), time.Second, &actualMsg)
	assert.NoError(t, err)
	assert.True(t, gotMsg)
	assert.Equal(t, msg, actualMsg)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Error(err)
	}
}

func TestGetMessageBadData(t *testing.T) {
	db, mock := redismock.NewClientMock()

	msg := stream.Message{
		ID:   "msgId",
		Data: []byte("hello"),
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
					ID: msg.ID,
					Values: map[string]interface{}{
						dataKey: 1,
					},
				},
			},
		}}
	mock.ExpectXRead(expectedArgs).SetVal(expStream)

	rs := NewRedisEventStream(db, "stream")

	var actualMsg stream.Message
	gotMsg, err := rs.GetMessage(context.TODO(), time.Second, &actualMsg)
	assert.Error(t, err)
	assert.Equal(t, fmt.Errorf("unknown msg data type"), err)
	assert.False(t, gotMsg)

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

	rs := NewRedisEventStream(db, "stream")
	var msg stream.Message
	gotMsg, err := rs.GetMessage(context.TODO(), time.Second, &msg)
	assert.NoError(t, err)
	assert.False(t, gotMsg)

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

	rs := NewRedisEventStream(db, "stream")
	var msg stream.Message
	gotMessage, err := rs.GetMessage(context.TODO(), time.Second, &msg)
	assert.Equal(t, expectedErr, err)
	assert.False(t, gotMessage)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Error(err)
	}
}

func TestGetMessageRangeBufferSizeEqualMsgCount(t *testing.T) {
	db, mock := redismock.NewClientMock()

	streamName := "stream"
	expectedData := "hello"
	expMsgs := []stream.Message{
		{
			ID:   "msgId",
			Data: []byte(expectedData),
		},
	}
	rawMsgs := []redis.XMessage{
		{
			ID:     expMsgs[0].ID,
			Values: map[string]interface{}{dataKey: expMsgs[0].Data},
		},
	}
	mock.ExpectXRangeN(streamName, "start", "end", int64(len(rawMsgs))).SetVal(rawMsgs)

	rs := NewRedisEventStream(db, streamName)
	actualMsgs := make([]stream.Message, 1)
	countRead, err := rs.GetMessageRange(context.TODO(), "start", "end", actualMsgs)
	assert.NoError(t, err)
	assert.Equal(t, len(expMsgs), countRead)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Error(err)
	}
}

func TestGetMessageRangeBadData(t *testing.T) {
	db, mock := redismock.NewClientMock()

	streamName := "stream"
	expectedData := "hello"
	expMsgs := []stream.Message{
		{
			ID:   "msgId",
			Data: []byte(expectedData),
		},
	}
	rawMsgs := []redis.XMessage{
		{
			ID:     expMsgs[0].ID,
			Values: map[string]interface{}{dataKey: 1},
		},
	}
	mock.ExpectXRangeN(streamName, "start", "end", int64(len(rawMsgs))).SetVal(rawMsgs)

	rs := NewRedisEventStream(db, streamName)
	actualMsgs := make([]stream.Message, 1)
	countRead, err := rs.GetMessageRange(context.TODO(), "start", "end", actualMsgs)
	assert.Error(t, err)
	assert.Equal(t, fmt.Errorf("unknown msg data type"), err)
	assert.Equal(t, 0, countRead)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Error(err)
	}
}

func TestGetMessageRangeEmptyBuffer(t *testing.T) {
	db, _ := redismock.NewClientMock()
	streamName := "stream"

	rs := NewRedisEventStream(db, streamName)
	// should read nothing because acutal has no capacity
	actualMsgs := make([]stream.Message, 0)
	countRead, err := rs.GetMessageRange(context.TODO(), "start", "end", actualMsgs)
	assert.Equal(t, "can't get message range with empty buffer", err.Error())
	assert.Equal(t, 0, countRead)
}

func TestGetMessageRangeError(t *testing.T) {
	db, mock := redismock.NewClientMock()

	expErr := fmt.Errorf("Fail")
	streamName := "stream"
	mock.ExpectXRangeN(streamName, "start", "end", 1).SetErr(expErr)

	rs := NewRedisEventStream(db, streamName)
	actualMessages := make([]stream.Message, 1)
	count, err := rs.GetMessageRange(context.TODO(), "start", "end", actualMessages)
	assert.Equal(t, expErr, err)
	assert.Equal(t, 0, count)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Error(err)
	}
}
