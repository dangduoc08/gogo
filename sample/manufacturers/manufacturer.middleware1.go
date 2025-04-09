package manufacturers

import (
	"fmt"

	"github.com/dangduoc08/gogo/ctx"
)

func ManufacturerMiddleware1(c *ctx.Context) {
	fmt.Println("[Module] Manufacturer middleware 1")

	c.Next()
}
