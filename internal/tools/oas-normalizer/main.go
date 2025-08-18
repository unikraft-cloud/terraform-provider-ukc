package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
)

type OpenAPISpec map[string]any

func main() {
	var inputFile, outputFile string
	flag.StringVar(&inputFile, "input", "", "Input OpenAPI spec file (required)")
	flag.StringVar(&outputFile, "output", "", "Output OpenAPI spec file (required)")
	flag.Parse()

	if inputFile == "" || outputFile == "" {
		fmt.Println("Usage: openapi-fixer -input <input.json> -output <output.json>")
		flag.PrintDefaults()
		os.Exit(1)
	}

	data, err := os.ReadFile(inputFile)
	if err != nil {
		log.Fatalf("Error reading input file: %v", err)
	}

	var spec OpenAPISpec
	if err := json.Unmarshal(data, &spec); err != nil {
		log.Fatalf("Error parsing JSON: %v", err)
	}

	fixedSpec := fixOpenAPISpec(spec)

	output, err := json.MarshalIndent(fixedSpec, "", "  ")
	if err != nil {
		log.Fatalf("Error marshaling JSON: %v", err)
	}

	if err := os.WriteFile(outputFile, output, 0644); err != nil {
		log.Fatalf("Error writing output file: %v", err)
	}

	fmt.Printf("Successfully fixed OpenAPI spec and wrote to %s\n", outputFile)
}

func fixOpenAPISpec(spec OpenAPISpec) OpenAPISpec {
	fixed := deepCopy(spec).(OpenAPISpec)
	fixDefaultResponses(fixed)
	fixed = flattenAllOfOneOf(fixed).(OpenAPISpec)
	return fixed
}

func fixDefaultResponses(spec map[string]any) {
	paths, ok := spec["paths"].(map[string]any)
	if !ok {
		return
	}

	for _, pathItem := range paths {
		pathMap, ok := pathItem.(map[string]any)
		if !ok {
			continue
		}

		for method, operation := range pathMap {
			if method == "parameters" || method == "summary" || method == "description" {
				continue
			}

			operationMap, ok := operation.(map[string]any)
			if !ok {
				continue
			}

			responses, ok := operationMap["responses"].(map[string]any)
			if !ok {
				continue
			}

			if defaultResp, hasDefault := responses["default"]; hasDefault {
				if _, has200 := responses["200"]; !has200 {
					responses["200"] = defaultResp
					delete(responses, "default")
				}
			}
		}
	}
}

func flattenAllOfOneOf(obj any) any {
	switch v := obj.(type) {
	case map[string]any:
		result := make(map[string]any)

		if allOf, hasAllOf := v["allOf"]; hasAllOf {
			result = flattenAllOf(allOf)
			for key, value := range v {
				if key != "allOf" {
					result[key] = flattenAllOfOneOf(value)
				}
			}
			return result
		}

		if oneOf, hasOneOf := v["oneOf"]; hasOneOf {
			result = flattenOneOf(oneOf)
			for key, value := range v {
				if key != "oneOf" {
					result[key] = flattenAllOfOneOf(value)
				}
			}
			return result
		}

		for key, value := range v {
			result[key] = flattenAllOfOneOf(value)
		}
		return result

	case []any:
		result := make([]any, len(v))
		for i, item := range v {
			result[i] = flattenAllOfOneOf(item)
		}
		return result

	default:
		return v
	}
}

func flattenAllOf(allOf any) map[string]any {
	result := make(map[string]any)

	allOfSlice, ok := allOf.([]any)
	if !ok {
		return result
	}

	// Initialize properties map
	result["properties"] = make(map[string]any)

	for _, item := range allOfSlice {
		itemMap, ok := item.(map[string]any)
		if !ok {
			continue
		}

		flattened := flattenAllOfOneOf(itemMap).(map[string]any)

		for key, value := range flattened {
			switch key {
			case "properties":
				resultProps := result["properties"].(map[string]any)
				if valueProps, ok := value.(map[string]any); ok {
					for propKey, propValue := range valueProps {
						resultProps[propKey] = propValue
					}
				}
			case "required":
				if resultRequired, hasRequired := result["required"].([]any); hasRequired {
					if valueRequired, ok := value.([]any); ok {
						result["required"] = append(resultRequired, valueRequired...)
					}
				} else {
					result[key] = value
				}
			default:
				result[key] = value
			}
		}
	}

	// Clean up empty properties
	if props := result["properties"].(map[string]any); len(props) == 0 {
		delete(result, "properties")
	}

	// Ensure we have a type if we have properties
	if _, hasType := result["type"]; !hasType {
		if _, hasProps := result["properties"]; hasProps {
			result["type"] = "object"
		}
	}

	return result
}

func flattenOneOf(oneOf any) map[string]any {
	result := make(map[string]any)

	oneOfSlice, ok := oneOf.([]any)
	if !ok {
		return result
	}

	allProperties := make(map[string]any)
	var allRequired []any

	for _, item := range oneOfSlice {
		itemMap, ok := item.(map[string]any)
		if !ok {
			continue
		}

		flattened := flattenAllOfOneOf(itemMap).(map[string]any)

		if props, hasProps := flattened["properties"].(map[string]any); hasProps {
			for propKey, propValue := range props {
				allProperties[propKey] = propValue
			}
		}

		if required, hasRequired := flattened["required"].([]any); hasRequired {
			allRequired = append(allRequired, required...)
		}

		for key, value := range flattened {
			switch key {
			case "properties", "required", "title":
				// Skip these as they're handled above
			default:
				result[key] = value
			}
		}
	}

	if len(allProperties) > 0 {
		result["properties"] = allProperties
	}
	if len(allRequired) > 0 {
		result["required"] = removeDuplicates(allRequired)
	}

	if _, hasType := result["type"]; !hasType && len(allProperties) > 0 {
		result["type"] = "object"
	}

	return result
}

func removeDuplicates(slice []any) []any {
	seen := make(map[string]bool)
	var result []any

	for _, item := range slice {
		if str, ok := item.(string); ok {
			if !seen[str] {
				seen[str] = true
				result = append(result, item)
			}
		}
	}

	return result
}

func deepCopy(obj any) any {
	if obj == nil {
		return nil
	}

	switch v := obj.(type) {
	case map[string]any:
		result := make(map[string]any)
		for key, value := range v {
			result[key] = deepCopy(value)
		}
		return result

	case []any:
		result := make([]any, len(v))
		for i, item := range v {
			result[i] = deepCopy(item)
		}
		return result

	default:
		return v
	}
}
