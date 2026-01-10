package request

import (
	"bytes"
	"fmt"
	"io"
	"ja_httpserver/headers"
	"log/slog"
	"strconv"
)

type parserState string

const (
	StateInit    parserState = "init"
	StateDone    parserState = "done"
	StateHeaders parserState = "headers"
	StateBody    parserState = "body"
	StateError   parserState = "error"
)

type Request struct {
	RequestLine RequestLine
	Headers     *headers.Headers
	Body        string
	state       parserState
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

func getInt(headers *headers.Headers, name string, defaultValue int) int {
	valueStr, ok := headers.Get(name)
	if !ok {
		return defaultValue
	}
	valueInt, err := strconv.Atoi(valueStr)
	if err != nil {
		return defaultValue
	}
	return valueInt
}

func (r *Request) hasBody() bool {
	return getInt(r.Headers, "Content-Length", 0) > 0
}

func (r *Request) parse(data []byte) (int, error) {
	read := 0
outer:
	for {
		currentData := data[read:]
		if len(currentData) == 0 {
			break outer
		}
		switch r.state {
		case StateError:
			return 0, ERROR_REQUEST_IN_ERROR_STATE
		case StateInit:
			rl, n, err := parseRequestLine(currentData)
			if err != nil {
				r.state = StateError
				return 0, err
			}
			if n == 0 {
				break outer
			}
			r.RequestLine = *rl
			read += n
			r.state = StateHeaders
		case StateHeaders:
			n, done, err := r.Headers.Parse(currentData)
			if err != nil {
				r.state = StateError
				return 0, err
			}
			if n == 0 {
				break outer
			}
			read += n
			if done {
				if r.hasBody() {
					r.state = StateBody
				} else {
					r.state = StateDone
				}
			}
		case StateBody:
			lenght := getInt(r.Headers, "Content-Length", 0)
			if lenght == 0 {
				// r.state = StateDone
				// break outer
				panic("Content-Length is 0")
			}
			remaining := min(lenght-len(r.Body), len(currentData))

			r.Body += string(currentData[:remaining])
			read += remaining
			if len(r.Body) == lenght {
				r.state = StateDone
			}
		case StateDone:
			break outer
			// default:
			// 	panic("mmmmmmmmmmmmmmmmmmmmmmoooooooooooooooo")
		}
	}

	return read, nil
}

func (r *Request) done() bool {
	return r.state == StateDone || r.state == StateError
}

func newRequest() *Request {
	return &Request{
		state:   StateInit,
		Headers: headers.NewHeaders(),
		Body:    "",
	}
}

var ERROR_MALFORMED_REQUEST_LINE = fmt.Errorf("malformed request-line")
var ERROR_UNSUPPORTED_HTTP_VERSION = fmt.Errorf("unsupported http version")
var ERROR_REQUEST_IN_ERROR_STATE = fmt.Errorf("request in error state")
var SEPERATOR = []byte("\r\n")

func parseRequestLine(b []byte) (*RequestLine, int, error) {
	idx := bytes.Index(b, SEPERATOR)
	if idx == -1 {
		return nil, 0, nil
	}

	startLine := b[:idx]
	read := idx + len(SEPERATOR)

	parts := bytes.Split(startLine, []byte(" "))
	if len(parts) != 3 {
		return nil, 0, ERROR_MALFORMED_REQUEST_LINE
	}

	httpParts := bytes.Split(parts[2], []byte("/"))
	if len(httpParts) != 2 || string(httpParts[0]) != "HTTP" || string(httpParts[1]) != "1.1" {
		return nil, 0, ERROR_MALFORMED_REQUEST_LINE
	}

	rl := &RequestLine{
		Method:        string(parts[0]),
		RequestTarget: string(parts[1]),
		HttpVersion:   string(httpParts[1]),
	}
	return rl, read, nil
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	request := newRequest()

	buff := make([]byte, 1024)
	buffLen := 0
	for !request.done() {
		n, err := reader.Read(buff[buffLen:])
		slog.Debug("reading request", "currentBufferLength", buffLen, "state", request.state)
		if err != nil {
			return nil, err
		}
		buffLen += n

		readN, err := request.parse(buff[:buffLen])
		if err != nil {
			return nil, err
		}
		copy(buff, buff[readN:buffLen])
		buffLen -= readN
	}
	// data, err := io.ReadAll(reader)
	// if err != nil {
	// 	return nil, errors.Join(fmt.Errorf("unable to io.readALl"), err)
	// }

	// str := string(data)

	// rl, _, err := parseRequestLine(str)
	// if err != nil {
	// 	return nil, err
	// }
	return request, nil
}
