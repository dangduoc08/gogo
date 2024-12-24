package dogs

import "github.com/dangduoc08/gogo/core"

var DogModule = func() *core.Module {
	var module = core.ModuleBuilder().
		Providers(
			DogProvider{},
		).
		Controllers(
			DogController{},
		).
		Build()

	module.
		Prefix("dogs")

	return module
}
