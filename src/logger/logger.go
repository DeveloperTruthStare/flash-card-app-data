package logger

import (
	"fmt"
	"strings"

	"github.com/fatih/color"
)

type LogLevel uint8

const (
	DEBUG LogLevel = 1 << iota // Used for debugging
	INFO                       // Used for information that the user might want, or real environment debugging
	WARN                       // Used when developer does something wrong
	ERROR                      // Used when program cannot continue normal operations
)

const ALL = INFO | DEBUG | WARN | ERROR

var logLevel LogLevel = ALL

func SetLogLevel(newLevel LogLevel) {
	logLevel = newLevel
}

func Debug(format string, args ...interface{}) {
	if logLevel&DEBUG == 0 {
		return
	}
	debug := color.New(color.FgHiWhite).SprintFunc()
	println(format, debug, args...)
}

func Info(format string, args ...interface{}) {
	if logLevel&INFO == 0 {
		return
	}
	info := color.New(color.FgBlue).SprintFunc()
	println(format, info, args...)
}

func Warn(format string, args ...interface{}) {
	if logLevel&WARN == 0 {
		return
	}
	warning := color.New(color.FgYellow).SprintFunc()
	println(format, warning, args...)
}

// Same Log Level as Info, but for when something has happened
func Success(format string, args ...interface{}) {
	if logLevel&INFO == 0 {
		return
	}
	success := color.New(color.FgGreen).SprintFunc()
	println(format, success, args...)
}

func Error(format string, args ...interface{}) {
	e := color.New(color.FgRed).SprintFunc()
	println(format, e, args...)
}

func println(format string, color func(a ...interface{}) string, args ...interface{}) {
	message := fmt.Sprintf(format, args...)
	fmt.Println(color(message))
}

// Progress prints a progress bar with the given message and progress percentage.
func Progress(message string, progress uint8) {
	if progress > 100 {
		progress = 100
	}

	barLength := 50 // Length of the progress bar
	filledLength := int(barLength * int(progress) / 100)
	emptyLength := barLength - filledLength

	// Construct the progress bar
	bar := fmt.Sprintf("[%s%s]", strings.Repeat("#", filledLength), strings.Repeat("-", emptyLength))

	// Adjust the output based on the platform

	fmt.Printf("\r%s: %s %d%%", message, bar, progress)

	// Finish with a newline when progress is 100%
	if progress == 100 {
		fmt.Println()
	}
}
