package document

// DocumentNode is the interface implemented by all nodes in the document
// tree. Each node carries a type tag, an optional list of children, and
// a style that governs its visual presentation.
type DocumentNode interface {
	// NodeType returns the kind of node (page, text, box, etc.).
	NodeType() NodeType
	// Children returns the direct child nodes, or nil for leaf nodes.
	Children() []DocumentNode
	// Style returns the visual style applied to this node.
	Style() Style
}

// NodeType enumerates the kinds of document nodes.
type NodeType int

const (
	// NodeDocument is the root node representing the entire document.
	NodeDocument NodeType = iota
	// NodePage represents a single page.
	NodePage
	// NodeBox is a generic container with CSS-like box model.
	NodeBox
	// NodeText is a leaf node containing text content.
	NodeText
	// NodeImage is a leaf node containing an image.
	NodeImage
	// NodeTable is a container for tabular data.
	NodeTable
	// NodeTableRow represents a single row within a table.
	NodeTableRow
	// NodeTableCell represents a single cell within a table row.
	NodeTableCell
	// NodeList is a container for ordered or unordered list items.
	NodeList
	// NodeListItem represents a single item in a list.
	NodeListItem
	// NodeRichText is an inline formatting context that places multiple
	// styled text fragments on a single line.
	NodeRichText
)

// BreakPolicy controls how page breaks interact with a node. It allows
// authors to force or suppress breaks before, after, or within a node.
type BreakPolicy struct {
	BreakBefore BreakValue
	BreakAfter  BreakValue
	BreakInside BreakValue
}

// BreakValue specifies the behavior of a page break control point.
type BreakValue int

const (
	// BreakAuto lets the layout engine decide whether to break.
	BreakAuto BreakValue = iota
	// BreakAvoid requests the engine to avoid a break at this point.
	BreakAvoid
	// BreakAlways forces a break at this point.
	BreakAlways
	// BreakPage forces a page break at this point.
	BreakPage
)
