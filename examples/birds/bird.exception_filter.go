package birds

import (
	"fmt"

	"github.com/dangduoc08/gogo"
	"github.com/dangduoc08/gogo/common"
	"github.com/dangduoc08/gogo/exception"
)

type BirdExceptionFilter struct {
	common.Logger
}

func (instance BirdExceptionFilter) NewExceptionFilter() BirdExceptionFilter {
	return instance
}

func (instance BirdExceptionFilter) Catch(c gogo.Context, ex gogo.Exception) {
	fmt.Println("catch 1")

	fmt.Println(ex.GetResponse())
	fmt.Println(ex.GetHTTPStatus())
	panic(exception.UnprocessableEntityException("from catch 1"))
	c.JSON(gogo.Map{
		"error": ex.GetResponse(),
	})
}

type BirdExceptionFilter2 struct {
	common.Logger
}

func (instance BirdExceptionFilter2) NewExceptionFilter() BirdExceptionFilter2 {
	return instance
}

func (instance BirdExceptionFilter2) Catch(c gogo.Context, ex gogo.Exception) {
	fmt.Println("catch 2")
	fmt.Println(ex.GetResponse())
	fmt.Println(ex.GetHTTPStatus())
	c.JSON(gogo.Map{
		"error": ex.GetResponse(),
	})
}
