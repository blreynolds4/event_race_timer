package repl

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

type ReadEvalPrintLoop interface {
	GetName() string
	Run()
}

type CommandLine func([]string) bool

type repl struct {
	name   string
	input  io.Reader
	cmdRun CommandLine
}

func NewReadEvalPrintLoop(name string, input io.Reader, runner CommandLine) ReadEvalPrintLoop {
	return &repl{
		name:   name,
		input:  input,
		cmdRun: runner,
	}
}

func (r *repl) GetName() string {
	return r.name
}

func (r *repl) Run() {
	scanner := bufio.NewScanner(r.input)
	done := false
	for !done {
		// read a line of input into an array of strings
		fmt.Printf("%s>", r.name)
		if scanner.Scan() {
			fmt.Println()
			line := scanner.Text()
			// look up the first string as the command and pass the rest to the command if one is found.
			cmdArgs := strings.Split(line, " ")
			done = r.cmdRun(cmdArgs)
		}
	}
}
