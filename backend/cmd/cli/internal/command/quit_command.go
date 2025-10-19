package command

import "fmt"

func NewQuitCommand() Command {
	return &noStateCommand{
		CmdFunc: func(args []string) (bool, error) {
			fmt.Println("quitting...")
			return true, nil
		},
	}
}
