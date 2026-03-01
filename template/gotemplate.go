package template

import (
	"bytes"
	"encoding/json"
	"fmt"
	gotemplate "text/template"
)

// FromJSON creates a Document from a JSON schema definition.
// The JSON may contain Go template expressions (e.g., {{.Field}}) that are
// resolved using the provided data. If data is nil, no template processing
// is performed and the JSON is parsed directly.
//
// Optional Option values (WithFont, WithDefaultFont, etc.) override any
// settings defined in the JSON schema.
//
// Example:
//
//	schema := []byte(`{
//	  "page": {"size": "A4", "margins": "20mm"},
//	  "body": [{"row": {"cols": [
//	    {"span": 12, "text": "Hello, {{.Name}}!", "style": {"size": 24, "bold": true}}
//	  ]}}]
//	}`)
//	doc, err := template.FromJSON(schema, map[string]string{"Name": "World"})
//	if err != nil { ... }
//	pdf, err := doc.Generate()
func FromJSON(schema []byte, data any, opts ...Option) (*Document, error) {
	jsonBytes := schema

	if data != nil {
		tmpl, err := gotemplate.New("schema").Funcs(defaultFuncMap()).Parse(string(schema))
		if err != nil {
			return nil, fmt.Errorf("gpdf: parsing template: %w", err)
		}
		var buf bytes.Buffer
		if err := tmpl.Execute(&buf, data); err != nil {
			return nil, fmt.Errorf("gpdf: executing template: %w", err)
		}
		jsonBytes = buf.Bytes()
	}

	var s Schema
	if err := json.Unmarshal(jsonBytes, &s); err != nil {
		return nil, fmt.Errorf("gpdf: parsing JSON schema: %w", err)
	}

	return buildFromSchema(&s, opts)
}

// FromTemplate executes a pre-parsed Go template with data and creates a
// Document from the resulting JSON output. The template must produce valid
// JSON conforming to the Schema format.
//
// Use TemplateFuncMap to get helper functions (like toJSON) when parsing
// templates:
//
//	tmpl := template.Must(gotemplate.New("").Funcs(template.TemplateFuncMap()).Parse(str))
//	doc, err := template.FromTemplate(tmpl, data)
func FromTemplate(tmpl *gotemplate.Template, data any, opts ...Option) (*Document, error) {
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return nil, fmt.Errorf("gpdf: executing template: %w", err)
	}

	var s Schema
	if err := json.Unmarshal(buf.Bytes(), &s); err != nil {
		return nil, fmt.Errorf("gpdf: parsing JSON from template output: %w", err)
	}

	return buildFromSchema(&s, opts)
}

// TemplateFuncMap returns a template.FuncMap containing helper functions
// for use with Go templates that produce JSON schemas for gpdf.
//
// Available functions:
//   - toJSON: Marshals a Go value to its JSON representation.
//     Usage in template: {{toJSON .Items}}
func TemplateFuncMap() gotemplate.FuncMap {
	return defaultFuncMap()
}

// defaultFuncMap returns the built-in template function map.
func defaultFuncMap() gotemplate.FuncMap {
	return gotemplate.FuncMap{
		"toJSON": toJSONFunc,
	}
}

// toJSONFunc marshals a value to a JSON string. This is essential in
// templates for embedding structured data like table rows.
func toJSONFunc(v any) (string, error) {
	b, err := json.Marshal(v)
	if err != nil {
		return "", err
	}
	return string(b), nil
}
