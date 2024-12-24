package main

import (
	"github.com/dangduoc08/gogo/core"
	"github.com/dangduoc08/gogo/examples/pets"
	"github.com/dangduoc08/gogo/examples/shared"
	"github.com/dangduoc08/gogo/log"
	"github.com/dangduoc08/gogo/middlewares"
	"github.com/dangduoc08/gogo/modules/config"
	"github.com/dangduoc08/gogo/versioning"
)

func main() {
	app := core.New()
	logger := log.NewLog(&log.LogOptions{
		Level:     log.DebugLevel,
		LogFormat: log.PrettyFormat,
	})

	// pattern := `^https?:\/\/([a-zA-Z0-9-]+\.)+[a-zA-Z]{2,}(:\d+)?(\/.*)?$`

	app.
		UseLogger(logger).
		Use(
			middlewares.RequestLogger(logger),
			middlewares.CORS(
				&middlewares.CORSOptions{
					// IsPreflightContinue: true,
					AllowOrigin: []string{
						"http://localhost:3000",
					},
				})).
		BindGlobalInterceptors(shared.LoggingInterceptor{}, shared.ResponseInterceptor{})

	app.EnableVersioning(versioning.Versioning{
		Type: versioning.HEADER,
		Key:  "X-API-Version",
	})

	app.Create(
		core.ModuleBuilder().
			Imports(
				pets.PetModule,
				config.Register(&config.ConfigModuleOptions{
					IsGlobal:          true,
					IsExpandVariables: true,
					Hooks: []config.ConfigHookFn{
						func(c config.ConfigService) {
							c.Set("PORT", 4000)
						}},
				}),
			).
			Build(),
	)

	configService := app.Get(config.ConfigService{}).(config.ConfigService)

	app.Logger.Fatal("AppError", "error", app.Listen(configService.Get("PORT").(int)))
}

// func goBookstore() error {
// 	return errors.New("bookstore was on fire!!!")
// }

// var errBookstore = errors.New("bookstore was on fire!!!")

// func readBook() error {
// 	return fmt.Errorf("not find any books due to\n wrapped error: %w", errBookstore)
// 	// return errors.Join(errors.New("not find any books due to"), errBookstore)
// }

// func main() {
// 	// bookstoreErr := goBookstore()

// 	readBookErr := readBook()

// 	fmt.Println(errors.Is(readBookErr, errBookstore))

// 	panic(readBookErr)
// }
