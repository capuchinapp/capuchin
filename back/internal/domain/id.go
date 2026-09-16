package domain

import (
	"crypto/rand"

	"github.com/oklog/ulid/v2"
)

// NewID возвращает новый идентификатор.
func NewID() string {
	entropy := ulid.Monotonic(rand.Reader, 0)

	return ulid.MustNew(ulid.Now(), entropy).String()
}

// IsID возвращает true, если id является допустимым идентификатором.
func IsID(id string) bool {
	_, err := ulid.Parse(id)

	return err == nil
}

// PartID возвращает последнюю часть идентификатора.
func PartID(id string) string {
	return ulid.MustParse(id).String()[18:]
}
