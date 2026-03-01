package document

// RichTextFragment holds a single span of text with its own style.
// Multiple fragments are combined within a RichText node to produce
// inline mixed-style text on a single logical line.
type RichTextFragment struct {
	// Content is the text string for this fragment.
	Content string
	// FragmentStyle controls font, color, and decoration for this fragment.
	FragmentStyle Style
}

// RichText is a block-level node that arranges multiple styled text
// fragments inline, allowing mixed fonts, sizes, and colors within a
// single paragraph. The BlockStyle governs paragraph-level properties
// such as TextAlign, LineHeight, and TextIndent.
type RichText struct {
	// Fragments is the ordered list of inline text spans.
	Fragments []RichTextFragment
	// BlockStyle holds paragraph-level style (alignment, line height, indent).
	BlockStyle Style
	// BreakPolicy controls page-break behavior for this node.
	BreakPolicy BreakPolicy
}

// NodeType returns NodeRichText.
func (rt *RichText) NodeType() NodeType { return NodeRichText }

// Children returns nil because RichText is a leaf-like node whose
// inline fragments are handled internally by the layout engine.
func (rt *RichText) Children() []DocumentNode { return nil }

// Style returns the paragraph-level block style.
func (rt *RichText) Style() Style { return rt.BlockStyle }
