package command

import (
	"github.com/mitchellh/cli"
)

var Commands map[string]cli.CommandFactory

func InitCommands(
	workingDir string,
) {
	meta := Meta{
		WorkingDir: workingDir,
	}

	Commands = map[string]cli.CommandFactory{
		"validate": func() (cli.Command, error) {
			return &ValidateCommand{
				Meta: meta,
			}, nil
		},
	}
}
