package model

type PropertyField interface{}

const (
	StringType  = "string"
	IntegerType = "integer"
	NumberType  = "number"
	BooleanType = "boolean"
	ObjectType  = "object"
	ArrayType   = "array"

	RootType = ObjectType
)

type RootParameter struct {
	Type       string                   `json:"type"`
	Properties map[string]PropertyField `json:"properties"`
	Required   []string                 `json:"required"`
}

type BaseProperty struct {
	Type        string `json:"type"`
	Description string `json:"description"`
}

type ObjectProperty struct {
	Type        string                   `json:"type"`
	Description string                   `json:"description"`
	Properties  map[string]PropertyField `json:"properties"`
	Required    []string                 `json:"required"`
}

type ArrayProperty struct {
	Type        string          `json:"type"`
	Description string          `json:"description"`
	Items       []PropertyField `json:"items"`
}
