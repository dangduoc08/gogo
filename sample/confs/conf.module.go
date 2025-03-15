package confs

import (
	"reflect"
	"strings"

	"github.com/dangduoc08/gogo/modules/config"
)

var ENV ConfModel

var ConfModule = config.Register(&config.ConfigModuleOptions{
	IsGlobal:          true,
	IsExpandVariables: true,
	Hooks: []config.ConfigHookFn{
		func(c config.ConfigService) {

			// transform to proper types
			dto, _ := c.Transform(ConfModel{})

			confDTO := dto.(ConfModel)

			if len(confDTO.DomainWhitelist) > 0 {
				confDTO.DomainWhitelist = strings.Split(confDTO.DomainWhitelist[0], ",")
			}
			ENV = confDTO

			// re-assign to config struct
			dtoConfigType := reflect.TypeOf(confDTO)
			for i := 0; i < dtoConfigType.NumField(); i++ {
				field := dtoConfigType.Field(i)
				fieldValue := reflect.ValueOf(confDTO).Field(i)
				envKey := field.Tag.Get("bind")

				if envKey != "" {
					c.Set(envKey, fieldValue.Interface())
				}
			}
		},
	},
})
