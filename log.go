package main

import (
	"log"
	"os"
	"time"
)

type logWriter struct{}

func (_ logWriter) Write(bytes []byte) (int, error) {
	return os.Stdout.WriteString(time.Now().Format("2006-01-02 15:04:05.999999 ") + string(bytes))
}

func init() {
	log.SetFlags(0)
	log.SetOutput(logWriter{})
}
