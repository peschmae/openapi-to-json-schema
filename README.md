# OpenAPI component to JSON Schema extraction

Minimal CLI to extract a single schema from the `component` section of a openAPI v3 schema, 
and convert it into a valid JSON schema.

The biggest change done to the schema, is that the `type` attribute is converted to a list,
and the `nullable: true` flag, is removed and placed into the `type` list as a specific value.

## Usage

```sh
openapi-to-json-schema openapi-petstore-schema.yaml -c Order
```

You can also pipe to the `openapi-to-json-schema` command, eg.

```sh
ytt --data-values-schema-inspect=true --output=openapi-v3 -f schema.yaml  | openapi-to-json-schema
```

The CLI will work with JSON or YAML as input, but only output JSON. By default, the generated
schema will be written to stdout, unless the `--output` flag is set.
