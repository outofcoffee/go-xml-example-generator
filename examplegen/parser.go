package examplegen

import (
	"fmt"

	"github.com/xuri/xgen"
)

// parseSchema parses an XSD schema file and returns its proto tree
func parseSchema(schemaPath string) ([]interface{}, error) {
	parser := xgen.NewParser(&xgen.Options{
		FilePath:            schemaPath,
		IncludeMap:          make(map[string]bool),
		LocalNameNSMap:      make(map[string]string),
		NSSchemaLocationMap: make(map[string]string),
		ParseFileList:       make(map[string]bool),
		ParseFileMap:        make(map[string][]interface{}),
		ProtoTree:           make([]interface{}, 0),
		RemoteSchema:        make(map[string][]byte),
		Extract:             true,
	})

	if err := parser.Parse(); err != nil {
		return nil, fmt.Errorf("failed to parse schema: %w", err)
	}

	return parser.ProtoTree, nil
}
