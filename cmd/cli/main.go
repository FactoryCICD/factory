package main

import (
	"log"
	"os"
	"path/filepath"
	"runtime"

	"github.com/factorycicd/factory/cmd/cli/command"
	"github.com/mitchellh/cli"
)

func main() {
	os.Exit(realMain())
}

func realMain() int {
	log.Printf("[INFO] Starting factory_agent")
	log.Printf("[INFO] Go runtime verion: %s", runtime.Version())

	// Load cliconfig using LoadConfig()
	// Handle any errors

	// Get the command line args
	binName := filepath.Base(os.Args[0])
	args := os.Args[1:]

	originalWd, err := os.Getwd()
	if err != nil {
		log.Printf("[ERROR] Failed to get working directory: %s", err)
		return 1
	}

	// Initialize the commands
	command.InitCommands(originalWd)

	// look into mitcheelh/cli package
	cliRunner := &cli.CLI{
		Name:       binName,
		Args:       args,
		Commands:   command.Commands,
		HelpFunc:   command.HelpFunc,
		HelpWriter: os.Stdout,
	}

	exitCode, err := cliRunner.Run()
	if err != nil {
		log.Printf("[ERROR] Error executing CLI: %s", err.Error())
		return 1
	}

	return exitCode
}
