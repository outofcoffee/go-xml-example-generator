# go-xml-example-generator

A Go tool that generates example XML snippets from XSD schemas. Given an XSD schema and an element name, it produces a valid XML example that matches the schema's structure.

## Features

- Generates example XML for any element defined in an XSD schema
- Handles complex types with nested elements
- Supports common XSD types (string, int, boolean, date, etc.)
- Uses British English sample text values
- Preserves XML structure and indentation

## Installation

As a command line tool:
```bash
go install github.com/outofcoffee/go-xml-example-generator@latest
```

As a dependency in your Go project:
```bash
go get github.com/outofcoffee/go-xml-example-generator
```

## Usage

### Command Line
```bash
go-xml-example-generator [xsd_path] [element_name]
```

For example:
```bash
go-xml-example-generator schemas/petstore.xsd getPetByIdResponse
```

This will output something like:
```xml
<getPetByIdResponse>
  <id>42</id>
  <name>colour</name>
</getPetByIdResponse>
```

### As a Library
```go
import "github.com/outofcoffee/go-xml-example-generator/examplegen"

// Generate XML from an XSD file
xml, err := examplegen.Generate("path/to/schema.xsd", "elementName")
if err != nil {
    log.Fatal(err)
}
fmt.Println(xml)
```

## Building from Source

```bash
git clone https://github.com/outofcoffee/go-xml-example-generator.git
cd go-xml-example-generator
go build
```

## Running Tests

```bash
go test ./...
```

## Dependencies

- github.com/xuri/xgen - XML schema parser

## Contributing

Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.

Please make sure to update tests as appropriate.

## License

[MIT](https://choosealicense.com/licenses/mit/)

## Author

[outofcoffee](https://github.com/outofcoffee) 