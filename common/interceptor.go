package common

import (
	"reflect"

	"github.com/dangduoc08/gooh/context"
	"github.com/dangduoc08/gooh/routing"
)

type Interceptable interface {
	Intercept(*context.Context, any) any
}

type InterceptorHandler struct {
	Interceptable Interceptable
	Handlers      []any
}

type InterceptorItem struct {
	Method  string
	Route   string
	Handler any
}

type Interceptor struct {
	InterceptorHandlers []InterceptorHandler
}

func (i *Interceptor) BindInterceptor(interceptable Interceptable, handlers ...any) *Interceptor {
	interceptorHandler := InterceptorHandler{
		Interceptable: interceptable,
		Handlers:      handlers,
	}

	i.InterceptorHandlers = append(i.InterceptorHandlers, interceptorHandler)

	return i
}

func (i *Interceptor) AddInterceptorsToModule(r *Rest, cb func(int, reflect.Type, reflect.Value, reflect.Value)) []InterceptorItem {
	interceptorItemArr := []InterceptorItem{}

	for _, interceptorHandler := range i.InterceptorHandlers {

		interceptableType := reflect.TypeOf(interceptorHandler.Interceptable)
		interceptableValue := reflect.ValueOf(interceptorHandler.Interceptable)
		newInterceptor := reflect.New(interceptableType)

		for i := 0; i < interceptableType.NumField(); i++ {

			// callback use to inject providers
			cb(i, interceptableType, interceptableValue, newInterceptor)
		}

		// invoke interceptor constructor
		// if NewInterceptor was declared
		newInterceptable := newInterceptor.Interface()
		newInterceptable = Construct(newInterceptable, "NewInterceptor")

		interceptorHandler.Interceptable = newInterceptable.(Interceptable)

		for pattern := range r.patternToFnNameMap {
			httpMethod, route := routing.SplitRoute(pattern)
			interceptorItemArr = append(interceptorItemArr, InterceptorItem{
				Method:  httpMethod,
				Route:   routing.ToEndpoint(route),
				Handler: interceptorHandler.Interceptable.Intercept,
			})
		}
	}

	return interceptorItemArr
}
