package cats

import (
	"github.com/dangduoc08/gogo"
	"github.com/dangduoc08/gogo/common"
)

type CatExceptionFilter struct {
	common.Logger
}

func (instance CatExceptionFilter) NewExceptionFilter() CatExceptionFilter {
	return instance
}

func (instance CatExceptionFilter) Catch(c gogo.Context, ex gogo.Exception) {
	c.JSON(gogo.Map{
		"error": ex.GetResponse(),
	})
}
