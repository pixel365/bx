package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/pixel365/bx/cmd"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	root := cmd.NewRootCmd(ctx)
	if err := root.Execute(); err != nil {
		stop()
		log.Fatalf("error: %v", err)
	}
}
