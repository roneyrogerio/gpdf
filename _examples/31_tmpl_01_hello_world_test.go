package examples_test

import (
	"testing"

	"github.com/gpdf-dev/gpdf/template"
)

func TestExample_31_Tmpl_01_HelloWorld(t *testing.T) {
	schema := []byte(`{
		"page": {"size": "A4", "margins": "20mm"},
		"body": [
			{"row": {"cols": [
				{"span": 12, "text": "{{.Greeting}}", "style": {"size": 24, "bold": true}}
			]}}
		]
	}`)

	data := map[string]any{
		"Greeting": "Hello, World!",
	}

	doc, err := template.FromJSON(schema, data)
	if err != nil {
		t.Fatalf("FromJSON error: %v", err)
	}
	generatePDF(t, "31_tmpl_01_hello_world.pdf", doc)
}
