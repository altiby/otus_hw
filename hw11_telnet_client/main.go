package main

import (
	"context"
	"flag"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	timeout := flag.Duration("timeout", time.Second*10, "timeout")
	flag.Parse()
	host := flag.Arg(0)
	port := flag.Arg(1)
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer cancel()

	client := NewTelnetClient(ctx, net.JoinHostPort(host, port), *timeout, os.Stdin, os.Stdout)
	if err := client.Connect(); err != nil {
		_ = fmt.Errorf("failed to connect %w ", err)
	}

	go func() {
		client.Send()
		cancel()
	}()

	go func() {
		client.Receive()
		fmt.Println("Connection was closed by peer")
		cancel()
	}()

	<-ctx.Done()
}
