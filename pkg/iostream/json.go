package iostream

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"reflect"

	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

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
