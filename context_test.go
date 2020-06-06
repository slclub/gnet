package gnet

import (
	"bytes"
	"fmt"
	"github.com/stretchr/testify/assert"
	"mime/multipart"
	"net"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func createMultipartRequest(t *testing.T) *http.Request {
	boundary := "--testboundary"
	body := new(bytes.Buffer)
	mw := multipart.NewWriter(body)
	defer mw.Close()

	(mw.SetBoundary(boundary))
	(mw.WriteField("foo", "bar"))
	(mw.WriteField("bar", "10"))
	(mw.WriteField("bar", "foo2"))
	(mw.WriteField("array[]", "first"))
	(mw.WriteField("array[]", "second"))
	(mw.WriteField("array1", "a1"))
	(mw.WriteField("array1", "a2"))
	(mw.WriteField("id", ""))
	(mw.WriteField("time_local", "31/12/2016 14:55"))
	(mw.WriteField("time_utc", "31/12/2016 14:55"))
	(mw.WriteField("time_location", "31/12/2016 14:55"))
	(mw.WriteField("names[a]", "thinkerou"))
	(mw.WriteField("names[b]", "tianou"))
	req, err := http.NewRequest("POST", "/", body)
	assert.Nil(t, err)
	req.Header.Set("Content-Type", "multipart/form-data; boundary="+boundary)
	//fmt.Println("TEST:create request form:", req.Header.Get("Content-Type"))
	return req
}

func createTestContext(res http.ResponseWriter, req *http.Request) *Context {
	ctx := NewContext()

	if req == nil {
		req, _ = http.NewRequest("GET", "/Ping/xiaoming", nil)
	}
	r := NewRequest()
	r.InitWithHttp(req)

	s := &Response{}
	s.InitSelf(res)

	ctx.SetRequest(nil)
	ctx.SetRequest(r)
	ctx.SetResponse(nil)
	ctx.SetResponse(s)

	ctx.SetAccess(nil)

	return ctx
}

func TestContextFormFile(t *testing.T) {
	buf := new(bytes.Buffer)
	mw := multipart.NewWriter(buf)
	w, err := mw.CreateFormFile("file", "test")
	if assert.NoError(t, err) {
		_, err = w.Write([]byte("test"))
		assert.NoError(t, err)
	}
	mw.Close()

	ctx := createTestContext(httptest.NewRecorder(), nil)

	req, _ := http.NewRequest("POST", "/", buf)
	req.Header.Set("Content-Type", mw.FormDataContentType())
	ctx.Request().InitWithHttp(req)

	f, err := ctx.Request().FormFile("file")

	assert.Nil(t, err)
	fmt.Println("file name", f.Filename)
	SaveUploadFile(f, "/tmp/glog/test")
}

func TestContextReset(t *testing.T) {

	ctx := createTestContext(httptest.NewRecorder(), nil)
	ctx.Reset()
	//ctx.Abort()
	assert.Nil(t, ctx.access)
	assert.Nil(t, ctx.GetAccess())

	//assert.Nil(t, ctx.stack_error)
	//assert.Nil(t, ctx.GetStackError())

	//assert.Nil(t, ctx.ext_values)

	assert.NotNil(t, ctx.Request())
	assert.NotNil(t, ctx.Response())

	assert.Panics(t, func() { ctx.Abort() })
	assert.Panics(t, func() { ctx.Exit() })

	assert.Panics(t, func() { ctx.AbortStatus(404) })
}

func TestContextSetterGetter(t *testing.T) {
	ctx := createTestContext(httptest.NewRecorder(), nil)
	ctx.Set("k1", "value1")
	ctx.Set("k2", 255255)
	ctx.Set("k3", int64(255255))
	assert.Equal(t, "value1", ctx.GetString("k1"))
	assert.Equal(t, 255255, ctx.GetInt("k2"))
	assert.Equal(t, 255255, (ctx.GetInt("k3")))
	assert.Equal(t, int64(255255), ctx.GetInt64("k3"))

	// test Get
	iv, ret := ctx.Get("k1")
	assert.True(t, ret)
	ivv, ok := iv.(string)
	assert.True(t, ok)
	assert.Equal(t, "value1", ivv)

	ctx.SetSameSite(1)

	assert.Nil(t, ctx.GetAccess())
}

func TestContextClientIP(t *testing.T) {

	ctx := createTestContext(httptest.NewRecorder(), nil)
	req := ctx.Request().GetHttpRequest()
	req.Header.Set("X-Real-Ip", "192.168.3.10")
	ip := ctx.ClientIP()
	//fmt.Println("test client ip", ip)
	assert.Equal(t, "192.168.3.10", ip)

	req.Header.Set("X-Forwarded-For", "192.168.3.20, 192.168.3.30")
	//X-Appengine-Remote-Addr
	req.Header.Set("X-Appengine-Remote-Addr", "192.168.3.50")
	req.RemoteAddr = "192.168.3.40:80"
	assert.Equal(t, "192.168.3.20", ctx.ClientIP())

	req.Header.Del("X-Forwarded-For")
	assert.Equal(t, "192.168.3.10", ctx.ClientIP())
	req.Header.Del("X-Real-Ip")
	assert.Equal(t, "192.168.3.50", ctx.ClientIP())
	req.Header.Del("X-Appengine-Remote-Addr")
	//ctx.Request().GetRemoteAddr()
	//assert.Equal(t, "192.168.3.40", ctx.Request().GetRemoteAddr())
	assert.Equal(t, "192.168.3.40", ctx.ClientIP())

	ip, _, err := net.SplitHostPort(strings.TrimSpace(ctx.Request().GetRemoteAddr()))
	fmt.Println("TEST:Context ClientIP", err)
}

func TestContextSetParam(t *testing.T) {
	ctx := createTestContext(httptest.NewRecorder(), nil)
	ctx.SetParam("xiaohai", "i lov")

	p1, _ := ctx.Request().GetString("xiaohai")
	assert.Equal(t, "i lov", p1)
}
