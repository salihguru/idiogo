package rest

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
)

type Cookie = fiber.Cookie

type Response struct {
	Data       any
	Cookies    []*Cookie
	Headers    map[string]string
	StatusCode int
}

type ReqHandler[I any] = func(ctx context.Context, input I) error
type HandlerWithRes[I any, O any] = func(ctx context.Context, input I) (O, error)

func Create[I any](h ReqHandler[I]) Handler[I] {
	return func(fctx *fiber.Ctx, payload I) error {
		if err := h(fctx.UserContext(), payload); err != nil {
			return err
		}
		return fctx.SendStatus(fiber.StatusCreated)
	}
}

func CreateResponds[I any, O any](h HandlerWithRes[I, O]) Handler[I] {
	return func(fctx *fiber.Ctx, payload I) error {
		res, err := h(fctx.UserContext(), payload)
		if err != nil {
			return err
		}
		return fctx.Status(fiber.StatusCreated).JSON(res)
	}
}

func Void[I any](h ReqHandler[I]) Handler[I] {
	return func(fctx *fiber.Ctx, payload I) error {
		if err := h(fctx.UserContext(), payload); err != nil {
			return err
		}
		return fctx.SendStatus(fiber.StatusNoContent)
	}
}

func Todo[I any](h ReqHandler[I]) Handler[I] {
	return func(fctx *fiber.Ctx, payload I) error {
		if err := h(fctx.UserContext(), payload); err != nil {
			return err
		}
		return fctx.SendStatus(fiber.StatusOK)
	}
}

func Data[I any, O any](h HandlerWithRes[I, O]) Handler[I] {
	return func(fctx *fiber.Ctx, payload I) error {
		res, err := h(fctx.UserContext(), payload)
		if err != nil {
			return err
		}
		return respond(fctx, res, fiber.StatusOK)
	}
}

func respond[O any](fctx *fiber.Ctx, res O, defStatus int) error {
	if response, ok := any(res).(*Response); ok {
		for k, v := range response.Headers {
			fctx.Set(k, v)
		}
		for _, cookie := range response.Cookies {
			fctx.Cookie(cookie)
		}
		status := response.StatusCode
		if status == 0 {
			status = defStatus
		}
		if response.Data == nil {
			return fctx.SendStatus(status)
		}
		return fctx.Status(status).JSON(response.Data)
	}
	return fctx.Status(defStatus).JSON(res)
}

type CookieOpts struct {
	Value   string
	Name    string
	Domain  string
	Expires time.Time
	IsDev   bool
}

func NewCookie(opts CookieOpts) *fiber.Cookie {
	return &fiber.Cookie{
		Name:     opts.Name,
		Value:    opts.Value,
		Domain:   opts.Domain,
		Expires:  opts.Expires,
		Secure:   !opts.IsDev,
		HTTPOnly: !opts.IsDev,
		SameSite: "Strict",
	}
}

func NewCookieExpired(opts CookieOpts) *fiber.Cookie {
	return &fiber.Cookie{
		Name:     opts.Name,
		Value:    "",
		Domain:   opts.Domain,
		Expires:  time.Now().Add(-1 * time.Hour),
		Secure:   !opts.IsDev,
		HTTPOnly: !opts.IsDev,
		SameSite: "Strict",
	}
}

func NewBatchCookieExpired(names []string, opts CookieOpts) []*fiber.Cookie {
	cookies := make([]*fiber.Cookie, 0, len(names))
	for _, name := range names {
		cookies = append(cookies, &fiber.Cookie{
			Name:     name,
			Value:    "",
			Domain:   opts.Domain,
			Expires:  time.Now().Add(-1 * time.Hour),
			Secure:   !opts.IsDev,
			HTTPOnly: !opts.IsDev,
			SameSite: "Strict",
		})
	}
	return cookies
}
