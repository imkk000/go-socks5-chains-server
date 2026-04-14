package main

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/things-go/go-socks5"
)

func main() {
	srv := socks5.NewServer(socks5.WithDialAndRequest(func(ctx context.Context, network, addr string, req *socks5.Request) (net.Conn, error) {
		fmt.Printf("from: %s -> %s\n", req.LocalAddr.String(), addr)
		return (&net.Dialer{}).DialContext(ctx, network, addr)
	}))

	go srv.ListenAndServe("tcp", "127.0.0.1:9001")

	done := make(chan os.Signal, 1)
	defer close(done)
	signal.Notify(done, os.Interrupt, syscall.SIGTERM)

	<-done
}
