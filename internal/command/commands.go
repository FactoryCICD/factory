package command

import (
	"os"

	"github.com/mitchellh/cli"
)

var Commands map[string]cli.CommandFactory

func InitCommands(
	workingDir string,
) {
	meta := Meta{
		WorkingDir: workingDir,
		Ui: &cli.BasicUi{
			Writer:      os.Stdout,
			ErrorWriter: os.Stderr,
			Reader:      os.Stdin,
		},
	}

	Commands = map[string]cli.CommandFactory{
		"validate": func() (cli.Command, error) {
			return &ValidateCommand{
				Meta: meta,
			}, nil
		},
	}
}
