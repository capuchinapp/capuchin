package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewID(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "new id",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewID()
			assert.True(t, IsID(got))
		})
	}
}

func TestIsID(t *testing.T) {
	tests := []struct {
		name string
		id   string
		want bool
	}{
		{
			name: "valid id",
			id:   "01J6GP155TZP3FMQ6F86KHE34Y",
			want: true,
		},
		{
			name: "invalid id",
			id:   "invalid",
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IsID(tt.id)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestPartID(t *testing.T) {
	tests := []struct {
		name string
		id   string
		want string
	}{
		{
			name: "valid id",
			id:   "01J6GP155TZP3FMQ6F86KHE34Y",
			want: "86KHE34Y",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := PartID(tt.id)
			assert.Equal(t, tt.want, got)
		})
	}
}
