package examplegen

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"strings"
	"testing"
	"time"
)

func TestGenerate_Examples(t *testing.T) {
	elements := []string{"getPetByIdRequest", "getPetByIdResponse", "fault"}

	t.Log("\nExample XML for each element type:")
	for _, element := range elements {
		xmlStr, err := Generate("../schemas/petstore.xsd", element)
		if err != nil {
			t.Errorf("Failed to generate XML for %s: %v", element, err)
			continue
		}
		t.Logf("\n%s:\n%s", element, xmlStr)
	}
}

func TestGenerate_GetPetByIdRequest(t *testing.T) {
	xmlStr, err := Generate("../schemas/petstore.xsd", "getPetByIdRequest")
	if err != nil {
		t.Fatalf("Failed to generate XML: %v", err)
	}

	// Verify it's valid XML
	if err := xml.Unmarshal([]byte(xmlStr), new(interface{})); err != nil {
		t.Errorf("Generated XML is not valid: %v\nXML: %s", err, xmlStr)
	}

	// Verify structure
	if !strings.Contains(xmlStr, "<getPetByIdRequest>") {
		t.Error("Missing root element getPetByIdRequest")
	}
	if !strings.Contains(xmlStr, "<id>") {
		t.Error("Missing id element")
	}
}

func TestGenerate_GetPetByIdResponse(t *testing.T) {
	xmlStr, err := Generate("../schemas/petstore.xsd", "getPetByIdResponse")
	if err != nil {
		t.Fatalf("Failed to generate XML: %v", err)
	}

	// Verify it's valid XML
	if err := xml.Unmarshal([]byte(xmlStr), new(interface{})); err != nil {
		t.Errorf("Generated XML is not valid: %v\nXML: %s", err, xmlStr)
	}

	// Verify structure
	if !strings.Contains(xmlStr, "<getPetByIdResponse>") {
		t.Error("Missing root element getPetByIdResponse")
	}
	if !strings.Contains(xmlStr, "<id>") {
		t.Error("Missing id element")
	}
	if !strings.Contains(xmlStr, "<name>") {
		t.Error("Missing name element")
	}
}

func TestGenerate_Fault(t *testing.T) {
	xmlStr, err := Generate("../schemas/petstore.xsd", "fault")
	if err != nil {
		t.Fatalf("Failed to generate XML: %v", err)
	}

	// Verify it's valid XML
	if err := xml.Unmarshal([]byte(xmlStr), new(interface{})); err != nil {
		t.Errorf("Generated XML is not valid: %v\nXML: %s", err, xmlStr)
	}

	// Verify structure
	if !strings.Contains(xmlStr, "<fault>") {
		t.Error("Missing root element fault")
	}
	// Since fault is a simple string type, verify it contains some content
	if strings.Contains(xmlStr, "<fault></fault>") {
		t.Error("Fault element is empty")
	}
}

func TestGenerate_NonExistentElement(t *testing.T) {
	xmlStr, err := Generate("../schemas/petstore.xsd", "nonexistent")
	if err != nil {
		t.Fatalf("Failed to generate XML: %v", err)
	}

	if !strings.Contains(xmlStr, "<!-- Element nonexistent not found -->") {
		t.Error("Expected comment for non-existent element")
	}
}

func TestGenerate_NonExistentSchema(t *testing.T) {
	_, err := Generate("nonexistent.xsd", "element")
	if err == nil {
		t.Error("Expected error for non-existent schema")
	}
}

func TestGenerateSimpleTypes(t *testing.T) {
	g := NewGenerator(nil) // nil is fine as we're not using the protoTree
	tests := []struct {
		typeName string
		validate func(string) bool
	}{
		{"xs:string", func(s string) bool { return len(s) > 0 }},
		{"xs:int", func(s string) bool {
			var i int
			_, err := fmt.Sscanf(s, "%d", &i)
			return err == nil && i >= 0 && i < 100
		}},
		{"xs:boolean", func(s string) bool {
			return s == "true" || s == "false"
		}},
		{"xs:date", func(s string) bool {
			_, err := time.Parse("2006-01-02", s)
			return err == nil
		}},
	}

	for _, tt := range tests {
		t.Run(tt.typeName, func(t *testing.T) {
			var buf bytes.Buffer
			g.generateSimpleTypeContent(&buf, tt.typeName)
			result := buf.String()
			if !tt.validate(result) {
				t.Errorf("Invalid value for %s: %s", tt.typeName, result)
			}
		})
	}
}

func TestGenerateWithNs_GetPetByIdResponse(t *testing.T) {
	xmlStr, err := GenerateWithNs("../schemas/petstore.xsd", "getPetByIdResponse", "urn:foo:bar", "foo")
	if err != nil {
		t.Fatalf("Failed to generate XML: %v", err)
	}

	// Verify it's valid XML
	if err := xml.Unmarshal([]byte(xmlStr), new(interface{})); err != nil {
		t.Errorf("Generated XML is not valid: %v\nXML: %s", err, xmlStr)
	}

	// Verify namespace declaration
	if !strings.Contains(xmlStr, `xmlns:foo="urn:foo:bar"`) {
		t.Error("Missing namespace declaration")
	}

	// Verify structure with prefixes
	expectedElements := []string{
		"<foo:getPetByIdResponse",
		"<foo:id>",
		"</foo:id>",
		"<foo:name>",
		"</foo:name>",
		"</foo:getPetByIdResponse>",
	}

	for _, expected := range expectedElements {
		if !strings.Contains(xmlStr, expected) {
			t.Errorf("Missing expected element: %s", expected)
		}
	}
}

func TestGenerateWithNs_EmptyPrefix(t *testing.T) {
	xmlStr, err := GenerateWithNs("../schemas/petstore.xsd", "getPetByIdResponse", "urn:foo:bar", "")
	if err != nil {
		t.Fatalf("Failed to generate XML: %v", err)
	}

	// Should not have any namespace declarations or prefixes
	if strings.Contains(xmlStr, "xmlns:") {
		t.Error("Should not have namespace declaration with empty prefix")
	}
	if strings.Contains(xmlStr, ":getPetByIdResponse") {
		t.Error("Should not have prefixed elements with empty prefix")
	}
}

func TestGenerateWithNs_EmptyNamespace(t *testing.T) {
	xmlStr, err := GenerateWithNs("../schemas/petstore.xsd", "getPetByIdResponse", "", "foo")
	if err != nil {
		t.Fatalf("Failed to generate XML: %v", err)
	}

	// Should not have any namespace declarations
	if strings.Contains(xmlStr, "xmlns:") {
		t.Error("Should not have namespace declaration with empty namespace")
	}
}
