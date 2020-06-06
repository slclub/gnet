package gnet

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	//"mime/multipart"
	//"net/http"
	"net/http/httptest"
	"testing"
)

func TestResponseReset(t *testing.T) {

	res := httptest.NewRecorder()
	ctx := createTestContext(res, nil)
	ctx.Response().Reset()
	assert.Equal(t, -1, ctx.Response().Size())
	assert.Equal(t, 200, ctx.Response().Status())

	fmt.Println("TEST:Response")

	// error write header
	ctx.Response().WriteHeader(-100)
	assert.Equal(t, 200, ctx.Response().Status())
	ctx.Response().WriteHeader(300)
	assert.Equal(t, 300, ctx.Response().Status())

	assert.False(t, ctx.Response().Written())

	ctx.Response().WriteHeader(404)
	assert.NotEqual(t, 404, res.Code)
	ctx.Response().Flush()
	assert.Equal(t, 404, res.Code)
}

func TestResponseFlushHeader(t *testing.T) {

	ctx := createTestContext(httptest.NewRecorder(), nil)
	ctx.Response().Reset()
	assert.Equal(t, -1, ctx.Response().Size())

	assert.False(t, ctx.Response().Written())
	ctx.Response().FlushHeader()
	assert.True(t, ctx.Response().Written())
	assert.True(t, 0 == ctx.Response().Size())
	fmt.Println("TEST:Response.FlushHeader size", ctx.Response().Size())
}

func TestResponseWrite(t *testing.T) {

	ctx := createTestContext(httptest.NewRecorder(), nil)
	// write byte
	ctx.Response().Reset()
	ctx.Response().Write([]byte{'a', 'b', 'c'})

	assert.Equal(t, 3, ctx.Response().Size())
	// write string
	ctx.Response().Reset()
	ctx.Response().WriteString("ILOVE")
	assert.Equal(t, 5, ctx.Response().Size())

	fmt.Println("PRINT:Response", ctx.Response())
}

func TestResponseHijacker(t *testing.T) {
	ctx := createTestContext(httptest.NewRecorder(), nil)
	ctx.Response().Reset()
	testWriter := httptest.NewRecorder()
	ctx.Response().InitSelf(testWriter)

	assert.Panics(t, func() {
		_, _, err := ctx.Response().Hijack()
		assert.NoError(t, err)
	})
	assert.True(t, ctx.Response().Written())

	assert.Panics(t, func() {
		ctx.Response().CloseNotify()
	})
	ctx.Response().Headers("http-1", "nothing")
	assert.Equal(t, "nothing", ctx.Response().Header().Get("http-1"))
	ctx.Response().Headers("http-1", "")
	assert.Empty(t, ctx.Response().Header().Get("http-1"))
	ctx.Response().Flush()

	pub := ctx.Response().Pusher()
	fmt.Println("TEST:Response.Pusher", pub)
}
