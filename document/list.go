package document

// ListType distinguishes ordered (numbered) from unordered (bulleted) lists.
type ListType int

const (
	// Unordered is a bulleted list.
	Unordered ListType = iota
	// Ordered is a numbered list.
	Ordered
)

// List is a document node representing an ordered or unordered list.
type List struct {
	// Items holds the list items.
	Items []ListItem
	// ListType selects bullet or numbered style.
	ListType ListType
	// ListStyle controls the list's visual properties.
	ListStyle Style
	// BreakPolicy controls page break behavior around and within this list.
	BreakPolicy BreakPolicy
	// MarkerIndent is the width reserved for the bullet or number marker
	// in points. If zero a default of 20pt is used.
	MarkerIndent float64
}

// NodeType returns NodeList.
func (l *List) NodeType() NodeType { return NodeList }

// Children returns a DocumentNode slice wrapping each ListItem.
func (l *List) Children() []DocumentNode {
	nodes := make([]DocumentNode, len(l.Items))
	for i := range l.Items {
		nodes[i] = &ListItemNode{Item: &l.Items[i], listStyle: l.ListStyle}
	}
	return nodes
}

// Style returns the list's style.
func (l *List) Style() Style { return l.ListStyle }

// ListItem represents a single item in a list.
type ListItem struct {
	// Content holds the child nodes rendered inside this item.
	Content []DocumentNode
	// ItemStyle controls the item's visual properties.
	ItemStyle Style
}

// ListItemNode wraps a ListItem to implement the DocumentNode interface.
type ListItemNode struct {
	Item      *ListItem
	listStyle Style
}

// NodeType returns NodeListItem.
func (n *ListItemNode) NodeType() NodeType { return NodeListItem }

// Children returns the item's content nodes.
func (n *ListItemNode) Children() []DocumentNode { return n.Item.Content }

// Style returns the item's style, falling back to the list style.
func (n *ListItemNode) Style() Style {
	s := n.listStyle
	if n.Item.ItemStyle.FontSize > 0 {
		s = n.Item.ItemStyle
	}
	return s
}
