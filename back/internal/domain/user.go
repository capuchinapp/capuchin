package domain

import (
	"time"
)

// User представляет пользователя.
type User struct {
	ID           string     `json:"id"`
	Email        string     `json:"email"`
	ActivateCode *string    `json:"-"`
	CreatedAt    time.Time  `json:"createdAt"`
	UpdatedAt    *time.Time `json:"updatedAt"`
}
