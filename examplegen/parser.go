package examplegen

import (
	"encoding/xml"
	"fmt"
	"os"
	"strings"

	"github.com/outofcoffee/xgen"
)

type xsdElement struct {
	XMLName xml.Name `xml:"element"`
	Ref     string   `xml:"ref,attr"`
	Name    string   `xml:"name,attr"`
}

type xsdComplexType struct {
	XMLName  xml.Name     `xml:"complexType"`
	Name     string       `xml:"name,attr"`
	All      []xsdElement `xml:"all>element"`
	Sequence []xsdElement `xml:"sequence>element"`
}

type xsdSchema struct {
	XMLName            xml.Name         `xml:"schema"`
	ElementFormDefault string           `xml:"elementFormDefault,attr"`
	TargetNamespace    string           `xml:"targetNamespace,attr"`
	Elements           []xsdElement     `xml:"element"`
	ComplexTypes       []xsdComplexType `xml:"complexType"`
	Attrs              []xml.Attr       `xml:",any,attr"`
}

// parseSchema parses an XSD schema file and returns its proto tree and elementFormDefault value
func parseSchema(schemaPath string) ([]interface{}, bool, error) {
	// Read the schema file to get elementFormDefault and parse references
	schemaBytes, err := os.ReadFile(schemaPath)
	if err != nil {
		return nil, false, fmt.Errorf("failed to read schema file: %w", err)
	}

	// Parse just enough of the XML to get the schema attributes
	var schema xsdSchema
	decoder := xml.NewDecoder(strings.NewReader(string(schemaBytes)))
	if err := decoder.Decode(&schema); err != nil {
		return nil, false, fmt.Errorf("failed to parse schema attributes: %w", err)
	}

	// Create the xgen parser
	parser := xgen.NewParser(&xgen.Options{
		FilePath:            schemaPath,
		IncludeMap:          make(map[string]bool),
		LocalNameNSMap:      make(map[string]string),
		NSSchemaLocationMap: make(map[string]string),
		ParseFileList:       make(map[string]bool),
		ParseFileMap:        make(map[string][]interface{}),
		ProtoTree:           make([]interface{}, 0),
		RemoteSchema:        make(map[string][]byte),
		Extract:             false,
		SkipGenerate:        true,
	})

	if err := parser.Parse(); err != nil {
		return nil, false, fmt.Errorf("failed to parse schema: %w", err)
	}

	// Default is "unqualified" if not specified
	elementFormQual := schema.ElementFormDefault == "qualified"

	return parser.ProtoTree, elementFormQual, nil
}
