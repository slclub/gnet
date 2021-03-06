package gnet

import (
	//"fmt"
	"github.com/slclub/link"
	"github.com/slclub/utils"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

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

// 32M
const default_multipart_memory_size = 32 << 20

var _ IRequest = &Request{}

type Request struct {
	http_request *http.Request

	// request type
	content_type string
	accepted     []string

	// come from form
	// here just quote relationship form.
	form_params url.Values
	form_state  bool

	// come from url query params.
	// just quote relationship.
	query_params url.Values
	query_state  bool

	// the param come from path
	params map[string]string
}

func NewRequest() IRequest {
	rq := &Request{
		//query_params: make(url.Values, 1024),
		form_state: false,
	}
	return rq
}

// impplement RequestParameter-----------------------------------------------------
func (req *Request) GetString(key string, args ...string) (value string, ret bool) {
	// first find from path param.
	// highest prioity.
	value, ret = req.params[key]
	if ret {
		return
	}
	// query form finded and return
	req.initForm()
	value, ret = utils.GetStringFromUrl(req.form_params, key)
	if ret {
		return
	}
	// query from query param.
	//req.once.Do(func() {
	req.initQuery()
	//link.DEBUG_PRINT("[GNET][REQUEST][GetString]", "initQuery")
	//})
	value, ret = utils.GetStringFromUrl(req.query_params, key)
	if ret {
		return
	}

	ret = false
	// Give value an default value.
	if len(args) >= 1 {
		value = args[0]
	}
	return
}

func (req *Request) GetArray(key string, args ...[]string) ([]string, bool) {
	req.initForm()
	value, ret := utils.GetArrayFromUrl(req.form_params, key)
	if ret {
		return value, ret
	}

	req.initQuery()
	value, ret = utils.GetArrayFromUrl(req.query_params, key)
	if ret {
		return value, ret
	}
	ret = false
	value = nil
	if len(args) >= 1 {
		value = args[0]
	}
	return value, ret
}

func (req *Request) GetMapString(key string, args ...map[string]string) (map[string]string, bool) {

	req.initForm()
	value, ret := utils.GetMapFromUrl(req.form_params, key)
	if ret {
		return value, ret
	}

	req.initQuery()
	value, ret = utils.GetMapFromUrl(req.query_params, key)
	if ret {
		return value, ret
	}

	ret = false
	value = nil
	if len(args) >= 1 {
		value = args[0]
	}
	return value, ret
}

func (req *Request) GetInt64(key string, args ...int64) (int64, bool) {
	value, ret := req.GetString(key)
	var default_val int64 = 0
	if len(args) >= 1 {
		default_val = args[0]
	}
	if !ret {
		return default_val, false
	}
	if rv, err := strconv.ParseInt(value, 10, 64); err == nil {
		return rv, true
	}

	return default_val, false
}

func (req *Request) GetInt(key string, args ...int) (int, bool) {
	rv, ret := req.GetInt64(key)
	return int(rv), ret
}

func (req *Request) BodyByte() ([]byte, error) {
	return ioutil.ReadAll(req.http_request.Body)
}

func (req *Request) SetParam(key, value string) {
	if req.params == nil {
		req.params = make(map[string]string)
	}
	req.params[key] = value
}

// impplement RequestParameter-----------------------------------------------------

// http---------------------------------------------------------------------------
func (req *Request) GetHttpRequest() *http.Request {
	return req.http_request
}

// http---------------------------------------------------------------------------
// init and reset-----------------------------------------------------------------
func (req *Request) InitWithHttp(hr *http.Request) {
	req.http_request = hr
	//req.initQuery()
	//req.initForm()
}

func (req *Request) Reset() {
	req.http_request = nil
	req.query_params = nil
	req.form_params = nil
	req.content_type = ""
	req.accepted = nil
	req.params = nil

	// reset query_param
	req.query_state = false
	//for k := range req.query_params {
	//	delete(req.query_params, k)
	//}
}

func (req *Request) initQuery() {

	if req.query_state {
		return
	}
	req.query_state = true
	req.query_params = req.http_request.URL.Query()
	//req.urlQuery()
	//link.DEBUG_PRINT("------request query param ", req.query_params == nil)
}

func (req *Request) initForm() {
	if req.form_params != nil || (req.http_request != nil && req.http_request.Method == http.MethodGet) {
		return
	}
	//size_valid := link.GetSize("form.multipart_memory", default_multipart_memory_size)
	size_valid := link.GetSizeInt64("form.multipart_memory", default_multipart_memory_size)

	req.form_params = make(url.Values)
	if err := req.http_request.ParseMultipartForm(size_valid); err != nil {
		link.ERROR("[FORM][MULTIPART_MEMORY][OVERFLOW]", "please check your form.multipart_memory config")
		//fmt.Println("error:init form", req.form_params, err, "content-type:", req.http_request.Header.Get("Content-Type"))
	}
	req.form_params = req.http_request.PostForm
	//link.DEBUG_PRINT("init form", req.form_params, req.http_request.PostForm, "\n")
}

// init and reset-----------------------------------------------------------------
// request header-----------------------------------------------------------------
func (req *Request) GetHeader(key string) string {
	// http
	return req.http_request.Header.Get(key)
}

func (req *Request) ContentType(args ...bool) string {
	force := false
	if len(args) >= 1 {
		force = args[0]
	}
	// http
	if req.content_type != "" && !force {
		return req.content_type
	}
	req.content_type = utils.GetPartFilterTrimOrSemicolon(req.GetHeader("Content-Type"))
	return req.content_type
}

// request header-----------------------------------------------------------------
// request file-------------------------------------------------------------------
func (req *Request) FormFile(key string) (*multipart.FileHeader, error) {
	if req.http_request.MultipartForm == nil {
		size_valid := link.GetSizeInt64("form.multipart_memory", default_multipart_memory_size)
		if err := req.http_request.ParseMultipartForm(size_valid); err != nil {
			return nil, err
		}
	}
	f, fh, err := req.http_request.FormFile(key)
	if err != nil {
		return nil, err
	}
	f.Close()
	return fh, err
}

// request file-------------------------------------------------------------------

func (req *Request) GetRemoteAddr() string {
	return req.http_request.RemoteAddr
}

func (req *Request) urlQuery() (err error) {
	query := req.GetHttpRequest().URL.RawQuery
	for query != "" {
		key := query
		if i := strings.IndexAny(key, "&;"); i >= 0 {
			key, query = key[:i], key[i+1:]
		} else {
			query = ""
		}
		if key == "" {
			continue
		}
		value := ""
		if i := strings.Index(key, "="); i >= 0 {
			key, value = key[:i], key[i+1:]
		}
		key, err1 := url.QueryUnescape(key)
		if err1 != nil {
			if err == nil {
				err = err1
			}
			continue
		}
		value, err1 = url.QueryUnescape(value)
		if err1 != nil {
			if err == nil {
				err = err1
			}
			continue
		}
		// here the performance is very bad.
		req.query_params[key] = append(req.query_params[key], value)
	}
	//link.DEBUG_PRINT(len(req.query_params["name"]))
	return
}
