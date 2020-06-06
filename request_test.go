package gnet

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	//"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRequestGetString(t *testing.T) {

	// query request
	req, _ := http.NewRequest("GET", "/Ping/xiaoming?myname=xiaoming&sex=girl", nil)
	ctx := createTestContext(httptest.NewRecorder(), req)
	//ctx.Request().InitWithHttp(req)

	p1, ret := ctx.Request().GetString("myname")
	fmt.Println("test Request.GetString", p1)
	assert.Equal(t, "xiaoming", p1)
	p1, ret = ctx.Request().GetString("myname not exist")
	assert.Equal(t, "", p1)
	assert.False(t, ret)

	p1, ret = ctx.Request().GetString("myname not exist", "xiaohu")
	assert.Equal(t, "xiaohu", p1)

	p1, ret = ctx.Request().GetString("sex")
	assert.Equal(t, "girl", p1)

	//form request
	req = createMultipartRequest(t)
	ctx.Request().InitWithHttp(req)
	p1, ret = ctx.Request().GetString("foo")
	//fmt.Println("TEST:GetString  PostForm", req.PostForm, "new-Content-Type:", req.Header.Get("Content-Type"))
	assert.True(t, ret)
	assert.Equal(t, "bar", p1)

	// convert int success!
	n1, ret1 := ctx.Request().GetInt("bar")
	assert.True(t, ret1)
	assert.Equal(t, 10, n1)
	n2, ret1 := ctx.Request().GetInt64("bar")
	assert.True(t, ret1)
	assert.Equal(t, int64(10), n2)

	// convret int fatal
	n1, ret1 = ctx.Request().GetInt("foo")
	assert.False(t, ret1)
	assert.Equal(t, 0, n1)

	vb, _ := ctx.Request().BodyByte()
	fmt.Println("TEST:Request.BodyByte", len(vb))
}

func TestRequestGetArray(t *testing.T) {
	// query request
	req, _ := http.NewRequest("GET", "/Ping/xiaoming?myname[]=xiaoming&myname[]=xiaohu&sex=girl", nil)
	ctx := createTestContext(httptest.NewRecorder(), req)

	names, ret := ctx.Request().GetArray("myname")
	assert.True(t, ret)
	assert.Equal(t, 2, len(names))
	assert.Equal(t, "xiaoming", names[0])
	assert.Equal(t, "xiaohu", names[1])

	names, ret = ctx.Request().GetArray("myname not exist", []string{"none"})
	assert.False(t, ret)
	assert.Equal(t, 1, len(names))
	assert.Equal(t, "none", names[0])

}

func TestRequestGetMap(t *testing.T) {
	// query request
	req, _ := http.NewRequest("GET", "/Ping/xiaoming?myname[a]=xiaoming&myname[b]=xiaohu&sex=girl", nil)
	ctx := createTestContext(httptest.NewRecorder(), req)

	names, ret := ctx.Request().GetMapString("myname")

	assert.True(t, ret)
	assert.Equal(t, 2, len(names))
	assert.Equal(t, "xiaoming", names["a"])
	assert.Equal(t, "xiaohu", names["b"])

	names, ret = ctx.Request().GetMapString("myname not exist", map[string]string{"none": "none"})
	assert.False(t, ret)
	assert.Equal(t, 1, len(names))
	assert.Equal(t, "none", names["none"])
}

func TestRequestContentType(t *testing.T) {
	// query request
	req, _ := http.NewRequest("GET", "/Ping/xiaoming?myname[a]=xiaoming&myname[b]=xiaohu&sex=girl", nil)
	req.Header.Set("Content-Type", "application/html")
	ctx := createTestContext(httptest.NewRecorder(), req)

	con := ctx.Request().ContentType()
	assert.Equal(t, "application/html", con)

	req.Header.Set("Content-Type", "application/json")
	con = ctx.Request().ContentType()
	assert.Equal(t, "application/html", con)

	con = ctx.Request().ContentType(true)
	assert.Equal(t, "application/json", con)

	req.RemoteAddr = "10.20.30.10"
	assert.Equal(t, "10.20.30.10", ctx.Request().GetRemoteAddr())
}
