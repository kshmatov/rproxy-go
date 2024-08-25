package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
)

func main() {
	dst := flag.String("d", "127.0.0.1:9999/test-it", "destination url")
	flag.Parse()

	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()

	url := *dst
	if !strings.HasPrefix(url, "http") {
		url = "http://" + url
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	req.AddCookie(&http.Cookie{Name: "srvName", Value: "first request"})
	if err != nil {
		log.Fatalf("create request: %v", err)
	}

	cli := http.Client{}
	resp, err := cli.Do(req)
	if err != nil {
		log.Fatalf("response: %v", err)
	}
	defer resp.Body.Close()

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("read response: %v", err)
	}

	fmt.Printf("[%s] -> %d [%s]\n", url, resp.StatusCode, string(b))
	for k, v := range resp.Header {
		fmt.Printf("\tHeader <%s> -> %v\n", k, v)
	}
	for _, v := range resp.Cookies() {
		fmt.Printf("\tCookie <%s> -> %s\n", v.Name, v.Value)
	}

	fmt.Print("-------------------------------------------------------\n\n")
}
