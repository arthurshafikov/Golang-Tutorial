package main

import (
	"io"
	"net"
	"time"
)

type TelnetClient interface {
	Connect() error
	io.Closer
	Send() error
	Receive() error
}

type TelClient struct {
	in      *io.ReadCloser
	out     *io.Writer
	address string
	timeout time.Duration
	conn    *net.Conn
}

func (t *TelClient) Connect() error {
	conn, err := net.DialTimeout("tcp", t.address, t.timeout)
	if err != nil {
		return err
	}
	t.conn = &conn
	return nil
}

func (t *TelClient) Close() error {
	return (*t.conn).Close()
}

func (t *TelClient) Send() error {
	_, err := io.Copy(*t.conn, *t.in)
	return err
}

func (t *TelClient) Receive() error {
	_, err := io.Copy(*t.out, *t.conn)
	return err
}

func NewTelnetClient(address string, timeout time.Duration, in io.ReadCloser, out io.Writer) TelnetClient {
	t := TelClient{
		in:      &in,
		out:     &out,
		address: address,
		timeout: timeout,
	}
	return &t
}
