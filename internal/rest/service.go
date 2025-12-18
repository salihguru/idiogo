package rest

import (
	"context"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/timeout"
	"github.com/restayway/rescode"
	"github.com/salihguru/idiogo/internal/rest/middleware"
	"github.com/salihguru/idiogo/pkg/i18np"
	"github.com/salihguru/idiogo/pkg/state"
	"github.com/salihguru/idiogo/pkg/validation"
)

type Service struct {
	i18n      i18np.I18n
	validator validation.Srv
	locales   []string
}

func NewService(i18n i18np.I18n, validator validation.Srv, locales []string) *Service {
	return &Service{
		i18n:      i18n,
		validator: validator,
		locales:   locales,
	}
}

func (s Service) ValidateStruct() func(ctx context.Context, sc interface{}) error {
	return s.validator.ValidateStruct
}

func (s Service) ErrorHandler() fiber.ErrorHandler {
	return func(c *fiber.Ctx, err error) error {
		code := fiber.StatusBadRequest
		if res, ok := err.(*rescode.RC); ok {
			res.Message = s.i18n.Translate(res.Message, state.LocaleStr(c.UserContext()))
			if res.Data != nil {
				return c.Status(res.HttpCode).JSON(res.JSON())
			}
			return c.Status(res.HttpCode).JSON(fiber.Map{
				"message": res.Message,
				"code":    res.Code,
			})
		}
		if e, ok := err.(*fiber.Error); ok {
			code = e.Code
		}
		fmt.Println(err)
		return c.Status(code).JSON(map[string]interface{}{})
	}
}

func (s Service) IpAddr() fiber.Handler {
	return middleware.IpAddr
}

func (s Service) I18n() fiber.Handler {
	return middleware.NewI18n(s.locales)
}

func (s Service) Recover() fiber.Handler {
	return recover.New(recover.Config{
		EnableStackTrace: true,
	})
}

func (h Service) RateLimit(limit int) fiber.Handler {
	return limiter.New(limiter.Config{
		Max:        limit,
		Expiration: 10 * time.Minute,
	})
}

func (h Service) Timeout(fn fiber.Handler) fiber.Handler {
	return timeout.NewWithContext(fn, 50*time.Second)
}
