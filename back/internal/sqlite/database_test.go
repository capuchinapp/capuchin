package sqlite

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBuildDSN(t *testing.T) {
	tests := []struct {
		name string
		path string
		want string
	}{
		{
			name: "простой путь к файлу",
			path: "capuchin.db",
			want: `file:capuchin.db?_pragma=busy_timeout(5000)&_pragma=journal_mode(WAL)&_pragma=synchronous(NORMAL)&_txlock=immediate`,
		},
		{
			name: "путь с подкаталогами",
			path: "data/capuchin.db",
			want: `file:data/capuchin.db?_pragma=busy_timeout(5000)&_pragma=journal_mode(WAL)&_pragma=synchronous(NORMAL)&_txlock=immediate`,
		},
		{
			name: "абсолютный путь",
			path: "/var/lib/capuchin/capuchin.db",
			want: `file:/var/lib/capuchin/capuchin.db?_pragma=busy_timeout(5000)&_pragma=journal_mode(WAL)&_pragma=synchronous(NORMAL)&_txlock=immediate`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := buildDSN(tt.path)
			assert.Equal(t, tt.want, got)
			assert.Contains(t, got, "file:"+tt.path+"?")
		})
	}
}
