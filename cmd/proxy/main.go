package main

import (
	"context"

	_ "net/http/pprof"

	"github.com/kshmatov/proxy-server/internal/config"
	"github.com/kshmatov/proxy-server/internal/proxy"
)

func main() {
	cfg := config.Get()
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	p := proxy.New(cfg)
	p.Start(ctx)

	<-ctx.Done()

	p.Close()
}
