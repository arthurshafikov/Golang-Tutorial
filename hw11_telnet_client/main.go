package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"net"
	"time"
)

func main() {
	timeout := flag.Duration("timeout", time.Duration(60), "some usage")
	flag.Parse()
	args := flag.Args()

	address := net.JoinHostPort(args[0], args[1])

	in := &bytes.Buffer{}
	out := &bytes.Buffer{}

	client := NewTelnetClient(address, *timeout, ioutil.NopCloser(in), out)

	client.Connect()
	fmt.Println(args)

	fmt.Println(*timeout)

	// host := net.JoinHostPort(address, "")

	// client := NewTelnetClient(host, timeout, ioutil.NopCloser(in), out)

	// Place your code here,
	// P.S. Do not rush to throw context down, think think if it is useful with blocking operation?
}
