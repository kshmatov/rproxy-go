package main

import (
	"flag"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"time"
)

func main() {
	p := flag.Int("p", 8081, "port")
	name := flag.String("n", fmt.Sprintf("srv_%d", time.Now().Unix()), "server ID name")
	flag.Parse()

	loger := slog.Default().With("ID", *name)

	originServerHandler := http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rc, err := req.Cookie("srvName")
		var cv string
		if err == nil {
			cv = rc.Value
		} else {
			cv = err.Error()
		}

		loger.Info("request ", "remoteAddr", req.RemoteAddr, "uri", req.RequestURI, "cookie", cv)
		rw.Header().Add("Srv-Timestamp", time.Now().String())
		rw.Header().Add("X-Origin", "Test server '"+*name+"'")
		rw.Header().Add("Srv-Origin", "Test server '"+*name+"'")
		c := http.Cookie{
			Name:  "srvName",
			Value: *name,
		}
		http.SetCookie(rw, &c)

		_, _ = fmt.Fprintf(rw, "origin server <%s> response", *name)
	})

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *p), originServerHandler))
}
