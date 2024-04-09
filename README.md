# OpenAPI component to JSON Schema extraction

Minimal CLI to extract a single schema from the `component` section of a OpenAPI schema, 
and convert it into a valid JSON schema.

The biggest change done to the schema, is that the `type` attribute is converted to a list,
and the `nullable: true` flag, is removed and placed into the `type` list as a specific value.
