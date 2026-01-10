package response

import (
	"fmt"
	"io"
	"ja_httpserver/internal/headers"
)

type Response struct {
	StatusCode StatusCode
}

type StatusCode int

const (
	StatusOK                  StatusCode = 200
	StatusBadRequest          StatusCode = 400
	StatusNotFound            StatusCode = 404
	StatusInternalServerError StatusCode = 500
)

func WriteStatusLine(w io.Writer, statusCode StatusCode) error {
	statusText := ""
	switch statusCode {
	case StatusOK:
		statusText = "OK"
	case StatusBadRequest:
		statusText = "Bad Request"
	case StatusNotFound:
		statusText = "Not Found"
	case StatusInternalServerError:
		statusText = "Internal Server Error"
	default:
		statusText = "Unknown Status"
	}
	_, err := fmt.Fprintf(w, "HTTP/1.1 %d %s\r\n", statusCode, statusText)
	return err
}

func GetDefaultHeaders(contentLen int) *headers.Headers {
	h := headers.NewHeaders()
	h.Set("Content-Length", fmt.Sprintf("%d", contentLen))
	h.Set("Connection", "close")
	h.Set("Content-Type", "text/plain")
	return h
}

func WriteHeaders(w io.Writer, h *headers.Headers) error {
	b := []byte{}
	h.ForEach(func(n, v string) {
		b = fmt.Appendf(b, "%s: %s\r\n", n, v)
	})

	b = fmt.Append(b, "\r\n")
	_, err := w.Write(b)
	return err
}
