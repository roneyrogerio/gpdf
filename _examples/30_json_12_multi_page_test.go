package examples_test

import (
	"testing"

	"github.com/gpdf-dev/gpdf/template"
)

func TestExample_30_JSON_12_MultiPage(t *testing.T) {
	schema := []byte(`{
		"page": {"size": "A4", "margins": "20mm"},
		"body": [
			{"row": {"cols": [
				{"span": 12, "text": "Multi-Page Document", "style": {"size": 20, "bold": true}}
			]}},
			{"row": {"cols": [
				{"span": 12, "spacer": "5mm"}
			]}},
			{"row": {"cols": [
				{"span": 12, "line": {"color": "#000000", "thickness": "1pt"}}
			]}},
			{"row": {"cols": [
				{"span": 12, "spacer": "10mm"}
			]}},
			{"row": {"cols": [
				{"span": 12, "text": "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Sed do eiusmod tempor incididunt ut labore et dolore magna aliqua."}
			]}},
			{"row": {"cols": [
				{"span": 12, "text": "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Sed do eiusmod tempor incididunt ut labore et dolore magna aliqua."}
			]}},
			{"row": {"cols": [
				{"span": 12, "text": "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Sed do eiusmod tempor incididunt ut labore et dolore magna aliqua."}
			]}},
			{"row": {"cols": [
				{"span": 12, "text": "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Sed do eiusmod tempor incididunt ut labore et dolore magna aliqua."}
			]}},
			{"row": {"cols": [
				{"span": 12, "text": "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Sed do eiusmod tempor incididunt ut labore et dolore magna aliqua."}
			]}},
			{"row": {"cols": [
				{"span": 12, "text": "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Sed do eiusmod tempor incididunt ut labore et dolore magna aliqua."}
			]}},
			{"row": {"cols": [
				{"span": 12, "text": "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Sed do eiusmod tempor incididunt ut labore et dolore magna aliqua."}
			]}},
			{"row": {"cols": [
				{"span": 12, "text": "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Sed do eiusmod tempor incididunt ut labore et dolore magna aliqua."}
			]}},
			{"row": {"cols": [
				{"span": 12, "text": "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Sed do eiusmod tempor incididunt ut labore et dolore magna aliqua."}
			]}},
			{"row": {"cols": [
				{"span": 12, "text": "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Sed do eiusmod tempor incididunt ut labore et dolore magna aliqua."}
			]}},
			{"row": {"cols": [
				{"span": 12, "text": "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Sed do eiusmod tempor incididunt ut labore et dolore magna aliqua."}
			]}},
			{"row": {"cols": [
				{"span": 12, "text": "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Sed do eiusmod tempor incididunt ut labore et dolore magna aliqua."}
			]}},
			{"row": {"cols": [
				{"span": 12, "text": "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Sed do eiusmod tempor incididunt ut labore et dolore magna aliqua."}
			]}},
			{"row": {"cols": [
				{"span": 12, "text": "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Sed do eiusmod tempor incididunt ut labore et dolore magna aliqua."}
			]}},
			{"row": {"cols": [
				{"span": 12, "text": "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Sed do eiusmod tempor incididunt ut labore et dolore magna aliqua."}
			]}},
			{"row": {"cols": [
				{"span": 12, "text": "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Sed do eiusmod tempor incididunt ut labore et dolore magna aliqua."}
			]}},
			{"row": {"cols": [
				{"span": 12, "text": "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Sed do eiusmod tempor incididunt ut labore et dolore magna aliqua."}
			]}},
			{"row": {"cols": [
				{"span": 12, "text": "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Sed do eiusmod tempor incididunt ut labore et dolore magna aliqua."}
			]}},
			{"row": {"cols": [
				{"span": 12, "text": "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Sed do eiusmod tempor incididunt ut labore et dolore magna aliqua."}
			]}},
			{"row": {"cols": [
				{"span": 12, "text": "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Sed do eiusmod tempor incididunt ut labore et dolore magna aliqua."}
			]}}
		]
	}`)

	doc, err := template.FromJSON(schema, nil)
	if err != nil {
		t.Fatalf("FromJSON error: %v", err)
	}
	generatePDF(t, "30_json_12_multi_page.pdf", doc)
}
