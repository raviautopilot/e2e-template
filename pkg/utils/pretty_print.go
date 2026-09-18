package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
)

// PrettyJSON serializes any Go value (struct, map, slice, primitive) into an indented JSON string.
func PrettyJSON(v interface{}) (string, error) {
	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to pretty format JSON: %w", err)
	}
	return string(b), nil
}

// MustPrettyJSON serializes a value to indented JSON, returning an error message string if marshaling fails.
func MustPrettyJSON(v interface{}) string {
	res, err := PrettyJSON(v)
	if err != nil {
		return fmt.Sprintf("<!error formatting JSON: %v>", err)
	}
	return res
}

// PrettyJSONString parses a raw JSON string and re-formats it with clean 2-space indentation.
func PrettyJSONString(raw string) (string, error) {
	var out bytes.Buffer
	if err := json.Indent(&out, []byte(raw), "", "  "); err != nil {
		return "", fmt.Errorf("failed to indent raw JSON string: %w", err)
	}
	return out.String(), nil
}

// PrintJSON prints the indented JSON representation of v to os.Stdout.
func PrintJSON(v interface{}) {
	FprintJSON(os.Stdout, v)
}

// FprintJSON writes the indented JSON representation of v to the given io.Writer.
func FprintJSON(w io.Writer, v interface{}) error {
	formatted, err := PrettyJSON(v)
	if err != nil {
		fmt.Fprintf(w, "error: %v\n", err)
		return err
	}
	fmt.Fprintln(w, formatted)
	return nil
}
