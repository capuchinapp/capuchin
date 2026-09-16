package sanitizer

import "github.com/microcosm-cc/bluemonday"

// Sanitizer представляет санитайзер для фильтрации входящих данных.
type Sanitizer struct {
	policy *bluemonday.Policy
}

// New возвращает новый экземпляр Sanitizer.
func New() *Sanitizer {
	return &Sanitizer{
		policy: bluemonday.UGCPolicy(),
	}
}

// Sanitize фильтрует входящие данные.
func (s *Sanitizer) Sanitize(input string) string {
	return s.policy.Sanitize(input)
}
