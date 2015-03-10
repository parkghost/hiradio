package main

import (
	"fmt"
	"log"
	"os"
)

var logger = log.New(os.Stderr, "", 0)

func print(levelText string, msg string) {
	logger.Printf("[%s]%s\n", levelText, msg)
}

func Fatalf(format string, args ...interface{}) {
	print("Error", fmt.Sprintf(format, args...))
	os.Exit(1)
}

func Fatal(args ...interface{}) {
	print("Error", fmt.Sprint(args...))
	os.Exit(1)
}

func Warnf(format string, args ...interface{}) {
	print("Warning", fmt.Sprintf(format, args...))
}

func Warn(args ...interface{}) {
	print("Warning", fmt.Sprint(args...))
}
