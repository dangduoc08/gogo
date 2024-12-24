package birds

import (
	"github.com/dangduoc08/gogo/core"
	"github.com/dangduoc08/gogo/examples/cats"
	"github.com/dangduoc08/gogo/examples/dogs"
)

var BirdModule = func() *core.Module {
	var module = core.ModuleBuilder().
		Providers(
			BirdProvider{},
			cats.CatProvider{},
			dogs.DogProvider{},
		).
		Controllers(
			BirdController{},
		).
		Build()

	module.
		Prefix("birds")

	return module
}
