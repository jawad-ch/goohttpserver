package main

import (
	"io"
	"ja_httpserver/internal/request"
	"ja_httpserver/internal/response"
	"ja_httpserver/internal/server"
	"log"
	"os"
	"os/signal"
	"syscall"
)

const port = 42069

func main() {
	server, err := server.Serve(port, func(w io.Writer, req *request.Request) *server.HandlerError {
		switch req.RequestLine.RequestTarget {
		case "/yourproblem":
			return &server.HandlerError{
				StatusCode: response.StatusInternalServerError,
				Message:    "Internal Server Error",
				Err:        nil,
			}
		case "/myproblem":
			return &server.HandlerError{
				StatusCode: response.StatusBadRequest,
				Message:    "Bad Request",
				Err:        nil,
			}
		default:
			w.Write([]byte("all good!!"))
		}
		return nil

	})
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
	defer server.Close()
	log.Println("Server started on port", port)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	log.Println("Server gracefully stopped")
}
