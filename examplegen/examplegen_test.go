package examplegen

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/xuri/xgen"
)

func setupTestSchema(t *testing.T) []interface{} {
	parser := xgen.NewParser(&xgen.Options{
		FilePath:            "../schemas/petstore.xsd",
		IncludeMap:          make(map[string]bool),
		LocalNameNSMap:      make(map[string]string),
		NSSchemaLocationMap: make(map[string]string),
		ParseFileList:       make(map[string]bool),
		ParseFileMap:        make(map[string][]interface{}),
		ProtoTree:           make([]interface{}, 0),
		RemoteSchema:        make(map[string][]byte),
		Extract:             true,
	})
	err := parser.Parse()
	if err != nil {
		t.Fatalf("Failed to parse schema: %v", err)
	}
	return parser.ProtoTree
}

func TestGenerate_Examples(t *testing.T) {
	protoTree := setupTestSchema(t)
	elements := []string{"getPetByIdRequest", "getPetByIdResponse", "fault"}

	t.Log("\nExample XML for each element type:")
	for _, element := range elements {
		xmlStr := Generate(protoTree, element)
		t.Logf("\n%s:\n%s", element, xmlStr)
	}
}

func TestGenerate_GetPetByIdRequest(t *testing.T) {
	protoTree := setupTestSchema(t)
	xmlStr := Generate(protoTree, "getPetByIdRequest")

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
	protoTree := setupTestSchema(t)
	xmlStr := Generate(protoTree, "getPetByIdResponse")

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
	protoTree := setupTestSchema(t)
	xmlStr := Generate(protoTree, "fault")

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
	protoTree := setupTestSchema(t)
	xmlStr := Generate(protoTree, "nonexistent")

	if !strings.Contains(xmlStr, "<!-- Element nonexistent not found -->") {
		t.Error("Expected comment for non-existent element")
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
