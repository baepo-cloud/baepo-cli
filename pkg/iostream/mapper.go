package iostream

// can be only func(T) string
type FormatterFunc any

type FieldConfig struct {
	DisplayName string
	Verbose     bool
	FormatFunc  FormatterFunc
}

type ObjectConfig struct {
	Path        string
	DisplayName string
	Full        bool
	Fields      []any // in fact it is a map where value is either FieldConfig, ObjectConfig, ArrayConfig
}

type ArrayConfig struct {
	Path         string
	DisplayName  string
	Verbose      bool
	FormatFunc   FormatterFunc // used if ObjectConfig is nil
	ObjectConfig *ObjectConfig // used if Not nil
}
