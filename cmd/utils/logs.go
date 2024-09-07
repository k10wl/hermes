package utils

import (
	"fmt"
	"io"
)

func LogError(w io.Writer, error error) {
	fmt.Fprintf(w, withGroupPrefix("Error", error))
}

func LogFail(w io.Writer, reason string) {
	fmt.Fprintf(w, withGroupPrefix("Fail", reason))
}

func withGroupPrefix(prefix string, content any) string {
	return fmt.Sprintf("%s: %v\n", prefix, content)
}
