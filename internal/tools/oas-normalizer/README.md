# oas-normalizer

A quick tool to clean up OpenAPI specs that don't play well with the terraform-plugin-codegen-openapi code generator.

## What it patches

Two main issues that trip up the the terraform-plugin-codegen-openapi generator:

1. **allOf/oneOf schemas** - Flattens these into regular object schemas by merging properties
2. **"default" responses** - Replaces them with proper "200" responses (when no 200 exists)

## Usage

```bash
go build -o openapi-fixer
./openapi-fixer -input spec.json -output spec-fixed.json
```

## Why this exists

terraform-plugin-codegen-openapi doesn't support allOf/oneOf constructs and doesn't handle "default" response codes properly. This tool preprocesses the OpenAPI spec to make it more generator-friendly without losing the structure.

## Testing

```bash
go test
```