package gnet

import (
	//"io"
	"github.com/slclub/gnet/permission"
	//"net/http"
)

// ====================================================================================
// should move this interface to an Single package. Dont confuse with router.
type Contexter interface {
	// Extend ParameterArray
	ParameterArray

	IContextRequest
	IContextResponse
	Reset()
	// server code http or rpc code
	Status(code int)
	Aborter
	// access
	// GetAccess() permission.Accesser
	// SetAccess(permission.Accesser)
	permission.IAccess
}

type Param interface {
	GetKey(string)
	GetValue() string
	Set(key string, value string)
}

type ParameterArray interface {
	Get(string) interface{}
	GetString(string) string
	GetAll() []Param
	SetParam(key string, value interface{})
}

// request from other place.
type IContextRequest interface {
	Request() IRequest
	SetRequest(IRequest) bool
}

// response to client or other server.
type IContextResponse interface {
	Response() IResponse
	SetResponse(IResponse)
}

// interrupt interface.
type Aborter interface {
	// abort current handle .
	Abort()
	IsAbort() bool
	// Jump out of the whole execution process.
	// break whole work flow.
	IsExist() bool
	Exit()
}

//====================================================================================
//
type Executer interface {
	Execute(Contexter)
}

//=====================================================================================

// router haddle func
type HandleFunc func(Contexter)
