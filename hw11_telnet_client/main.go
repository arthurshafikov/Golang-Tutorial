package main

import (
	"io/ioutil"
	"net"
)

func main() {
	host := net.JoinHostPort(address, "")

	client := NewTelnetClient(host, timeout, ioutil.NopCloser(in), out)

	// Place your code here,
	// P.S. Do not rush to throw context down, think think if it is useful with blocking operation?
}
