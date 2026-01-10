package server

import (
	"bytes"
	"fmt"
	"io"
	"ja_httpserver/internal/request"
	"ja_httpserver/internal/response"
	"net"
)

type HandlerError struct {
	StatusCode response.StatusCode
	Message    string
	Err        error
}

type Handler func(w io.Writer, r *request.Request) *HandlerError

type Server struct {
	closed  bool
	handler Handler
	Port    uint16
}

func runConnection(s *Server, conn io.ReadWriteCloser) {
	// out := []byte("HTTP/1.1 200 OK\r\nContent-Length: 13\r\n\r\nHello, world!")
	// conn.Write(out)
	defer conn.Close()

	headers := response.GetDefaultHeaders(0)
	r, err := request.RequestFromReader(conn)
	if err != nil {
		response.WriteStatusLine(conn, response.StatusBadRequest)
		response.WriteHeaders(conn, headers)
		return
	}
	writer := bytes.NewBuffer(nil)
	handlerError := s.handler(writer, r)

	var body []byte = nil
	var status response.StatusCode = response.StatusOK
	if handlerError != nil {
		status = handlerError.StatusCode
		body = []byte(handlerError.Message)
	} else {
		body = writer.Bytes()
	}

	headers.Replace("Content-Length", fmt.Sprintf("%d", len(body)))
	response.WriteStatusLine(conn, status)
	response.WriteHeaders(conn, headers)
	conn.Write(body)

}

func runServer(s *Server, listener net.Listener) {
	for {
		conn, err := listener.Accept()
		if s.closed {
			return
		}
		if err != nil {
			return
		}

		go runConnection(s, conn)
	}
}

func Serve(port uint16, handler Handler) (*Server, error) {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return nil, err
	}

	server := &Server{Port: port, closed: false}
	go runServer(server, listener)
	return server, nil
}

func (s *Server) Close() error {
	s.closed = true
	return nil
}
