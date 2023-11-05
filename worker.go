package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"time"
)

func mainWorker(ctx context.Context, cli *CLI) int {
	terminal, err := NewTerminal(cli.outStream, cli.errStream, cli.outFd)
	if err != nil {
		fmt.Fprintln(cli.errStream, "Failed to initialize terminal:", err)
		return 1
	}

	inputCh := watchInput(ctx, cli.inStream)

	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return 1
		case <-ticker.C:
			terminal.Render()
		case newLine, ok := <-inputCh:
			if !ok {
				terminal.Render()
				return 0
			}

			terminal.AddLine(newLine)
		}
	}
}

func watchInput(ctx context.Context, inStream io.Reader) <-chan string {
	out := make(chan string)
	go func() {
		defer close(out)

		scanner := bufio.NewScanner(inStream)

		for {
			select {
			case <-ctx.Done():
				return
			default:
				if !scanner.Scan() {
					return
				}
				out <- scanner.Text()
			}
		}
	}()

	return out
}
