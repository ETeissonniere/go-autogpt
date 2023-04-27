package helpers

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSplitCommand(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  []string
	}{
		{
			name:  "simple",
			input: "ls",
			want:  []string{"ls"},
		},
		{
			name:  "with args",
			input: "ls -l",
			want:  []string{"ls", "-l"},
		},
		{
			name:  "with more args",
			input: "ls -l -a",
			want:  []string{"ls", "-l", "-a"},
		},
		{
			name:  "with quoted args",
			input: "ls -l -a \"-h\"",
			want:  []string{"ls", "-l", "-a", "-h"},
		},
		{
			name:  "with quoted args and spaces",
			input: "ls -l -a \"t e s t\"",
			want:  []string{"ls", "-l", "-a", "t e s t"},
		},
		{
			name:  "with single quoted args",
			input: "ls -l -a '-h'",
			want:  []string{"ls", "-l", "-a", "-h"},
		},
		{
			name:  "with single quoted args and spaces",
			input: "ls -l -a 't e s t'",
			want:  []string{"ls", "-l", "-a", "t e s t"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := SplitCommand(tt.input)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestWithEscapeCharacters(t *testing.T) {
	assert.Equal(t, "test", WithEscapeCharacters("test"))
	assert.Equal(t, "test\ntest", WithEscapeCharacters("test\\ntest"))
	assert.Equal(t, "test\ntest", WithEscapeCharacters("test\ntest"))
}
