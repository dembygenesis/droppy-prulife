package validators

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
)

type Validator struct {
	c *fiber.Ctx
}

func (v *Validator) SetContext(c *fiber.Ctx) {
	v.c = c
}

// Parameters passed here must not be an empty string
func (v *Validator) HasParams([]string)  {

}