package shared

import (
	"encoding/json"
	"net/http"

	"github.com/dangduoc08/gogo"
	"github.com/dangduoc08/gogo/common"
	"github.com/dangduoc08/gogo/ctx"
	"github.com/dangduoc08/gogo/exception"
)

type LoggingInterceptor struct {
	Logger common.Logger
}

func (instance LoggingInterceptor) Intercept(c gogo.Context, aggregation gogo.Aggregation) any {
	datas := []any{}

	if c.Method == http.MethodPost || c.Method == http.MethodPut || c.Method == http.MethodPatch {
		body := c.Body()
		if len(body) > 0 {
			bodyJSON, _ := json.Marshal(body)
			bodyJSONStr := string(bodyJSON)
			datas = append(datas, "body", bodyJSONStr)
		} else {
			formMap := c.Form()
			if len(formMap) > 0 {
				formJSON, _ := json.Marshal(formMap)
				formJSONStr := string(formJSON)
				datas = append(datas, "body", formJSONStr)
			} else {
				datas = append(datas, "body", nil)
			}
		}
	}

	queryMap := c.Query()
	if len(queryMap) > 0 {
		queryJSON, _ := json.Marshal(queryMap)
		queryJSONStr := string(queryJSON)
		datas = append(datas, "query", queryJSONStr)
	} else {
		datas = append(datas, "query", nil)
	}

	paramMap := c.Param()
	if len(paramMap) > 0 {
		paramJSON, _ := json.Marshal(c.Param())
		paramJSONStr := string(paramJSON)
		datas = append(datas, "param", paramJSONStr)
	} else {
		datas = append(datas, "param", nil)
	}

	datas = append(datas, ctx.REQUEST_ID, c.GetID())
	instance.Logger.Info(
		"RequestData",
		datas...,
	)

	return aggregation.Pipe(
		aggregation.Consume(func(c gogo.Context, data any) any {
			resJSON, _ := json.Marshal(data)
			resJSONStr := string(resJSON)
			if resJSONStr != "{}" {
				instance.Logger.Info(
					"SuccessResponse",
					"data", resJSONStr,
					ctx.REQUEST_ID, c.GetID(),
				)
			} else {
				instance.Logger.Info(
					"SuccessResponse",
					"data", nil,
					ctx.REQUEST_ID, c.GetID(),
				)
			}
			return data
		}),
		aggregation.Error(func(c gogo.Context, err any) any {
			if exception, ok := err.(exception.Exception); ok {
				instance.Logger.Debug(
					"ErrorResponse",
					"data", err,
					"message", exception.GetResponse(),
					ctx.REQUEST_ID, c.GetID(),
				)
			} else {
				instance.Logger.Debug(
					"ErrorResponse",
					"data", err,
					ctx.REQUEST_ID, c.GetID(),
				)
			}

			return nil
		}),
	)
}
