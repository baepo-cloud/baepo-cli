// Package iostream provides a helper for CLI tools to enable output to stdout and stderr
// in either plain text or JSON format depending on configuration.
package iostream

import (
	"fmt"
	"io"
	"os"
	"reflect"
	"strings"
)

// IOStream is the core structure for handling CLI output in either plain text or JSON format.
type IOStream struct {
	// JSONOutput determines whether the output is in JSON format (true) or plain text (false)
	JSONOutput bool
	// Stdout is the writer for standard output
	Stdout io.Writer
	// Stderr is the writer for error output
	Stderr io.Writer
}

// ErrorMessage represents an error message with optional details
type errorMessage struct {
	Error   string `json:"error"`
	Details string `json:"details,omitempty"`
}

// New creates a new IOStream with the specified JSON output flag
func New(jsonOutput bool) *IOStream {
	return &IOStream{
		JSONOutput: jsonOutput,
		Stdout:     os.Stdout,
		Stderr:     os.Stderr,
	}
}

// Message outputs a general message to stdout
func (s *IOStream) Message(str string, args ...interface{}) {
	msg := fmt.Sprintf(str, args...)
	if s.JSONOutput {
		s.writeJSON(s.Stdout, map[string]string{"message": msg})
	} else {
		fmt.Fprintln(s.Stdout, msg)
	}
}

// Error outputs an error message to stderr
func (s *IOStream) Error(str string, args ...interface{}) {
	msg := fmt.Sprintf(str, args...)
	if s.JSONOutput {
		s.writeJSON(s.Stderr, errorMessage{Error: msg})
	} else {
		fmt.Fprintln(s.Stderr, "Error: "+msg)
	}
}

// ErrorOptions represents options for customizing error output
type ErrorOptions struct {
	Error   string `json:"error"`
	Details string `json:"details,omitempty"`
	Code    string `json:"code,omitempty"`
}

// ErrorWithDetails outputs an error message with additional details to stderr
func (s *IOStream) ErrorWithDetails(opts ErrorOptions, args ...interface{}) {
	if s.JSONOutput {
		opts.Error = fmt.Sprintf(opts.Error, args...)
		s.writeJSON(s.Stderr, opts)
	} else {
		sb := strings.Builder{}
		sb.WriteString(fmt.Sprintf(opts.Error, args...))
		if opts.Details != "" {
			sb.WriteString(fmt.Sprintf(" (%s)", opts.Details))
		}
		if opts.Code != "" {
			sb.WriteString(fmt.Sprintf(" [code: %s]", opts.Code))
		}
		fmt.Fprintln(s.Stderr, sb.String())
	}
}

// ObjectOptions provides configuration options for the Object function
type ObjectOptions struct {
	Full bool
}

// Array processes and displays a slice of objects of type T based on the provided configuration
func (s *IOStream) Array(data interface{}, config []any, opts ObjectOptions) {
	if data == nil {
		return
	}

	if s.JSONOutput {
		s.writeJSON(s.Stdout, data)
		return
	}

	// Get the slice value
	sliceVal := reflect.ValueOf(data)
	if sliceVal.Kind() != reflect.Slice {
		fmt.Fprintln(s.Stderr, "Error: data is not a slice")
		return
	}

	// Extract headers from config
	headers := make([]string, 0)
	for _, cfg := range config {
		switch c := cfg.(type) {
		case FieldConfig:
			if !c.Verbose || opts.Full {
				headers = append(headers, c.DisplayName)
			}
		}
	}

	// Docker CLI-like table
	// Calculate column widths (minimum width = length of header)
	colWidths := make([]int, len(headers))
	for i, h := range headers {
		colWidths[i] = len(h)
	}

	// Build rows
	rows := make([][]string, sliceVal.Len())
	for i := 0; i < sliceVal.Len(); i++ {
		obj := sliceVal.Index(i).Interface()
		row := make([]string, 0, len(headers))

		headerIdx := 0
		for _, cfg := range config {
			switch c := cfg.(type) {
			case FieldConfig:
				if !c.Verbose || opts.Full {
					formatterVal := reflect.ValueOf(c.FormatFunc)
					args := []reflect.Value{reflect.ValueOf(obj)}
					result := formatterVal.Call(args)
					value := result[0].String()

					row = append(row, value)
					if len(value) > colWidths[headerIdx] {
						colWidths[headerIdx] = len(value)
					}
					headerIdx++
				}
			}
		}
		rows[i] = row
	}

	// Print headers
	for i, h := range headers {
		if i > 0 {
			fmt.Fprint(s.Stdout, "  ")
		}
		fmt.Fprintf(s.Stdout, "%-*s", colWidths[i], h)
	}
	fmt.Fprintln(s.Stdout)

	// Print rows
	for _, row := range rows {
		for i, cell := range row {
			if i > 0 {
				fmt.Fprint(s.Stdout, "  ")
			}
			fmt.Fprintf(s.Stdout, "%-*s", colWidths[i], cell)
		}
		fmt.Fprintln(s.Stdout)
	}
}

// Object processes and displays an object of type T based on the provided configuration
func (s *IOStream) Object(data interface{}, config []any, opts ObjectOptions) {
	if data == nil {
		return
	}

	if s.JSONOutput {
		s.writeJSON(s.Stdout, data)
		return
	}

	// For a single object, use processObject
	s.processObject(data, config, "", opts)
}

// processObject handles the formatting of a single object with a tree-style layout
func (s *IOStream) processObject(obj interface{}, config []any, prefix string, opts ObjectOptions) {
	if obj == nil {
		return
	}

	totalConfigs := len(config)
	for i, cfg := range config {
		isLast := i == totalConfigs-1
		var currentPrefix, childPrefix string

		if prefix == "" {
			// Root level elements
			if isLast {
				currentPrefix = "└─ "
			} else {
				currentPrefix = "├─ "
			}

			if isLast {
				childPrefix = "    " // Space for child items of last element
			} else {
				childPrefix = "│   " // Vertical line for child items
			}
		} else {
			// Already inside a nested structure, maintain the tree
			currentPrefix = prefix
			childPrefix = strings.Replace(prefix, "└─", "    ", 1)
			childPrefix = strings.Replace(childPrefix, "├─", "│   ", 1)
		}

		switch c := cfg.(type) {
		case FieldConfig:
			if !c.Verbose || opts.Full {
				// Call the format function
				formatterVal := reflect.ValueOf(c.FormatFunc)
				args := []reflect.Value{reflect.ValueOf(obj)}
				result := formatterVal.Call(args)
				value := result[0].String()
				if value != "" {
					fmt.Fprintf(s.Stdout, "%s%s: %s\n", currentPrefix, c.DisplayName, value)
				}
			}
		case ObjectConfig:
			if opts.Full {
				// Get the nested object using reflection
				objVal := reflect.ValueOf(obj)
				var nestedObj interface{}

				// Handle path if specified
				if c.Path != "" {
					// Find the field based on path
					if objVal.Kind() == reflect.Ptr {
						objVal = objVal.Elem()
					}

					field := objVal.FieldByName(c.Path)
					if !field.IsValid() {
						continue
					}
					nestedObj = field.Interface()
				} else {
					nestedObj = obj
				}

				if nestedObj != nil {
					// Print the object header
					if c.DisplayName != "" {
						fmt.Fprintf(s.Stdout, "%s%s:\n", currentPrefix, c.DisplayName)
					}

					// Process the nested object with the tree-style prefix
					s.processObject(nestedObj, c.Fields, childPrefix, opts)
				}
			}
		case ArrayConfig:
			if !c.Verbose || opts.Full {
				// Get the array using reflection
				objVal := reflect.ValueOf(obj)
				if objVal.Kind() == reflect.Ptr {
					objVal = objVal.Elem()
				}

				// Find the field based on path
				field := objVal.FieldByName(c.Path)
				if !field.IsValid() || field.IsNil() {
					continue
				}

				// Print the array header
				fmt.Fprintf(s.Stdout, "%s%s:\n", currentPrefix, c.DisplayName)

				// Process each item in the array
				arrayLen := field.Len()
				for j := 0; j < arrayLen; j++ {
					item := field.Index(j).Interface()

					// Determine if this is the last item in the array
					isLastItem := j == arrayLen-1

					// Create an index label
					indexLabel := fmt.Sprintf("[%d]", j)

					// Choose the appropriate prefix for array items
					var itemPrefix string
					if isLastItem {
						itemPrefix = childPrefix + "└─ " + indexLabel
					} else {
						itemPrefix = childPrefix + "├─ " + indexLabel
					}

					fmt.Fprintf(s.Stdout, "%s\n", itemPrefix)

					// Determine child prefix for array item's fields
					var itemChildPrefix string
					if isLastItem {
						itemChildPrefix = childPrefix + "    ├─ "
					} else {
						itemChildPrefix = childPrefix + "│   ├─ "
					}

					// For the last field in each object, use └─ instead of ├─
					var itemChildLastPrefix string
					if isLastItem {
						itemChildLastPrefix = childPrefix + "    └─ "
					} else {
						itemChildLastPrefix = childPrefix + "│   └─ "
					}

					if c.ObjectConfig != nil {
						// Process each field of the item
						fieldsLen := len(c.ObjectConfig.Fields)
						for k, fieldCfg := range c.ObjectConfig.Fields {
							isLastField := k == fieldsLen-1

							if fc, ok := fieldCfg.(FieldConfig); ok {
								formatterVal := reflect.ValueOf(fc.FormatFunc)
								args := []reflect.Value{reflect.ValueOf(item)}
								result := formatterVal.Call(args)
								value := result[0].String()

								if value != "" {
									fieldPrefix := itemChildPrefix
									if isLastField {
										fieldPrefix = itemChildLastPrefix
									}
									fmt.Fprintf(s.Stdout, "%s%s: %s\n", fieldPrefix, fc.DisplayName, value)
								}
							}
						}
					} else if c.FormatFunc != nil {
						// Format the item using the provided function
						formatterVal := reflect.ValueOf(c.FormatFunc)
						args := []reflect.Value{reflect.ValueOf(item)}
						result := formatterVal.Call(args)
						value := result[0].String()
						fmt.Fprintf(s.Stdout, "%s%s\n", itemChildPrefix, value)
					} else {
						// Just print the item as a string
						fmt.Fprintf(s.Stdout, "%s%v\n", itemChildPrefix, item)
					}
				}
			}
		}
	}
}
