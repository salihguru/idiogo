package rest

import "github.com/gofiber/fiber/v2"

type EmptyReq struct{}

type Handler[T any] func(c *fiber.Ctx, payload T) error

func Handle[T any](h Handler[T]) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var payload T
		return h(c, payload)
	}
}

func WithBody[T any](h Handler[T]) Handler[T] {
	return func(c *fiber.Ctx, payload T) error {
		if err := c.BodyParser(&payload); err != nil {
			return err
		}
		return h(c, payload)
	}
}

func WithQuery[T any](h Handler[T]) Handler[T] {
	return func(c *fiber.Ctx, payload T) error {
		if err := c.QueryParser(&payload); err != nil {
			return err
		}
		return h(c, payload)
	}
}

func WithParams[T any](h Handler[T]) Handler[T] {
	return func(c *fiber.Ctx, payload T) error {
		if err := c.ParamsParser(&payload); err != nil {
			return err
		}
		return h(c, payload)
	}
}

func WithHeaders[T any](h Handler[T]) Handler[T] {
	return func(c *fiber.Ctx, payload T) error {
		if err := c.ReqHeaderParser(&payload); err != nil {
			return err
		}
		return h(c, payload)
	}
}

func WithCookies[T any](h Handler[T]) Handler[T] {
	return func(c *fiber.Ctx, payload T) error {
		if err := c.CookieParser(&payload); err != nil {
			return err
		}
		return h(c, payload)
	}
}

func WithValidation[T any](fn ValidatorFn, h Handler[T]) Handler[T] {
	return func(c *fiber.Ctx, payload T) error {
		if err := fn(c.UserContext(), &payload); err != nil {
			return err
		}
		return h(c, payload)
	}
}
