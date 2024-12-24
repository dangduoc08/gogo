package birds

import (
	"github.com/dangduoc08/gogo"
	"github.com/dangduoc08/gogo/common"
)

type BirdGuard struct {
	common.Logger
}

func (instance BirdGuard) NewGuard() BirdGuard {
	return instance
}

func (instance BirdGuard) CanActivate(c gogo.Context) bool {
	return true
}
