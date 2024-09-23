/*
Copyright © 2024 Mathias Petermann <mathias.petermann@gmail.com>
*/
package jsonschema

// Schema is a representation of a JSON schema
type Schema struct {
	Schema               string            `json:"$schema,omitempty"`
	Title                string            `json:"title,omitempty"`
	Id                   string            `json:"$id,omitempty"`
	Type                 []string          `json:"type"`
	Items                *Schema           `json:"items,omitempty"`
	Default              interface{}       `json:"default,omitempty"`
	Description          string            `json:"description,omitempty"`
	Properties           map[string]Schema `json:"properties,omitempty"`
	AdditionalProperties *bool             `json:"additionalProperties,omitempty"`
	MinLength            int               `json:"minLength,omitempty"`
	MaxLength            int               `json:"maxLength,omitempty"`
	Min                  int               `json:"minimum,omitempty"`
	Max                  int               `json:"maximum,omitempty"`
	MinItems             int               `json:"minItems,omitempty"`
	MaxItems             int               `json:"maxItems,omitempty"`
	MinProperties        int               `json:"minProperties,omitempty"`
	MaxProperties        int               `json:"maxProperties,omitempty"`
	Enum                 []interface{}     `json:"enum,omitempty"`
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
