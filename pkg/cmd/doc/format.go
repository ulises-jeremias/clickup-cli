package doc

import (
	"strconv"
	"time"

	"github.com/triptechtravel/clickup-cli/internal/text"
)

// formatDocTimestamp formats ClickUp v3 unix-ms string timestamps for CLI output.
func formatDocTimestamp(ms string) string {
	if ms == "" {
		return ""
	}
	n, err := strconv.ParseInt(ms, 10, 64)
	if err != nil {
		return ms
	}
	t := time.UnixMilli(n)
	return t.Format("2006-01-02 15:04") + " (" + text.RelativeTime(t) + ")"
}
