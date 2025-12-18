package port

import (
	"context"

	"github.com/gofiber/fiber/v2"
)

type RestService interface {
	IpAddr() fiber.Handler
	I18n() fiber.Handler
	RateLimit(limit int) fiber.Handler
	Timeout(fn fiber.Handler) fiber.Handler
	ValidateStruct() ValidatorFn
}

type ValidatorFn = func(ctx context.Context, sc interface{}) error
