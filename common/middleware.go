package common

import (
	"fmt"

	"github.com/dangduoc08/gogo/ctx"
)

type MiddlewareFn = ctx.Handler

type RESTMiddlewareItem struct {
	Method  string
	Route   string
	Version string
	Pattern string
	Common  CommonItem
}

type WSMiddlewareItem struct {
	EventName string
	Common    CommonItem
}

type MiddlewareItem struct {
	REST RESTMiddlewareItem
	WS   WSMiddlewareItem
}

type MiddlewareHandler struct {
	MiddlewareFn MiddlewareFn
	Handlers     []any
}

type Middleware struct {
	MiddlewareHandlers []MiddlewareHandler
}

func (mw *Middleware) Apply(middleware ctx.Handler, handlers ...any) *Middleware {
	middlewareHandler := MiddlewareHandler{
		MiddlewareFn: middleware,
		Handlers:     handlers, // handlers = nil => apply for all handler
	}

	mw.MiddlewareHandlers = append(mw.MiddlewareHandlers, middlewareHandler)

	return mw
}

func (mw *Middleware) GenerateRESTMiddlewares(r *REST) []MiddlewareItem {
	middlewareItemArr := []MiddlewareItem{}
	shouldAddMiddleware := map[string]bool{}
	fmt.Println(" r.FnNameToPatternMap ", r.FnNameToPatternMap, "\n")
	for _, middlewareHandler := range mw.MiddlewareHandlers {
		// middlewareName := GetFnName(middlewareHandler.MiddlewareFn)

		for _, mainHandler := range middlewareHandler.Handlers {
			mainHandlerName := GetFnName(mainHandler)
			// fmt.Println("this is route", r.FnNameToPatternMap[mainHandlerName])
			if pattern, ok := r.FnNameToPatternMap[mainHandlerName]; ok {
				shouldAddMiddleware[pattern] = true
			}

			// for pattern := range r.PatternToFnNameMap {
			// 	fmt.Println("this is pattern", pattern, shouldAddMiddleware)
			// }

			// if pattern, ok := r.FnNameToPatternMap[mainHandlerName]; ok {
			// 	method, route, version := routing.PatternToMethodRouteVersion(pattern)
			// 	httpMethod := routing.OperationsMapHTTPMethods[method]

			// 	middlewareItemArr = append(middlewareItemArr, MiddlewareItem{
			// 		REST: RESTMiddlewareItem{
			// 			Method:  httpMethod,
			// 			Route:   routing.ToEndpoint(route),
			// 			Version: version,
			// 			Pattern: pattern,
			// 			Common: CommonItem{
			// 				Handler:         middlewareHandler.MiddlewareFn,
			// 				Name:            middlewareName,
			// 				MainHandlerName: mainHandlerName,
			// 			},
			// 		},
			// 	})
			// }
		}
	}

	return middlewareItemArr
}
