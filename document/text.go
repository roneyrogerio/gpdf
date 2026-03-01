package document

// Text is a leaf document node that holds a string of text content
// and the style used to render it (font, size, color, alignment, etc.).
type Text struct {
	// Content is the text string to render.
	Content string
	// TextStyle controls font, color, alignment, and spacing.
	TextStyle Style
}

// NodeType returns NodeText.
func (t *Text) NodeType() NodeType { return NodeText }

// Children returns nil because text is a leaf node.
func (t *Text) Children() []DocumentNode { return nil }

// Style returns the text's visual style.
func (t *Text) Style() Style { return t.TextStyle }
