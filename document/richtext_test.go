package document

import "testing"

func TestRichTextNodeType(t *testing.T) {
	rt := &RichText{}
	if rt.NodeType() != NodeRichText {
		t.Errorf("NodeType() = %v, want NodeRichText", rt.NodeType())
	}
}

func TestRichTextChildren(t *testing.T) {
	rt := &RichText{
		Fragments: []RichTextFragment{
			{Content: "hello"},
		},
	}
	if rt.Children() != nil {
		t.Errorf("Children() should return nil, got %v", rt.Children())
	}
}

func TestRichTextStyle(t *testing.T) {
	style := Style{FontSize: 14, TextAlign: AlignCenter}
	rt := &RichText{BlockStyle: style}
	if rt.Style() != style {
		t.Errorf("Style() = %v, want %v", rt.Style(), style)
	}
}
