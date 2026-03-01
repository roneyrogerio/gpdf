package document

// PathOp identifies the type of a path drawing operation.
type PathOp int

const (
	// PathMoveTo moves the current point without drawing.
	PathMoveTo PathOp = iota
	// PathLineTo draws a straight line from the current point.
	PathLineTo
	// PathCurveTo draws a cubic Bezier curve using two control points.
	PathCurveTo
	// PathClose closes the current subpath with a straight line back
	// to its starting point.
	PathClose
)

// PathSegment is a single drawing instruction within a path.
type PathSegment struct {
	Op PathOp
	// Points contains the operand points. The number of points depends
	// on the operation:
	//   MoveTo:  1 (destination)
	//   LineTo:  1 (endpoint)
	//   CurveTo: 3 (control1, control2, endpoint)
	//   Close:   0
	Points []Point
}

// Path is an ordered sequence of path segments that together describe
// a shape suitable for filling, stroking, or both.
type Path struct {
	Segments []PathSegment
}
