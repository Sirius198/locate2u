package main

import (
	"errors"
	"fmt"
	"log"
)

type TestLogger struct {
	// For testing purposes only - to be replaced with production logger
}

func NewTestLogger() *TestLogger {
	return &TestLogger{}
}

func (r *TestLogger) log(severity string, message string, params ...interface{}) error {
	e := fmt.Sprintf("[%s] %s", severity, fmt.Sprintf(message, params...))
	log.Println(e)
	return errors.New(e)
}

func (r *TestLogger) Error(message string, params ...interface{}) error {
	return r.log("ERROR", message, params...)
}

func (r *TestLogger) Info(message string, params ...interface{}) error {
	return r.log("INFO", message, params...)
}

func (r *TestLogger) Warning(message string, params ...interface{}) error {
	return r.log("WARNING", message, params...)
}

func (r *TestLogger) Debug(message string, params ...interface{}) error {
	return r.log("DEBUG", message, params...)
}
