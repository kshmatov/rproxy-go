package config

import (
	"flag"
	"log"
	"log/slog"
	"net/url"
	"strings"
	"sync"
)

type Config struct {
	h            string
	pos          int
	destinations []*url.URL
	l            int
	m            *sync.Mutex
	log          *slog.Logger
}

func Get() *Config {
	l := flag.String("l", "info", "log level (info, debug, warn, error)")
	h := flag.String("h", "127.0.0.1:9999", "host:port to listen")
	flag.Parse()

	hosts := flag.Args()

	if len(hosts) == 0 {
		log.Fatal("empty destination list")
	}

	aList := make([]*url.URL, len(hosts))
	for i, v := range hosts {
		slog.Info("add dst", "url", v)
		if !strings.Contains(v, "//") {
			v = "http://" + v
		}
		u, err := url.Parse(v)
		if err != nil {
			log.Fatalf("can't parse destination <%s>: %v", v, err)
		}
		aList[i] = u
	}

	level := slog.LevelInfo
	switch *l {
	case "debug", "d", "dbg":
		level = slog.LevelDebug
	case "err", "error", "e":
		level = slog.LevelError
	case "warn", "warning", "w":
		level = slog.LevelWarn
	}

	slog.SetLogLoggerLevel(level)
	return &Config{
		destinations: aList,
		m:            new(sync.Mutex),
		l:            len(hosts),
		h:            *h,
	}
}

func (c *Config) Destination() url.URL {
	if c.l == 1 {
		return *c.destinations[0]
	}

	c.m.Lock()
	defer func() {
		c.pos++
		c.m.Unlock()
	}()

	if c.pos >= c.l {
		c.pos = 0
	}
	return *c.destinations[c.pos]
}

func (c *Config) Host() string {
	return c.h
}

func (c *Config) Logger() *slog.Logger {
	return c.log
}
