package helpers

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMatch(t *testing.T) {
	tests := []struct {
		name     string
		data 	string
		regex 	string
		expected [][]string
	}{
		{
			name: fmt.Sprintf("%s: HappyDay", GetTestName()),
			data: "i HAVE alot of DATA that i need to CUT out using REGEX",
			regex: "i (.+) alot of (.+) that i need to (.+) out using (.+)",
			expected: [][]string{
				{
					"i HAVE alot of DATA that i need to CUT out using REGEX",
					"HAVE",
					"DATA",
					"CUT",
					"REGEX",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual, err := Match(tt.data, tt.regex)
			if err != nil {
			    assert.Equal(t, "", err.Error())
			    return
			}

			assert.Equal(t, len(tt.expected), len(actual))
			if len(tt.expected) > 0 {
				for i, a := range tt.expected {
					assert.Equal(t, len(a), len(actual[i]))
				}
			}

			for i, a := range tt.expected {
				for i2, b := range a {
					assert.Equal(t, b, actual[i][i2])
				}
			}
		})
	}
}
