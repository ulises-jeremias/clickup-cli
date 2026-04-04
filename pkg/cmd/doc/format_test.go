package doc

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFormatDocTimestamp(t *testing.T) {
	assert.Equal(t, "", formatDocTimestamp(""))

	s := formatDocTimestamp("not-a-number")
	assert.Equal(t, "not-a-number", s)

	// 2006-01-02 15:04:05 UTC in ms
	got := formatDocTimestamp("1136214245000")
	assert.Contains(t, got, "2006-01-02")
	assert.Contains(t, got, "(")
	assert.Contains(t, got, ")")
}
