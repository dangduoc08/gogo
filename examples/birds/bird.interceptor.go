package birds

import (
	"reflect"

	"github.com/dangduoc08/gogo"
	"github.com/dangduoc08/gogo/common"
)

type BirdInterceptor struct {
	common.Logger
}

func (instance BirdInterceptor) Intercept(c gogo.Context, aggregation gogo.Aggregation) any {
	return aggregation.Pipe(
		aggregation.Consume(func(c gogo.Context, data any) any {
			return data
		}),
		aggregation.Error(func(c gogo.Context, data any) any {
			reflect.TypeOf(data)
			// c.JSON(gogo.Map{
			// 	"data": data,
			// })
			return nil
		}),
	)
}
