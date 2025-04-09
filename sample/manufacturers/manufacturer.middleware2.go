package manufacturers

import (
	"fmt"

	"github.com/dangduoc08/gogo/ctx"
)

func ManufacturerMiddleware2(c *ctx.Context) {
	fmt.Println("[Module] Manufacturer middleware 2")

	c.Next()
}
