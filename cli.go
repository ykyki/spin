package main

import (
	"context"
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

func (cli *CLI) Run(args []string) int {
	if len(args) > 1 {
		fmt.Fprintln(cli.errStream, "option not supported")

		return 1
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	for {
		select {
		case sig := <-watchSignal(ctx):
			fmt.Fprintf(cli.errStream, "Received signal: %s\n", sig)

			return 1
		case result := <-doMain(ctx, cli):
			return result
		}
	}
}

func doMain(ctx context.Context, cli *CLI) <-chan int {
	out := make(chan int, 1)
	go func() {
		defer close(out)

		terminal, err := NewTerminal(cli.outStream, cli.outFd)
		if err != nil {
			fmt.Fprintln(cli.errStream, "Failed to initialize terminal:", err)
			out <- 1

			return
		}

		height, err := terminal.getHeight()
		if err != nil {
			fmt.Fprintln(cli.errStream, "Failed to get terminal height:", err)
			out <- 1

			return
		}

		linesInit := min(5, height/2)

		lines := linesInit
		frame := 0

		for {
			select {
			case <-ctx.Done():
				return
			case <-time.After(50 * time.Millisecond):
				if frame >= 100 {
					fmt.Fprintln(cli.outStream, "Done!")

					return
				}

				if frame > 0 {
					for i := 0; i < lines; i++ {
						terminal.clearCurrentLine()
						terminal.moveCursorUp(1)
					}

					height, err := terminal.getHeight()
					if err != nil {
						out <- 0

						return
					}

					lines = min(linesInit+frame%200, height-2)
				}

				for i := 0; i < lines; i++ {
					terminal.writeSpinnerLine(frame, i%5)
				}

				frame++
			}
		}
	}()

	return out
}
