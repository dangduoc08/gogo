package dogs

import (
	"github.com/dangduoc08/gogo"
	"github.com/dangduoc08/gogo/common"
)

type DogGuard struct {
	common.Logger
}

func (instance DogGuard) NewGuard() DogGuard {
	return instance
}

func (instance DogGuard) CanActivate(c gogo.Context) bool {
	return true
}
