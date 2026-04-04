package api

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestClient_APIV3BaseURL(t *testing.T) {
	u, err := url.Parse("https://api.clickup.com/api/v2/")
	require.NoError(t, err)
	c := NewClient("token")
	c.Clickup.BaseURL = u

	got := c.APIV3BaseURL()
	assert.Equal(t, "https://api.clickup.com/api/v3/", got)
}
