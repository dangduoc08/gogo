package common

import (
	"reflect"

	"github.com/dangduoc08/gooh/ctx"
	"github.com/dangduoc08/gooh/exception"
	"github.com/dangduoc08/gooh/routing"
)

type Catch = func(*ctx.Context, *exception.HTTPException)

type ExceptionFilterable interface {
	Catch(*ctx.Context, *exception.HTTPException)
}

type ExceptionFilterHandler struct {
	ExceptionFilterable ExceptionFilterable
	Handlers            []any
}

type ExceptionFilterItem struct {
	// for REST
	Method string
	Route  string

	// for WS
	EventName string

	Handler any
}

type ExceptionFilter struct {
	ExceptionFilterHandlers []ExceptionFilterHandler
}

func (e *ExceptionFilter) BindExceptionFilter(exceptionFilterable ExceptionFilterable, handlers ...any) *ExceptionFilter {
	exceptionFilterHandler := ExceptionFilterHandler{
		ExceptionFilterable: exceptionFilterable,
		Handlers:            handlers,
	}

	e.ExceptionFilterHandlers = append(e.ExceptionFilterHandlers, exceptionFilterHandler)
	return e
}

func (e *ExceptionFilter) InjectProvidersIntoRESTExceptionFilters(r *REST, cb func(int, reflect.Type, reflect.Value, reflect.Value)) []ExceptionFilterItem {
	exceptionFilterItemArr := []ExceptionFilterItem{}

	for _, exceptionFilterHandler := range e.ExceptionFilterHandlers {

		exceptionFilterableType := reflect.TypeOf(exceptionFilterHandler.ExceptionFilterable)
		exceptionFilterableValue := reflect.ValueOf(exceptionFilterHandler.ExceptionFilterable)
		newExceptionFilter := reflect.New(exceptionFilterableType)

		for i := 0; i < exceptionFilterableType.NumField(); i++ {

			// callback use to inject providers
			cb(i, exceptionFilterableType, exceptionFilterableValue, newExceptionFilter)
		}

		// invoke exceptionFilter constructor
		// if NewExceptionFilter was declared
		newExceptionFilterable := newExceptionFilter.Interface()
		newExceptionFilterable = Construct(newExceptionFilterable, "NewExceptionFilter")

		exceptionFilterHandler.ExceptionFilterable = newExceptionFilterable.(ExceptionFilterable)

		shouldAddExceptionFilter := map[string]bool{}
		for _, handler := range exceptionFilterHandler.Handlers {
			fnName := GetFnName(handler)
			httpMethod, route := ParseFnNameToURL(fnName, RESTOperations)
			route = r.addPrefixesToRoute(route, fnName, r.GetPrefixes())
			shouldAddExceptionFilter[routing.AddMethodToRoute(route, httpMethod)] = true
		}

		for pattern := range r.PatternToFnNameMap {
			if _, ok := shouldAddExceptionFilter[pattern]; ok || len(shouldAddExceptionFilter) == 0 {
				httpMethod, route := routing.SplitRoute(pattern)
				exceptionFilterItemArr = append(exceptionFilterItemArr, ExceptionFilterItem{
					Method:  httpMethod,
					Route:   routing.ToEndpoint(route),
					Handler: exceptionFilterHandler.ExceptionFilterable.Catch,
				})
			}
		}
	}

	return exceptionFilterItemArr
}

func (e *ExceptionFilter) InjectProvidersIntoWSExceptionFilters(ws *WS, cb func(int, reflect.Type, reflect.Value, reflect.Value)) []ExceptionFilterItem {
	exceptionFilterItemArr := []ExceptionFilterItem{}

	for _, exceptionFilterHandler := range e.ExceptionFilterHandlers {

		exceptionFilterableType := reflect.TypeOf(exceptionFilterHandler.ExceptionFilterable)
		exceptionFilterableValue := reflect.ValueOf(exceptionFilterHandler.ExceptionFilterable)
		newExceptionFilter := reflect.New(exceptionFilterableType)

		for i := 0; i < exceptionFilterableType.NumField(); i++ {

			// callback use to inject providers
			cb(i, exceptionFilterableType, exceptionFilterableValue, newExceptionFilter)
		}

		// invoke exceptionFilter constructor
		// if NewExceptionFilter was declared
		newExceptionFilterable := newExceptionFilter.Interface()
		newExceptionFilterable = Construct(newExceptionFilterable, "NewExceptionFilter")

		exceptionFilterHandler.ExceptionFilterable = newExceptionFilterable.(ExceptionFilterable)

		shouldAddExceptionFilter := map[string]bool{}
		for _, handler := range exceptionFilterHandler.Handlers {
			fnName := GetFnName(handler)
			_, eventName := ParseFnNameToURL(fnName, WSOperations)
			eventName = ToWSEventName(ws.subprotocol, eventName)
			shouldAddExceptionFilter[eventName] = true
		}

		for pattern := range ws.patternToFnNameMap {
			if _, ok := shouldAddExceptionFilter[pattern]; ok || len(shouldAddExceptionFilter) == 0 {
				exceptionFilterItemArr = append(exceptionFilterItemArr, ExceptionFilterItem{
					EventName: pattern,
					Handler:   exceptionFilterHandler.ExceptionFilterable.Catch,
				})
			}
		}
	}

	return exceptionFilterItemArr
}
