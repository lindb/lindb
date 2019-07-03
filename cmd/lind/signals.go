package lind

import (
	"context"
	"os"
	"os/signal"
	"syscall"
)

// newCtxWithSignals returns a context which will can be canceled by sending signal.
func newCtxWithSignals() context.Context {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		defer cancel()
		select {
		case <-ctx.Done():
			return
		case <-c:
			return
		}
	}()
	return ctx
}
