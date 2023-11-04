package main

import (
	"bufio"
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

		inputCh := make(chan string, 10)
		go func() {
			defer close(inputCh)

			scanner := bufio.NewScanner(cli.inStream)

			for {
				select {
				case <-ctx.Done():
					return
				case <-time.After(5 * time.Millisecond):
					if !scanner.Scan() {
						return
					}
					inputCh <- scanner.Text()
				}
			}
		}()

		var inputBuf []string

		renderedLineCount := 0

		for frame := 0; ; {
			select {
			case <-ctx.Done():
				return
			case newLine, ok := <-inputCh:
				if !ok {
					return
				}

				inputBuf = append(inputBuf, newLine)

				for i := 0; i < renderedLineCount; i++ {
					terminal.clearCurrentLine()
					terminal.moveCursorUp(1)
				}

				var maxLineCount int
				{
					height, err := terminal.getHeight()
					if err != nil {
						out <- 1

						fmt.Fprintln(cli.errStream, "Failed to get terminal height:", err)

						return
					}
					maxLineCount = height - 2
				}

				l := len(inputBuf)
				lineCount := min(maxLineCount, l)

				terminal.writeSpinnerLine(frame, 3)

				for i := 0; i < lineCount; i++ {
					fmt.Fprintf(terminal.outStream, "%s\n", inputBuf[l-lineCount+i])
				}

				renderedLineCount = lineCount + 1
				frame++
			}
		}
	}()

	return out
}
