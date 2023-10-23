package main

import "net/http"

func newMux(file string) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/", rootHandler)
	return mux
}

func replyTextContent(writer http.ResponseWriter, request *http.Request, status int, content string) {
	writer.Header().Set("Content-Type", "text/plain")
	writer.WriteHeader(status)
	writer.Write([]byte(content))
}
