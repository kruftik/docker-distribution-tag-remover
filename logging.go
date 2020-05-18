package main

import "fmt"

type Log struct {}

var (
	log *Log
	DebugEnabled = false
)

func (l *Log) Debug(args ...interface{}) {
	if DebugEnabled {
		msg := fmt.Sprintf("%v", args[0])
		fmt.Printf("[DEBUG] " + msg + "\n", args[1:]...)
	}
}

func (l *Log) Info(args ...interface{}) {
	msg := fmt.Sprintf("%v", args[0])
	fmt.Printf("[INFO] " + msg + "\n", args[1:]...)
}

func (l *Log) Warn(args ...interface{}) {
	msg := fmt.Sprintf("%v", args[0])
	fmt.Printf("[WARN] " +msg + "\n", args[1:]...)
}

func (l *Log) Error(args ...interface{}) {
	msg := fmt.Sprintf("%v", args[0])
	fmt.Printf("[ERROR] " + msg + "\n", args[1:]...)
}