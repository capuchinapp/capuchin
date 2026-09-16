package sqlite

import (
	"database/sql"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestTimeToUnix(t *testing.T) {
	ny, err := time.LoadLocation("America/New_York")
	if !assert.NoError(t, err) {
		return
	}

	tests := []struct {
		name string
		in   time.Time
		want int64
	}{
		{
			name: "UTC время",
			in:   time.Date(2025, 2, 3, 9, 11, 25, 0, time.UTC),
			want: 1738573885,
		},
		{
			name: "время с субсекундной точностью отбрасывается",
			in:   time.Date(2025, 2, 3, 9, 11, 25, 999_999_999, time.UTC),
			want: 1738573885,
		},
		{
			name: "время в ненулевом таймзоне конвертируется в UTC-мгновение",
			in:   time.Date(2025, 2, 3, 4, 11, 25, 0, ny),
			want: 1738573885,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, timeToUnix(tt.in))
		})
	}
}

func TestUnixToTime(t *testing.T) {
	tests := []struct {
		name string
		in   int64
		want time.Time
	}{
		{
			name: "целые секунды",
			in:   1738573885,
			want: time.Date(2025, 2, 3, 9, 11, 25, 0, time.UTC),
		},
		{
			name: "нулевое значение",
			in:   0,
			want: time.Unix(0, 0).UTC(),
		},
		{
			name: "часовой пояс всегда UTC",
			in:   1738573885,
			want: time.Date(2025, 2, 3, 9, 11, 25, 0, time.UTC),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := unixToTime(tt.in)
			assert.Equal(t, tt.want, got)
			assert.Equal(t, time.UTC, got.Location())
		})
	}
}

func TestTimeRoundTrip(t *testing.T) {
	tests := []struct {
		name string
		in   time.Time
	}{
		{name: "UTC", in: time.Date(2025, 2, 3, 9, 11, 25, 0, time.UTC)},
		{name: "новая эпоха", in: time.Date(2026, 9, 14, 14, 38, 0, 0, time.UTC)},
		{name: "далёкое будущее", in: time.Date(2100, 12, 31, 23, 59, 59, 0, time.UTC)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := unixToTime(timeToUnix(tt.in))
			assert.Equal(t, tt.in.UTC().Truncate(time.Second), got)
		})
	}
}

func TestNullableUnixFromDomain(t *testing.T) {
	ny, err := time.LoadLocation("America/New_York")
	if !assert.NoError(t, err) {
		return
	}

	value := time.Date(2025, 2, 3, 4, 11, 25, 0, ny)

	tests := []struct {
		name string
		in   *time.Time
		want sql.NullInt64
	}{
		{
			name: "указатель на время",
			in:   &value,
			want: sql.NullInt64{Int64: 1738573885, Valid: true},
		},
		{
			name: "nil",
			in:   nil,
			want: sql.NullInt64{Valid: false},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, nullableUnixFromDomain(tt.in))
		})
	}
}

func TestNullableUnixToDomain(t *testing.T) {
	tests := []struct {
		name string
		in   sql.NullInt64
		want *time.Time
	}{
		{
			name: "валидное значение",
			in:   sql.NullInt64{Int64: 1738573885, Valid: true},
			want: func() *time.Time {
				t := time.Date(2025, 2, 3, 9, 11, 25, 0, time.UTC)
				return &t
			}(),
		},
		{
			name: "невалидное значение",
			in:   sql.NullInt64{Valid: false},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := nullableUnixToDomain(tt.in)
			assert.Equal(t, tt.want, got)
			if got != nil {
				assert.Equal(t, time.UTC, got.Location())
			}
		})
	}
}
