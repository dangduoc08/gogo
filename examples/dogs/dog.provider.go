package dogs

import (
	"github.com/dangduoc08/gogo/core"
)

type DogProvider struct{}

func (instance DogProvider) NewProvider() core.Provider {
	return instance
}
