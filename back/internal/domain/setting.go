package domain

// Setting представляет настройку.
type Setting struct {
	UserID string `json:"-"`
	Key    string `json:"key"`
	Value  string `json:"value"`
}
