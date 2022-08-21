package internal

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseVersion(t *testing.T) {
	tests := []struct {
		inp   string
		ver   int
		isErr bool
	}{
		{"v1", 1, false},
		{"1", 1, false},
		{"v0", -1, true},
		{"0", -1, true},
		{"v-1", -1, true},
		{"-1", -1, true},
		{"v128", 128, false},
		{"128", 128, false},
	}
	for _, test := range tests {
		i, err := ParseVersion(test.inp)
		if test.isErr {
			assert.Error(t, err)
		} else {
			assert.Equal(t, i, test.ver)
			assert.NoError(t, err)
		}
	}
}
