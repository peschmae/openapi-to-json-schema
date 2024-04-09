package openapi

import (
	"os"

	"github.com/peschmae/opanpi-to-json-schema/pkg/jsonschema"
	"gopkg.in/yaml.v3"
)

// OpenAPI is a representation of an OpenAPI document
type OpenAPI struct {
	Openapi    string     `json:"openapi" yaml:"openapi"`
	Info       Info       `json:"info" yaml:"info"`
	Components Components `json:"components" yaml:"components"`
}

// Info is a representation of the info section of an OpenAPI document
type Info struct {
	Title       string `json:"title" yaml:"title"`
	Description string `json:"description" yaml:"description"`
	Version     string `json:"version" yaml:"version"`
}

// Schema is a representation of a schema in a parameter in an operation in a path in the paths section of an OpenAPI document
type Schema struct {
	Type                 string            `json:"type" yaml:"type"`
	AdditionalProperties bool              `json:"additionalProperties" yaml:"additionalProperties"`
	Items                *Schema           `json:"items" yaml:"items"`
	MinLength            int               `json:"minLength" yaml:"minLength"`
	MaxLength            int               `json:"maxLength" yaml:"maxLength"`
	Min                  int               `json:"min" yaml:"min"`
	Max                  int               `json:"max" yaml:"max"`
	Enum                 []interface{}     `json:"enum" yaml:"enum"`
	Nullable             bool              `json:"nullable" yaml:"nullable"`
	Properties           map[string]Schema `json:"properties" yaml:"properties"`
}

// Components is a representation of the components section of an OpenAPI document
type Components struct {
	Schemas map[string]Schema `json:"schemas" yaml:"schemas"`
}

// LoadOpenApiSchema loads an OpenAPI schema from a file, assumes the file exists and is in YAML format
func LoadOpenApiSchema(file string) (*OpenAPI, error) {
	schema, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}

	openApi := OpenAPI{}

	err = yaml.Unmarshal(schema, &openApi)
	if err != nil {
		return nil, err
	}

	return &openApi, nil
}

func (s *Schema) UnmarshalYAML(unmarshal func(interface{}) error) error {
	// indirection to avoid infinite recursion when unmarshaling
	type rawSchema Schema
	raw := rawSchema{
		AdditionalProperties: true,
	} // Put your defaults here
	if err := unmarshal(&raw); err != nil {
		return err
	}

	*s = Schema(raw)
	return nil
}

func (o *OpenAPI) ConvertToJsonSchema(component string) (*jsonschema.Schema, error) {
	schema := jsonschema.Schema{
		Schema:               "https://json-schema.org/draft/2020-12/schema",
		Id:                   "https://example.biz/schema/ytt/data-values.json",
		Type:                 []string{"object"},
		AdditionalProperties: false,
		Properties:           make(map[string]jsonschema.Schema),
	}
	for k, s := range o.Components.Schemas[component].Properties {
		property := convertProperty(s)
		schema.Properties[k] = property
	}
	return &schema, nil
}

func convertProperty(s Schema) jsonschema.Schema {
	property := jsonschema.Schema{
		Type:                 []string{s.Type},
		AdditionalProperties: s.AdditionalProperties,
		MinLength:            s.MinLength,
		MaxLength:            s.MaxLength,
		Min:                  s.Min,
		Max:                  s.Max,
	}
	if s.Nullable {
		property.Type = append(property.Type, "null")
	}
	if len(s.Enum) > 0 {
		property.Enum = s.Enum
	}
	if s.Properties != nil {
		property.Properties = make(map[string]jsonschema.Schema)
		for k, v := range s.Properties {
			property.Properties[k] = convertProperty(v)
		}
	}
	return property
}
