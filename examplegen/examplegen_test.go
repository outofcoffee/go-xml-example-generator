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
	tests := []struct {
		name       string
		xsdPath    string
		element    string
		wantErr    bool
		validation func(string) bool
	}{
		{
			name:    "getPetByIdRequest",
			xsdPath: "../schemas/simple/petstore.xsd",
			element: "getPetByIdRequest",
			validation: func(xml string) bool {
				return strings.Contains(xml, "<getPetByIdRequest>") &&
					strings.Contains(xml, "<id>")
			},
		},
		{
			name:    "getPetByIdResponse",
			xsdPath: "../schemas/simple/petstore.xsd",
			element: "getPetByIdResponse",
			validation: func(xml string) bool {
				return strings.Contains(xml, "<getPetByIdResponse>") &&
					strings.Contains(xml, "<id>") &&
					strings.Contains(xml, "<name>")
			},
		},
		{
			name:    "fault",
			xsdPath: "../schemas/simple/petstore.xsd",
			element: "fault",
			validation: func(xml string) bool {
				return strings.Contains(xml, "<fault>") &&
					!strings.Contains(xml, "<fault></fault>")
			},
		},
		{
			name:    "non-existent element",
			xsdPath: "../schemas/simple/petstore.xsd",
			element: "nonexistent",
			validation: func(xml string) bool {
				return strings.Contains(xml, "<!-- Element nonexistent not found -->")
			},
		},
		{
			name:    "non-existent schema",
			xsdPath: "nonexistent.xsd",
			element: "element",
			wantErr: true,
		},
		{
			name:    "element-ref getPetByIdRequest",
			xsdPath: "../schemas/element-ref/petstore.xsd",
			element: "getPetByIdRequest",
			validation: func(xml string) bool {
				return strings.Contains(xml, "<getPetByIdRequest>") &&
					strings.Contains(xml, "<id>") &&
					!strings.Contains(xml, "tns:id") &&
					!strings.Contains(xml, "</id></id>")
			},
		},
		{
			name:    "element-ref getPetByIdResponse",
			xsdPath: "../schemas/element-ref/petstore.xsd",
			element: "getPetByIdResponse",
			validation: func(xml string) bool {
				return strings.Contains(xml, "<getPetByIdResponse>") &&
					strings.Contains(xml, "<id>") &&
					!strings.Contains(xml, "tns:id") &&
					strings.Contains(xml, "<name>") &&
					!strings.Contains(xml, "</id></id>")
			},
		},
		{
			name:    "element-ref fault",
			xsdPath: "../schemas/element-ref/petstore.xsd",
			element: "fault",
			validation: func(xml string) bool {
				return strings.Contains(xml, "<fault>") &&
					!strings.Contains(xml, "<fault></fault>")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			xmlStr, err := Generate(tt.xsdPath, tt.element)
			if (err != nil) != tt.wantErr {
				t.Errorf("Generate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil {
				return
			}

			// Verify it's valid XML unless it's a non-existent element case
			if !strings.Contains(xmlStr, "<!-- Element") {
				if err := xml.Unmarshal([]byte(xmlStr), new(interface{})); err != nil {
					t.Errorf("Generated XML is not valid: %v\nXML: %s", err, xmlStr)
				}
			}

			// Run custom validation if provided
			if tt.validation != nil && !tt.validation(xmlStr) {
				t.Errorf("Generated XML failed validation\nXML: %s", xmlStr)
			}
		})
	}
}

func TestGenerateSimpleTypes(t *testing.T) {
	g := NewGenerator(nil, false) // nil is fine as we're not using the protoTree, and elementFormQual doesn't matter for simple types
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

func TestGenerateWithNs(t *testing.T) {
	tests := []struct {
		name      string
		xsdPath   string
		element   string
		namespace string
		prefix    string
		validate  func(string) bool
	}{
		{
			name:      "getPetByIdResponse with namespace",
			xsdPath:   "../schemas/simple/petstore.xsd",
			element:   "getPetByIdResponse",
			namespace: "urn:foo:bar",
			prefix:    "foo",
			validate: func(xml string) bool {
				return strings.Contains(xml, `xmlns:foo="urn:foo:bar"`) &&
					strings.Contains(xml, "<foo:getPetByIdResponse") &&
					strings.Contains(xml, "<foo:id>") &&
					strings.Contains(xml, "</foo:id>") &&
					strings.Contains(xml, "<foo:name>") &&
					strings.Contains(xml, "</foo:name>") &&
					strings.Contains(xml, "</foo:getPetByIdResponse>")
			},
		},
		{
			name:      "empty prefix",
			xsdPath:   "../schemas/simple/petstore.xsd",
			element:   "getPetByIdResponse",
			namespace: "urn:foo:bar",
			prefix:    "",
			validate: func(xml string) bool {
				return !strings.Contains(xml, "xmlns:") &&
					!strings.Contains(xml, ":getPetByIdResponse")
			},
		},
		{
			name:      "empty namespace",
			xsdPath:   "../schemas/simple/petstore.xsd",
			element:   "getPetByIdResponse",
			namespace: "",
			prefix:    "foo",
			validate: func(xml string) bool {
				return !strings.Contains(xml, "xmlns:")
			},
		},
		{
			name:      "element-ref getPetByIdResponse with namespace",
			xsdPath:   "../schemas/element-ref/petstore.xsd",
			element:   "getPetByIdResponse",
			namespace: "urn:com:example:petstore",
			prefix:    "tns",
			validate: func(xml string) bool {
				return strings.Contains(xml, `xmlns:tns="urn:com:example:petstore"`) &&
					strings.Contains(xml, "<tns:getPetByIdResponse") &&
					strings.Contains(xml, "<tns:id>") &&
					strings.Contains(xml, "<tns:name>") &&
					strings.Contains(xml, "</tns:getPetByIdResponse>") &&
					!strings.Contains(xml, "</tns:id></tns:id>")
			},
		},
		{
			name:      "element-ref getPetByIdRequest with namespace",
			xsdPath:   "../schemas/element-ref/petstore.xsd",
			element:   "getPetByIdRequest",
			namespace: "urn:com:example:petstore",
			prefix:    "tns",
			validate: func(xml string) bool {
				return strings.Contains(xml, `xmlns:tns="urn:com:example:petstore"`) &&
					strings.Contains(xml, "<tns:getPetByIdRequest") &&
					strings.Contains(xml, "<tns:id>") &&
					strings.Contains(xml, "</tns:getPetByIdRequest>") &&
					!strings.Contains(xml, "</tns:id></tns:id>")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			xmlStr, err := GenerateWithNs(tt.xsdPath, tt.element, tt.namespace, tt.prefix)
			if err != nil {
				t.Fatalf("Failed to generate XML: %v", err)
			}

			// Verify it's valid XML
			if err := xml.Unmarshal([]byte(xmlStr), new(interface{})); err != nil {
				t.Errorf("Generated XML is not valid: %v\nXML: %s", err, xmlStr)
			}

			// Run custom validation
			if !tt.validate(xmlStr) {
				t.Errorf("Generated XML failed validation\nXML: %s", xmlStr)
			}
		})
	}
}
