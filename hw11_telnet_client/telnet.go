package main

import (
	"bufio"
	"context"
	"io"
	"time"
)

type TelnetClient interface {
	Connect() error
	io.Closer
	Send() error
	Receive() error
}

type TelClient struct {
	in  io.ReadCloser
	out io.Writer
}

func (t TelClient) Connect() error {

	return nil
}

func (t TelClient) Close() error {

	return nil
}

func (t TelClient) Send() error {

	return nil
}

func (t TelClient) Receive() error {

	return nil
}

/*
$ go-telnet --timeout=10s host port
$ go-telnet mysite.ru 8080
$ go-telnet --timeout=3s 1.1.1.1 123

*/
func NewTelnetClient(address string, timeout time.Duration, in io.ReadCloser, out io.Writer) TelnetClient {
	t := TelClient{
		in:  in,
		out: out,
	}

	scanner := bufio.NewScanner(t)

	ctx, _ := context.WithTimeout(context.Background(), timeout)

	// Place your code here.
	return nil
}

// Place your code here.
// P.S. Author's solution takes no more than 50 lines.
