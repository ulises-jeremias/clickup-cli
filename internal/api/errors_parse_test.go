package api

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAPIErrorFromResponseBody_JSONMessage(t *testing.T) {
	body := []byte(`{"err":"OAUTH_027","message":"Team not authorized","ECODE":"OAUTH_027"}`)
	err := APIErrorFromResponseBody(403, body)
	require.Error(t, err)
	apiErr, ok := err.(*APIError)
	require.True(t, ok)
	assert.Equal(t, 403, apiErr.StatusCode)
	assert.Contains(t, apiErr.Error(), "Team not authorized")
}
