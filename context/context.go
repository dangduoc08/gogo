package context

import (
	"fmt"
	"net/http"
	"time"

	"github.com/dangduoc08/gooh/utils"
)

type (
	Map     map[string]any
	ErrFn   func(error)
	Handler = func(c *Context)
)

type Responser interface {
	Set(map[string]string) Responser
	Status(int) Responser
	Text(string, ...any)
	JSONP(...any)
	JSON(...any)
	Param() Values
	Error(ErrFn)
}

type Context struct {
	*http.Request
	http.ResponseWriter

	param     Values
	paramKeys map[string][]int
	paramVals []string

	Next      func()
	Event     *event
	Code      int
	Timestamp time.Time
}

func NewContext() *Context {
	return &Context{
		Code:      http.StatusOK,
		Timestamp: time.Now(),
	}
}

func (c *Context) SetParamKeys(paramKeys map[string][]int) {
	c.paramKeys = paramKeys
}

func (c *Context) SetParamVals(paramVals []string) {
	c.paramVals = paramVals
}

func (c *Context) Set(pair map[string]string) Responser {
	for key, value := range pair {
		c.ResponseWriter.Header().Set(key, value)
	}

	return c
}

func (c *Context) Status(code int) Responser {
	c.Code = code

	return c
}

func (c *Context) Text(content string, args ...any) {
	c.WriteHeader(c.Code)
	fmt.Fprintf(c.ResponseWriter, content, args...)
	c.Event.Emit(REQUEST_FINISHED)
}

func (c *Context) JSON(args ...any) {
	buf, err := handleJSON(args...)
	if err != nil {
		panic(err.Error())
	}

	c.Set(map[string]string{
		"Content-Type": "application/json",
	})
	c.WriteHeader(c.Code)
	c.ResponseWriter.Write(buf)
	c.Event.Emit(REQUEST_FINISHED)
}

func (c *Context) JSONP(args ...any) {
	cb := utils.StrRemoveSpace(c.URL.Query().Get("callback"))
	if cb == "" {
		c.JSON(args...)
		return
	}

	buf, err := handleJSON(args...)
	if err != nil {
		panic(err.Error())
	}
	c.Set(map[string]string{
		"Content-Type": "text/javascript; charset=utf-8",
	})
	c.Text(buildJSONP(string(buf), cb))
}

func (c *Context) Error(cb ErrFn) {
	if rec := recover(); rec != nil {
		c.Status(http.StatusInternalServerError)
		cb(rec.(error))
		c.Event.Emit(REQUEST_FINISHED)
	}
}

func (c *Context) Reset() {
	c.Code = http.StatusOK
	c.Next = nil
}