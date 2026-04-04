package doc

import (
	"bytes"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/triptechtravel/clickup-cli/internal/apiv3"
	"github.com/triptechtravel/clickup-cli/internal/iostreams"
)

func TestPrintPageTree(t *testing.T) {
	var buf bytes.Buffer
	ios := iostreams.TestWithWriters(&buf, nil)
	cs := ios.ColorScheme()

	pages := []apiv3.PageRef{
		{ID: "1", Name: "Alpha", Pages: []apiv3.PageRef{
			{ID: "2", Name: "Beta"},
		}},
	}
	printPageTree(&buf, pages, 0, cs)

	out := buf.String()
	require.Contains(t, out, "Alpha")
	require.Contains(t, out, "Beta")
	lines := strings.Split(strings.TrimSpace(out), "\n")
	assert.GreaterOrEqual(t, len(lines), 2)
	assert.Contains(t, lines[0], "Alpha")
	assert.True(t, strings.HasPrefix(lines[1], "  "), "nested page should be indented")
}
