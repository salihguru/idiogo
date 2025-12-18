package todo

import (
	"github.com/gofiber/fiber/v2"
	"github.com/salihguru/idiogo/internal/port"
	"github.com/salihguru/idiogo/internal/rest"
)

type Handler struct {
	srv Service
}

func NewHandler(srv Service) *Handler {
	return &Handler{srv: srv}
}

func (h *Handler) RegisterRoutes(srv port.RestService, router fiber.Router) {
	group := router.Group("/todos")

	group.Post("/",
		srv.Timeout(rest.Handle(rest.WithBody(rest.WithValidation(srv.ValidateStruct(), rest.CreateResponds(h.srv.Create))))))

	group.Get("/",
		srv.Timeout(rest.Handle(rest.WithQuery(rest.Data(h.srv.Find)))))
	group.Get("/:id",
		srv.Timeout(rest.Handle(rest.WithParams(rest.WithValidation(srv.ValidateStruct(), rest.Data(h.srv.View))))))

	group.Patch("/:id",
		srv.Timeout(rest.Handle(rest.WithParams(rest.WithBody(rest.WithValidation(srv.ValidateStruct(), rest.Data(h.srv.Update)))))))

	group.Delete("/:id",
		srv.Timeout(rest.Handle(rest.WithParams(rest.WithValidation(srv.ValidateStruct(), rest.Void(h.srv.Delete))))))
}
