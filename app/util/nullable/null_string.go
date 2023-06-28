package nullable

import (
	"capuchin/app/util"
	"database/sql"
	"encoding/json"
)

// Nullable String that overrides sql.NullString
type NullString struct {
	sql.NullString
}

func NewNullString(s string, valid bool) NullString {
	return NullString{
		sql.NullString{
			String: s,
			Valid:  valid,
		},
	}
}

func (ns NullString) MarshalJSON() ([]byte, error) {
	if ns.Valid {
		return json.Marshal(ns.String)
	}

	return json.Marshal(nil)
}

func (ns *NullString) UnmarshalJSON(data []byte) error {
	var s *string

	if err := json.Unmarshal(data, &s); err != nil {
		return util.ErrTrace("json.Unmarshal", err)
	}

	if s != nil {
		ns.Valid = true
		ns.String = *s
	} else {
		ns.Valid = false
	}

	return nil
}
