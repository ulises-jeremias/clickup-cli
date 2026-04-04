package doc

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/triptechtravel/clickup-cli/internal/apiv3"
	"github.com/triptechtravel/clickup-cli/internal/iostreams"
	"github.com/triptechtravel/clickup-cli/pkg/cmdutil"
)

func TestPrintDocView(t *testing.T) {
	var outBuf bytes.Buffer
	ios := iostreams.TestWithWriters(&outBuf, nil)
	f := &cmdutil.Factory{IOStreams: ios}

	d := &apiv3.DocCore{}
	d.ID = "doc1"
	d.Name = "My Doc"
	d.Visibility = "PRIVATE"
	d.DateCreated = "1700000000000"

	err := printDocView(f, d)
	require.NoError(t, err)

	out := outBuf.String()
	assert.Contains(t, out, "My Doc")
	assert.Contains(t, out, "#doc1")
	assert.Contains(t, out, "private")
	assert.Contains(t, out, "Created:")
}
