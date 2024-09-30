package parser

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseSize(t *testing.T) {
	tests := map[string]struct {
		sizeStr     string
		expectedRes int
		expectError bool
	}{
		"parse size in B": {
			sizeStr:     "4B",
			expectedRes: 4,
			expectError: false,
		},
		"parse size in KB": {
			sizeStr:     "4KB",
			expectedRes: 4096,
			expectError: false,
		},
		"parse size in MB": {
			sizeStr:     "4MB",
			expectedRes: 4194304,
			expectError: false,
		},
		"parse size in GB": {
			sizeStr:     "4GB",
			expectedRes: 4294967296,
			expectError: false,
		},
		"parse invalid size suffix": {
			sizeStr:     "4G",
			expectedRes: 0,
			expectError: true,
		},
		"parse null size": {
			sizeStr:     "0KB",
			expectedRes: 0,
			expectError: true,
		},
		"parse negative size": {
			sizeStr:     "-4KB",
			expectedRes: 0,
			expectError: true,
		},
		"parse empty string": {
			sizeStr:     "",
			expectedRes: 0,
			expectError: true,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			result, err := ParseSize(test.sizeStr)
			assert.Equal(t, test.expectedRes, result)
			if test.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
