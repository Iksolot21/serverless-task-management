package logger

import (
	"fmt"
	"log"
	"runtime"
)

func Error(message string, err error) {
	_, file, line, ok := runtime.Caller(1)
	if ok {
		log.Printf("[ERROR] %s:%d %s: %v \n", file, line, message, err)
	} else {
		log.Printf("[ERROR] %s: %v \n", message, err)
	}

}
func Info(message string, params ...interface{}) {
	_, file, line, ok := runtime.Caller(1)
	if ok {
		log.Printf("[INFO] %s:%d %s \n", file, line, fmt.Sprintf(message, params...))
	} else {
		log.Printf("[INFO] %s \n", fmt.Sprintf(message, params...))
	}
}
func Debug(message string, params ...interface{}) {
	_, file, line, ok := runtime.Caller(1)
	if ok {
		log.Printf("[DEBUG] %s:%d %s \n", file, line, fmt.Sprintf(message, params...))
	} else {
		log.Printf("[DEBUG] %s \n", fmt.Sprintf(message, params...))
	}

}
func Trace(message string, params ...interface{}) {
	_, file, line, ok := runtime.Caller(1)
	if ok {
		log.Printf("[TRACE] %s:%d %s \n", file, line, fmt.Sprintf(message, params...))
	} else {
		log.Printf("[TRACE] %s \n", fmt.Sprintf(message, params...))
	}

}
func Warn(message string, params ...interface{}) {
	_, file, line, ok := runtime.Caller(1)
	if ok {
		log.Printf("[WARN] %s:%d %s \n", file, line, fmt.Sprintf(message, params...))
	} else {
		log.Printf("[WARN] %s \n", fmt.Sprintf(message, params...))
	}
}
