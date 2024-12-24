package cats

import (
	"github.com/dangduoc08/gogo/core"
	"github.com/dangduoc08/gogo/examples/dogs"
)

var CatModule = func() *core.Module {
	var module = core.ModuleBuilder().
		Providers(
			CatProvider{},
			dogs.DogProvider{},
		).
		Controllers(
			CatController{},
		).
		Build()

	module.
		Prefix("cats")

	return module
}
