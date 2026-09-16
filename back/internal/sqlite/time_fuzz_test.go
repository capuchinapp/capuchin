package sqlite

import (
	"testing"
)

func FuzzTimeRoundTrip(f *testing.F) {
	seeds := []int64{
		0,
		1,
		-1,
		1738573885,
		253402300799, // 9999-12-31T23:59:59Z
	}
	for _, s := range seeds {
		f.Add(s)
	}

	f.Fuzz(func(t *testing.T, ts int64) {
		got := timeToUnix(unixToTime(ts))
		if got != ts {
			t.Fatalf("round trip imbalance: in=%d out=%d", ts, got)
		}
	})
}
