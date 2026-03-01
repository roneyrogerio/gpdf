package main

import (
	"fmt"
	"os"

	"github.com/gpdf-dev/gpdf/document"
	"github.com/gpdf-dev/gpdf/pdf"
	"github.com/gpdf-dev/gpdf/template"
)

func main() {
	doc := template.New(
		template.WithPageSize(document.A4),
		template.WithMargins(document.UniformEdges(document.Mm(15))),
		template.WithMetadata(document.DocumentMetadata{
			Title:  "gpdf Example",
			Author: "gpdf",
		}),
	)

	// Header on every page.
	doc.Header(func(p *template.PageBuilder) {
		p.Row(document.Mm(10), func(r *template.RowBuilder) {
			r.Col(6, func(c *template.ColBuilder) {
				c.Text("gpdf", template.FontSize(14), template.Bold(),
					template.TextColor(pdf.RGB(0.2, 0.4, 0.8)))
			})
			r.Col(6, func(c *template.ColBuilder) {
				c.Text("Example Document", template.AlignRight(), template.FontSize(10))
			})
		})
		p.AutoRow(func(r *template.RowBuilder) {
			r.Col(12, func(c *template.ColBuilder) {
				c.Line(template.LineColor(pdf.Gray(0.7)))
			})
		})
	})

	// Page 1: Title and intro.
	page := doc.AddPage()
	page.AutoRow(func(r *template.RowBuilder) {
		r.Col(12, func(c *template.ColBuilder) {
			c.Spacer(document.Mm(10))
			c.Text("Hello, gpdf!", template.FontSize(28), template.Bold())
			c.Spacer(document.Mm(5))
			c.Text("A pure Go, zero-dependency PDF generation library.")
			c.Spacer(document.Mm(10))
		})
	})

	// Two-column layout.
	page.AutoRow(func(r *template.RowBuilder) {
		r.Col(6, func(c *template.ColBuilder) {
			c.Text("Features", template.FontSize(16), template.Bold())
			c.Spacer(document.Mm(3))
			c.Text("- Pure Go, zero dependencies")
			c.Text("- 4-layer architecture")
			c.Text("- CJK Day 1 support")
			c.Text("- Standard library patterns")
		})
		r.Col(6, func(c *template.ColBuilder) {
			c.Text("Architecture", template.FontSize(16), template.Bold())
			c.Spacer(document.Mm(3))
			c.Text("Layer 1: PDF Primitives")
			c.Text("Layer 2: Document Model")
			c.Text("Layer 3: Template API")
			c.Text("Layer 4: CSS Renderer (future)")
		})
	})

	// Table.
	page.AutoRow(func(r *template.RowBuilder) {
		r.Col(12, func(c *template.ColBuilder) {
			c.Spacer(document.Mm(10))
			c.Text("Comparison", template.FontSize(16), template.Bold())
			c.Spacer(document.Mm(3))
			c.Table(
				[]string{"Library", "Status", "Dependencies", "CJK"},
				[][]string{
					{"gpdf", "Active", "Zero", "Day 1"},
					{"gofpdf", "Archived", "Zero", "Partial"},
					{"gopdf", "Active", "Zero", "Yes"},
					{"maroto", "Active", "gofpdf", "Partial"},
				},
				template.ColumnWidths(30, 25, 25, 20),
			)
		})
	})

	data, err := doc.Generate()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	if err := os.WriteFile("example.pdf", data, 0644); err != nil {
		fmt.Fprintf(os.Stderr, "Error writing file: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Generated example.pdf (%d bytes)\n", len(data))
}
