// +build debug

package main

import (
	"net/http"
	_ "net/http/pprof"
	"os"
)

func init() {
	listen := os.Getenv("SYGGO_HTTP")
	if listen == "" {
		listen = "[::]:8082"
	}
	stdout.Printf("profiling is enabled, HTTP server is attached to %s\n", listen)
	go func() {
		if err := http.ListenAndServe(listen, nil); err != nil {
			stderr.Fatalf("http failed: %v\n", err)
		}
	}()
}
