package document

import (
	"testing"

	"github.com/gpdf-dev/gpdf/pdf"
)

func TestInheritStyle_FontFamily(t *testing.T) {
	parent := Style{FontFamily: "Helvetica"}
	child := Style{}
	got := InheritStyle(parent, child)
	if got.FontFamily != "Helvetica" {
		t.Errorf("FontFamily = %q, want %q", got.FontFamily, "Helvetica")
	}

	// Child override takes precedence.
	child2 := Style{FontFamily: "Courier"}
	got2 := InheritStyle(parent, child2)
	if got2.FontFamily != "Courier" {
		t.Errorf("FontFamily = %q, want %q", got2.FontFamily, "Courier")
	}
}

func TestInheritStyle_FontSize(t *testing.T) {
	parent := Style{FontSize: 16}
	child := Style{}
	got := InheritStyle(parent, child)
	if got.FontSize != 16 {
		t.Errorf("FontSize = %v, want 16", got.FontSize)
	}

	child2 := Style{FontSize: 10}
	got2 := InheritStyle(parent, child2)
	if got2.FontSize != 10 {
		t.Errorf("FontSize = %v, want 10", got2.FontSize)
	}
}

func TestInheritStyle_FontWeight(t *testing.T) {
	parent := Style{FontWeight: WeightBold}
	child := Style{}
	got := InheritStyle(parent, child)
	if got.FontWeight != WeightBold {
		t.Errorf("FontWeight = %v, want %v", got.FontWeight, WeightBold)
	}

	child2 := Style{FontWeight: WeightNormal}
	got2 := InheritStyle(parent, child2)
	if got2.FontWeight != WeightNormal {
		t.Errorf("FontWeight = %v, want %v", got2.FontWeight, WeightNormal)
	}
}

func TestInheritStyle_Color(t *testing.T) {
	red := pdf.RGB(1, 0, 0)
	parent := Style{Color: red}
	child := Style{}
	got := InheritStyle(parent, child)
	if got.Color != red {
		t.Errorf("Color = %v, want %v", got.Color, red)
	}

	blue := pdf.RGB(0, 0, 1)
	child2 := Style{Color: blue}
	got2 := InheritStyle(parent, child2)
	if got2.Color != blue {
		t.Errorf("Color = %v, want %v", got2.Color, blue)
	}
}

func TestInheritStyle_LineHeight(t *testing.T) {
	parent := Style{LineHeight: 1.5}
	child := Style{}
	got := InheritStyle(parent, child)
	if got.LineHeight != 1.5 {
		t.Errorf("LineHeight = %v, want 1.5", got.LineHeight)
	}
}

func TestInheritStyle_TextIndent(t *testing.T) {
	parent := Style{TextIndent: Pt(20)}
	child := Style{}
	got := InheritStyle(parent, child)
	if got.TextIndent != Pt(20) {
		t.Errorf("TextIndent = %v, want Pt(20)", got.TextIndent)
	}
}

func TestInheritStyle_NonInheritable(t *testing.T) {
	parent := Style{
		TextAlign:      AlignCenter,
		LetterSpacing:  2,
		WordSpacing:    3,
		TextDecoration: DecorationUnderline,
		VerticalAlign:  VAlignMiddle,
		Background:     &pdf.Color{R: 1, A: 1},
		Margin:         UniformEdges(Pt(10)),
		Padding:        UniformEdges(Pt(5)),
	}
	child := Style{}
	got := InheritStyle(parent, child)

	// Non-inheritable: child's zero values should remain.
	if got.TextAlign != AlignLeft {
		t.Errorf("TextAlign = %v, should not inherit", got.TextAlign)
	}
	if got.LetterSpacing != 0 {
		t.Errorf("LetterSpacing = %v, should not inherit", got.LetterSpacing)
	}
	if got.WordSpacing != 0 {
		t.Errorf("WordSpacing = %v, should not inherit", got.WordSpacing)
	}
	if got.TextDecoration != DecorationNone {
		t.Errorf("TextDecoration = %v, should not inherit", got.TextDecoration)
	}
	if got.VerticalAlign != VAlignTop {
		t.Errorf("VerticalAlign = %v, should not inherit", got.VerticalAlign)
	}
	if got.Background != nil {
		t.Errorf("Background = %v, should not inherit", got.Background)
	}
	if got.Margin != (Edges{}) {
		t.Errorf("Margin = %v, should not inherit", got.Margin)
	}
	if got.Padding != (Edges{}) {
		t.Errorf("Padding = %v, should not inherit", got.Padding)
	}
}
