package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/outofcoffee/go-xml-example-generator/examplegen"
)

func main() {
	// Parse command line arguments
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [xsd_path] [element_name]\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "\nGenerates example XML for the specified element from an XSD schema.\n\n")
		flag.PrintDefaults()
	}
	flag.Parse()

	if flag.NArg() != 2 {
		flag.Usage()
		os.Exit(1)
	}

	xsdPath := flag.Arg(0)
	elementName := flag.Arg(1)

	// Generate the XML
	xml, err := examplegen.Generate(xsdPath, elementName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error generating XML: %v\n", err)
		os.Exit(1)
	}

	fmt.Println(xml)
}
