package utils

import (
	"log"
)

// LogInfo prints informational logs
func LogInfo(message string) {
	log.Println("[INFO]:", message)
}

// LogError prints error logs
func LogError(message string) {
	log.Println("[ERROR]:", message)
}

// LogDebug prints debug logs
func LogDebug(message string) {
	log.Println("[DEBUG]:", message)
}

// LogWarning prints warning logs
func LogWarning(message string) {
	log.Println("[WARNING]:", message)
}

// LogFatal prints fatal logs and stops the program
func LogFatal(message string) {
	log.Fatal("[FATAL]:", message)
}

// LogInfof prints formatted info logs
func LogInfof(format string, args ...interface{}) {
	log.Printf("[INFO]: "+format, args...)
}

// LogErrorf prints formatted error logs
func LogErrorf(format string, args ...interface{}) {
	log.Printf("[ERROR]: "+format, args...)
}

// LogDebugf prints formatted debug logs
func LogDebugf(format string, args ...interface{}) {
	log.Printf("[DEBUG]: "+format, args...)
}

// LogWarningf prints formatted warning logs
func LogWarningf(format string, args ...interface{}) {
	log.Printf("[WARNING]: "+format, args...)
}

// LogErrorWithErr logs a message along with an error object
func LogErrorWithErr(message string, err error) {
	if err != nil {
		log.Println("[ERROR]:", message, "|", err.Error())
		return
	}
	log.Println("[ERROR]:", message)
}