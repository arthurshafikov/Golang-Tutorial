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
	timeout := flag.Duration("timeout", time.Duration(60), "timeout for connection")
	flag.Parse()
	args := flag.Args()
	address := net.JoinHostPort(args[0], args[1])

	client := NewTelnetClient(address, *timeout, os.Stdin, os.Stdout)
	if err := client.Connect(); err != nil {
		fmt.Fprintln(os.Stderr, err)
	}

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGKILL)
	ctx, stop := context.WithTimeout(ctx, *timeout)

	go readRoutine(client, cancel)
	go writeRoutine(client, cancel)

	<-ctx.Done()
	stop()
	cancel()
	client.Close()
}

func readRoutine(client TelnetClient, cancel context.CancelFunc) {
	if err := client.Receive(); err != nil {
		cancel()
	}
}

func writeRoutine(client TelnetClient, cancel context.CancelFunc) {
	if err := client.Send(); err != nil {
		cancel()
	}
}
