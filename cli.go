package main

import (
	"context"
	"fmt"
	"io"
	"os"
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

	signalCh := watchSignal(ctx)

	mainCh := make(chan int)
	go func() {
		defer close(mainCh)

		mainCh <- mainWorker(ctx, cli)
	}()

	for {
		select {
		case sig := <-signalCh:
			fmt.Fprintf(cli.errStream, "Received signal: %s\n", sig)

			return 1
		case result := <-mainCh:
			return result
		}
	}
}
