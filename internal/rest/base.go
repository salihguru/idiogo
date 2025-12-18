package rest

import (
	"context"

	"github.com/salihguru/idiogo/pkg/i18np"
)

type ValidatorFn = func(ctx context.Context, sc interface{}) error

type BaseHandler interface {
	Rest() Service
	Validator() ValidatorFn
	I18n() i18np.I18n
}

type handler struct {
	rest      Service
	validator ValidatorFn
	i18n      i18np.I18n
}

func NewBaseHandler(rest Service, validator ValidatorFn, i18n i18np.I18n) BaseHandler {
	return &handler{
		rest:      rest,
		validator: validator,
		i18n:      i18n,
	}
}

func (h *handler) Rest() Service {
	return h.rest
}

func (h *handler) Validator() ValidatorFn {
	return h.validator
}

func (h *handler) I18n() i18np.I18n {
	return h.i18n
}
