package nullable

import (
	"capuchin/app/util"
	"database/sql"
	"encoding/json"
	"time"
)

// Nullable Time that overrides sql.NullTime
type NullTime struct {
	sql.NullTime
}

func NewNullTime(t time.Time, valid bool) NullTime {
	return NullTime{
		sql.NullTime{
			Time:  t,
			Valid: valid,
		},
	}
}

func (nt NullTime) MarshalJSON() ([]byte, error) {
	if nt.Valid {
		return json.Marshal(nt.Time)
	}

	return json.Marshal(nil)
}

func (nt *NullTime) UnmarshalJSON(data []byte) error {
	var t *time.Time

	if err := json.Unmarshal(data, &t); err != nil {
		return util.ErrTrace("json.Unmarshal", err)
	}

	if t != nil {
		nt.Valid = true
		nt.Time = *t
	} else {
		nt.Valid = false
	}

	return nil
}
