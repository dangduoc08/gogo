package routing

import (
	"fmt"
	"net/http"
	"os"
	"testing"

	"github.com/dangduoc08/gogo/ctx"
	"github.com/dangduoc08/gogo/utils"
)

func TestRouteAdd(t *testing.T) {
	cases := []string{
		"/users/{userId}",
		"/feeds/all",
		"/users/{userId}/friends/all",
		"/schools/{schoolId}/subjects/{subjectId}/{schoolId}",
	}
	r := NewRouter()

	for _, path := range cases {
		r.Add(path, http.MethodGet, nil)
	}

	expected1 := 15
	actual1 := r.Trie.len()
	if actual1 != expected1 {
		t.Errorf(utils.ErrorMessage(actual1, expected1, "trie length should be equal"))
	}

	expected2 := map[string][]int{
		"schoolId":  {0, 2},
		"subjectId": {1},
	}
	actual2 := r.Children["schools"].Children["$"].Children["subjects"].Children["$"].Children["$"].Children[fromMethodtoPattern(http.MethodGet)].ParamKeys

	for key, indexs := range expected2 {
		if actual2[key] == nil {
			t.Errorf(utils.ErrorMessage(actual2[key], expected2, "params should not be null"))
		}

		for i, index := range indexs {
			if actual2[key][i] != index {
				t.Errorf(utils.ErrorMessage(actual2[key][i], index, "params index should be equal"))
			}
		}
	}
}

func TestRouterMatch(t *testing.T) {
	cases := map[string]string{
		"/lv1/ping":                             "/lv1/ping",
		"/lv1/642e948adef44c303cdd2df3":         "/lv1/{id1}",
		"/lv1/foo/bar":                          "/lv1/*",
		"/lv1":                                  "/lv1/*",
		"/lv1/lv2/pong":                         "/lv1/lv2/pong",
		"/lv1/lv2/642e951525714c1ec609b338":     "/lv1/lv2/{id2}",
		"/lv1/lv2/foo/bar":                      "/lv1/lv2/*",
		"/lv1/lv2/file/index.html":              "lv1/lv2/file/in*.html",
		"/lv1/lv2/file/in.html":                 "lv1/lv2/file/in*.html",
		"/lv1/lv2/file/image.png":               "lv1/lv2/file/image.*",
		"/lv1/lv2/lv3/peng":                     "/lv1/lv2/lv3/peng",
		"/lv1/lv2/lv3/642e95c4fbb2ad847ca96840": "/lv1/lv2/lv3/{id3}",
		"/lv1/lv2/lv3/foo/bar":                  "/lv1/lv2/lv3/*",
		"/lv1/lv2/lv3":                          "/lv1/lv2/lv3/*",
		"/lv1/lv2/lv3/file/index.html":          "lv1/lv2/lv3/file/in*.html",
		"/lv1/lv2/lv3/file/in.html":             "lv1/lv2/lv3/file/in*.html",
		"/lv1/lv2/lv3/file/image.jpeg":          "lv1/lv2/lv3/file/image.*",
		"/api/feeds/{feedApiId}/next/files/index.html/endpoint/any/things/after": "/*/feeds/{feed*Id}/*/files/*.html/*/",
		"/users/633b0aa5d7fc3578b655b9bd/friends/633b0af45f4fe7d45b00fba5":       "/users/{userId}/friends/{friendId}",
	}

	r := NewRouter()
	for _, path := range cases {
		for _, httpMethod := range HTTPMethods {
			r.Add(path, httpMethod, nil)
		}
	}

	for requestedRoute, expectedRoute := range cases {
		expectedRoute = AddMethodToRoute(expectedRoute, http.MethodPost)
		_, actualRoute, _, _, _ := r.Match(requestedRoute, http.MethodPost)

		if actualRoute != expectedRoute {
			t.Errorf(utils.ErrorMessage(actualRoute, expectedRoute, "routes should be matched"))
		}
	}
}

func TestRouterGroup(t *testing.T) {
	r1 := NewRouter()
	case1 := []string{
		"/users/get",
		"/users/get/{userId}",
	}
	for _, route := range case1 {
		for _, httpMethod := range HTTPMethods {
			r1.Add(route, httpMethod, nil)
		}
	}

	r2 := NewRouter()
	case2 := []string{
		"/users/update/{userId}",
		"/users/delete/{userId}",
	}
	for _, route := range case2 {
		for _, httpMethod := range HTTPMethods {
			r2.Add(route, httpMethod, nil)
		}
	}

	gr := NewRouter()
	gr.Group("/v1", r1, r2)

	_, actualRoute1, _, _, _ := gr.Match("/v1/users/update/123", http.MethodPatch)
	expectedRoute1 := AddMethodToRoute("/v1"+case2[0], http.MethodPatch)
	if actualRoute1 != expectedRoute1 {
		t.Errorf(utils.ErrorMessage(actualRoute1, expectedRoute1, "routes should be matched"))
	}
}

func TestRouterMiddleware(t *testing.T) {
	counter := 0

	handler1 := func(c *ctx.Context) {
		counter++
		c.Next()
	}

	handler2 := func(c *ctx.Context) {
		counter += 2
		c.Next()
	}

	handler3 := func(c *ctx.Context) {
		counter += 3
		c.Next()
	}

	handler4 := func(c *ctx.Context) {
		counter += 4
		c.Next()
	}

	r0 := NewRouter()
	r0.Use(handler1)
	r0.For("/test0", HTTPMethods)(handler1)
	for _, httpMethod := range HTTPMethods {
		r0.Add("/test0", httpMethod, handler1)
	}

	_, _, _, _, handlers := r0.Match("/test0", http.MethodHead)

	if len(handlers) != 3 {
		t.Errorf(utils.ErrorMessage(len(handlers), 3, "router 0 handlers total should be equal"))
	}

	isNext := true
	c := ctx.NewContext()
	c.Next = func() {
		isNext = true
	}
	for i, handler := range handlers {
		if isNext {
			isNext = false
			handler(c)

			if i == 0 && counter != 1 {
				t.Errorf(utils.ErrorMessage(counter, 1, "router 0 handlers increase counter should be equal"))
			}

			if i == 1 && counter != 2 {
				t.Errorf(utils.ErrorMessage(counter, 2, "router 0 handlers increase counter should be equal"))
			}

			if i == 2 && counter != 3 {
				t.Errorf(utils.ErrorMessage(counter, 3, "router 0 handlers increase counter should be equal"))
			}
		}
	}

	r1 := NewRouter()
	r1.Use(handler1)
	r1.Use(handler2)
	for _, httpMethod := range HTTPMethods {
		r1.Add("/test1", httpMethod, handler4)
	}
	r1.For("/test1", HTTPMethods)(handler3)
	r1.Use(handler1)

	_, _, _, _, handlers = r1.Match("/test1", http.MethodPatch)

	if len(handlers) != 5 {
		t.Errorf(utils.ErrorMessage(len(handlers), 5, "router 1 handlers total should be equal"))
	}

	isNext = true
	c = ctx.NewContext()
	c.Next = func() {
		isNext = true
	}
	for i, handler := range handlers {
		if isNext {
			isNext = false
			handler(c)

			if i == 0 && counter != 4 {
				t.Errorf(utils.ErrorMessage(counter, 4, "router 1 handlers increase counter should be equal"))
			}

			if i == 1 && counter != 6 {
				t.Errorf(utils.ErrorMessage(counter, 6, "router 1 handlers increase counter should be equal"))
			}

			if i == 2 && counter != 10 {
				t.Errorf(utils.ErrorMessage(counter, 10, "router 1 handlers increase counter should be equal"))
			}

			if i == 3 && counter != 13 {
				t.Errorf(utils.ErrorMessage(counter, 13, "router 1 handlers increase counter should be equal"))
			}

			if i == 4 && counter != 14 {
				t.Errorf(utils.ErrorMessage(counter, 14, "router 1 handlers increase counter should be equal"))
			}
		}
	}

	r2 := NewRouter()
	r2.For("/test2/{param}", HTTPMethods)(handler1)
	r2.Use(handler2)
	for _, httpMethod := range HTTPMethods {
		r2.Add("/test2/{param}", httpMethod, handler3)
	}
	r2.For("/test2/{param}", HTTPMethods)(handler4)

	_, _, _, _, handlers = r2.Match("/test2/123", http.MethodOptions)

	if len(handlers) != 4 {
		t.Errorf(utils.ErrorMessage(len(handlers), 4, "router 2 handlers total should be equal"))
	}

	isNext = true
	c = ctx.NewContext()
	c.Next = func() {
		isNext = true
	}
	for i, handler := range handlers {
		if isNext {
			isNext = false
			handler(c)

			if i == 0 && counter != 15 {
				t.Errorf(utils.ErrorMessage(counter, 15, "router 2 handlers increase counter should be equal"))
			}

			if i == 1 && counter != 17 {
				t.Errorf(utils.ErrorMessage(counter, 17, "router 2 handlers increase counter should be equal"))
			}

			if i == 2 && counter != 20 {
				t.Errorf(utils.ErrorMessage(counter, 20, "router 2 handlers increase counter should be equal"))
			}

			if i == 3 && counter != 24 {
				t.Errorf(utils.ErrorMessage(counter, 24, "router 2 handlers increase counter should be equal"))
			}
		}
	}

	gr := NewRouter()
	for _, httpMethod := range HTTPMethods {
		gr.Add("/group/test1", httpMethod, handler3)
	}
	gr.Use(handler4).Use(handler2).Use(handler1)
	gr.Group("/group", r1, r2)
	gr.For("/group/test2/{param}", HTTPMethods)(handler3)

	_, _, _, _, handlers = gr.Match("/group/test2/123", http.MethodOptions)

	if len(handlers) != 8 {
		t.Errorf(utils.ErrorMessage(len(handlers), 8, "router group handlers total should be equal"))
	}

	isNext = true
	c = ctx.NewContext()
	c.Next = func() {
		isNext = true
	}
	for i, handler := range handlers {
		if isNext {
			isNext = false
			handler(c)

			if i == 0 && counter != 28 {
				t.Errorf(utils.ErrorMessage(counter, 28, "router group handlers increase counter should be equal"))
			}

			if i == 1 && counter != 30 {
				t.Errorf(utils.ErrorMessage(counter, 30, "router group handlers increase counter should be equal"))
			}

			if i == 2 && counter != 31 {
				t.Errorf(utils.ErrorMessage(counter, 31, "router group handlers increase counter should be equal"))
			}

			if i == 3 && counter != 32 {
				t.Errorf(utils.ErrorMessage(counter, 32, "router group handlers increase counter should be equal"))
			}

			if i == 4 && counter != 34 {
				t.Errorf(utils.ErrorMessage(counter, 34, "router group handlers increase counter should be equal"))
			}

			if i == 5 && counter != 37 {
				t.Errorf(utils.ErrorMessage(counter, 37, "router group handlers increase counter should be equal"))
			}

			if i == 6 && counter != 41 {
				t.Errorf(utils.ErrorMessage(counter, 41, "router group handlers increase counter should be equal"))
			}

			if i == 7 && counter != 44 {
				t.Errorf(utils.ErrorMessage(counter, 44, "router group handlers increase counter should be equal"))
			}
		}
	}

	if counter != 44 {
		t.Errorf(utils.ErrorMessage(counter, 44, "final counter should be equal"))
	}
}

func TestRouteToJSON(t *testing.T) {
	if os.Getenv("v") == "true" {
		paths := []string{
			"/users/{userId}",
			"/feeds/all",
			"/users/{userId}/friends/{friendId}",
			"/schools/{schoolId}/subjects/{subjectId}/{subjectId}",
			"/schools/*",
			"/*/feeds/{feed***Id}/**/files/*.html/***",
		}
		r := NewRouter()

		for _, path := range paths {
			r.Add(path, "", func(c *ctx.Context) {})
		}

		json, err := r.ToJSON()
		if err != nil {
			fmt.Printf("Error: %s", err.Error())
		} else {
			fmt.Println(json)
		}
	} else {
		t.Skip()
	}

}
