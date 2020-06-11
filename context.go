package gnet

import (
	//"io"
	"github.com/slclub/gerror"
	"github.com/slclub/gnet/defined"
	"github.com/slclub/gnet/permission"
	"github.com/slclub/utils"
	"net"
	"net/http"
	"net/url"
	"strings"
)

// ====================================================================================
// should move this interface to an Single package. Dont confuse with router.
type Contexter interface {
	// Extend ParameterArray
	IParameter

	IContextRequest
	IContextResponse
	// server code http or rpc code
	//Status(code int)
	Aborter
	// Running flow run  by permission that access offered.
	// GetAccess() permission.Accesser
	// SetAccess(permission.Accesser)
	permission.IAccess

	// Add defined object by yourself
	// Extend interface.
	SetterGetter

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
}

type SetterGetter interface {
	// setter
	Set(key string, val interface{})
	// getter
	Get(key string) (interface{}, bool)
	GetString(key string) string
	GetInt(key string) int
	GetInt64(key string) int64
}

type IParameter interface {
	SetParam(key string, value string)
}

// request from other place.
type IContextRequest interface {
	Request() IRequest
	SetRequest(IRequest) bool

	// cookie
	SetCookie(name, value string, args ...interface{})
	Cookie(string) (string, error)
}

// response to client or other server.
type IContextResponse interface {
	Response() IResponse
	SetResponse(IResponse) bool
}

// interrupt interface.
type Aborter interface {
	// abort current handle .
	Abort()
	AbortStatus(int)
	// Jump out of the whole execution process.
	// break whole work flow.
	Exit()
}

//=====================================================================================

// router haddle func
type HandleFunc func(Contexter)

var _ Contexter = &Context{}

//var _ IParameter = &Context{}
//var _ IContextRequest = &Context{}

// Context struct *********************************************************************
type Context struct {
	request     IRequest
	response    IResponse
	ext_values  map[string]interface{}
	same_site   http.SameSite
	stack_error gerror.StackError
	access      permission.Accesser
	handle      HandleFunc
	exec        Executer
}

func NewContext() *Context {
	return &Context{
		ext_values: make(map[string]interface{}),
	}
}

// ---------------------------------Parameter ----------------------------------------
func (ctx *Context) SetParam(key, val string) {
	ctx.Request().SetParam(key, val)
}

// ---------------------------------Parameter ----------------------------------------
// ---------------------------------Request Response----------------------------------
func (ctx *Context) Response() IResponse {
	return ctx.response
}

func (ctx *Context) SetResponse(res IResponse) bool {
	if res == nil {
		return false
	}
	ctx.response = res
	return true
}

func (ctx *Context) Request() IRequest {
	return ctx.request
}

func (ctx *Context) SetRequest(rq IRequest) bool {
	if rq == nil {
		return false
	}
	ctx.request = rq
	return true
}

func (ctx *Context) SetSameSite(st http.SameSite) {
	ctx.same_site = st
}

// ---------------------------------Request Response----------------------------------
// ---------------------------------Aborter-------------------------------------------
// like try catch.
func (ctx *Context) Abort() {
	ctx.Response().Flush()
	gerror.Panic(defined.CODE_JUMP_CURRENT_NODE, "[USER][ABORT]")
}

func (ctx *Context) AbortStatus(code int) {
	ctx.Response().WriteHeader(code)
	ctx.Response().FlushHeader()
	ctx.Abort()
}

// like try catch. throw an speical code
func (ctx *Context) Exit() {
	gerror.Panic(defined.CODE_JUMP_CURRENT_FLOW, "[USER][ABORT]")
}

// ---------------------------------Aborter-------------------------------------------
// ---------------------------------SetterGetter--------------------------------------
func (ctx *Context) Set(key string, obj interface{}) {
	ctx.ext_values[key] = obj
}

func (ctx *Context) Get(key string) (interface{}, bool) {
	val, ok := ctx.ext_values[key]
	return val, ok
}

func (ctx *Context) GetString(key string) (val string) {
	if value, ok := ctx.Get(key); ok && value != nil {
		val, _ = value.(string)
	}
	return
}

func (ctx *Context) GetInt(key string) (val int) {
	if value, ok := ctx.Get(key); ok && value != nil {
		val, ok = value.(int)
		if !ok {
			val64, _ := utils.ForceInt64(value)
			val = int(val64)
		}
	}
	return
}
func (ctx *Context) GetInt64(key string) (val int64) {
	if value, ok := ctx.Get(key); ok && value != nil {
		val, ok = value.(int64)
		if !ok {
			val, _ = utils.ForceInt64(value)
		}
	}
	return
}

// ---------------------------------SetterGetter--------------------------------------
// -----------------------------------access------------------------------------------
func (ctx *Context) GetAccess() permission.Accesser {
	return ctx.access
}
func (ctx *Context) SetAccess(access permission.Accesser) {
	ctx.access = access
}

// -----------------------------------access------------------------------------------

// *********************************cookie ******************************************
/** Set cookie
 * @param	name	string
 * @param	value	string
 * @param	MaxAge	int			args[0]
 * @param	path	string		args[1]
 * @param	domain	string		args[2]
 * @param	secure	bool		args[3]
 * @param	http_only	bool	args[4]
 * @return nil
 */
func (ctx *Context) SetCookie(name, value string, args ...interface{}) {
	max_age := 0
	path, domain := "", ""
	secure, http_only := true, true
	la := len(args)

	if la >= 1 {
		max_age, _ = args[0].(int)
	}
	if la >= 2 {
		path, _ = args[1].(string)
	}
	if la >= 3 {
		domain, _ = args[2].(string)
	}
	if la >= 4 {
		secure, _ = args[3].(bool)
	}
	if la >= 5 {
		http_only, _ = args[4].(bool)
	}

	if path == "" {
		path = "/"
	}

	http.SetCookie(ctx.Response(), &http.Cookie{
		Name:     name,
		Value:    value,
		MaxAge:   max_age,
		Path:     path,
		Domain:   domain,
		SameSite: ctx.same_site,
		Secure:   secure,
		HttpOnly: http_only,
	})
}

func (ctx *Context) Cookie(name string) (string, error) {
	cookie, err := ctx.Request().GetHttpRequest().Cookie(name)
	if err != nil {
		return "", err
	}
	val, _ := url.QueryUnescape(cookie.Value)
	return val, nil
}

// *********************************cookie ******************************************
func (ctx *Context) GetStackError() gerror.StackError {
	return ctx.stack_error
}

func (ctx *Context) Reset() {
	ctx.Request().Reset()
	ctx.Response().Reset()
	ctx.stack_error = make(gerror.StackError, 0)
	ctx.ext_values = make(map[string]interface{})
	ctx.access = nil
	//ctx.request = nil
	//ctx.response = nil
}

func (ctx *Context) ClientIP() string {
	clientIP := ctx.Request().GetHeader("X-Forwarded-For")
	clientIP = strings.TrimSpace(strings.Split(clientIP, ",")[0])
	if clientIP == "" {
		clientIP = strings.TrimSpace(ctx.Request().GetHeader("X-Real-Ip"))
	}
	if clientIP != "" {
		return clientIP
	}

	if addr := ctx.Request().GetHeader("X-Appengine-Remote-Addr"); addr != "" {
		return addr
	}

	if ip, _, err := net.SplitHostPort(strings.TrimSpace(ctx.Request().GetRemoteAddr())); err == nil {
		return ip
	}

	return ""
}

// *********************************handle and execute ******************************************
func (ctx *Context) GetHandler() HandleFunc {
	return ctx.handle
}

func (ctx *Context) SetHandler(handle HandleFunc) {
	ctx.handle = handle
}

func (ctx *Context) GetExecute() Executer {
	return ctx.exec
}
func (ctx *Context) SetExecute(exe Executer) {
	ctx.exec = exe
}

// Context struct *********************************************************************

//************************************************************************************

//====================================================================================
//
type Executer interface {
	Execute(Contexter)
}
