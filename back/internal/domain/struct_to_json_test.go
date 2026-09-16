package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStructToJSON(t *testing.T) {
	type args struct {
		s any
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Should convert struct to json",
			args: args{
				s: struct {
					Name string `json:"name"`
					Age  int    `json:"age"`
				}{
					Name: "John",
					Age:  30,
				},
			},
			want: `{"name":"John","age":30}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := StructToJSON(tt.args.s)
			assert.Equal(t, tt.want, got)
		})
	}
}
