package template

import (
	"testing"
	gotemplate "text/template"
)

func TestFromJSON_Static(t *testing.T) {
	schema := []byte(`{
		"page": {"size": "A4", "margins": "20mm"},
		"metadata": {"title": "Static Doc", "author": "test"},
		"body": [
			{"row": {"cols": [
				{"span": 12, "text": "Hello, World!", "style": {"size": 24, "bold": true}}
			]}}
		]
	}`)

	doc, err := FromJSON(schema, nil)
	if err != nil {
		t.Fatalf("FromJSON error: %v", err)
	}

	data, err := doc.Generate()
	if err != nil {
		t.Fatalf("Generate error: %v", err)
	}
	assertValidPDFBytes(t, data)
}

func TestFromJSON_WithTemplateData(t *testing.T) {
	schema := []byte(`{
		"page": {"size": "A4", "margins": "15mm"},
		"metadata": {"title": "{{.Title}}"},
		"body": [
			{"row": {"cols": [
				{"span": 12, "text": "{{.Greeting}}", "style": {"size": 24, "bold": true}}
			]}},
			{"row": {"cols": [
				{"span": 12, "text": "Author: {{.Author}}"}
			]}}
		]
	}`)

	data := map[string]string{
		"Title":    "Dynamic Document",
		"Greeting": "Hello, gpdf!",
		"Author":   "Test User",
	}

	doc, err := FromJSON(schema, data)
	if err != nil {
		t.Fatalf("FromJSON error: %v", err)
	}

	pdfBytes, err := doc.Generate()
	if err != nil {
		t.Fatalf("Generate error: %v", err)
	}
	assertValidPDFBytes(t, pdfBytes)
}

func TestFromJSON_WithTableData(t *testing.T) {
	schema := []byte(`{
		"page": {"size": "A4"},
		"body": [
			{"row": {"cols": [
				{"span": 12, "text": "{{.Title}}", "style": {"size": 18, "bold": true}}
			]}},
			{"row": {"cols": [
				{"span": 12, "table": {
					"header": ["Item", "Qty", "Price"],
					"rows": {{toJSON .Items}},
					"columnWidths": [50, 25, 25],
					"headerStyle": {"bold": true, "color": "white", "background": "#1A237E"},
					"stripeColor": "#F5F5F5"
				}}
			]}}
		]
	}`)

	data := map[string]any{
		"Title": "Invoice",
		"Items": [][]string{
			{"Widget", "10", "$50.00"},
			{"Gadget", "5", "$25.00"},
			{"Doohickey", "3", "$15.00"},
		},
	}

	doc, err := FromJSON(schema, data)
	if err != nil {
		t.Fatalf("FromJSON error: %v", err)
	}

	pdfBytes, err := doc.Generate()
	if err != nil {
		t.Fatalf("Generate error: %v", err)
	}
	assertValidPDFBytes(t, pdfBytes)
}

func TestFromJSON_WithListData(t *testing.T) {
	schema := []byte(`{
		"page": {"size": "A4"},
		"body": [
			{"row": {"cols": [
				{"span": 12, "elements": [
					{"type": "text", "content": "Features:", "style": {"bold": true}},
					{"type": "list", "list": {"items": {{toJSON .Features}}}}
				]}
			]}}
		]
	}`)

	data := map[string]any{
		"Features": []string{"Zero dependencies", "CJK support", "10-30x faster"},
	}

	doc, err := FromJSON(schema, data)
	if err != nil {
		t.Fatalf("FromJSON error: %v", err)
	}

	pdfBytes, err := doc.Generate()
	if err != nil {
		t.Fatalf("Generate error: %v", err)
	}
	assertValidPDFBytes(t, pdfBytes)
}

func TestFromJSON_TemplateError(t *testing.T) {
	// Invalid template syntax.
	schema := []byte(`{
		"page": {"size": "A4"},
		"body": [{"row": {"cols": [{"span": 12, "text": "{{.Oops"}]}}]
	}`)

	if _, err := FromJSON(schema, map[string]string{}); err == nil {
		t.Error("expected template parse error")
	}
}

func TestFromJSON_TemplateExecutionError(t *testing.T) {
	// call is a function that will fail at execution time.
	schema := []byte(`{
		"page": {"size": "A4"},
		"body": [{"row": {"cols": [{"span": 12, "text": "{{call .Func}}"}]}}]
	}`)

	// Passing a non-callable value for .Func triggers an execution error.
	if _, err := FromJSON(schema, map[string]string{"Func": "notfunc"}); err == nil {
		t.Error("expected template execution error")
	}
}

func TestFromJSON_InvalidJSON(t *testing.T) {
	schema := []byte(`not json at all`)
	if _, err := FromJSON(schema, nil); err == nil {
		t.Error("expected JSON parse error")
	}
}

func TestFromTemplate_Basic(t *testing.T) {
	tmplStr := `{
		"page": {"size": "{{.PageSize}}", "margins": "{{.Margins}}"},
		"body": [
			{"row": {"cols": [
				{"span": 12, "text": "{{.Title}}", "style": {"size": 24, "bold": true}}
			]}}
		]
	}`

	tmpl, err := gotemplate.New("test").Funcs(TemplateFuncMap()).Parse(tmplStr)
	if err != nil {
		t.Fatalf("template parse error: %v", err)
	}

	data := map[string]string{
		"PageSize": "A4",
		"Margins":  "20mm",
		"Title":    "From Template",
	}

	doc, err := FromTemplate(tmpl, data)
	if err != nil {
		t.Fatalf("FromTemplate error: %v", err)
	}

	pdfBytes, err := doc.Generate()
	if err != nil {
		t.Fatalf("Generate error: %v", err)
	}
	assertValidPDFBytes(t, pdfBytes)
}

func TestFromTemplate_WithToJSON(t *testing.T) {
	tmplStr := `{
		"page": {"size": "A4"},
		"body": [
			{"row": {"cols": [
				{"span": 12, "table": {
					"header": {{toJSON .Headers}},
					"rows": {{toJSON .Rows}}
				}}
			]}}
		]
	}`

	tmpl, err := gotemplate.New("test").Funcs(TemplateFuncMap()).Parse(tmplStr)
	if err != nil {
		t.Fatalf("template parse error: %v", err)
	}

	data := map[string]any{
		"Headers": []string{"A", "B", "C"},
		"Rows":    [][]string{{"1", "2", "3"}, {"4", "5", "6"}},
	}

	doc, err := FromTemplate(tmpl, data)
	if err != nil {
		t.Fatalf("FromTemplate error: %v", err)
	}

	pdfBytes, err := doc.Generate()
	if err != nil {
		t.Fatalf("Generate error: %v", err)
	}
	assertValidPDFBytes(t, pdfBytes)
}

func TestFromTemplate_ConditionalRows(t *testing.T) {
	tmplStr := `{
		"page": {"size": "A4"},
		"body": [
			{"row": {"cols": [
				{"span": 12, "text": "Title", "style": {"size": 24, "bold": true}}
			]}}
			{{- if .ShowSubtitle}},
			{"row": {"cols": [
				{"span": 12, "text": "{{.Subtitle}}"}
			]}}
			{{- end}}
		]
	}`

	tmpl, err := gotemplate.New("test").Funcs(TemplateFuncMap()).Parse(tmplStr)
	if err != nil {
		t.Fatalf("template parse error: %v", err)
	}

	// With subtitle.
	data := map[string]any{
		"ShowSubtitle": true,
		"Subtitle":     "A dynamic subtitle",
	}
	doc, err := FromTemplate(tmpl, data)
	if err != nil {
		t.Fatalf("FromTemplate error: %v", err)
	}
	pdfBytes, err := doc.Generate()
	if err != nil {
		t.Fatalf("Generate error: %v", err)
	}
	assertValidPDFBytes(t, pdfBytes)

	// Without subtitle.
	data["ShowSubtitle"] = false
	doc, err = FromTemplate(tmpl, data)
	if err != nil {
		t.Fatalf("FromTemplate (no subtitle) error: %v", err)
	}
	pdfBytes, err = doc.Generate()
	if err != nil {
		t.Fatalf("Generate (no subtitle) error: %v", err)
	}
	assertValidPDFBytes(t, pdfBytes)
}

func TestFromTemplate_LoopRows(t *testing.T) {
	tmplStr := `{
		"page": {"size": "A4"},
		"body": [
			{{- range $i, $sec := .Sections}}
			{{- if $i}},{{end}}
			{"row": {"cols": [
				{"span": 12, "text": "{{$sec}}", "style": {"size": 16, "bold": true}}
			]}}
			{{- end}}
		]
	}`

	tmpl, err := gotemplate.New("test").Funcs(TemplateFuncMap()).Parse(tmplStr)
	if err != nil {
		t.Fatalf("template parse error: %v", err)
	}

	data := map[string]any{
		"Sections": []string{"Introduction", "Methods", "Results", "Conclusion"},
	}

	doc, err := FromTemplate(tmpl, data)
	if err != nil {
		t.Fatalf("FromTemplate error: %v", err)
	}
	pdfBytes, err := doc.Generate()
	if err != nil {
		t.Fatalf("Generate error: %v", err)
	}
	assertValidPDFBytes(t, pdfBytes)
}

func TestFromTemplate_ExecutionError(t *testing.T) {
	tmpl, _ := gotemplate.New("test").Parse(`{"page": {"size": "A4"}, "body": []}`)
	// Passing a channel which cannot be accessed by template.
	type badData struct {
		Ch chan int
	}
	// This should not error since the template doesn't reference Ch.
	doc, err := FromTemplate(tmpl, badData{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, err := doc.Generate(); err != nil {
		t.Fatalf("Generate error: %v", err)
	}
}

func TestTemplateFuncMap_ToJSON(t *testing.T) {
	fm := TemplateFuncMap()
	fn, ok := fm["toJSON"]
	if !ok {
		t.Fatal("toJSON not found in FuncMap")
	}

	toJSON, ok := fn.(func(any) (string, error))
	if !ok {
		t.Fatal("toJSON has unexpected signature")
	}

	result, err := toJSON([]string{"a", "b"})
	if err != nil {
		t.Fatalf("toJSON error: %v", err)
	}
	if result != `["a","b"]` {
		t.Errorf("toJSON = %q, want %q", result, `["a","b"]`)
	}
}

func TestFromJSON_ComplexInvoice(t *testing.T) {
	schema := []byte(`{
		"page": {"size": "A4", "margins": "20mm"},
		"metadata": {"title": "Invoice #001", "author": "ACME Corp"},
		"header": [
			{"row": {"cols": [
				{"span": 6, "text": "ACME Corporation", "style": {"size": 20, "bold": true, "color": "#1A237E"}},
				{"span": 6, "text": "INVOICE", "style": {"size": 24, "bold": true, "align": "right", "color": "#1A237E"}}
			]}},
			{"row": {"cols": [
				{"span": 12, "line": {"color": "#1A237E", "thickness": "2pt"}}
			]}}
		],
		"footer": [
			{"row": {"cols": [
				{"span": 12, "elements": [
					{"type": "pageNumber", "style": {"align": "center"}}
				]}
			]}}
		],
		"body": [
			{"row": {"height": "25mm", "cols": [
				{"span": 6, "elements": [
					{"type": "text", "content": "Bill To:", "style": {"bold": true}},
					{"type": "text", "content": "John Smith"},
					{"type": "text", "content": "123 Main St"}
				]},
				{"span": 6, "elements": [
					{"type": "text", "content": "Date: 2026-03-01", "style": {"align": "right"}},
					{"type": "text", "content": "#INV-001", "style": {"align": "right"}}
				]}
			]}},
			{"row": {"cols": [
				{"span": 12, "spacer": "10mm"}
			]}},
			{"row": {"cols": [
				{"span": 12, "table": {
					"header": ["Description", "Qty", "Price", "Amount"],
					"rows": [
						["Web Development", "40", "$150", "$6,000"],
						["UI Design", "20", "$120", "$2,400"]
					],
					"columnWidths": [40, 15, 20, 25],
					"headerStyle": {"bold": true, "color": "white", "background": "#1A237E"},
					"stripeColor": "#F5F5F5"
				}}
			]}},
			{"row": {"cols": [
				{"span": 8, "text": ""},
				{"span": 4, "elements": [
					{"type": "text", "content": "Total: $8,400.00", "style": {"bold": true, "size": 14, "align": "right"}}
				]}
			]}}
		]
	}`)

	doc, err := FromJSON(schema, nil)
	if err != nil {
		t.Fatalf("FromJSON error: %v", err)
	}

	pdfBytes, err := doc.Generate()
	if err != nil {
		t.Fatalf("Generate error: %v", err)
	}
	assertValidPDFBytes(t, pdfBytes)
}

// assertValidPDFBytes checks that data starts with the PDF magic header.
func assertValidPDFBytes(t *testing.T, data []byte) {
	t.Helper()
	if len(data) == 0 {
		t.Fatal("generated PDF is empty")
	}
	if string(data[:5]) != "%PDF-" {
		t.Fatalf("invalid PDF header: %q", string(data[:5]))
	}
}
