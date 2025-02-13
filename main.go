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

	if err := cmd.Execute(ctx); err != nil {
		log.Fatal(err)
	}
}
