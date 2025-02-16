package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/pixel365/bx/internal"

	cfg "github.com/pixel365/bx/internal/config"

	"github.com/pixel365/bx/cmd"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	conf, err := cfg.NewConfig()
	if err != nil {
		log.Fatal(err)
	}
	var configManager internal.ConfigManager = conf

	root := cmd.NewRootCmd(ctx, configManager)
	if err := root.ExecuteContext(ctx); err != nil {
		log.Fatal(err)
	}
}
