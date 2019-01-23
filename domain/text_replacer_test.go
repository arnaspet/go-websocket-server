package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReplaceBytes(t *testing.T) {
	assert.Equal(t, []byte("!"), ReplaceBytes([]byte("?")))
	assert.Equal(t, []byte("!"), ReplaceBytes([]byte("!")))
	assert.Equal(t, []byte("hello!"), ReplaceBytes([]byte("hello?")))
	assert.Equal(t, []byte("!a!b!c!d!"), ReplaceBytes([]byte("?a?b?c?d?")))
}
