package domain

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestReplaceBytes(t *testing.T) {
	assert.Equal(t, []byte("!"), ReplaceBytes([]byte("?")))
	assert.Equal(t, []byte("!"), ReplaceBytes([]byte("!")))
	assert.Equal(t, []byte("hello!"), ReplaceBytes([]byte("hello?")))
	assert.Equal(t, []byte("!a!b!c!d!"), ReplaceBytes([]byte("?a?b?c?d?")))
}
