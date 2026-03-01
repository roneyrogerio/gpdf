package template

import (
	"testing"

	"github.com/gpdf-dev/gpdf/document"
	"github.com/gpdf-dev/gpdf/pdf"
)

func TestParseValue(t *testing.T) {
	tests := []struct {
		input string
		want  document.Value
	}{
		{"15mm", document.Mm(15)},
		{"20pt", document.Pt(20)},
		{"2.5cm", document.Cm(2.5)},
		{"1in", document.In(1)},
		{"1.5em", document.Em(1.5)},
		{"50%", document.Pct(50)},
		{"72", document.Pt(72)},
		{"auto", document.Auto},
		{"", document.Auto},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got, err := parseValue(tt.input)
			if err != nil {
				t.Fatalf("parseValue(%q) error: %v", tt.input, err)
			}
			if got != tt.want {
				t.Errorf("parseValue(%q) = %+v, want %+v", tt.input, got, tt.want)
			}
		})
	}
}

func TestParseValueError(t *testing.T) {
	bad := []string{"abc", "12xyz", "mm15"}
	for _, s := range bad {
		if _, err := parseValue(s); err == nil {
			t.Errorf("parseValue(%q) expected error, got nil", s)
		}
	}
}

func TestParsePageSize(t *testing.T) {
	tests := []struct {
		input string
		want  document.Size
	}{
		{"A4", document.A4},
		{"a4", document.A4},
		{"A3", document.A3},
		{"Letter", document.Letter},
		{"legal", document.Legal},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got, err := parsePageSize(tt.input)
			if err != nil {
				t.Fatalf("parsePageSize(%q) error: %v", tt.input, err)
			}
			if got != tt.want {
				t.Errorf("parsePageSize(%q) = %+v, want %+v", tt.input, got, tt.want)
			}
		})
	}

	if _, err := parsePageSize("Unknown"); err == nil {
		t.Error("parsePageSize(Unknown) expected error")
	}
}

func TestParseColor(t *testing.T) {
	tests := []struct {
		input string
		want  pdf.Color
	}{
		{"black", pdf.Black},
		{"white", pdf.White},
		{"red", pdf.Red},
		{"#FF0000", pdf.RGBHex(0xFF0000)},
		{"#1A237E", pdf.RGBHex(0x1A237E)},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got, err := parseColor(tt.input)
			if err != nil {
				t.Fatalf("parseColor(%q) error: %v", tt.input, err)
			}
			if got != tt.want {
				t.Errorf("parseColor(%q) = %+v, want %+v", tt.input, got, tt.want)
			}
		})
	}

	badColors := []string{"#GG0000", "unknown", "#FFF"}
	for _, s := range badColors {
		if _, err := parseColor(s); err == nil {
			t.Errorf("parseColor(%q) expected error, got nil", s)
		}
	}
}

func TestApplySchemaStyle(t *testing.T) {
	// nil style should return nil.
	if opts := applySchemaStyle(nil); opts != nil {
		t.Error("applySchemaStyle(nil) should return nil")
	}

	ss := &SchemaStyle{
		Size:          18,
		Bold:          true,
		Italic:        true,
		Align:         "center",
		Color:         "red",
		Background:    "#F5F5F5",
		Underline:     true,
		Strikethrough: true,
		LetterSpacing: 1.5,
	}
	opts := applySchemaStyle(ss)
	if len(opts) == 0 {
		t.Fatal("expected non-empty options from full SchemaStyle")
	}

	// Apply to a default style and verify some fields.
	style := document.DefaultStyle()
	for _, opt := range opts {
		opt(&style)
	}
	if style.FontSize != 18 {
		t.Errorf("FontSize = %v, want 18", style.FontSize)
	}
	if style.FontWeight != document.WeightBold {
		t.Error("expected bold")
	}
	if style.FontStyle != document.StyleItalic {
		t.Error("expected italic")
	}
	if style.TextAlign != document.AlignCenter {
		t.Error("expected center align")
	}
	if style.TextDecoration&document.DecorationUnderline == 0 {
		t.Error("expected underline")
	}
	if style.TextDecoration&document.DecorationStrikethrough == 0 {
		t.Error("expected strikethrough")
	}
	if style.LetterSpacing != 1.5 {
		t.Errorf("LetterSpacing = %v, want 1.5", style.LetterSpacing)
	}
}

func TestBuildFromSchema_Basic(t *testing.T) {
	schema := &Schema{
		Page: SchemaPage{
			Size:    "A4",
			Margins: "20mm",
		},
		Metadata: &SchemaMeta{
			Title:  "Test Document",
			Author: "gpdf",
		},
		Body: []SchemaRow{
			{Row: SchemaRowDef{
				Cols: []SchemaCol{
					{Span: 12, Text: "Hello, World!", Style: &SchemaStyle{Size: 24, Bold: true}},
				},
			}},
		},
	}

	doc, err := buildFromSchema(schema, nil)
	if err != nil {
		t.Fatalf("buildFromSchema error: %v", err)
	}

	data, err := doc.Generate()
	if err != nil {
		t.Fatalf("Generate error: %v", err)
	}
	if len(data) == 0 {
		t.Fatal("generated PDF is empty")
	}
	if string(data[:5]) != "%PDF-" {
		t.Fatalf("invalid PDF header: %q", string(data[:5]))
	}
}

func TestBuildFromSchema_FixedHeightRow(t *testing.T) {
	schema := &Schema{
		Page: SchemaPage{Size: "Letter"},
		Body: []SchemaRow{
			{Row: SchemaRowDef{
				Height: "25mm",
				Cols: []SchemaCol{
					{Span: 6, Text: "Left"},
					{Span: 6, Text: "Right"},
				},
			}},
		},
	}

	doc, err := buildFromSchema(schema, nil)
	if err != nil {
		t.Fatalf("buildFromSchema error: %v", err)
	}
	if _, err := doc.Generate(); err != nil {
		t.Fatalf("Generate error: %v", err)
	}
}

func TestBuildFromSchema_HeaderFooter(t *testing.T) {
	schema := &Schema{
		Page: SchemaPage{Size: "A4", Margins: "15mm"},
		Header: []SchemaRow{
			{Row: SchemaRowDef{
				Cols: []SchemaCol{
					{Span: 12, Text: "Page Header", Style: &SchemaStyle{Bold: true}},
				},
			}},
		},
		Footer: []SchemaRow{
			{Row: SchemaRowDef{
				Cols: []SchemaCol{
					{Span: 12, Elements: []SchemaElement{
						{Type: "pageNumber", Style: &SchemaStyle{Align: "center"}},
					}},
				},
			}},
		},
		Body: []SchemaRow{
			{Row: SchemaRowDef{
				Cols: []SchemaCol{
					{Span: 12, Text: "Body content"},
				},
			}},
		},
	}

	doc, err := buildFromSchema(schema, nil)
	if err != nil {
		t.Fatalf("buildFromSchema error: %v", err)
	}
	if _, err := doc.Generate(); err != nil {
		t.Fatalf("Generate error: %v", err)
	}
}

func TestBuildFromSchema_Elements(t *testing.T) {
	schema := &Schema{
		Page: SchemaPage{Size: "A4"},
		Body: []SchemaRow{
			{Row: SchemaRowDef{
				Cols: []SchemaCol{
					{Span: 12, Elements: []SchemaElement{
						{Type: "text", Content: "Title", Style: &SchemaStyle{Size: 24, Bold: true}},
						{Type: "spacer", Height: "10mm"},
						{Type: "line", Line: &SchemaLine{Color: "#CCCCCC", Thickness: "2pt"}},
						{Type: "spacer", Height: "5mm"},
						{Type: "text", Content: "Body text"},
						{Type: "list", List: &SchemaList{Items: []string{"Item 1", "Item 2", "Item 3"}}},
						{Type: "list", List: &SchemaList{Type: "ordered", Items: []string{"First", "Second"}}},
					}},
				},
			}},
		},
	}

	doc, err := buildFromSchema(schema, nil)
	if err != nil {
		t.Fatalf("buildFromSchema error: %v", err)
	}
	if _, err := doc.Generate(); err != nil {
		t.Fatalf("Generate error: %v", err)
	}
}

func TestBuildFromSchema_Table(t *testing.T) {
	schema := &Schema{
		Page: SchemaPage{Size: "A4"},
		Body: []SchemaRow{
			{Row: SchemaRowDef{
				Cols: []SchemaCol{
					{Span: 12, Table: &SchemaTable{
						Header:       []string{"Name", "Qty", "Price"},
						Rows:         [][]string{{"Widget", "10", "$5.00"}, {"Gadget", "3", "$15.00"}},
						ColumnWidths: []float64{50, 25, 25},
						HeaderStyle:  &SchemaStyle{Bold: true, Color: "white", Background: "#1A237E"},
						StripeColor:  "#F5F5F5",
					}},
				},
			}},
		},
	}

	doc, err := buildFromSchema(schema, nil)
	if err != nil {
		t.Fatalf("buildFromSchema error: %v", err)
	}
	if _, err := doc.Generate(); err != nil {
		t.Fatalf("Generate error: %v", err)
	}
}

func TestBuildFromSchema_QRCodeBarcode(t *testing.T) {
	schema := &Schema{
		Page: SchemaPage{Size: "A4"},
		Body: []SchemaRow{
			{Row: SchemaRowDef{
				Cols: []SchemaCol{
					{Span: 6, QRCode: &SchemaQRCode{Data: "https://gpdf.dev", Size: "30mm"}},
					{Span: 6, Barcode: &SchemaBarcode{Data: "ABC-123", Width: "50mm", Height: "20mm", Format: "code128"}},
				},
			}},
		},
	}

	doc, err := buildFromSchema(schema, nil)
	if err != nil {
		t.Fatalf("buildFromSchema error: %v", err)
	}
	if _, err := doc.Generate(); err != nil {
		t.Fatalf("Generate error: %v", err)
	}
}

func TestBuildFromSchema_Line(t *testing.T) {
	schema := &Schema{
		Page: SchemaPage{Size: "A4"},
		Body: []SchemaRow{
			{Row: SchemaRowDef{
				Cols: []SchemaCol{
					{Span: 12, Line: &SchemaLine{Color: "red", Thickness: "2pt"}},
				},
			}},
			{Row: SchemaRowDef{
				Cols: []SchemaCol{
					{Span: 12, Elements: []SchemaElement{
						{Type: "line"}, // default line
					}},
				},
			}},
		},
	}

	doc, err := buildFromSchema(schema, nil)
	if err != nil {
		t.Fatalf("buildFromSchema error: %v", err)
	}
	if _, err := doc.Generate(); err != nil {
		t.Fatalf("Generate error: %v", err)
	}
}

func TestBuildFromSchema_WithOptions(t *testing.T) {
	schema := &Schema{
		Page: SchemaPage{Size: "A4"},
		Body: []SchemaRow{
			{Row: SchemaRowDef{
				Cols: []SchemaCol{
					{Span: 12, Text: "With custom metadata"},
				},
			}},
		},
	}

	// Options passed to buildFromSchema override schema settings.
	doc, err := buildFromSchema(schema, []Option{
		WithMetadata(document.DocumentMetadata{Title: "Override Title"}),
	})
	if err != nil {
		t.Fatalf("buildFromSchema error: %v", err)
	}
	if _, err := doc.Generate(); err != nil {
		t.Fatalf("Generate error: %v", err)
	}
}

func TestBuildFromSchema_InvalidPageSize(t *testing.T) {
	schema := &Schema{
		Page: SchemaPage{Size: "B5"},
		Body: []SchemaRow{},
	}
	if _, err := buildFromSchema(schema, nil); err == nil {
		t.Error("expected error for unknown page size")
	}
}

func TestBuildFromSchema_InvalidMargins(t *testing.T) {
	schema := &Schema{
		Page: SchemaPage{Size: "A4", Margins: "abc"},
		Body: []SchemaRow{},
	}
	if _, err := buildFromSchema(schema, nil); err == nil {
		t.Error("expected error for invalid margins")
	}
}
