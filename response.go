package gnet

/********************************************************
 * study and copy some part code from gin.ResponseWriter.
 ********************************************************/

import (
	"bufio"
	"io"
	"net"
	"net/http"
)

type IResponse interface {
	http.ResponseWriter
	http.Hijacker
	http.Flusher
	http.CloseNotifier

	// Returns the HTTP response status code of the current request.
	Status() int

	// Returns the number of bytes already written into the response http body.
	// See Written()
	Size() int

	// Writes the string into the response body.
	WriteString(string) (int, error)

	// Returns true if the response body was already written.
	Written() bool

	// Forces to write the http header (status code + headers).
	FlushHeader()

	// get the http.Pusher for server push
	Pusher() http.Pusher

	// init reset
	Reset()
	InitSelf(http.ResponseWriter)

	// update ResponseWriter.Header()
	Headers(key, value string)
}

const zero_written = -1

type Response struct {
	http.ResponseWriter
	size   int
	status int
}

var _ IResponse = &Response{}

func (w *Response) Reset() {
	w.size = zero_written
	w.status = http.StatusOK
}

func (w *Response) InitSelf(writer http.ResponseWriter) {
	w.ResponseWriter = writer
}

func (w *Response) WriteHeader(code int) {
	if code <= 0 || w.status == code {
		return
	}
	w.status = code
}

func (w *Response) Written() bool {
	return w.size != zero_written
}

func (w *Response) Size() int {
	return w.size
}

func (w *Response) Status() int {
	return w.status
}

// safe write header
// only write once.
// http.ResponseWriter write head twice will panic
func (w *Response) FlushHeader() {
	if w.Written() {
		return
	}
	w.size = 0
	w.ResponseWriter.WriteHeader(w.status)
}

func (w *Response) Write(data []byte) (int, error) {
	w.FlushHeader()
	n, err := w.ResponseWriter.Write(data)
	w.size += n
	return n, err
}

func (w *Response) WriteString(data string) (int, error) {
	w.FlushHeader()
	n, err := io.WriteString(w.ResponseWriter, data)
	w.size += n
	return n, err
}

// Hijack implements the http.Hijacker interface.
// On the other hand. manage conn by your self.
// Can write response content dose not follow the http schame.
func (w *Response) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	if w.size < 0 {
		w.size = 0
	}
	return w.ResponseWriter.(http.Hijacker).Hijack()
}

// CloseNotify implements the http.CloseNotify interface.
func (w *Response) CloseNotify() <-chan bool {
	return w.ResponseWriter.(http.CloseNotifier).CloseNotify()
}

// Flush implements the http.Flush interface.
func (w *Response) Flush() {
	w.FlushHeader()
	w.ResponseWriter.(http.Flusher).Flush()
}

func (w *Response) Pusher() (pusher http.Pusher) {
	if pusher, ok := w.ResponseWriter.(http.Pusher); ok {
		return pusher
	}
	return nil
}

func (w *Response) Headers(key, value string) {
	// delete key from header()
	if value == "" {
		w.ResponseWriter.Header().Del(key)
		return
	}
	w.ResponseWriter.Header().Set(key, value)
}
