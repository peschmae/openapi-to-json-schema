/*
Copyright © 2024 Mathias Petermann <mathias.petermann@gmail.com>
*/
package openapi

import (
	"os"
	"slices"

	"encoding/json"

	"github.com/peschmae/openapi-to-json-schema/pkg/jsonschema"
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
	Title                string            `json:"title" yaml:"title"`
	Description          string            `json:"description" yaml:"description"`
	AdditionalProperties *bool             `json:"additionalProperties,omitempty" yaml:"additionalProperties,omitempty"`
	Items                *Schema           `json:"items" yaml:"items"`
	Default              interface{}       `json:"default" yaml:"default"`
	MinLength            int               `json:"minLength" yaml:"minLength"`
	MaxLength            int               `json:"maxLength" yaml:"maxLength"`
	Min                  int               `json:"minimum" yaml:"minimum"`
	Max                  int               `json:"maximum" yaml:"maximum"`
	MinItems             int               `json:"minItems,omitempty" yaml:"minItems,omitempty"`
	MaxItems             int               `json:"maxItems,omitempty" yaml:"maxItems,omitempty"`
	MinProperties        int               `json:"minProperties,omitempty" yaml:"minProperties,omitempty"`
	MaxProperties        int               `json:"maxProperties,omitempty" yaml:"maxProperties,omitempty"`
	Enum                 []interface{}     `json:"enum" yaml:"enum"`
	Nullable             bool              `json:"nullable" yaml:"nullable"`
	Properties           map[string]Schema `json:"properties" yaml:"properties"`
}

// Components is a representation of the components section of an OpenAPI document
type Components struct {
	Schemas map[string]Schema `json:"schemas" yaml:"schemas"`
}

// LoadOpenApiYamlSchema loads an OpenAPI schema from a byte array, assumes the schema is in YAML format
func LoadOpenApiYamlSchema(schema []byte) (*OpenAPI, error) {
	openApi := OpenAPI{}

	err := yaml.Unmarshal(schema, &openApi)
	if err != nil {
		return nil, err
	}

	return &openApi, nil
}

// LoadOpenApiYamlSchemaFromFile loads an OpenAPI schema from a file, assumes the file exists and is in YAML format
func LoadOpenApiYamlSchemaFromFile(file string) (*OpenAPI, error) {
	schema, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}

	return LoadOpenApiYamlSchema(schema)

}

// LoadOpenApiJsonSchema loads an OpenAPI schema from a byte array, assumes the schema is in JSON format
func LoadOpenApiJsonSchema(schema []byte) (*OpenAPI, error) {
	openApi := OpenAPI{}

	err := json.Unmarshal(schema, &openApi)
	if err != nil {
		return nil, err
	}

	return &openApi, nil
}

// LoadOpenApiJsonSchemaFromFile loads an OpenAPI schema from a file, assumes the file exists and is in YAML format
func LoadOpenApiJsonSchemaFromFile(file string) (*OpenAPI, error) {
	schema, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}

	return LoadOpenApiJsonSchema(schema)
}

func (s *Schema) UnmarshalYAML(unmarshal func(interface{}) error) error {
	// indirection to avoid infinite recursion when unmarshaling
	type rawSchema Schema
	raw := rawSchema{} // Put your defaults here

	if err := unmarshal(&raw); err != nil {
		return err
	}

	*s = Schema(raw)
	return nil
}

func (s *Schema) UnmarshalJSON(unmarshal func(interface{}) error) error {
	// indirection to avoid infinite recursion when unmarshaling
	type rawSchema Schema
	raw := rawSchema{} // Put your defaults here

	if err := unmarshal(&raw); err != nil {
		return err
	}

	*s = Schema(raw)
	return nil
}

func (s *Schema) isRequired() bool {
	return s.MinItems > 0 || s.MinLength > 0 || s.Min > 0 || s.MinProperties > 0 || (len(s.Enum) > 0 && !slices.Contains(s.Enum, s.Default))
}

func (o *OpenAPI) ConvertToJsonSchema(component string) (*jsonschema.Schema, error) {
	schema := jsonschema.Schema{
		Schema:               "https://json-schema.org/draft/2020-12/schema",
		Id:                   "https://example.biz/schema/ytt/data-values.json",
		Type:                 []string{"object"},
		AdditionalProperties: new(bool),
		Properties:           make(map[string]jsonschema.Schema),
	}
	*schema.AdditionalProperties = false
	for k, s := range o.Components.Schemas[component].Properties {
		property := convertProperty(s)
		if property.Title == "" {
			property.Title = k
		}
		schema.Properties[k] = property
	}
	return &schema, nil
}

func convertProperty(s Schema) jsonschema.Schema {
	property := jsonschema.Schema{
		Type:                 []string{s.Type},
		Title:                s.Title,
		Description:          s.Description,
		Default:              s.Default,
		AdditionalProperties: s.AdditionalProperties,
		MinLength:            s.MinLength,
		MaxLength:            s.MaxLength,
		Min:                  s.Min,
		Max:                  s.Max,
		MinItems:             s.MinItems,
		MaxItems:             s.MaxItems,
		MinProperties:        s.MinProperties,
		MaxProperties:        s.MaxProperties,
		Required:             []string{},
	}
	if s.Nullable {
		property.Type = append(property.Type, "null")
	}
	if len(s.Enum) > 0 {
		property.Enum = s.Enum
	}
	if s.Items != nil {
		i := convertProperty(*s.Items)
		property.Items = &i
	}
	if s.Properties != nil {
		property.Properties = make(map[string]jsonschema.Schema)
		for k, v := range s.Properties {
			p := convertProperty(v)
			if p.Title == "" {
				p.Title = k
			}
			property.Properties[k] = p
			if v.isRequired() || p.IsRequired() {
				property.Required = append(property.Required, k)
			}
		}
	}
	return property
}
