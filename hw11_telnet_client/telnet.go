package main

import (
	"bufio"
	"errors"
	"io"
	"log"
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
	debug("Connect....")
	return nil
}

func (t TelClient) Close() error {
	debug("Close...")
	return (*t.conn).Close()
}

func (t TelClient) Send() error {
	debug("Send...")
	return t.ReadAndWrite((*t.in), (*t.conn))
}

func (t TelClient) Receive() error {
	debug("Receive...")
	return t.ReadAndWrite((*t.conn), (*t.out))
}

func (t TelClient) ReadAndWrite(readFrom io.Reader, writeTo io.Writer) error {
	debug("Read and write...")
	reader := bufio.NewReader(readFrom)
	text, err := reader.ReadBytes('\n')
	debug("=====================================", text)
	debug("Read and write...", text, err)
	if err != nil && !errors.Is(err, io.EOF) {
		return err
	}
	debug("Read and writed...")
	text = append(text, '\n')
	n, err := writeTo.Write(text)
	debug("Read and writedd err", n, err)
	if err != nil {
		return err
	}
	return nil
}

func NewTelnetClient(address string, timeout time.Duration, in io.ReadCloser, out io.Writer) TelnetClient {
	debug("New client....")
	conn, err := net.DialTimeout("tcp", address, timeout)
	if err != nil {
		log.Fatalln(err)
	}
	t := TelClient{
		in:   &in,
		out:  &out,
		conn: &conn,
	}
	return t
}
