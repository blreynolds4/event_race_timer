package command

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

func NewPingCommand(rdb *redis.Client) Command {
	return &noStateCommand{
		CmdFunc: func(args []string) (bool, error) {
			cmdResult := rdb.Ping(context.TODO())
			fmt.Println(cmdResult.String())
			return false, cmdResult.Err()
		},
	}
}
