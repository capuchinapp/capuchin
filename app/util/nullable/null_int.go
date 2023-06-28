package nullable

import (
	"capuchin/app/util"
	"database/sql"
	"encoding/json"
)

// Nullable Int64 that overrides sql.NullInt64
type NullInt64 struct {
	sql.NullInt64
}

func NewNullInt64(i int64, valid bool) NullInt64 {
	return NullInt64{
		sql.NullInt64{
			Int64: i,
			Valid: valid,
		},
	}
}

func (ns NullInt64) MarshalJSON() ([]byte, error) {
	if ns.Valid {
		return json.Marshal(ns.Int64)
	}

	return json.Marshal(nil)
}

func (ns *NullInt64) UnmarshalJSON(data []byte) error {
	var i *int64

	if err := json.Unmarshal(data, &i); err != nil {
		return util.ErrTrace("json.Unmarshal", err)
	}

	if i != nil {
		ns.Valid = true
		ns.Int64 = *i
	} else {
		ns.Valid = false
	}

	return nil
}
