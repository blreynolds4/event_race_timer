package command

import (
	"testing"

	"github.com/go-redis/redismock/v9"
	"github.com/stretchr/testify/assert"
)

func TestPingCommand(t *testing.T) {
	db, mock := redismock.NewClientMock()

	// set up expectations
	mock.ExpectPing().SetVal("pong")

	ping := NewPingCommand(db)
	q, err := ping.Run([]string{})
	assert.NoError(t, err)
	assert.False(t, q)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Error(err)
	}
}
