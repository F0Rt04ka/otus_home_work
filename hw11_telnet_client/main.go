package main

import (
	"context"
	"flag"
	"fmt"
	"net"
	"os"
	"os/signal"
	"time"
)

func main() {
	var timeout time.Duration
	flag.DurationVar(&timeout, "timeout", 10*time.Second, "connection timeout")

	flag.Parse()
	address := net.JoinHostPort(flag.Arg(0), flag.Arg(1))
	client := NewTelnetClient(address, timeout, os.Stdin, os.Stdout)

	if err := client.Connect(); err != nil {
		fmt.Fprintln(os.Stderr, "Connection error: ", err)
		os.Exit(1)
	}

	fmt.Fprintln(os.Stderr, "Connected to "+address)
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	defer func() {
		if err := client.Close(); err != nil {
			fmt.Fprintln(os.Stderr, "Close error: ", err)
		}
	}()

	go func() {
		if err := client.Send(); err != nil {
			fmt.Fprintln(os.Stderr, "Send error: ", err)
		}
	}()

	go func() {
		if err := client.Receive(); err != nil {
			fmt.Fprintln(os.Stderr, "Receive error: ", err)
		}
	}()

	<-ctx.Done()
}
