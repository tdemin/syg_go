// +build debug

package main

import (
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
)

func init() {
	listen := os.Getenv("SYGGO_HTTP")
	if listen == "" {
		listen = "[::]:8082"
	}
	log.Printf("profiling is enabled, HTTP server is attached to %s\n", listen)
	go func() {
		if err := http.ListenAndServe(listen, nil); err != nil {
			log.Fatalf("http failed: %v\n", err)
		}
	}()
}
