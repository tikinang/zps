package main

import (
	"fmt"
	"time"
	"zps/pkg/graceful"
)

func main() {
	ctx, cancel := graceful.Context()
	defer cancel()

	fmt.Println("temporal-worker started")
	for {
		select {
		case <-ctx.Done():
			fmt.Println("graceful shutdown")
			return
		case <-time.After(time.Second):
			fmt.Println("running")
			continue
		}
	}
}
