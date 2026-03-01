package template

import "github.com/gpdf-dev/gpdf/document"

// RichTextBuilder accumulates inline text fragments with individual
// styles. It is used inside ColBuilder.RichText to construct a
// document.RichText node.
type RichTextBuilder struct {
	defaultStyle document.Style
	fragments    []document.RichTextFragment
}

// Span appends a text fragment. Options modify a copy of the default
// style for this fragment only.
func (b *RichTextBuilder) Span(text string, opts ...TextOption) {
	style := b.defaultStyle
	for _, opt := range opts {
		opt(&style)
	}
	b.fragments = append(b.fragments, document.RichTextFragment{
		Content:       text,
		FragmentStyle: style,
	})
}
