// Package gpdf is a pure Go, zero-dependency PDF generation library.
//
// gpdf provides a layered architecture for PDF creation:
//
//   - pdf (Layer 1): Low-level PDF primitives — objects, streams, fonts, images
//   - document (Layer 2): Document model — nodes, box model, layout engine, renderer
//   - template (Layer 3): High-level declarative API — builders, grids, components
//
// Most users should start with the template package for the simplest API:
//
//	doc := gpdf.NewDocument(
//	    gpdf.WithPageSize(document.A4),
//	    gpdf.WithMargins(document.UniformEdges(document.Mm(15))),
//	)
//	page := doc.AddPage()
//	page.AutoRow(func(r *template.RowBuilder) {
//	    r.Col(12, func(c *template.ColBuilder) {
//	        c.Text("Hello, World!", template.FontSize(24))
//	    })
//	})
//	data, err := doc.Generate()
package gpdf

import (
	"github.com/gpdf-dev/gpdf/document"
	"github.com/gpdf-dev/gpdf/template"
)

// NewDocument creates a new PDF document builder with the given options.
// This is the primary entry point for creating PDFs using the high-level API.
func NewDocument(opts ...template.Option) *template.Document {
	return template.New(opts...)
}

// Re-export commonly used option functions for convenience.
var (
	// WithPageSize sets the page dimensions.
	WithPageSize = template.WithPageSize
	// WithMargins sets the page margins.
	WithMargins = template.WithMargins
	// WithFont registers a TrueType font.
	WithFont = template.WithFont
	// WithDefaultFont sets the default font family and size.
	WithDefaultFont = template.WithDefaultFont
	// WithMetadata sets document metadata (title, author, etc.).
	WithMetadata = template.WithMetadata
)

// Re-export commonly used page sizes for convenience.
var (
	A4     = document.A4
	A3     = document.A3
	Letter = document.Letter
	Legal  = document.Legal
)
