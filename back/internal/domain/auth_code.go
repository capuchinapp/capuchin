package domain

import (
	"time"
)

// AuthCode представляет код авторизации.
type AuthCode struct {
	ID        string
	UserID    string
	UserEmail string
	Code      string
	CreatedAt time.Time
}
