package logging

import (
	"fmt"
)

func LogInfo(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	fmt.Printf("[INFO] %s\n", msg)
}

func LogError(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	fmt.Printf("[ERROR] %s\n", msg)
}

func LogWarn(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	fmt.Printf("[WARN] %s\n", msg)
}
