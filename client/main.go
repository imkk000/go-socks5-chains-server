package main

import (
	"fmt"
	"io"
	"log"
	"math/rand/v2"
	"net"
	"net/http"

	"golang.org/x/net/proxy"
)

func main() {
	doRequest := func(url string) (*http.Response, error) {
		chains := []string{
			"127.0.0.1:9001",
			"127.0.0.1:9002",
			"127.0.0.1:9003",
			"127.0.0.1:9004",
			"127.0.0.1:9050",
		}
		dialer, err := createProxyChains(chains)
		if err != nil {
			return nil, err
		}

		client := http.Client{
			Transport: &http.Transport{
				Dial: func(network, addr string) (net.Conn, error) {
					return dialer.Dial(network, addr)
				},
			},
		}
		return client.Get(url)
	}

	resp, err := doRequest("https://ifconfig.me")
	if err != nil {
		log.Fatal("request:", err)
	}
	defer resp.Body.Close()

	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(raw))
}

func createProxyChains(chains []string) (proxy.Dialer, error) {
	l := len(chains)
	if l == 1 {
		return proxy.SOCKS5("tcp", chains[0], nil, proxy.Direct)
	}
	previousDialer, err := proxy.SOCKS5("tcp", chains[l-1], nil, proxy.Direct)
	if err != nil {
		return nil, fmt.Errorf("create exit proxy: %w", err)
	}
	// random chains
	midChains := chains[1 : l-1]
	rand.Shuffle(len(midChains), func(i int, j int) {
		midChains[i], midChains[j] = midChains[j], midChains[i]
	})
	for _, addr := range midChains {
		previousDialer, err = proxy.SOCKS5("tcp", addr, nil, previousDialer)
		if err != nil {
			return nil, fmt.Errorf("create middle proxy: %w", err)
		}
	}
	entryDialer, err := proxy.SOCKS5("tcp", chains[0], nil, previousDialer)
	if err != nil {
		return nil, fmt.Errorf("create entry proxy: %w", err)
	}

	return entryDialer, nil
}
