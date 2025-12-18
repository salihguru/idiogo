package ptr

import (
	"strings"
	"time"

	"github.com/google/uuid"
)

// UUID returns a pointer to the given uuid.UUID.
// Example: ptr.UUID(uuid.New())
// Returns: *uuid.UUID
func UUID(id uuid.UUID) *uuid.UUID {
	if id == uuid.Nil {
		return nil
	}
	return &id
}

// String returns a pointer to the given string.
// Example: ptr.String("hello")
// Returns: *string
func String(s string) *string {
	return &s
}

// Time returns a pointer to the given time.
// Example: ptr.Time(time.Now())
// Returns: *time.Time
func Time(t time.Time) *time.Time {
	if t.IsZero() {
		return nil
	}
	return &t
}

// Int returns a pointer to the given int.
// Example: ptr.Int(42)
// Returns: *int
func Int(i int) *int {
	return &i
}

// IntRef returns the value of the given int pointer.
// Example: ptr.IntRef(ptr.Int(42))
// Returns: int
func IntRef(i *int) int {
	if i == nil {
		return 0
	}
	return *i
}

// Int returns a pointer to the given int.
// Example: ptr.Int64(42)
// Returns: *int
func Int64(i int64) *int64 {
	return &i
}

func Float64(f float64) *float64 {
	return &f
}

// Bool returns a pointer to the given bool.
// Example: ptr.Bool(true)
// Returns: *bool
func Bool(b bool) *bool {
	return &b
}

func BoolRef(b *bool) bool {
	if b == nil {
		return false
	}
	return *b
}

func StringRef(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

func Float64Ref(f *float64) float64 {
	if f == nil {
		return 0
	}
	return *f
}

func StrToSlice(s *string) []string {
	if s == nil {
		return nil
	}
	return strings.Split(*s, ",")
}

func UUIDRef(id *uuid.UUID) uuid.UUID {
	if id == nil {
		return uuid.Nil
	}
	return *id
}

func UUIDRefStr(id *string) uuid.UUID {
	if id == nil || *id == "" {
		return uuid.Nil
	}
	u, err := uuid.Parse(*id)
	if err != nil {
		return uuid.Nil
	}
	return u
}
