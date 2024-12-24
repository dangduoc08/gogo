package dogs

import (
	"fmt"

	"github.com/dangduoc08/gogo"
	"github.com/dangduoc08/gogo/common"
	"github.com/dangduoc08/gogo/core"
)

type DogController struct {
	common.REST
	common.Guard
	common.Interceptor
	common.ExceptionFilter
}

func (instance DogController) NewController() core.Controller {
	instance.BindGuard(
		DogGuard{},
		instance.CREATE_VERSION_1,
	)

	instance.BindInterceptor(
		DogInterceptor{},
		instance.CREATE_VERSION_1,
	)

	instance.BindExceptionFilter(
		DogExceptionFilter{},
		instance.CREATE_VERSION_1,
	)

	return instance
}

func (instance DogController) CREATE_VERSION_1() any {
	return gogo.Map{
		"dog": "Tom",
	}
}

func (instance DogController) UPDATE_VERSION_1() any {
	return gogo.Map{
		"dog": "Tom",
	}
}

func (instance DogController) PREFLIGHT_VERSION_1() any {
	fmt.Println("zo")
	return gogo.Map{
		"dog": "Tom",
	}
}
