package graceful

import (
	"context"
	"os"
	"os/signal"
	"syscall"
)

func Context() (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(context.Background())

	signals := make(chan os.Signal)
	signal.Notify(signals, TerminationSignals()...)

	go func() {
		select {
		case <-ctx.Done():
		case <-signals:
			cancel()
		}
	}()

	return ctx, cancel
}

func TerminationSignals() []os.Signal {
	return []os.Signal{
		syscall.SIGTERM,
		syscall.SIGINT,
		syscall.SIGQUIT,
		syscall.SIGKILL,
		syscall.SIGHUP,
	}
}
