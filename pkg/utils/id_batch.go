package utils

import (
	"fmt"
	"time"
)

// GenerateBatchID creates a readable unique identifier for background tasks.
// Format: #PREFIX_DD_MM_YY_HH_MM_SS (e.g., #IMA_GEN_02_03_26_1453)
func GenerateBatchID(prefix string) string {
	now := time.Now()
	timestamp := now.Format("02_01_06_1504")
	// Add seconds if needed for higher resolution, but for now simple is better
	return fmt.Sprintf("#%s_%s", prefix, timestamp)
}
