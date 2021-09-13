package main

import (
	"bufio"
	"bytes"
	"context"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

func debug(text ...interface{}) {
	if false {
		log.Println(text...)
	}
}

/*
	короче ебучий test.sh не проходит и ебал я его в рот почему то пропуск строки не пропускает
	тест проходит и сам мейн норм работает блять заебался короче с днём рождения
*/
func main() {
	timeout := flag.Duration("timeout", time.Duration(60), "some usage")
	flag.Parse()
	args := flag.Args()

	address := net.JoinHostPort(args[0], args[1])

	var wg sync.WaitGroup

	wg.Add(2)
	in := &bytes.Buffer{}
	out := &bytes.Buffer{}
	client := runClient(address, *timeout, ioutil.NopCloser(in), out)

	// ctx, _ := context.WithTimeout(context.Background(), *timeout)
	ctx, close := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGINT, syscall.SIGKILL)

	go readRoutine(ctx, client, &wg)
	go writeRoutine(ctx, client, in, &wg)

	debug("wg.Wait()!")
	wg.Wait()
	debug("Client close!")

	close()
	// cancel()
	client.Close()
	// Place your code here,
	// P.S. Do not rush to throw context down, think think if it is useful with blocking operation?
}

func runClient(address string, timeout time.Duration, in io.ReadCloser, out io.Writer) TelClient {
	debug("Run client!")
	client := NewTelnetClient(address, timeout, ioutil.NopCloser(in), out)
	client.Connect()

	return client.(TelClient)
}

func readRoutine(ctx context.Context, client TelClient, wg *sync.WaitGroup) {
	defer wg.Done()
	debug("readRoutine!")
	scanner := bufio.NewScanner((*client.conn))
OUTER:
	for {
		debug("readfor!")
		select {
		case <-ctx.Done():
			debug("read <-ctx.Done()!")
			break OUTER
		default:
			if !scanner.Scan() {
				debug("!read scanner.Scan()!")
				break OUTER
			}
			text := scanner.Text()
			fmt.Println(text)
			debug("Received from channel: ", text)
		}
	}
	debug("Connection lost readRoutine")
}

func writeRoutine(ctx context.Context, client TelClient, in io.Writer, wg *sync.WaitGroup) {
	defer wg.Done()
	debug("writeRoutine!")
	scanner := bufio.NewScanner(os.Stdin)
OUTER:
	for {
		debug("writefor!")
		select {
		case <-ctx.Done():
			debug("write <-ctx.Done()!")
			break OUTER
		default:
			if !scanner.Scan() {
				debug("!write scanner.Scan()!")
				break OUTER
			}
			text := scanner.Text()
			in.Write([]byte(text))
			err := client.Send()
			if err != nil {
				debug("Connection lost writeRoutine")
				return
			}
			debug("Write to channel: ", text)
		}
	}
	debug("Connection lost writeRoutine")
}

/*
export GOROOT=/usr/local/go
export GOPATH=~/.go
export PATH=$PATH:$GOROOT/bin:$GOPATH/bin
*/
