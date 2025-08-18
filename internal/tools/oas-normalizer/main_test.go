package main

import (
	"reflect"
	"testing"
)

func TestFixDefaultResponses(t *testing.T) {
	spec := map[string]any{
		"paths": map[string]any{
			"/test": map[string]any{
				"get": map[string]any{
					"responses": map[string]any{
						"default": map[string]any{
							"description": "Success",
						},
					},
				},
			},
		},
	}

	fixDefaultResponses(spec)

	paths := spec["paths"].(map[string]any)
	testPath := paths["/test"].(map[string]any)
	getOp := testPath["get"].(map[string]any)
	responses := getOp["responses"].(map[string]any)

	if _, hasDefault := responses["default"]; hasDefault {
		t.Error("Expected default response to be removed")
	}

	if _, has200 := responses["200"]; !has200 {
		t.Error("Expected 200 response to be added")
	}
}

func TestFixDefaultResponsesPreserves200(t *testing.T) {
	spec := map[string]any{
		"paths": map[string]any{
			"/test": map[string]any{
				"get": map[string]any{
					"responses": map[string]any{
						"200": map[string]any{
							"description": "OK",
						},
						"default": map[string]any{
							"description": "Error",
						},
					},
				},
			},
		},
	}

	fixDefaultResponses(spec)

	paths := spec["paths"].(map[string]any)
	testPath := paths["/test"].(map[string]any)
	getOp := testPath["get"].(map[string]any)
	responses := getOp["responses"].(map[string]any)

	if _, hasDefault := responses["default"]; !hasDefault {
		t.Error("Expected default response to be preserved when 200 exists")
	}

	resp200 := responses["200"].(map[string]any)
	if resp200["description"] != "OK" {
		t.Error("Expected 200 response to remain unchanged")
	}
}

func TestFlattenAllOf(t *testing.T) {
	allOf := []any{
		map[string]any{
			"type": "object",
			"properties": map[string]any{
				"name": map[string]any{
					"type": "string",
				},
			},
			"required": []any{"name"},
		},
		map[string]any{
			"properties": map[string]any{
				"age": map[string]any{
					"type": "integer",
				},
			},
			"required": []any{"age"},
		},
	}

	result := flattenAllOf(allOf)

	if result["type"] != "object" {
		t.Error("Expected type to be object")
	}

	props := result["properties"].(map[string]any)
	if len(props) != 2 {
		t.Errorf("Expected 2 properties, got %d", len(props))
	}

	if _, hasName := props["name"]; !hasName {
		t.Error("Expected name property")
	}

	if _, hasAge := props["age"]; !hasAge {
		t.Error("Expected age property")
	}

	required := result["required"].([]any)
	if len(required) != 2 {
		t.Errorf("Expected 2 required fields, got %d", len(required))
	}
}

func TestFlattenOneOf(t *testing.T) {
	oneOf := []any{
		map[string]any{
			"type": "object",
			"properties": map[string]any{
				"name": map[string]any{
					"type": "string",
				},
			},
		},
		map[string]any{
			"type": "object",
			"properties": map[string]any{
				"id": map[string]any{
					"type": "integer",
				},
			},
		},
	}

	result := flattenOneOf(oneOf)

	if result["type"] != "object" {
		t.Error("Expected type to be object")
	}

	props := result["properties"].(map[string]any)
	if len(props) != 2 {
		t.Errorf("Expected 2 properties, got %d", len(props))
	}

	if _, hasName := props["name"]; !hasName {
		t.Error("Expected name property")
	}

	if _, hasID := props["id"]; !hasID {
		t.Error("Expected id property")
	}
}

func TestFlattenAllOfOneOfSimple(t *testing.T) {
	input := map[string]any{
		"allOf": []any{
			map[string]any{
				"type": "object",
				"properties": map[string]any{
					"name": map[string]any{"type": "string"},
				},
			},
			map[string]any{
				"properties": map[string]any{
					"age": map[string]any{"type": "integer"},
				},
			},
		},
	}

	result := flattenAllOfOneOf(input).(map[string]any)

	if _, hasAllOf := result["allOf"]; hasAllOf {
		t.Error("Expected allOf to be removed")
	}

	if result["type"] != "object" {
		t.Errorf("Expected type to be 'object', got %v", result["type"])
	}

	props, hasProps := result["properties"].(map[string]any)
	if !hasProps {
		t.Fatal("Expected properties to exist")
	}

	if len(props) != 2 {
		t.Errorf("Expected 2 properties, got %d", len(props))
	}

	if _, hasName := props["name"]; !hasName {
		t.Error("Expected name property")
	}

	if _, hasAge := props["age"]; !hasAge {
		t.Error("Expected age property")
	}
}

func TestRemoveDuplicates(t *testing.T) {
	input := []any{"name", "age", "name", "email", "age"}
	result := removeDuplicates(input)

	expected := []any{"name", "age", "email"}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}

func TestDeepCopy(t *testing.T) {
	original := map[string]any{
		"string": "value",
		"number": 42,
		"nested": map[string]any{
			"array": []any{1, 2, 3},
		},
	}

	copied := deepCopy(original).(map[string]any)

	// Modify the copy
	copied["string"] = "modified"
	nestedCopy := copied["nested"].(map[string]any)
	arrayCopy := nestedCopy["array"].([]any)
	arrayCopy[0] = 99

	// Original should be unchanged
	if original["string"] != "value" {
		t.Error("Original string was modified")
	}

	originalNested := original["nested"].(map[string]any)
	originalArray := originalNested["array"].([]any)
	if originalArray[0] != 1 {
		t.Error("Original array was modified")
	}
}

func TestFixOpenAPISpec(t *testing.T) {
	// TODO
}
