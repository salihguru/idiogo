package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/salihguru/idiogo/pkg/state"
	"github.com/salihguru/idiogo/pkg/xip"
)

func IpAddr(c *fiber.Ctx) error {
	ip := xip.ClaimRealIP(func(name string) string {
		return c.Get(name)
	}, c.IP())
	c.SetUserContext(state.SetIP(c.UserContext(), ip))
	return c.Next()
}
