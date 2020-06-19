# gnet

[![GoTest](https://api.travis-ci.org/slclub/gnet.svg)](https://github.com/slclub/gnet/actions)
[![codecov](https://codecov.io/gh/slclub/gnet/branch/master/graph/badge.svg)](https://codecov.io/gh/slclub/gnet)

An go server web framework.

## Summary

Context object. Include request body response body and some common methods


## Request
	
Request object. Obtain all kinds of requested information and parameter routes in the object here

[Source Code](https://github.com/slclub/gnet/blob/master/request.go), You can read it from this link.


### Request Paramters

You can use these methods get paramters from path param, query, form, and so on.

You don't care about  where the paramters come from.

```go
type RequestParameter interface {
	// Get param inetface ----------------------------------------------------------
	// just get string by key string.
	// q=a
	// return a
	GetString(key string, args ...string) (value string, ret bool)
	// q[]=a&q[]=b
	// return []string{a, b}
	GetArray(key string, args ...[]string) ([]string, bool)
	// q[a]=a&q[b]=b
	// return map[string]string{"a":"a", "b":"b"}
	GetMapString(key string, args ...map[string]string) (map[string]string, bool)
	GetInt64(key string, args ...int64) (int64, bool)
	GetInt(key string, args ...int) (int, bool)

	// set
	SetParam(key, value string)

	// input :// data
	BodyByte() ([]byte, error)
}
```

```go
	// example
	boy.R.GET(url, func(ctx gnet.Contexter){
		s1, ok := ctx.Request().GetString(key, default_value string)
	})
```

### Request Interface

```go
type IRequest interface {
	GetHttpRequest() *http.Request
	RequestParameter

	// init and reset
	InitWithHttp(*http.Request)
	Reset()

	// header
	GetHeader(key string) string
	ContentType(args ...bool) string
	GetRemoteAddr() string
	//file
	FormFile(key string) (*multipart.FileHeader, error)
}
```

```go
	// example
	boy.R.GET(url, func(ctx gnet.Contexter){
		s1 := ctx.Request().ContentType()
	})
```


## Response

Rewrite the interface of http.ResponseWriter. 

```go
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

```

```go
	// example
	boy.R.GET(url, func(ctx gnet.Contexter){
		ctx.Response().Write([]bytes{"hello girls"})
		ctx.Response().WriteString("hello girls")
	})

```

## Contexter Api

[Source Code](https://github.com/slclub/gnet/blob/master/context.go)

These methods been used in the same way.

```go
    // example
    boy.R.GET(url, func(ctx gnet.Contexter){
        ctx.Xxx(args ...)
    })

```


```go
	Reset()
	// gerror.StackError([]error) support:
	// Push(err error)
	// Pop() error
	// Size() int
	GetStackError() gerror.StackError
	SetSameSite(st http.SameSite)

	ClientIP() string
	//
	GetHandler() HandleFunc
	SetHandler(HandleFunc)

	GetExecute() Executer
	SetExecute(exe Executer)

	//redirect
	Redirect(location string, args ...int)
```

### Contexter Setter Getter

For custome key-value pairs extension.

```go
type SetterGetter interface {
	// setter
	Set(key string, val interface{})
	// getter
	Get(key string) (interface{}, bool)
	GetString(key string) string
	GetInt(key string) int
	GetInt64(key string) int64
}
```

### Contexter Get Set Request 

```go
// request from other place.
type IContextRequest interface {
	Request() IRequest
	SetRequest(IRequest) bool

	// cookie
	SetCookie(name, value string, args ...interface{})
	Cookie(string) (string, error)
}
```

### Contexter Get Set Response

```go
// response to client or other server.
type IContextResponse interface {
	Response() IResponse
	SetResponse(IResponse) bool
}
```

### Contexter Abort

Execution process jump control.

```go
// interrupt interface.
type Aborter interface {
	// abort current handle .
	Abort()
	AbortStatus(int)
	// Jump out of the whole execution process.
	// break whole work flow.
	Exit()
}
```

## Save File

Upload file.

```go
f1 = func(ctx gnet.Contexter) {
    f, err := ctx.Request().FormFile("file")
    gnet.SaveUploadFile(f, "/tmp/glog/test")
}
```
