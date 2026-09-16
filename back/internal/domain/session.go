package domain

import (
	"time"
)

// Session представляет сессию.
type Session struct {
	ID        string    `json:"id"`
	CookieID  string    `json:"-"`
	UserID    string    `json:"-"`
	CheckedAt time.Time `json:"checkedAt"`
	IsCurrent bool      `json:"isCurrent"`
}
