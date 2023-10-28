package main

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"
)

func newMux(file string) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/", rootHandler)
	todoHandler := todoRouter(file, &sync.Mutex{})
	mux.Handle("/todo", http.StripPrefix("/todo", todoHandler))
	mux.Handle("/todo/", http.StripPrefix("/todo/", todoHandler))
	return mux
}

func replyTextContent(writer http.ResponseWriter, request *http.Request, status int, content string) {
	writer.Header().Set("Content-Type", "text/plain")
	writer.WriteHeader(status)
	writer.Write([]byte(content))
}

func replyJSONContent(writer http.ResponseWriter, request *http.Request, status int, content any) {
	body, err := json.Marshal(content)
	if err != nil {
		replyError(writer, request, http.StatusInternalServerError, err.Error())
		return
	}
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(status)
	writer.Write(body)
}

func replyError(writer http.ResponseWriter, request *http.Request, status int, message string) {
	log.Printf("%s %s: Error: %d %s", request.Method, request.URL, status, message)
	http.Error(writer, http.StatusText(status), status)
}
