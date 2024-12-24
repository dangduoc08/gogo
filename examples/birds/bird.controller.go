package birds

import (
	"errors"
	"fmt"

	"github.com/dangduoc08/gogo"
	"github.com/dangduoc08/gogo/common"
	"github.com/dangduoc08/gogo/core"
	"github.com/dangduoc08/gogo/exception"
)

type BirdController struct {
	common.REST
	common.Guard
	common.Interceptor
	common.ExceptionFilter
	Name string `json:"name"`
}

func (instance BirdController) NewController() core.Controller {
	instance.BindGuard(
		BirdGuard{},
		instance.CREATE_VERSION_NEUTRAL,
	)

	// instance.BindInterceptor(
	// 	BirdInterceptor{},
	// 	instance.CREATE_VERSION_1,
	// )

	// instance.BindExceptionFilter(
	// 	BirdExceptionFilter2{},
	// 	instance.CREATE_VERSION_1,
	// )

	// instance.BindExceptionFilter(
	// 	BirdExceptionFilter{},
	// 	instance.CREATE_VERSION_1,
	// )

	return instance
}

var fileError = errors.New("File error")

func (instance BirdController) CREATE_VERSION_NEUTRAL(
	dto CREATE_VERSION_1_DTO,
	bodyDto CREATE_VERSION_1_Body_DTO,
	body gogo.Body,
	query gogo.Query,
	param gogo.Param,
) any {
	instance.Name = "bird controller"
	badGateway := exception.BadGatewayException(
		[]string{"very bad way", "Asdas"},
		// "customize errror",
		// fileError,
		exception.ExceptionOptions{
			// Description: "Error by",
			// Cause:       fileError,
		},
	)

	fmt.Println("try unwrap", badGateway, badGateway.Unwrap())
	fmt.Println("errors.Is", errors.Is(badGateway, fileError))

	panic(badGateway)
}
