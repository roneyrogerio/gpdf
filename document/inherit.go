package document

import "github.com/gpdf-dev/gpdf/pdf"

// InheritStyle produces a merged style where inheritable CSS properties
// are copied from parent to child when the child's value is unset (zero
// value). Non-inheritable properties (TextAlign, FontStyle, LetterSpacing,
// WordSpacing, TextDecoration, VerticalAlign, Background, Margin, Padding,
// Border) are always taken from the child unchanged.
//
// Limitation: properties whose zero value is a valid setting (FontStyle,
// TextAlign, LetterSpacing, WordSpacing) cannot be distinguished from
// "unset" and are therefore treated as non-inheritable here. Callers
// that need explicit reset of these properties should set them on the
// child style directly.
func InheritStyle(parent, child Style) Style {
	result := child

	// FontFamily — inheritable, unset = ""
	if result.FontFamily == "" {
		result.FontFamily = parent.FontFamily
	}

	// FontSize — inheritable, unset = 0
	if result.FontSize == 0 {
		result.FontSize = parent.FontSize
	}

	// FontWeight — inheritable, unset = 0 (valid weights start at 400)
	if result.FontWeight == 0 {
		result.FontWeight = parent.FontWeight
	}

	// Color — inheritable, unset = pdf.Color{} (all constructors set A=1.0)
	if result.Color == (pdf.Color{}) {
		result.Color = parent.Color
	}

	// LineHeight — inheritable, unset = 0
	if result.LineHeight == 0 {
		result.LineHeight = parent.LineHeight
	}

	// TextIndent — inheritable, unset = Value{}
	if result.TextIndent == (Value{}) {
		result.TextIndent = parent.TextIndent
	}

	return result
}
