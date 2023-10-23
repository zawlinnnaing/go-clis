package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"time"
)

func main() {
	host := flag.String("h", "localhost", "Server host")
	port := flag.Int("p", 8080, "Server port")
	file := flag.String("f", "todo-server.json", "Todo json file")
	flag.Parse()

	addr := fmt.Sprintf("%s:%d", *host, *port)

	server := &http.Server{
		Addr:         addr,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		Handler:      newMux(*file),
	}

	fmt.Fprintf(os.Stdout, "Server listening on: %s", addr)
	if err := server.ListenAndServe(); err != nil {
		fmt.Fprint(os.Stderr, err)
		os.Exit(1)
	}
}
