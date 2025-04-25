// Package iostream provides a helper for CLI tools to enable output to stdout and stderr
// in either plain text or JSON format depending on configuration.
package iostream

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"reflect"
	"strings"
	"text/tabwriter"

	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
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

// ObjectFormatter is a custom formatter for an entire object
type ObjectFormatter interface {
	// GetFields returns the list of field names to be included in the output
	GetFields() []string
	// FormatField handles the formatting of a specific field
	FormatField(fieldName string, value interface{}) (string, string)
}

// FieldDefinition describes how to access and format a specific field
type FieldDefinition struct {
	// DisplayName is the name to use in output (e.g., "RAM" instead of "memory_mb")
	DisplayName string
	// FormatFunc takes the entire object and returns a formatted string for this field
	FormatFunc func(obj interface{}) string
}

// ObjectOptions contains options for customizing Object output
type ObjectOptions struct {
	// Fields specifies which fields to include and how to format them
	Fields []FieldDefinition
}

// ArrayOptions contains options for customizing Array output
type ArrayOptions struct {
	// Fields specifies which fields to include and how to format them
	Fields []FieldDefinition
}

// Object outputs a single object/map with customizable field selection and formatting
func (s *IOStream) Object(data interface{}, opts ...ObjectOptions) {
	if s.JSONOutput {
		s.writeJSON(s.Stdout, data)
		return
	}

	var options ObjectOptions
	if len(opts) > 0 {
		options = opts[0]
	}

	// For plain text output, display in key-value format
	w := tabwriter.NewWriter(s.Stdout, 0, 0, 3, ' ', 0)
	defer w.Flush()

	if len(options.Fields) > 0 {
		// Use field definitions
		for _, field := range options.Fields {
			if field.FormatFunc != nil {
				displayValue := field.FormatFunc(data)
				fmt.Fprintf(w, "%s:\t%s\n", field.DisplayName, displayValue)
			}
		}
		return
	}

	// Extract the actual message from wrapper structures like response messages
	// This handles cases where data might be in a nested field like Msg.machine
	extractedData := extractNestedObject(data)

	// Default behavior - show all fields
	iterateObject(extractedData, func(key string, value interface{}) {
		var strValue string
		switch v := value.(type) {
		case string:
			strValue = v
		case fmt.Stringer:
			strValue = v.String()
		default:
			strValue = fmt.Sprintf("%v", v)
		}

		fmt.Fprintf(w, "%s:\t%s\n", key, strValue)
	})
}

// extractNestedObject attempts to extract the actual data object from wrapper structures
// like response messages that have fields like "Msg" containing the actual data
func extractNestedObject(data interface{}) interface{} {
	// If data is nil, return it as is
	if data == nil {
		return data
	}

	v := reflect.ValueOf(data)

	// Handle pointers
	if v.Kind() == reflect.Ptr {
		if v.IsNil() {
			return data
		}
		v = v.Elem()
	}

	// Only process structs
	if v.Kind() != reflect.Struct {
		return data
	}

	// Look for common wrapper fields like "Msg" or "Response"
	for _, fieldName := range []string{"Msg", "Response", "Data", "Result"} {
		field := v.FieldByName(fieldName)
		if field.IsValid() && field.CanInterface() {
			msgValue := field.Interface()

			// If this field is itself a struct or ptr to struct, check for common data fields
			msgReflect := reflect.ValueOf(msgValue)
			if msgReflect.Kind() == reflect.Ptr {
				if msgReflect.IsNil() {
					continue
				}
				msgReflect = msgReflect.Elem()
			}

			if msgReflect.Kind() == reflect.Struct {
				// Look for common data fields like "machine", "item", etc.
				for _, dataFieldName := range []string{"machine", "item", "user", "resource", "object"} {
					dataField := msgReflect.FieldByName(dataFieldName)
					if dataField.IsValid() && dataField.CanInterface() {
						return dataField.Interface()
					}
				}

				// If we didn't find a specific data field but found a wrapper field, return that
				return msgValue
			}

			// If Msg/Response is not a struct but some other value, return it
			return msgValue
		}
	}

	// If we didn't find any known wrapper fields, return the original data
	return data
}

// Array outputs an array/slice of objects with customizable field selection and formatting
func (s *IOStream) Array(data interface{}, opts ...ArrayOptions) {
	if s.JSONOutput {
		s.writeJSON(s.Stdout, data)
		return
	}

	// For plain text, create a tabulated output
	w := tabwriter.NewWriter(s.Stdout, 0, 0, 3, ' ', 0)
	defer w.Flush()

	// Process options
	var options ArrayOptions
	if len(opts) > 0 {
		options = opts[0]
	}

	// Convert to slice
	slice, ok := convertToSlice(data)
	if !ok || len(slice) == 0 {
		fmt.Fprintln(s.Stdout, "No data to display")
		return
	}

	// Determine headers and format data based on options
	var headers []string
	var rows [][]string

	if len(options.Fields) > 0 {
		// Use field definitions
		for _, field := range options.Fields {
			headers = append(headers, field.DisplayName)
		}

		for _, item := range slice {
			row := make([]string, len(headers))
			for i, field := range options.Fields {
				if field.FormatFunc != nil {
					row[i] = field.FormatFunc(item)
				} else {
					row[i] = "-"
				}
			}
			rows = append(rows, row)
		}
	} else {
		// Default behavior - extract headers from first item
		headerMap := make(map[string]bool)
		iterateObject(slice[0], func(key string, value interface{}) {
			if !headerMap[key] {
				headers = append(headers, key)
				headerMap[key] = true
			}
		})

		// Process each item
		for _, item := range slice {
			valueMap := make(map[string]string)

			iterateObject(item, func(key string, value interface{}) {
				var strValue string
				switch v := value.(type) {
				case string:
					strValue = v
				case fmt.Stringer:
					strValue = v.String()
				default:
					strValue = fmt.Sprintf("%v", v)
				}
				valueMap[key] = strValue
			})

			row := make([]string, len(headers))
			for i, header := range headers {
				row[i] = valueMap[header]
			}
			rows = append(rows, row)
		}
	}

	// Print headers
	fmt.Fprintln(w, strings.Join(headers, "\t"))

	// Print rows
	for _, row := range rows {
		fmt.Fprintln(w, strings.Join(row, "\t"))
	}
}

// These helper functions for path extraction have been removed as they're no longer needed
// with the simplified field definition approach that takes the whole object

func (s *IOStream) writeJSON(w io.Writer, data interface{}) {
	// Check if data is a proto.Message or a slice of proto.Message
	switch v := data.(type) {
	case proto.Message:
		marshaler := protojson.MarshalOptions{
			Indent:          "  ",
			EmitUnpopulated: false,
		}
		jsonBytes, err := marshaler.Marshal(v)
		if err != nil {
			fmt.Fprintf(s.Stderr, "Error encoding proto to JSON: %v\n", err)
			return
		}
		if _, err := w.Write(jsonBytes); err != nil {
			fmt.Fprintf(s.Stderr, "Error writing JSON: %v\n", err)
		}
		if _, err := w.Write([]byte("\n")); err != nil {
			fmt.Fprintf(s.Stderr, "Error writing newline: %v\n", err)
		}
	default:
		// Check if it's a slice of proto.Message
		if isSliceOfProtoMessages(data) {
			jsonBytes, err := marshalSliceOfProtoMessages(data)
			if err != nil {
				fmt.Fprintf(s.Stderr, "Error encoding proto slice to JSON: %v\n", err)
				return
			}
			if _, err := w.Write(jsonBytes); err != nil {
				fmt.Fprintf(s.Stderr, "Error writing JSON: %v\n", err)
			}
			if _, err := w.Write([]byte("\n")); err != nil {
				fmt.Fprintf(s.Stderr, "Error writing newline: %v\n", err)
			}
		} else {
			// Use standard json for regular data
			encoder := json.NewEncoder(w)
			encoder.SetIndent("", "  ")
			if err := encoder.Encode(data); err != nil {
				fmt.Fprintf(s.Stderr, "Error encoding JSON: %v\n", err)
			}
		}
	}
}

// Helper functions for handling various data types (unchanged from original)
func isSliceOfProtoMessages(data interface{}) bool {
	val := reflect.ValueOf(data)
	if val.Kind() != reflect.Slice {
		return false
	}

	// Empty slice cannot be determined
	if val.Len() == 0 {
		return false
	}

	// Check if the first element is a proto.Message
	firstElem := val.Index(0).Interface()
	_, ok := firstElem.(proto.Message)
	return ok
}

func marshalSliceOfProtoMessages(data interface{}) ([]byte, error) {
	val := reflect.ValueOf(data)
	marshaler := protojson.MarshalOptions{
		Indent:          "  ",
		EmitUnpopulated: false,
	}

	var result bytes.Buffer
	result.WriteString("[\n")

	for i := 0; i < val.Len(); i++ {
		elem := val.Index(i).Interface()
		if protoMsg, ok := elem.(proto.Message); ok {
			elemBytes, err := marshaler.Marshal(protoMsg)
			if err != nil {
				return nil, err
			}

			if i > 0 {
				result.WriteString(",\n")
			}
			result.WriteString("  ")
			result.Write(elemBytes)
		}
	}

	result.WriteString("\n]")
	return result.Bytes(), nil
}

// convertToSlice attempts to convert data to a slice of interfaces
func convertToSlice(data interface{}) ([]interface{}, bool) {
	switch v := data.(type) {
	case []interface{}:
		return v, true
	case []map[string]interface{}:
		result := make([]interface{}, len(v))
		for i, item := range v {
			result[i] = item
		}
		return result, true
	}

	// Use reflection for other slice types
	result, ok := reflectToSlice(data)
	return result, ok
}

// iterateObject iterates over fields/keys of an object and calls the provided function
func iterateObject(data interface{}, fn func(key string, value interface{})) {
	switch v := data.(type) {
	case map[string]interface{}:
		for k, val := range v {
			fn(k, val)
		}
	case map[string]string:
		for k, val := range v {
			fn(k, val)
		}
	default:
		// Use reflection for structs and other types
		reflectOverObject(data, fn)
	}
}

// reflectToSlice uses reflection to convert data to a slice
func reflectToSlice(data interface{}) ([]interface{}, bool) {
	v := reflect.ValueOf(data)

	// Check if it's a nil interface
	if !v.IsValid() {
		return nil, false
	}

	// Handle pointers
	if v.Kind() == reflect.Ptr {
		if v.IsNil() {
			return nil, false
		}
		v = v.Elem()
	}

	// Check if it's a slice or array
	if v.Kind() != reflect.Slice && v.Kind() != reflect.Array {
		return nil, false
	}

	// Convert each element to interface{}
	length := v.Len()
	result := make([]interface{}, length)

	for i := 0; i < length; i++ {
		result[i] = v.Index(i).Interface()
	}

	return result, true
}

// reflectOverObject uses reflection to iterate over an object's fields
func reflectOverObject(data interface{}, fn func(key string, value interface{})) {
	v := reflect.ValueOf(data)

	// Check if it's a nil interface
	if !v.IsValid() {
		return
	}

	// Handle pointers
	if v.Kind() == reflect.Ptr {
		if v.IsNil() {
			return
		}
		v = v.Elem()
	}

	// Handle different kinds of objects
	switch v.Kind() {
	case reflect.Struct:
		t := v.Type()
		for i := 0; i < v.NumField(); i++ {
			field := t.Field(i)
			// Skip unexported fields
			if field.PkgPath != "" {
				continue
			}

			// Use JSON tag name if available, otherwise use field name
			name := field.Name
			if tag, ok := field.Tag.Lookup("json"); ok {
				parts := strings.Split(tag, ",")
				if parts[0] != "" && parts[0] != "-" {
					name = parts[0]
				}
			}

			fn(name, v.Field(i).Interface())
		}
	case reflect.Map:
		for _, key := range v.MapKeys() {
			// Convert the key to string
			var keyStr string
			switch k := key.Interface().(type) {
			case string:
				keyStr = k
			case fmt.Stringer:
				keyStr = k.String()
			default:
				keyStr = fmt.Sprintf("%v", k)
			}

			fn(keyStr, v.MapIndex(key).Interface())
		}
	}
}
