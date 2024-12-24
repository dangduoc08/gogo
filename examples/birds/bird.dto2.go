package birds

import (
	"fmt"
	"time"

	"github.com/dangduoc08/gogo"
	"github.com/dangduoc08/gogo/common"
)

type CREATE_VERSION_1_Body_DTO struct {
	Limit     int       `bind:"Name"`
	Offset    int       `bind:"offset"`
	CreatedAt time.Time `bind:"date"`
}

func (instance CREATE_VERSION_1_Body_DTO) Transform(query gogo.Body, meta common.ArgumentMetadata) any {

	_dto, fielLevels := query.Bind(instance)

	dto := _dto.(CREATE_VERSION_1_Body_DTO)

	for i, e := range fielLevels {
		fmt.Println(i, e.Namespace(), e.NestedTag())
	}

	return dto
}
