package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
)

var trapSignals = []os.Signal{
	syscall.SIGHUP,
	syscall.SIGINT,
	syscall.SIGTERM,
	syscall.SIGQUIT,
}

func watchSignal(ctx context.Context) <-chan os.Signal {
	out := make(chan os.Signal, 1)
	go func() {
		defer close(out)

		sigCh := make(chan os.Signal, 1)
		defer close(sigCh)
		signal.Notify(sigCh, trapSignals...)

		select {
		case <-ctx.Done():
		case sig := <-sigCh:
			out <- sig
		}
	}()

	return out
}
