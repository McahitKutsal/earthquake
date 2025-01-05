package utils

import (
	"fmt"
	"time"
)

func FormatDuration(d time.Duration) string {
	hours := int(d.Hours())
	minutes := int(d.Minutes()) % 60
	seconds := int(d.Seconds()) % 60
	milliseconds := int(d.Milliseconds()) % 1000

	// Return formatted string with milliseconds
	return fmt.Sprintf("%02d:%02d:%02d:%03d", hours, minutes, seconds, milliseconds)
}
