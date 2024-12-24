package birds

import (
	"fmt"

	"github.com/dangduoc08/gogo"
	"github.com/dangduoc08/gogo/common"
)

type CREATE_VERSION_1_DTO struct {
	Limit int `bind:"limit"`
}

func (instance CREATE_VERSION_1_DTO) Transform(query gogo.Query, meta common.ArgumentMetadata) any {

	_dto, fielLevels := query.Bind(instance)

	dto := _dto.(CREATE_VERSION_1_DTO)

	for i, e := range fielLevels {
		fmt.Println(i, e.Namespace(), e.NestedTag())
	}

	return dto
}
