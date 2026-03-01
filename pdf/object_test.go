package pdf_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/gpdf-dev/gpdf/pdf"
)

// ---------------------------------------------------------------------------
// Name
// ---------------------------------------------------------------------------

func TestName_WriteTo_Simple(t *testing.T) {
	var buf bytes.Buffer
	n, err := pdf.Name("Type").WriteTo(&buf)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := "/Type"
	if got := buf.String(); got != want {
		t.Errorf("got %q, want %q", got, want)
	}
	if int(n) != len(want) {
		t.Errorf("byte count = %d, want %d", n, len(want))
	}
}

func TestName_WriteTo_SpecialChars(t *testing.T) {
	tests := []struct {
		name string
		in   pdf.Name
		want string
	}{
		{"hash", pdf.Name("A#B"), "/A#23B"},
		{"slash", pdf.Name("A/B"), "/A#2FB"},
		{"parens", pdf.Name("A(B)"), "/A#28B#29"},
		{"angle brackets", pdf.Name("A<B>"), "/A#3CB#3E"},
		{"square brackets", pdf.Name("A[B]"), "/A#5BB#5D"},
		{"curly braces", pdf.Name("A{B}"), "/A#7BB#7D"},
		{"percent", pdf.Name("A%B"), "/A#25B"},
		{"space", pdf.Name("A B"), "/A#20B"},
		{"empty", pdf.Name(""), "/"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			_, err := tt.in.WriteTo(&buf)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got := buf.String(); got != tt.want {
				t.Errorf("got %q, want %q", got, tt.want)
			}
		})
	}
}

// ---------------------------------------------------------------------------
// LiteralString
// ---------------------------------------------------------------------------

func TestLiteralString_WriteTo(t *testing.T) {
	tests := []struct {
		name string
		in   pdf.LiteralString
		want string
	}{
		{"simple", pdf.LiteralString("Hello"), "(Hello)"},
		{"empty", pdf.LiteralString(""), "()"},
		{"escape_parens", pdf.LiteralString("A(B)C"), `(A\(B\)C)`},
		{"escape_backslash", pdf.LiteralString(`A\B`), `(A\\B)`},
		{"escape_cr", pdf.LiteralString("A\rB"), `(A\rB)`},
		{"escape_lf", pdf.LiteralString("A\nB"), `(A\nB)`},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			n, err := tt.in.WriteTo(&buf)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got := buf.String(); got != tt.want {
				t.Errorf("got %q, want %q", got, tt.want)
			}
			if int(n) != len(tt.want) {
				t.Errorf("byte count = %d, want %d", n, len(tt.want))
			}
		})
	}
}

// ---------------------------------------------------------------------------
// HexString
// ---------------------------------------------------------------------------

func TestHexString_WriteTo(t *testing.T) {
	tests := []struct {
		name string
		in   pdf.HexString
		want string
	}{
		{"hello", pdf.HexString("Hello"), "<48656C6C6F>"},
		{"empty", pdf.HexString(""), "<>"},
		{"binary", pdf.HexString("\x00\xFF"), "<00FF>"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			n, err := tt.in.WriteTo(&buf)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got := buf.String(); got != tt.want {
				t.Errorf("got %q, want %q", got, tt.want)
			}
			if int(n) != len(tt.want) {
				t.Errorf("byte count = %d, want %d", n, len(tt.want))
			}
		})
	}
}

// ---------------------------------------------------------------------------
// Integer
// ---------------------------------------------------------------------------

func TestInteger_WriteTo(t *testing.T) {
	tests := []struct {
		in   pdf.Integer
		want string
	}{
		{0, "0"},
		{42, "42"},
		{-7, "-7"},
		{1000000, "1000000"},
	}
	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			var buf bytes.Buffer
			n, err := tt.in.WriteTo(&buf)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got := buf.String(); got != tt.want {
				t.Errorf("got %q, want %q", got, tt.want)
			}
			if int(n) != len(tt.want) {
				t.Errorf("byte count = %d, want %d", n, len(tt.want))
			}
		})
	}
}

// ---------------------------------------------------------------------------
// Real
// ---------------------------------------------------------------------------

func TestReal_WriteTo(t *testing.T) {
	tests := []struct {
		name string
		in   pdf.Real
		want string
	}{
		{"zero", 0.0, "0"},
		{"positive", 3.14, "3.14"},
		{"negative", -2.5, "-2.5"},
		{"integer_value", 10.0, "10"},
		{"small", 0.001, "0.001"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			_, err := tt.in.WriteTo(&buf)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got := buf.String(); got != tt.want {
				t.Errorf("got %q, want %q", got, tt.want)
			}
		})
	}
}

// ---------------------------------------------------------------------------
// Boolean
// ---------------------------------------------------------------------------

func TestBoolean_WriteTo(t *testing.T) {
	tests := []struct {
		in   pdf.Boolean
		want string
	}{
		{true, "true"},
		{false, "false"},
	}
	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			var buf bytes.Buffer
			n, err := tt.in.WriteTo(&buf)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got := buf.String(); got != tt.want {
				t.Errorf("got %q, want %q", got, tt.want)
			}
			if int(n) != len(tt.want) {
				t.Errorf("byte count = %d, want %d", n, len(tt.want))
			}
		})
	}
}

// ---------------------------------------------------------------------------
// Null
// ---------------------------------------------------------------------------

func TestNull_WriteTo(t *testing.T) {
	var buf bytes.Buffer
	n, err := pdf.Null{}.WriteTo(&buf)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := "null"
	if got := buf.String(); got != want {
		t.Errorf("got %q, want %q", got, want)
	}
	if int(n) != len(want) {
		t.Errorf("byte count = %d, want %d", n, len(want))
	}
}

// ---------------------------------------------------------------------------
// Dict
// ---------------------------------------------------------------------------

func TestDict_WriteTo_Empty(t *testing.T) {
	var buf bytes.Buffer
	d := pdf.Dict{}
	_, err := d.WriteTo(&buf)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := "<< >>"
	if got := buf.String(); got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestDict_WriteTo_SingleEntry(t *testing.T) {
	var buf bytes.Buffer
	d := pdf.Dict{
		pdf.Name("Type"): pdf.Name("Catalog"),
	}
	_, err := d.WriteTo(&buf)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := "<< /Type /Catalog >>"
	if got := buf.String(); got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestDict_WriteTo_SortedKeys(t *testing.T) {
	var buf bytes.Buffer
	d := pdf.Dict{
		pdf.Name("Zebra"): pdf.Integer(1),
		pdf.Name("Apple"): pdf.Integer(2),
		pdf.Name("Mango"): pdf.Integer(3),
	}
	_, err := d.WriteTo(&buf)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got := buf.String()
	// Keys must be sorted: Apple, Mango, Zebra
	appleIdx := strings.Index(got, "/Apple")
	mangoIdx := strings.Index(got, "/Mango")
	zebraIdx := strings.Index(got, "/Zebra")
	if appleIdx > mangoIdx || mangoIdx > zebraIdx {
		t.Errorf("keys not sorted; got %q", got)
	}
}

func TestDict_WriteTo_NestedDict(t *testing.T) {
	var buf bytes.Buffer
	d := pdf.Dict{
		pdf.Name("Inner"): pdf.Dict{
			pdf.Name("Key"): pdf.LiteralString("Value"),
		},
	}
	_, err := d.WriteTo(&buf)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got := buf.String()
	if !strings.Contains(got, "<< /Key (Value) >>") {
		t.Errorf("nested dict not found in output: %q", got)
	}
}

// ---------------------------------------------------------------------------
// Array
// ---------------------------------------------------------------------------

func TestArray_WriteTo_Empty(t *testing.T) {
	var buf bytes.Buffer
	a := pdf.Array{}
	_, err := a.WriteTo(&buf)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := "[]"
	if got := buf.String(); got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestArray_WriteTo_Multiple(t *testing.T) {
	var buf bytes.Buffer
	a := pdf.Array{pdf.Integer(1), pdf.Integer(2), pdf.Integer(3)}
	_, err := a.WriteTo(&buf)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := "[1 2 3]"
	if got := buf.String(); got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestArray_WriteTo_MixedTypes(t *testing.T) {
	var buf bytes.Buffer
	a := pdf.Array{
		pdf.Name("Type"),
		pdf.Integer(42),
		pdf.LiteralString("hello"),
		pdf.Boolean(true),
		pdf.Null{},
	}
	_, err := a.WriteTo(&buf)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := "[/Type 42 (hello) true null]"
	if got := buf.String(); got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

// ---------------------------------------------------------------------------
// Stream
// ---------------------------------------------------------------------------

func TestStream_WriteTo(t *testing.T) {
	var buf bytes.Buffer
	s := pdf.Stream{
		Dict:    pdf.Dict{pdf.Name("Filter"): pdf.Name("FlateDecode")},
		Content: []byte("hello stream"),
	}
	_, err := s.WriteTo(&buf)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got := buf.String()

	// Must contain dict with Length
	if !strings.Contains(got, "/Length 12") {
		t.Errorf("missing /Length in output: %q", got)
	}
	// Must contain stream markers
	if !strings.Contains(got, "\nstream\n") {
		t.Errorf("missing stream keyword in output: %q", got)
	}
	if !strings.Contains(got, "\nendstream") {
		t.Errorf("missing endstream keyword in output: %q", got)
	}
	// Must contain the content
	if !strings.Contains(got, "hello stream") {
		t.Errorf("missing content in output: %q", got)
	}
	// Must contain filter
	if !strings.Contains(got, "/Filter") {
		t.Errorf("missing /Filter in output: %q", got)
	}
}

func TestStream_WriteTo_EmptyContent(t *testing.T) {
	var buf bytes.Buffer
	s := pdf.Stream{
		Dict:    pdf.Dict{},
		Content: []byte{},
	}
	_, err := s.WriteTo(&buf)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got := buf.String()
	if !strings.Contains(got, "/Length 0") {
		t.Errorf("missing /Length 0 in output: %q", got)
	}
}

func TestStream_WriteTo_OverridesLength(t *testing.T) {
	// If user sets Length in Dict, it should be overridden by the actual content length.
	var buf bytes.Buffer
	s := pdf.Stream{
		Dict:    pdf.Dict{pdf.Name("Length"): pdf.Integer(9999)},
		Content: []byte("abc"),
	}
	_, err := s.WriteTo(&buf)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got := buf.String()
	if !strings.Contains(got, "/Length 3") {
		t.Errorf("Length not overridden to 3: %q", got)
	}
	if strings.Contains(got, "9999") {
		t.Errorf("original Length 9999 should have been overridden: %q", got)
	}
}

// ---------------------------------------------------------------------------
// ObjectRef
// ---------------------------------------------------------------------------

func TestObjectRef_WriteTo(t *testing.T) {
	tests := []struct {
		ref  pdf.ObjectRef
		want string
	}{
		{pdf.ObjectRef{Number: 1, Generation: 0}, "1 0 R"},
		{pdf.ObjectRef{Number: 42, Generation: 3}, "42 3 R"},
		{pdf.ObjectRef{Number: 0, Generation: 0}, "0 0 R"},
	}
	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			var buf bytes.Buffer
			n, err := tt.ref.WriteTo(&buf)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got := buf.String(); got != tt.want {
				t.Errorf("got %q, want %q", got, tt.want)
			}
			if int(n) != len(tt.want) {
				t.Errorf("byte count = %d, want %d", n, len(tt.want))
			}
		})
	}
}

// ---------------------------------------------------------------------------
// Rectangle
// ---------------------------------------------------------------------------

func TestRectangle_WriteTo(t *testing.T) {
	var buf bytes.Buffer
	r := pdf.Rectangle{LLX: 0, LLY: 0, URX: 612, URY: 792}
	_, err := r.WriteTo(&buf)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := "[0 0 612 792]"
	if got := buf.String(); got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestRectangle_WriteTo_Fractional(t *testing.T) {
	var buf bytes.Buffer
	r := pdf.Rectangle{LLX: 0.5, LLY: 1.5, URX: 100.25, URY: 200.75}
	_, err := r.WriteTo(&buf)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := "[0.5 1.5 100.25 200.75]"
	if got := buf.String(); got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

// ---------------------------------------------------------------------------
// ResourceDict
// ---------------------------------------------------------------------------

func TestResourceDict_ToDict_Empty(t *testing.T) {
	rd := pdf.ResourceDict{}
	d := rd.ToDict()
	if len(d) != 0 {
		t.Errorf("expected empty dict, got %d entries", len(d))
	}
}

func TestResourceDict_ToDict_AllFields(t *testing.T) {
	rd := pdf.ResourceDict{
		Font:       pdf.Dict{pdf.Name("F1"): pdf.ObjectRef{Number: 1}},
		XObject:    pdf.Dict{pdf.Name("Im1"): pdf.ObjectRef{Number: 2}},
		ExtGState:  pdf.Dict{pdf.Name("GS1"): pdf.ObjectRef{Number: 3}},
		ColorSpace: pdf.Dict{pdf.Name("CS1"): pdf.ObjectRef{Number: 4}},
		Pattern:    pdf.Dict{pdf.Name("P1"): pdf.ObjectRef{Number: 5}},
	}
	d := rd.ToDict()
	if len(d) != 5 {
		t.Errorf("expected 5 entries, got %d", len(d))
	}
	for _, key := range []pdf.Name{"Font", "XObject", "ExtGState", "ColorSpace", "Pattern"} {
		if _, ok := d[key]; !ok {
			t.Errorf("missing key %q in dict", key)
		}
	}
}

func TestResourceDict_ToDict_PartialFields(t *testing.T) {
	rd := pdf.ResourceDict{
		Font: pdf.Dict{pdf.Name("F1"): pdf.ObjectRef{Number: 1}},
	}
	d := rd.ToDict()
	if len(d) != 1 {
		t.Errorf("expected 1 entry, got %d", len(d))
	}
	if _, ok := d[pdf.Name("Font")]; !ok {
		t.Errorf("missing key Font in dict")
	}
}

// ---------------------------------------------------------------------------
// DocumentInfo
// ---------------------------------------------------------------------------

func TestDocumentInfo_ToDict_Empty(t *testing.T) {
	di := pdf.DocumentInfo{}
	d := di.ToDict()
	if len(d) != 0 {
		t.Errorf("expected empty dict, got %d entries", len(d))
	}
}

func TestDocumentInfo_ToDict_AllFields(t *testing.T) {
	di := pdf.DocumentInfo{
		Title:    "My Document",
		Author:   "John Doe",
		Subject:  "Testing",
		Creator:  "TestApp",
		Producer: "gpdf",
	}
	d := di.ToDict()
	if len(d) != 5 {
		t.Errorf("expected 5 entries, got %d", len(d))
	}
	for _, key := range []pdf.Name{"Title", "Author", "Subject", "Creator", "Producer"} {
		if _, ok := d[key]; !ok {
			t.Errorf("missing key %q in dict", key)
		}
	}

	// Verify that the values are LiteralStrings with the correct content.
	var buf bytes.Buffer
	_, _ = d[pdf.Name("Title")].WriteTo(&buf)
	if got := buf.String(); got != "(My Document)" {
		t.Errorf("Title = %q, want %q", got, "(My Document)")
	}
}

func TestDocumentInfo_ToDict_PartialFields(t *testing.T) {
	di := pdf.DocumentInfo{
		Title:  "Title Only",
		Author: "Author Only",
	}
	d := di.ToDict()
	if len(d) != 2 {
		t.Errorf("expected 2 entries, got %d", len(d))
	}
	if _, ok := d[pdf.Name("Subject")]; ok {
		t.Errorf("Subject should not be present when empty")
	}
}

// ---------------------------------------------------------------------------
// Object interface compliance (compile-time checks)
// ---------------------------------------------------------------------------

func TestObjectInterfaceCompliance(t *testing.T) {
	// Verify that all types implement the Object interface at compile time.
	var _ pdf.Object = pdf.Name("")
	var _ pdf.Object = pdf.LiteralString("")
	var _ pdf.Object = pdf.HexString("")
	var _ pdf.Object = pdf.Integer(0)
	var _ pdf.Object = pdf.Real(0)
	var _ pdf.Object = pdf.Boolean(false)
	var _ pdf.Object = pdf.Null{}
	var _ pdf.Object = pdf.Dict{}
	var _ pdf.Object = pdf.Array{}
	var _ pdf.Object = pdf.Stream{}
	var _ pdf.Object = pdf.ObjectRef{}
	var _ pdf.Object = pdf.Rectangle{}
}
