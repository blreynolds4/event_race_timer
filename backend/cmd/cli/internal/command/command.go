package command

type Command interface {
	Run(args []string) (bool, error)
}

type noStateCommand struct {
	CmdFunc func(args []string) (bool, error)
}

func (nsc *noStateCommand) Run(args []string) (bool, error) {
	return nsc.CmdFunc(args)
}
