package main

import (
	"bufio"
	"fmt"
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
	in   *io.ReadCloser
	out  *io.Writer
	conn *net.Conn
}

func (t TelClient) Connect() error {
	fmt.Println("Connect....")

	return nil
}

func (t TelClient) Close() error {
	return (*t.conn).Close()
}

func (t TelClient) Send() error {
	return t.ReadAndWrite((*t.in), (*t.conn))
}

func (t TelClient) Receive() error {
	return t.ReadAndWrite((*t.conn), (*t.out))
}

func (t TelClient) ReadAndWrite(readFrom io.Reader, writeTo io.Writer) error {
	reader := bufio.NewReader(readFrom)
	text, err := reader.ReadBytes('\n')
	if err != nil {
		return err
	}
	writeTo.Write(text)

	return nil
}

/*
$ go-telnet --timeout=10s host port
$ go-telnet mysite.ru 8080
$ go-telnet --timeout=3s 1.1.1.1 123
*/
func NewTelnetClient(address string, timeout time.Duration, in io.ReadCloser, out io.Writer) TelnetClient {
	fmt.Println("New client....")
	conn, _ := net.DialTimeout("tcp", address, timeout)
	t := TelClient{
		in:   &in,
		out:  &out,
		conn: &conn,
	}
	return t
}
