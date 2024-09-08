package utils

import (
	"fmt"
	"strings"
	"unicode"

)

func cleanup(errString string, err error) string {
	errString = strings.Map(func(r rune) rune {
		if unicode.IsGraphic(r) {
			return r
		}
		return -1
	}, errString)
	return fmt.Sprintf(strings.Replace(errString, "\n ", "", 1), err.Error())
}

// Debug returns a debug message to the logger 
func Debug(debugString string) {
		Elilogger.Debug(debugString)
}

// Error returns an error message to the logger and terminates
func Error(errString string, err error) error {
	errString = cleanup(errString, err)
	if err != nil {
		Elilogger.Error(errString)
	}
	return nil
}

// Fatal returns a fatal msg to the logger and terminates
func Fatal(errString string, err error) {
	errString = cleanup(errString, err)
	if err != nil {
		Elilogger.Fatal(errString)
	}
}

// Info returns an info msg to the logger 
// is info or debug
func Info(infoString string) {
		Elilogger.Info(infoString)
}

// Warn returns a warning msg to the logger
func Warn(warnString string) {
	Elilogger.Warn(warnString)
}
