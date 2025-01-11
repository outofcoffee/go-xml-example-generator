package examplegen

import (
	"bytes"
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/xuri/xgen"
)

// Generator holds the state for XML snippet generation
type Generator struct {
	typeStack []string
	nElements int
	r         *rand.Rand
	protoTree []interface{}
	namespace string
	prefix    string
}

// words is a slice of sample words for text content
var words = []string{
	"colour", "flavour", "behaviour", "humour", "labour",
	"neighbour", "rumour", "splendour", "vigour", "valour",
}

// NewGenerator creates a new XML snippet generator
func NewGenerator(protoTree []interface{}) *Generator {
	return &Generator{
		typeStack: make([]string, 0),
		r:         rand.New(rand.NewSource(time.Now().UnixNano())),
		protoTree: protoTree,
	}
}

// Generate generates an example XML snippet for the given element from an XSD schema file
func Generate(schemaPath string, elementName string) (string, error) {
	protoTree, err := parseSchema(schemaPath)
	if err != nil {
		return "", err
	}

	g := NewGenerator(protoTree)
	return g.generateXML(elementName), nil
}

// GenerateWithNs generates an example XML snippet with the specified namespace and prefix
func GenerateWithNs(schemaPath string, elementName string, namespace string, prefix string) (string, error) {
	protoTree, err := parseSchema(schemaPath)
	if err != nil {
		return "", err
	}

	g := NewGenerator(protoTree)
	g.namespace = namespace
	g.prefix = prefix
	return g.generateXML(elementName), nil
}

// generateXML is the main entry point for XML generation
func (g *Generator) generateXML(elementName string) string {
	element := g.findElement(elementName)
	if element == nil {
		return fmt.Sprintf("<!-- Element %s not found -->", elementName)
	}

	var buf bytes.Buffer
	g.generateElement(&buf, element, 0)
	return buf.String()
}

// findElement searches for an element by name in the proto tree
func (g *Generator) findElement(elementName string) *xgen.Element {
	for _, item := range g.protoTree {
		if element, ok := item.(*xgen.Element); ok {
			if element.Name == elementName {
				return element
			}
		}
	}
	return nil
}

// findComplexType searches for a complex type by name
func (g *Generator) findComplexType(typeName string) *xgen.ComplexType {
	// Remove any namespace prefix (e.g., "tns:petType" -> "petType")
	if idx := strings.Index(typeName, ":"); idx >= 0 {
		typeName = typeName[idx+1:]
	}

	for _, item := range g.protoTree {
		if ct, ok := item.(*xgen.ComplexType); ok {
			if ct.Name == typeName {
				return ct
			}
		}
	}
	return nil
}

// generateElement handles generation of a single element
func (g *Generator) generateElement(buf *bytes.Buffer, element *xgen.Element, indent int) {
	if g.nElements > 1000 { // Limit to prevent infinite recursion
		return
	}
	g.nElements++

	// Write opening tag with indentation
	g.writeIndent(buf, indent)
	buf.WriteString("<")
	if g.prefix != "" {
		buf.WriteString(g.prefix)
		buf.WriteString(":")
	}
	buf.WriteString(element.Name)

	// Write namespace declaration for root element (indent == 0)
	if indent == 0 && g.namespace != "" && g.prefix != "" {
		buf.WriteString(fmt.Sprintf(" xmlns:%s=\"%s\"", g.prefix, g.namespace))
	}

	// Write closing tag
	if element.Type == "" {
		buf.WriteString("/>\n")
		return
	}

	buf.WriteString(">")

	// Generate content based on type
	complexType := g.findComplexType(element.Type)
	if complexType != nil {
		buf.WriteString("\n")
		g.generateComplexTypeContent(buf, complexType, indent+2)
		g.writeIndent(buf, indent)
	} else {
		// Remove any namespace prefix for simple types
		simpleType := element.Type
		if idx := strings.Index(simpleType, ":"); idx >= 0 {
			simpleType = simpleType[idx+1:]
		}
		g.generateSimpleTypeContent(buf, simpleType)
	}

	buf.WriteString("</")
	if g.prefix != "" {
		buf.WriteString(g.prefix)
		buf.WriteString(":")
	}
	buf.WriteString(element.Name)
	buf.WriteString(">\n")
}

// generateComplexTypeContent generates content for complex types
func (g *Generator) generateComplexTypeContent(buf *bytes.Buffer, complexType *xgen.ComplexType, indent int) {
	// Generate each element in the complex type
	for _, element := range complexType.Elements {
		g.writeIndent(buf, indent)
		buf.WriteString("<")
		if g.prefix != "" {
			buf.WriteString(g.prefix)
			buf.WriteString(":")
		}
		buf.WriteString(element.Name)
		buf.WriteString(">")
		// Remove any namespace prefix for simple types
		simpleType := element.Type
		if idx := strings.Index(simpleType, ":"); idx >= 0 {
			simpleType = simpleType[idx+1:]
		}
		g.generateSimpleTypeContent(buf, simpleType)
		buf.WriteString("</")
		if g.prefix != "" {
			buf.WriteString(g.prefix)
			buf.WriteString(":")
		}
		buf.WriteString(element.Name)
		buf.WriteString(">\n")
	}
}

// generateSimpleTypeContent generates content for simple types
func (g *Generator) generateSimpleTypeContent(buf *bytes.Buffer, typeName string) {
	switch typeName {
	case "string", "xs:string":
		buf.WriteString(g.generateSampleString())
	case "int", "xs:int", "integer", "xs:integer":
		buf.WriteString(fmt.Sprintf("%d", g.r.Intn(100)))
	case "decimal", "xs:decimal", "float", "xs:float", "double", "xs:double":
		buf.WriteString(fmt.Sprintf("%.2f", g.r.Float64()*100))
	case "boolean", "xs:boolean":
		if g.r.Intn(2) == 0 {
			buf.WriteString("true")
		} else {
			buf.WriteString("false")
		}
	case "date", "xs:date":
		now := time.Now()
		buf.WriteString(now.Format("2006-01-02"))
	case "time", "xs:time":
		now := time.Now()
		buf.WriteString(now.Format("15:04:05"))
	case "dateTime", "xs:dateTime":
		now := time.Now()
		buf.WriteString(now.Format("2006-01-02T15:04:05"))
	default:
		buf.WriteString(g.generateSampleString())
	}
}

// generateSampleString generates a sample string value
func (g *Generator) generateSampleString() string {
	return words[g.r.Intn(len(words))]
}

// writeIndent writes the specified number of spaces for indentation
func (g *Generator) writeIndent(buf *bytes.Buffer, indent int) {
	buf.WriteString(strings.Repeat(" ", indent))
}
