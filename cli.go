package main

import (
	"fmt"
	"io"
	"os"
	"time"
)

type CLI struct {
	inStream  io.Reader
	outStream io.Writer
	errStream io.Writer
	outFd     int // UNIX file descriptor for outStream
}

func NewCLIStd() *CLI {
	return &CLI{
		inStream:  os.Stdin,
		outStream: os.Stdout,
		errStream: os.Stderr,
		outFd:     int(os.Stdout.Fd()),
	}
}

func (c *CLI) Run(args []string) int {
	if len(args) > 1 {
		fmt.Fprintln(c.errStream, "option not supported")
		return 1
	}

	terminal := NewTerminal(c.outStream, c.outFd)

	height, err := terminal.getHeight()
	if err != nil {
		return 1
	}
	linesInit := min(5, height/2)

	lines := linesInit
	for frame := 0; frame < 100; frame++ {
		if frame > 0 {
			for i := 0; i < lines; i++ {
				terminal.clearCurrentLine()
				terminal.moveCursorUp(1)
			}

			height, err := terminal.getHeight()
			if err != nil {
				return 1
			}
			lines = min(linesInit+frame%200, height-2)
		}

		for i := 0; i < lines; i++ {
			terminal.writeSpinnerLine(frame, i%5)
		}

		time.Sleep(50 * time.Millisecond)
	}
	return 0
}
