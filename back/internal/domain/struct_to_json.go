package domain

import "encoding/json"

// StructToJSON конвертирует структуру в JSON.
func StructToJSON(s any) string {
	b, err := json.Marshal(s)
	if err != nil {
		return ""
	}

	return string(b)
}
