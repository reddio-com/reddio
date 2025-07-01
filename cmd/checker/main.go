package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"syscall"

	"github.com/reddio-com/reddio/checker"
)

var (
	host string
)

func init() {
	flag.StringVar(&host, "host", "http://localhost:9092", "server host")
}

func main() {
	flag.Parse()
	c := checker.NewChecker(host)
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		c.Run(ctx)
	}()
	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, os.Interrupt, syscall.SIGTERM)
	select {
	case <-sigint:
		cancel()
	}
}
