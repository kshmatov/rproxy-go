package proxy

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"log/slog"
	"net/http"
	"sync"

	"github.com/kshmatov/proxy-server/internal/config"
)

type Proxy struct {
	c  *config.Config
	wg *sync.WaitGroup
}

func New(cfg *config.Config) *Proxy {
	return &Proxy{
		c:  cfg,
		wg: new(sync.WaitGroup),
	}
}

func (p *Proxy) Start(ctx context.Context) {
	proxyHandler := http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		p.wg.Add(1)
		defer p.wg.Done()

		dst := p.c.Destination()
		slog.Info("incoming", "remote", req.RemoteAddr, "req", req.URL.String(), "destination", dst)
		req.Host = dst.Host
		req.URL.Host = dst.Host
		req.URL.Scheme = dst.Scheme
		req.RequestURI = ""

		cli := http.Client{}
		req.Header.Add("Origin", "http://127.0.0.1")
		resp, err := cli.Do(req.WithContext(ctx))
		if err != nil {
			slog.Error("do request", "url", req.URL.String(), "error", err.Error())
			rw.WriteHeader(http.StatusInternalServerError)
			_, _ = fmt.Fprint(rw, err)
			return
		}

		//		defer resp.Body.Close()
		data, err := io.ReadAll(resp.Body)
		if err != nil {
			slog.Error("read original response", "url", req.URL.String(), "error", err.Error())
			rw.WriteHeader(http.StatusInternalServerError)
			_, _ = fmt.Fprint(rw, err)
			return
		}

		head := resp.Header.Clone()
		for k, v := range head {
			for _, item := range v {
				slog.Debug("header", "key", k, "value", item)
				rw.Header().Set(k, item)
			}
		}

		rw.WriteHeader(resp.StatusCode)
		_, err = rw.Write(data)
		if err != nil {
			slog.Error("send boy", "error", err)
		}
	})

	server := &http.Server{Addr: p.c.Host(), Handler: proxyHandler}
	go func() {
		err := server.ListenAndServe()
		if err != nil && errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("can't start proxy: %v", err)
		}
	}()
}

func (p *Proxy) Close() {
	p.wg.Wait()
	slog.Info("done")
}
