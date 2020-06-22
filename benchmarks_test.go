package gnet

import (
	//"github.com/slclub/link"
	"net/http"
	"sync"
	"testing"
)

func BenchmarkContext(B *testing.B) {
	en := NewEngine()
	run_request(B, en, "GET", "/list/me?name=xiaohu")
}

// -=================================================================================

type Engine struct {
	pool sync.Pool
}

func NewEngine() *Engine {
	eg := &Engine{}
	eg.pool.New = func() interface{} {
		return eg.allocateContext()
	}
	return eg
}

func (eg *Engine) allocateContext() Contexter {
	ctx := NewContext()
	r := NewRequest()
	s := &Response{}
	ctx.SetRequest(r)
	ctx.SetResponse(s)

	return ctx
}

func (eg *Engine) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	ctx := eg.pool.Get().(Contexter)
	ctx.Request().InitWithHttp(req)
	ctx.Response().InitSelf(res)
	_, _ = ctx.Request().GetString("name")
	//link.DEBUG_PRINT("name:", name)

	ctx.Reset()
	eg.pool.Put(ctx)
}

// mock http

type header_writer struct {
	header http.Header
}

func new_mock_writer() *header_writer {
	return &header_writer{
		http.Header{},
	}
}

func (m *header_writer) Write(p []byte) (n int, err error) {
	return len(p), nil
}

func (m *header_writer) WriteString(s string) (n int, err error) {
	return len(s), nil
}

func (m *header_writer) Header() http.Header {
	return m.header
}

func (m *header_writer) WriteHeader(int) {}

func run_request(B *testing.B, en *Engine, method, path string) {
	// create fake request
	req, err := http.NewRequest(method, path, nil)
	if err != nil {
		panic(err)
	}
	w := new_mock_writer()
	B.ReportAllocs()
	B.ResetTimer()
	for i := 0; i < B.N; i++ {
		en.ServeHTTP(w, req)
	}
}
