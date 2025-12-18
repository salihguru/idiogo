package rest

import (
	"context"
	"encoding/json"
	"fmt"
	"slices"

	"github.com/gofiber/fiber/v2"
	"github.com/salihguru/idiogo/internal/config"
	"github.com/salihguru/idiogo/internal/port"
	"github.com/salihguru/idiogo/pkg/i18np"
	"github.com/salihguru/idiogo/pkg/validation"
	"github.com/salihguru/idiogo/pkg/xascii"
	"github.com/salihguru/idiogo/pkg/xip"
)

// Module represents a complete domain module with repository, service, and router
type Module[R any, S any] struct {
	Repo    R
	Service S
	Router  Router
}

// Router is the interface that all domain handlers must implement
type Router interface {
	RegisterRoutes(srv port.RestService, router fiber.Router)
}

type Server struct {
	app *fiber.App
	cnf Config
	srv Service
}

type Config struct {
	Rest      config.Rest
	I18n      i18np.I18n
	Validator validation.Srv
	Routers   []Router
	Locales   []string
}

func New(cnf Config) *Server {
	srv := Service{
		i18n:      cnf.I18n,
		validator: cnf.Validator,
		locales:   cnf.Locales,
	}
	return &Server{
		cnf: cnf,
		srv: srv,
		app: fiber.New(fiber.Config{
			ErrorHandler:            srv.ErrorHandler(),
			DisableStartupMessage:   true,
			AppName:                 "idiogo",
			ServerHeader:            "idiogo",
			JSONEncoder:             json.Marshal,
			JSONDecoder:             json.Unmarshal,
			CaseSensitive:           true,
			BodyLimit:               100 * 1024 * 1024,
			ReadBufferSize:          100 * 1024 * 1024,
			ProxyHeader:             fiber.HeaderXForwardedFor,
			EnableTrustedProxyCheck: true,
			TrustedProxies:          slices.Concat(xip.LocalIPs, xip.CloudflareIPv4, xip.CloudflareIPv6),
		}),
	}
}

func (s *Server) Listen() error {
	s.app.Use(s.srv.Recover(), s.srv.I18n(), s.srv.IpAddr())
	for _, r := range s.cnf.Routers {
		r.RegisterRoutes(s.srv, s.app)
	}
	xascii.Log()
	fmt.Printf("idiogo api is running on %s:%s\n", s.cnf.Rest.Host, s.cnf.Rest.Port)
	return s.app.Listen(fmt.Sprintf(":%v", s.cnf.Rest.Port))
}

func (s *Server) Shutdown(ctx context.Context) error {
	if s.app != nil {
		return s.app.ShutdownWithContext(ctx)
	}
	return nil
}
