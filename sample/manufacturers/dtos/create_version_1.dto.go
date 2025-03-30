package dtos

import (
	"fmt"

	"github.com/dangduoc08/gogo"
	"github.com/dangduoc08/gogo/common"
)

type CREATE_VERSION_1_Data_Body_TO struct {
	Name   string   `bind:"name"`
	Scores []string `bind:"scores"`
}

type CREATE_VERSION_1_Body_DTO struct {
	Data CREATE_VERSION_1_Data_Body_TO `bind:"data"`
}

func (instance CREATE_VERSION_1_Body_DTO) Transform(body gogo.Body, medata common.ArgumentMetadata) any {
	fmt.Println("[Module] CREATE_VERSION_1_Body dto")
	dto, _ := body.Bind(instance)

	return dto
}
