package sanitizer

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSanitizer_Sanitize(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "Should sanitize string with links",
			input: `<a onblur="alert(secret)" href="http://www.google.com">Google</a>`,
			want:  `<a href="http://www.google.com" rel="nofollow">Google</a>`,
		},
		{
			name:  "Should sanitize string with scripts",
			input: `<a href="javascript:alert('XSS1')" onmouseover="alert('XSS2')">XSS<a>`,
			want:  `XSS`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := New().Sanitize(tt.input)
			assert.Equal(t, tt.want, got)
		})
	}
}
