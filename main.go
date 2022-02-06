package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	signals := make(chan os.Signal)
	signal.Notify(signals, syscall.SIGTERM, syscall.SIGINT)
	go func() {
		select {
		case <-ctx.Done():
		case <-signals:
			cancel()
		}
	}()
	fmt.Println("waiting for sigkill")
	<-ctx.Done()
	fmt.Println("app terminated")
}
