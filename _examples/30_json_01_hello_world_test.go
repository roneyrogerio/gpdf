package examples_test

import (
	"testing"

	"github.com/gpdf-dev/gpdf/template"
)

func TestExample_30_JSON_01_HelloWorld(t *testing.T) {
	schema := []byte(`{
		"page": {"size": "A4", "margins": "20mm"},
		"body": [
			{"row": {"cols": [
				{"span": 12, "text": "Hello, World!", "style": {"size": 24, "bold": true}}
			]}}
		]
	}`)

	doc, err := template.FromJSON(schema, nil)
	if err != nil {
		t.Fatalf("FromJSON error: %v", err)
	}
	generatePDF(t, "30_json_01_hello_world.pdf", doc)
}
