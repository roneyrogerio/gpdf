package gpdf_test

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/png"

	gpdf "github.com/gpdf-dev/gpdf"
	"github.com/gpdf-dev/gpdf/document"
	"github.com/gpdf-dev/gpdf/pdf"
	"github.com/gpdf-dev/gpdf/template"
)

func ExampleNewDocument() {
	doc := gpdf.NewDocument(
		gpdf.WithPageSize(document.A4),
		gpdf.WithMargins(document.UniformEdges(document.Mm(20))),
	)
	page := doc.AddPage()
	page.AutoRow(func(r *template.RowBuilder) {
		r.Col(12, func(c *template.ColBuilder) {
			c.Text("Hello, World!", template.FontSize(24))
		})
	})
	data, _ := doc.Generate()
	fmt.Println("PDF starts with:", string(data[:5]))
	// Output: PDF starts with: %PDF-
}

func ExampleNewDocument_withOptions() {
	doc := gpdf.NewDocument(
		gpdf.WithPageSize(document.Letter),
		gpdf.WithMargins(document.Edges{
			Top: document.Mm(25), Right: document.Mm(20),
			Bottom: document.Mm(25), Left: document.Mm(20),
		}),
		gpdf.WithMetadata(document.DocumentMetadata{
			Title:  "Annual Report",
			Author: "gpdf",
		}),
	)
	page := doc.AddPage()
	page.AutoRow(func(r *template.RowBuilder) {
		r.Col(12, func(c *template.ColBuilder) {
			c.Text("Custom page size, margins, and metadata")
		})
	})
	data, _ := doc.Generate()
	fmt.Println("PDF starts with:", string(data[:5]))
	// Output: PDF starts with: %PDF-
}

func Example_textStyling() {
	doc := gpdf.NewDocument()
	page := doc.AddPage()
	page.AutoRow(func(r *template.RowBuilder) {
		r.Col(12, func(c *template.ColBuilder) {
			c.Text("Large bold title", template.FontSize(24), template.Bold())
			c.Text("Italic subtitle", template.Italic())
			c.Text("Red text", template.TextColor(pdf.Red))
			c.Text("Centered", template.AlignCenter())
			c.Text("Underlined", template.Underline())
			c.Text("Yellow background", template.BgColor(pdf.Yellow))
		})
	})
	data, _ := doc.Generate()
	fmt.Println("PDF starts with:", string(data[:5]))
	// Output: PDF starts with: %PDF-
}

func Example_gridLayout() {
	doc := gpdf.NewDocument()
	page := doc.AddPage()
	// Two-column layout (6+6 = 12)
	page.AutoRow(func(r *template.RowBuilder) {
		r.Col(6, func(c *template.ColBuilder) {
			c.Text("Left column")
		})
		r.Col(6, func(c *template.ColBuilder) {
			c.Text("Right column")
		})
	})
	// Three-column layout (4+4+4 = 12)
	page.AutoRow(func(r *template.RowBuilder) {
		r.Col(4, func(c *template.ColBuilder) { c.Text("Col 1") })
		r.Col(4, func(c *template.ColBuilder) { c.Text("Col 2") })
		r.Col(4, func(c *template.ColBuilder) { c.Text("Col 3") })
	})
	data, _ := doc.Generate()
	fmt.Println("PDF starts with:", string(data[:5]))
	// Output: PDF starts with: %PDF-
}

func Example_table() {
	doc := gpdf.NewDocument()
	page := doc.AddPage()
	page.AutoRow(func(r *template.RowBuilder) {
		r.Col(12, func(c *template.ColBuilder) {
			c.Table(
				[]string{"Name", "Age", "City"},
				[][]string{
					{"Alice", "30", "Tokyo"},
					{"Bob", "25", "New York"},
					{"Charlie", "35", "London"},
				},
			)
		})
	})
	data, _ := doc.Generate()
	fmt.Println("PDF starts with:", string(data[:5]))
	// Output: PDF starts with: %PDF-
}

func Example_tableStyled() {
	doc := gpdf.NewDocument()
	page := doc.AddPage()
	page.AutoRow(func(r *template.RowBuilder) {
		r.Col(12, func(c *template.ColBuilder) {
			c.Table(
				[]string{"Product", "Qty", "Price"},
				[][]string{
					{"Widget", "10", "$9.99"},
					{"Gadget", "5", "$24.99"},
				},
				template.ColumnWidths(50, 25, 25),
				template.TableHeaderStyle(
					template.TextColor(pdf.White),
					template.BgColor(pdf.RGBHex(0x1A237E)),
				),
				template.TableStripe(pdf.RGBHex(0xF5F5F5)),
			)
		})
	})
	data, _ := doc.Generate()
	fmt.Println("PDF starts with:", string(data[:5]))
	// Output: PDF starts with: %PDF-
}

func Example_headerFooter() {
	doc := gpdf.NewDocument()
	doc.Header(func(p *template.PageBuilder) {
		p.AutoRow(func(r *template.RowBuilder) {
			r.Col(6, func(c *template.ColBuilder) {
				c.Text("Company Name", template.Bold())
			})
			r.Col(6, func(c *template.ColBuilder) {
				c.Text("Report", template.AlignRight())
			})
		})
	})
	doc.Footer(func(p *template.PageBuilder) {
		p.AutoRow(func(r *template.RowBuilder) {
			r.Col(12, func(c *template.ColBuilder) {
				c.Text("Confidential", template.AlignCenter(), template.FontSize(8))
			})
		})
	})
	page := doc.AddPage()
	page.AutoRow(func(r *template.RowBuilder) {
		r.Col(12, func(c *template.ColBuilder) {
			c.Text("Page content goes here")
		})
	})
	data, _ := doc.Generate()
	fmt.Println("PDF starts with:", string(data[:5]))
	// Output: PDF starts with: %PDF-
}

func Example_list() {
	doc := gpdf.NewDocument()
	page := doc.AddPage()
	page.AutoRow(func(r *template.RowBuilder) {
		r.Col(12, func(c *template.ColBuilder) {
			c.Text("Bullet list:", template.Bold())
			c.List([]string{"First item", "Second item", "Third item"})
			c.Spacer(document.Mm(5))
			c.Text("Numbered list:", template.Bold())
			c.OrderedList([]string{"Step one", "Step two", "Step three"})
		})
	})
	data, _ := doc.Generate()
	fmt.Println("PDF starts with:", string(data[:5]))
	// Output: PDF starts with: %PDF-
}

func Example_image() {
	// Create a small 2x2 test PNG image.
	img := image.NewRGBA(image.Rect(0, 0, 2, 2))
	for y := range 2 {
		for x := range 2 {
			img.Set(x, y, color.RGBA{R: 0, G: 100, B: 200, A: 255})
		}
	}
	var buf bytes.Buffer
	_ = png.Encode(&buf, img)

	doc := gpdf.NewDocument()
	page := doc.AddPage()
	page.AutoRow(func(r *template.RowBuilder) {
		r.Col(12, func(c *template.ColBuilder) {
			c.Text("Embedded image:")
			c.Image(buf.Bytes(), template.FitWidth(document.Mm(30)))
		})
	})
	data, _ := doc.Generate()
	fmt.Println("PDF starts with:", string(data[:5]))
	// Output: PDF starts with: %PDF-
}

func Example_richText() {
	doc := gpdf.NewDocument()
	page := doc.AddPage()
	page.AutoRow(func(r *template.RowBuilder) {
		r.Col(12, func(c *template.ColBuilder) {
			c.RichText(func(rt *template.RichTextBuilder) {
				rt.Span("This is ")
				rt.Span("bold", template.Bold())
				rt.Span(" and this is ")
				rt.Span("red italic", template.Italic(), template.TextColor(pdf.Red))
				rt.Span(" in one line.")
			})
		})
	})
	data, _ := doc.Generate()
	fmt.Println("PDF starts with:", string(data[:5]))
	// Output: PDF starts with: %PDF-
}

func Example_pageNumbers() {
	doc := gpdf.NewDocument()
	doc.Footer(func(p *template.PageBuilder) {
		p.AutoRow(func(r *template.RowBuilder) {
			r.Col(6, func(c *template.ColBuilder) {
				c.PageNumber(template.FontSize(8))
			})
			r.Col(6, func(c *template.ColBuilder) {
				c.TotalPages(template.AlignRight(), template.FontSize(8))
			})
		})
	})
	for range 3 {
		page := doc.AddPage()
		page.AutoRow(func(r *template.RowBuilder) {
			r.Col(12, func(c *template.ColBuilder) {
				c.Text("Content page")
			})
		})
	}
	data, _ := doc.Generate()
	fmt.Println("PDF starts with:", string(data[:5]))
	// Output: PDF starts with: %PDF-
}

func Example_renderToWriter() {
	doc := gpdf.NewDocument()
	page := doc.AddPage()
	page.AutoRow(func(r *template.RowBuilder) {
		r.Col(12, func(c *template.ColBuilder) {
			c.Text("Written via Render(io.Writer)")
		})
	})
	var buf bytes.Buffer
	_ = doc.Render(&buf)
	fmt.Println("PDF starts with:", string(buf.Bytes()[:5]))
	// Output: PDF starts with: %PDF-
}

func Example_qrCode() {
	doc := gpdf.NewDocument()
	page := doc.AddPage()
	page.AutoRow(func(r *template.RowBuilder) {
		r.Col(12, func(c *template.ColBuilder) {
			c.Text("Scan to visit our site:")
			c.QRCode("https://gpdf.dev", template.QRSize(document.Mm(30)))
		})
	})
	data, _ := doc.Generate()
	fmt.Println("PDF starts with:", string(data[:5]))
	// Output: PDF starts with: %PDF-
}

func Example_barcode() {
	doc := gpdf.NewDocument()
	page := doc.AddPage()
	page.AutoRow(func(r *template.RowBuilder) {
		r.Col(12, func(c *template.ColBuilder) {
			c.Text("Order barcode:")
			c.Barcode("INV-2026-0001", template.BarcodeWidth(document.Mm(80)))
		})
	})
	data, _ := doc.Generate()
	fmt.Println("PDF starts with:", string(data[:5]))
	// Output: PDF starts with: %PDF-
}
