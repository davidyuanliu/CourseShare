package utils

import "log"

// LogInfo prints informational logs
func LogInfo(message string) {
	log.Println("[INFO]:", message)
}

// LogError prints error logs
func LogError(message string) {
	log.Println("[ERROR]:", message)
}