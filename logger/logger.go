package logger

import (
	"os"

	"github.com/charlesbases/colors"
)

// Fatal .
func Fatal(v ...interface{}) {
	os.Stderr.WriteString(colors.RedSprint(v...))
	os.Stderr.WriteString("\n")
	os.Exit(1)
}

// Fatalf .
func Fatalf(format string, v ...interface{}) {
	os.Stderr.WriteString(colors.RedSprintf(format, v...))
	os.Stderr.WriteString("\n")
	os.Exit(1)
}
