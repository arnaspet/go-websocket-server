package domain

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReplaceBytes(t *testing.T) {
	cases := []struct{input, output []byte} {
		{[]byte("?"), []byte("!")},
		{[]byte("!"), []byte("!")},
		{[]byte("hello?"), []byte("hello!")},
		{[]byte("?a?b?c?d?"), []byte("!a!b!c!d!")},
	}

	for _, tc := range cases {
		t.Run(fmt.Sprintf("%s replaces to %s", tc.input, tc.output), func(t *testing.T) {
			assert.Equal(t, tc.output, ReplaceBytes(tc.input))
		})
	}
}
