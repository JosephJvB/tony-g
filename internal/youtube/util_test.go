package youtube

import (
	"testing"
	"time"
)

func TestUtil(t *testing.T) {
	t.Run("test getting year from date string", func(t *testing.T) {
		dateStr := "2025-04-07T19:50:25Z"

		parsed, err := time.Parse(time.RFC3339, dateStr)
		if err != nil {
			t.Errorf("failed to parse dateStr:%v", err)
		}

		if parsed.Year() != 2025 {
			t.Errorf("expected parsed year to be 2025, received:%d", parsed.Year())
		}
	})
}
