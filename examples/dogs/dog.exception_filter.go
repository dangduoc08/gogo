package dogs

import (
	"github.com/dangduoc08/gogo"
	"github.com/dangduoc08/gogo/common"
)

type DogExceptionFilter struct {
	common.Logger
}

func (instance DogExceptionFilter) NewExceptionFilter() DogExceptionFilter {
	return instance
}

func (instance DogExceptionFilter) Catch(c gogo.Context, ex gogo.Exception) {
	c.JSON(gogo.Map{
		"error": ex.GetResponse(),
	})
}
