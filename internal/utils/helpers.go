package utils

import (
	"fmt"
	"time"
)

func FormatBytes(bytes uint64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := uint64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.2f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}

func FormatSpeed(bytesPerSec float64) string {
	const unit = 1024
	if bytesPerSec < unit {
		return fmt.Sprintf("%.2f B/s", bytesPerSec)
	}
	div, exp := float64(unit), 0
	for n := bytesPerSec / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.2f %cB/s", bytesPerSec/div, "KMGTPE"[exp])
}

func FormatDuration(d time.Duration) string {
	if d < time.Minute {
		return fmt.Sprintf("%ds", int(d.Seconds()))
	}
	if d < time.Hour {
		return fmt.Sprintf("%dm %ds", int(d.Minutes()), int(d.Seconds())%60)
	}
	return fmt.Sprintf("%dh %dm", int(d.Hours()), int(d.Minutes())%60)
}

func GetSignalQuality(signalBar int) string {
	switch signalBar {
	case 0:
		return "No Signal"
	case 1:
		return "Very Poor"
	case 2:
		return "Poor"
	case 3:
		return "Fair"
	case 4:
		return "Good"
	case 5:
		return "Excellent"
	default:
		return "Unknown"
	}
}

func GetBatteryStatus(level int) string {
	if level < 0 || level > 100 {
		return "Unknown"
	}
	if level < 10 {
		return "Critical"
	}
	if level < 20 {
		return "Low"
	}
	if level < 50 {
		return "Medium"
	}
	return "Good"
}
