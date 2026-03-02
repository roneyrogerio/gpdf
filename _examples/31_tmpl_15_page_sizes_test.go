package examples_test

import (
	"testing"

	"github.com/gpdf-dev/gpdf/template"
)

func TestExample_31_Tmpl_15_PageSizes(t *testing.T) {
	sizes := []struct {
		name string
		size string
		file string
	}{
		{"A4 (210mm x 297mm)", "A4", "31_tmpl_15a_pagesize_a4.pdf"},
		{"A3 (297mm x 420mm)", "A3", "31_tmpl_15b_pagesize_a3.pdf"},
		{"Letter (8.5in x 11in)", "Letter", "31_tmpl_15c_pagesize_letter.pdf"},
		{"Legal (8.5in x 14in)", "Legal", "31_tmpl_15d_pagesize_legal.pdf"},
	}
	for _, s := range sizes {
		t.Run(s.name, func(t *testing.T) {
			schema := []byte(`{
				"page": {"size": "{{.PageSize}}", "margins": "20mm"},
				"body": [
					{"row": {"cols": [
						{"span": 12, "text": "Page Size: {{.SizeName}}", "style": {"size": 20, "bold": true}}
					]}},
					{"row": {"cols": [
						{"span": 12, "spacer": "10mm"}
					]}},
					{"row": {"cols": [
						{"span": 12, "text": "This page demonstrates the {{.SizeName}} page format."}
					]}}
				]
			}`)
			data := map[string]any{
				"PageSize": s.size,
				"SizeName": s.name,
			}
			doc, err := template.FromJSON(schema, data)
			if err != nil {
				t.Fatalf("FromJSON error: %v", err)
			}
			generatePDF(t, s.file, doc)
		})
	}
}
