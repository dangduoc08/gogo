package birds

import (
	"github.com/dangduoc08/gogo/core"
	"github.com/dangduoc08/gogo/examples/cats"
)

type BirdProvider struct {
	cats.CatProvider
}

func (instance BirdProvider) NewProvider() core.Provider {
	return instance
}
