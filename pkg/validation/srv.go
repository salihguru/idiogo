package validation

import (
	"context"
	"regexp"
	"strings"

	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/tr"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"github.com/salihguru/idiogo/pkg/i18np"
	"github.com/salihguru/idiogo/pkg/state"
	"github.com/salihguru/idiogo/pkg/xrescode"
)

type Srv struct {
	validator *validator.Validate
	uni       *ut.UniversalTranslator
	i18n      *i18np.I18n
}

func New(i18n *i18np.I18n) *Srv {
	v := validator.New()
	
	// Register custom validators
	v.RegisterValidation("username", validateUsername)
	v.RegisterValidation("password", validatePassword)
	v.RegisterValidation("locale", validateLocale)
	v.RegisterValidation("slug", validateSlug)
	v.RegisterValidation("gender", validateGender)
	v.RegisterValidation("phone", validatePhone)
	
	return &Srv{validator: v, uni: ut.New(tr.New(), en.New()), i18n: i18n}
}

// Custom validation functions
func validateUsername(fl validator.FieldLevel) bool {
	username := fl.Field().String()
	if len(username) < 3 || len(username) > 20 {
		return false
	}
	matched, _ := regexp.MatchString("^[a-zA-Z0-9_]+$", username)
	return matched
}

func validatePassword(fl validator.FieldLevel) bool {
	password := fl.Field().String()
	if len(password) < 8 {
		return false
	}
	// Check for at least one uppercase, one lowercase, one digit, and one special char
	hasUpper := regexp.MustCompile(`[A-Z]`).MatchString(password)
	hasLower := regexp.MustCompile(`[a-z]`).MatchString(password)
	hasDigit := regexp.MustCompile(`[0-9]`).MatchString(password)
	hasSpecial := regexp.MustCompile(`[!@#$%^&*(),.?":{}|<>]`).MatchString(password)
	return hasUpper && hasLower && hasDigit && hasSpecial
}

func validateLocale(fl validator.FieldLevel) bool {
	locale := fl.Field().String()
	matched, _ := regexp.MatchString(localeRegexp, locale)
	return matched
}

func validateSlug(fl validator.FieldLevel) bool {
	slug := fl.Field().String()
	matched, _ := regexp.MatchString(slugRegexp, slug)
	return matched
}

func validateGender(fl validator.FieldLevel) bool {
	gender := fl.Field().String()
	matched, _ := regexp.MatchString(genderRegexp, gender)
	return matched
}

func validatePhone(fl validator.FieldLevel) bool {
	phone := fl.Field().String()
	matched, _ := regexp.MatchString(phoneWithCountryCodeRegexp, phone)
	return matched
}

func (s *Srv) translate(ctx context.Context, err validator.FieldError) string {
	if s.i18n == nil {
		return err.Translate(s.getTranslator(ctx))
	}
	msg := s.i18n.TranslateWithParams("validation_"+err.Tag(), map[string]interface{}{
		"Value": err.Value(),
		"Field": err.Field(),
		"Param": err.Param(),
	}, state.LocaleStr(ctx))
	if msg == "" || msg == "validation_"+err.Tag() {
		msg = err.Translate(s.getTranslator(ctx))
	}
	return msg
}

// ValidateStruct validates the given struct.
func (s *Srv) ValidateStruct(ctx context.Context, sc interface{}) error {
	var errs []*ErrorResponse
	err := s.validator.StructCtx(ctx, sc)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			var element ErrorResponse
			ns := s.mapStructNamespace(err.Namespace())
			if ns != "" {
				element.Namespace = ns
			}
			element.Field = err.Field()
			element.Value = err.Value()
			element.Message = s.translate(ctx, err)
			errs = append(errs, &element)
		}
	}
	if len(errs) > 0 {
		return xrescode.ValidationFailed().SetData(errs)
	}
	return nil
}

// ValidateMap validates the giveb struct.
func (s *Srv) ValidateMap(ctx context.Context, m map[string]interface{}, rules map[string]interface{}) error {
	var errs []*ErrorResponse
	errMap := s.validator.ValidateMapCtx(ctx, m, rules)
	for key, err := range errMap {
		var element ErrorResponse
		if _err, ok := err.(validator.ValidationErrors); ok {
			for _, err := range _err {
				element.Namespace = err.Namespace()
				element.Field = err.Field()
				if element.Field == "" {
					element.Field = key
				}
				element.Value = err.Value()
				element.Message = s.translate(ctx, err)
				errs = append(errs, &element)
			}
			continue
		}
	}
	if len(errs) > 0 {
		return xrescode.ValidationFailed().SetData(errs)
	}
	return nil
}

func (s *Srv) getTranslator(ctx context.Context) ut.Translator {
	locale := state.Locale(ctx)
	translator, found := s.uni.GetTranslator(string(locale))
	if !found {
		translator = s.uni.GetFallback()
	}
	return translator
}

func (s *Srv) mapStructNamespace(ns string) string {
	str := strings.Split(ns, ".")
	return strings.Join(str[1:], ".")
}
