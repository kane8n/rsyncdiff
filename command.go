package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"syscall"
)

type Command struct {
	Command       string
	Options       []string
	commandString string
}

func (c *Command) initCommand() {
	var com string
	switch c.Command {
	case "rsync":
		com = strings.Join(append([]string{c.Command, "-rzilcpogn", "--delete"}, c.Options...), " ")
		if rsyncExclude != "" {
			com = fmt.Sprintf("%s --exclude %s", com, rsyncExclude)
		}
		if rsyncExcludeFrom != "" {
			com = fmt.Sprintf("%s --exclude-from=%s", com, rsyncExcludeFrom)
		}
	case "vimdiff":
		com = strings.Join(append([]string{c.Command, "-R"}, c.Options...), " ")
	case "cat":
		com = strings.Join(append([]string{c.Command}, c.Options...), " ")
		if isColor {
			com = fmt.Sprintf("%s %s", com, " | colordiff")
		}
		if isLess {
			com = fmt.Sprintf("%s %s", com, " | less -Rr")
		}
	default:
		panic(fmt.Sprintf("%s is not supported", c.Command))
	}
	c.commandString = com
}

func (c *Command) Run() (int, error) {
	if c.commandString == "" {
		c.initCommand()
	}
	cmd := exec.Command(os.Getenv("SHELL"), "-c", c.commandString)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	err := cmd.Start()
	if err != nil {
		return 1, err
	}
	err = cmd.Wait()
	if err != nil {
		if e2, ok := err.(*exec.ExitError); ok {
			if s, ok := e2.Sys().(syscall.WaitStatus); ok {
				return s.ExitStatus(), err
			}
			panic(errors.New("Unimplemented for system where exec.ExitError.Sys() is not syscall.WaitStatus."))
		}
	}
	return 0, nil
}

func (c *Command) Output() ([]byte, error) {
	if c.commandString == "" {
		c.initCommand()
	}
	cmd := exec.Command(os.Getenv("SHELL"), "-c", c.commandString)
	cmd.Stderr = os.Stderr
	return cmd.Output()
}
