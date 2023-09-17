package main

import "os"

func getEnv(key string, defaultVal string) string {
	val, exists := os.LookupEnv(key)
	if exists {
		return val
	}
	return defaultVal
}
