package query

import (
	"github.com/google/uuid"
)

func SkipStrPtr(s *string) bool {
	return s == nil || *s == ""
}

func SkipUUIDPtr(u *uuid.UUID) bool {
	return u == nil || *u == uuid.Nil
}
