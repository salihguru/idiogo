package locale

import (
	"database/sql/driver"
	"errors"

	"github.com/salihguru/idiogo/pkg/entity"
)

type Map map[Locale]string
type Translation[T any] map[Locale]T

type Locale string

const (
	EN Locale = "en" // English
	TR Locale = "tr" // Turkish
	ZH Locale = "zh" // Chinese
	RU Locale = "ru" // Russian
	AZ Locale = "az" // Azerbaijani
	KZ Locale = "kz" // Kazakh
	UZ Locale = "uz" // Uzbek
	KK Locale = "kk" // Kazakh
	DE Locale = "de" // German
)

func IsLocale(locale string) bool {
	for _, l := range []Locale{EN, TR} {
		if l == Locale(locale) {
			return true
		}
	}
	return false
}

func IsLocaleList(locales []string) bool {
	for _, l := range locales {
		if !IsLocale(l) {
			return false
		}
	}
	return true
}

func ParseLocale(locale string) (Locale, error) {
	if !IsLocale(locale) {
		return "", errors.New("invalid locale: " + locale)
	}
	return Locale(locale), nil
}

func (l Locale) String() string {
	return string(l)
}

func (m *Map) Scan(value interface{}) error {
	return entity.JsonbObjScan(value, m)
}

func (m Map) Value() (driver.Value, error) {
	return entity.JsonbObjValue(m)
}

func (t *Translation[T]) Scan(value interface{}) error {
	return entity.JsonbObjScan(value, t)
}

func (t Translation[T]) Value() (driver.Value, error) {
	return entity.JsonbObjValue(t)
}
