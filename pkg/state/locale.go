package state

import (
	"context"

	"github.com/salihguru/idiogo/pkg/locale"
)

// SetLocale sets the locale in the context
func SetLocale(ctx context.Context, locale string) context.Context {
	return context.WithValue(ctx, KeyLocale, locale)
}

// GetLocale gets the locale from the context
func Locale(ctx context.Context) locale.Locale {
	if l, ok := ctx.Value(KeyLocale).(string); ok {
		loc, err := locale.ParseLocale(l)
		if err != nil {
			return locale.EN
		}
		return loc
	}
	return locale.EN
}

func LocaleStr(ctx context.Context) string {
	return string(Locale(ctx))
}
