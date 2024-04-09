package jsonschema

// Schema is a representation of a JSON schema
type Schema struct {
	Schema               string            `json:"$schema,omitempty"`
	Id                   string            `json:"$id,omitempty"`
	Type                 []string          `json:"type"`
	Description          string            `json:"description,omitempty"`
	Properties           map[string]Schema `json:"properties,omitempty"`
	AdditionalProperties bool              `json:"additionalProperties"`
	MinLength            int               `json:"minLength,omitempty"`
	MaxLength            int               `json:"maxLength,omitempty"`
	Min                  int               `json:"min,omitempty"`
	Max                  int               `json:"max,omitempty"`
	Enum                 []interface{}     `json:"enum,omitempty"`
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

func (s *Schema) UnmarshalJSON(unmarshal func(interface{}) error) error {
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
