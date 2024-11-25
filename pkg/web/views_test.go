package web

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSplitCommand(t *testing.T) {
	tests := map[string]struct {
		input        string
		expectedCmd  string
		expectedArgs []string
	}{
		"empty": {
			input:        "",
			expectedCmd:  "",
			expectedArgs: nil,
		},
		"simple": {
			input:        "!test",
			expectedCmd:  "test",
			expectedArgs: nil,
		},
		"complex": {
			input:        "!test1!test2!test3",
			expectedCmd:  "test1",
			expectedArgs: []string{"test2", "test3"},
		},
		"prefixed": {
			input:        "test1!test2!test3",
			expectedCmd:  "test1",
			expectedArgs: []string{"test2", "test3"},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			actualCmd, actualArgs := splitCommand(tc.input)
			assert.Equal(t, tc.expectedCmd, actualCmd)
			assert.Equal(t, tc.expectedArgs, actualArgs)
		})
	}
}
