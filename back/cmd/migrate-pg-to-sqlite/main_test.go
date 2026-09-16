//go:build migration

package main

import (
	"database/sql"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestTimeToUnix(t *testing.T) {
	tests := []struct {
		name     string
		input    time.Time
		expected int64
	}{
		{
			name:     "unix epoch",
			input:    time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC),
			expected: 0,
		},
		{
			name:     "known date",
			input:    time.Date(2025, 6, 15, 12, 30, 45, 0, time.UTC),
			expected: 1749990645,
		},
		{
			name:     "non-utc input is converted to utc",
			input:    time.Date(2025, 6, 15, 15, 30, 45, 0, time.FixedZone("MSK", 3*60*60)),
			expected: 1749990645,
		},
		{
			name:     "zero time",
			input:    time.Time{},
			expected: -62135596800,
		},
		{
			name:     "sub-second precision truncated",
			input:    time.Date(2025, 1, 1, 0, 0, 0, 999999999, time.UTC),
			expected: 1735689600,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := timeToUnix(tt.input)
			assert.Equal(t, tt.expected, got)
		})
	}
}

func TestNullableTimeToUnix(t *testing.T) {
	tests := []struct {
		name     string
		input    sql.NullTime
		expected sql.NullInt64
	}{
		{
			name:     "valid time",
			input:    sql.NullTime{Time: time.Date(2025, 6, 15, 12, 0, 0, 0, time.UTC), Valid: true},
			expected: sql.NullInt64{Int64: 1749988800, Valid: true},
		},
		{
			name:     "null time",
			input:    sql.NullTime{Valid: false},
			expected: sql.NullInt64{Valid: false},
		},
		{
			name:     "valid zero time",
			input:    sql.NullTime{Time: time.Time{}, Valid: true},
			expected: sql.NullInt64{Int64: -62135596800, Valid: true},
		},
		{
			name:     "valid time with non-utc timezone",
			input:    sql.NullTime{Time: time.Date(2025, 1, 1, 3, 0, 0, 0, time.FixedZone("UTC+3", 3*3600)), Valid: true},
			expected: sql.NullInt64{Int64: 1735689600, Valid: true},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := nullableTimeToUnix(tt.input)
			assert.Equal(t, tt.expected, got)
		})
	}
}

func TestTimeToUnixRoundTrip(t *testing.T) {
	original := time.Date(2025, 6, 15, 12, 30, 45, 0, time.UTC)
	unix := timeToUnix(original)
	restored := time.Unix(unix, 0).UTC()

	assert.Equal(t, original, restored)
}
