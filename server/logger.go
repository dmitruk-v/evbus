package server

import (
	"fmt"
	"os"
	"time"
)

func MustGetLogFile() *os.File {
	fname := time.Now().Format("02-01-2006")
	file, err := os.OpenFile(fname, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		panic(fmt.Sprintf("Could not open log file: %v", err))
	}
	return file
}
