package cats

import (
	"github.com/dangduoc08/gogo/core"
	"github.com/dangduoc08/gogo/examples/dogs"
)

type CatProvider struct {
	dogs.DogProvider
}

func (instance CatProvider) NewProvider() core.Provider {
	return instance
}
