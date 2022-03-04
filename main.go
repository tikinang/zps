package main

import (
	"fmt"
	"time"
	"zps/pkg/graceful"
)

func main() {
	ctx, cancel := graceful.Context()
	defer cancel()
	for {
		if ctx.Err() != nil {
			break
		}
		fmt.Printf("logging time: %s\n", time.Now())
	}
	fmt.Println("shutting down")
}
