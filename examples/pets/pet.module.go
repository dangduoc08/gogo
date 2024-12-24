package pets

import (
	"github.com/dangduoc08/gogo/core"
	"github.com/dangduoc08/gogo/examples/birds"
	"github.com/dangduoc08/gogo/examples/cats"
	"github.com/dangduoc08/gogo/examples/dogs"
)

var PetModule = func() *core.Module {
	var module = core.ModuleBuilder().
		Imports(
			cats.CatModule,
			dogs.DogModule,
			birds.BirdModule,
		).
		Build()

	module.
		Prefix("pets")

	return module
}
