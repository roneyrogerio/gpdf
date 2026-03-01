package benchmark_test

import (
	"bytes"
	"fmt"
	"testing"

	gpdf "github.com/gpdf-dev/gpdf/document"
	"github.com/gpdf-dev/gpdf/template"

	"github.com/go-pdf/fpdf"
	"github.com/johnfercher/maroto/v2"
	"github.com/johnfercher/maroto/v2/pkg/components/page"
	"github.com/johnfercher/maroto/v2/pkg/components/row"
	"github.com/johnfercher/maroto/v2/pkg/components/text"
	"github.com/johnfercher/maroto/v2/pkg/config"
	"github.com/johnfercher/maroto/v2/pkg/consts/align"
	"github.com/johnfercher/maroto/v2/pkg/consts/fontstyle"
	"github.com/johnfercher/maroto/v2/pkg/consts/pagesize"
	"github.com/johnfercher/maroto/v2/pkg/core"
	"github.com/johnfercher/maroto/v2/pkg/props"
	"github.com/signintech/gopdf"
	"golang.org/x/image/font/gofont/goregular"
)

// ---------------------------------------------------------------------------
// SinglePage: A4 1 page + 1 line of text
// ---------------------------------------------------------------------------

func BenchmarkGpdf_SinglePage(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for b.Loop() {
		doc := template.New(
			template.WithPageSize(gpdf.A4),
			template.WithMargins(gpdf.UniformEdges(gpdf.Mm(20))),
		)
		page := doc.AddPage()
		page.AutoRow(func(r *template.RowBuilder) {
			r.Col(12, func(c *template.ColBuilder) {
				c.Text("Hello, World!")
			})
		})
		data, err := doc.Generate()
		if err != nil {
			b.Fatal(err)
		}
		_ = data
	}
}

func BenchmarkMaroto_SinglePage(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for b.Loop() {
		cfg := config.NewBuilder().
			WithPageSize(pagesize.A4).
			WithLeftMargin(20).WithTopMargin(20).WithRightMargin(20).WithBottomMargin(20).
			Build()
		m := maroto.New(cfg)
		m.AddRow(10, text.NewCol(12, "Hello, World!"))
		doc, err := m.Generate()
		if err != nil {
			b.Fatal(err)
		}
		_ = doc.GetBytes()
	}
}

// ---------------------------------------------------------------------------
// Table: 4 columns x 10 rows
// ---------------------------------------------------------------------------

func BenchmarkGpdf_Table(b *testing.B) {
	header := []string{"Name", "Age", "City", "Score"}
	rows := make([][]string, 10)
	for i := range rows {
		rows[i] = []string{
			fmt.Sprintf("User %d", i+1),
			fmt.Sprintf("%d", 20+i),
			fmt.Sprintf("City %d", i+1),
			fmt.Sprintf("%d", 80+i),
		}
	}

	b.ReportAllocs()
	b.ResetTimer()
	for b.Loop() {
		doc := template.New(
			template.WithPageSize(gpdf.A4),
			template.WithMargins(gpdf.UniformEdges(gpdf.Mm(15))),
		)
		page := doc.AddPage()
		page.AutoRow(func(r *template.RowBuilder) {
			r.Col(12, func(c *template.ColBuilder) {
				c.Table(header, rows)
			})
		})
		data, err := doc.Generate()
		if err != nil {
			b.Fatal(err)
		}
		_ = data
	}
}

func BenchmarkMaroto_Table(b *testing.B) {
	header := []string{"Name", "Age", "City", "Score"}
	rows := make([][]string, 10)
	for i := range rows {
		rows[i] = []string{
			fmt.Sprintf("User %d", i+1),
			fmt.Sprintf("%d", 20+i),
			fmt.Sprintf("City %d", i+1),
			fmt.Sprintf("%d", 80+i),
		}
	}

	b.ReportAllocs()
	b.ResetTimer()
	for b.Loop() {
		cfg := config.NewBuilder().
			WithPageSize(pagesize.A4).
			WithLeftMargin(15).WithTopMargin(15).WithRightMargin(15).WithBottomMargin(15).
			Build()
		m := maroto.New(cfg)

		// Header row
		headerCols := make([]core.Col, len(header))
		for i, h := range header {
			headerCols[i] = text.NewCol(3, h, props.Text{Style: fontstyle.Bold})
		}
		m.AddRow(8, headerCols...)

		// Data rows
		for _, r := range rows {
			dataCols := make([]core.Col, len(r))
			for j, cell := range r {
				dataCols[j] = text.NewCol(3, cell)
			}
			m.AddRow(8, dataCols...)
		}

		doc, err := m.Generate()
		if err != nil {
			b.Fatal(err)
		}
		_ = doc.GetBytes()
	}
}

// ---------------------------------------------------------------------------
// 100Pages: 100-page generation
// ---------------------------------------------------------------------------

func BenchmarkGpdf_100Pages(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for b.Loop() {
		doc := template.New(
			template.WithPageSize(gpdf.A4),
			template.WithMargins(gpdf.UniformEdges(gpdf.Mm(20))),
		)
		for i := range 100 {
			pg := doc.AddPage()
			pg.AutoRow(func(r *template.RowBuilder) {
				r.Col(12, func(c *template.ColBuilder) {
					c.Text(fmt.Sprintf("Page %d", i+1))
				})
			})
		}
		data, err := doc.Generate()
		if err != nil {
			b.Fatal(err)
		}
		_ = data
	}
}

func BenchmarkMaroto_100Pages(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for b.Loop() {
		cfg := config.NewBuilder().
			WithPageSize(pagesize.A4).
			WithLeftMargin(20).WithTopMargin(20).WithRightMargin(20).WithBottomMargin(20).
			Build()
		m := maroto.New(cfg)

		for i := range 100 {
			pg := page.New()
			pg.Add(row.New(10).Add(text.NewCol(12, fmt.Sprintf("Page %d", i+1))))
			m.AddPages(pg)
		}

		doc, err := m.Generate()
		if err != nil {
			b.Fatal(err)
		}
		_ = doc.GetBytes()
	}
}

// ---------------------------------------------------------------------------
// ComplexDocument: header + table + multiple text blocks
// ---------------------------------------------------------------------------

func BenchmarkGpdf_ComplexDocument(b *testing.B) {
	header := []string{"Product", "Qty", "Price", "Total"}
	rows := make([][]string, 10)
	for i := range rows {
		rows[i] = []string{
			fmt.Sprintf("Product %d", i+1),
			fmt.Sprintf("%d", (i+1)*2),
			fmt.Sprintf("$%d.00", 10+i),
			fmt.Sprintf("$%d.00", (10+i)*(i+1)*2),
		}
	}

	b.ReportAllocs()
	b.ResetTimer()
	for b.Loop() {
		doc := template.New(
			template.WithPageSize(gpdf.A4),
			template.WithMargins(gpdf.UniformEdges(gpdf.Mm(15))),
		)

		// Header
		doc.Header(func(p *template.PageBuilder) {
			p.AutoRow(func(r *template.RowBuilder) {
				r.Col(6, func(c *template.ColBuilder) {
					c.Text("ACME Corp", template.FontSize(16), template.Bold())
				})
				r.Col(6, func(c *template.ColBuilder) {
					c.Text("Invoice #001", template.FontSize(12), template.AlignRight())
				})
			})
		})

		pg := doc.AddPage()

		// Info section
		pg.AutoRow(func(r *template.RowBuilder) {
			r.Col(6, func(c *template.ColBuilder) {
				c.Text("Bill To: John Doe")
				c.Text("123 Main St")
				c.Text("Anytown, ST 12345")
			})
			r.Col(6, func(c *template.ColBuilder) {
				c.Text("Date: 2026-01-15", template.AlignRight())
				c.Text("Due: 2026-02-15", template.AlignRight())
			})
		})

		// Table
		pg.AutoRow(func(r *template.RowBuilder) {
			r.Col(12, func(c *template.ColBuilder) {
				c.Table(header, rows)
			})
		})

		// Footer text
		pg.AutoRow(func(r *template.RowBuilder) {
			r.Col(12, func(c *template.ColBuilder) {
				c.Text("Thank you for your business!")
				c.Text("Payment terms: Net 30")
			})
		})

		data, err := doc.Generate()
		if err != nil {
			b.Fatal(err)
		}
		_ = data
	}
}

func BenchmarkMaroto_ComplexDocument(b *testing.B) {
	header := []string{"Product", "Qty", "Price", "Total"}
	rows := make([][]string, 10)
	for i := range rows {
		rows[i] = []string{
			fmt.Sprintf("Product %d", i+1),
			fmt.Sprintf("%d", (i+1)*2),
			fmt.Sprintf("$%d.00", 10+i),
			fmt.Sprintf("$%d.00", (10+i)*(i+1)*2),
		}
	}

	b.ReportAllocs()
	b.ResetTimer()
	for b.Loop() {
		cfg := config.NewBuilder().
			WithPageSize(pagesize.A4).
			WithLeftMargin(15).WithTopMargin(15).WithRightMargin(15).WithBottomMargin(15).
			Build()
		m := maroto.New(cfg)

		// Header
		_ = m.RegisterHeader(
			row.New(10).Add(
				text.NewCol(6, "ACME Corp", props.Text{Size: 16, Style: fontstyle.Bold}),
				text.NewCol(6, "Invoice #001", props.Text{Size: 12, Align: align.Right}),
			),
		)

		// Info section
		m.AddRow(8, text.NewCol(6, "Bill To: John Doe"), text.NewCol(6, "Date: 2026-01-15", props.Text{Align: align.Right}))
		m.AddRow(8, text.NewCol(6, "123 Main St"), text.NewCol(6, "Due: 2026-02-15", props.Text{Align: align.Right}))
		m.AddRow(8, text.NewCol(6, "Anytown, ST 12345"))

		// Table header
		headerCols := make([]core.Col, len(header))
		for i, h := range header {
			headerCols[i] = text.NewCol(3, h, props.Text{Style: fontstyle.Bold})
		}
		m.AddRow(8, headerCols...)

		// Table data
		for _, r := range rows {
			dataCols := make([]core.Col, len(r))
			for j, cell := range r {
				dataCols[j] = text.NewCol(3, cell)
			}
			m.AddRow(8, dataCols...)
		}

		// Footer
		m.AddRow(8, text.NewCol(12, "Thank you for your business!"))
		m.AddRow(8, text.NewCol(12, "Payment terms: Net 30"))

		doc, err := m.Generate()
		if err != nil {
			b.Fatal(err)
		}
		_ = doc.GetBytes()
	}
}

// ===========================================================================
// go-pdf/fpdf benchmarks (built-in Helvetica font)
// ===========================================================================

func BenchmarkFpdf_SinglePage(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for b.Loop() {
		pdf := fpdf.New("P", "mm", "A4", "")
		pdf.SetMargins(20, 20, 20)
		pdf.AddPage()
		pdf.SetFont("Helvetica", "", 12)
		pdf.Cell(0, 10, "Hello, World!")
		var buf bytes.Buffer
		if err := pdf.Output(&buf); err != nil {
			b.Fatal(err)
		}
		_ = buf.Bytes()
	}
}

func BenchmarkFpdf_Table(b *testing.B) {
	header := []string{"Name", "Age", "City", "Score"}
	rows := make([][]string, 10)
	for i := range rows {
		rows[i] = []string{
			fmt.Sprintf("User %d", i+1),
			fmt.Sprintf("%d", 20+i),
			fmt.Sprintf("City %d", i+1),
			fmt.Sprintf("%d", 80+i),
		}
	}
	colW := []float64{45, 30, 45, 30}

	b.ReportAllocs()
	b.ResetTimer()
	for b.Loop() {
		pdf := fpdf.New("P", "mm", "A4", "")
		pdf.SetMargins(15, 15, 15)
		pdf.AddPage()

		// Header
		pdf.SetFont("Helvetica", "B", 10)
		for i, h := range header {
			pdf.CellFormat(colW[i], 8, h, "1", 0, "C", false, 0, "")
		}
		pdf.Ln(-1)

		// Data
		pdf.SetFont("Helvetica", "", 10)
		for _, r := range rows {
			for j, cell := range r {
				pdf.CellFormat(colW[j], 8, cell, "1", 0, "L", false, 0, "")
			}
			pdf.Ln(-1)
		}

		var buf bytes.Buffer
		if err := pdf.Output(&buf); err != nil {
			b.Fatal(err)
		}
		_ = buf.Bytes()
	}
}

func BenchmarkFpdf_100Pages(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for b.Loop() {
		pdf := fpdf.New("P", "mm", "A4", "")
		pdf.SetMargins(20, 20, 20)
		pdf.SetFont("Helvetica", "", 12)
		for i := range 100 {
			pdf.AddPage()
			pdf.Text(20, 30, fmt.Sprintf("Page %d", i+1))
		}
		var buf bytes.Buffer
		if err := pdf.Output(&buf); err != nil {
			b.Fatal(err)
		}
		_ = buf.Bytes()
	}
}

func BenchmarkFpdf_ComplexDocument(b *testing.B) {
	header := []string{"Product", "Qty", "Price", "Total"}
	rows := make([][]string, 10)
	for i := range rows {
		rows[i] = []string{
			fmt.Sprintf("Product %d", i+1),
			fmt.Sprintf("%d", (i+1)*2),
			fmt.Sprintf("$%d.00", 10+i),
			fmt.Sprintf("$%d.00", (10+i)*(i+1)*2),
		}
	}
	colW := []float64{60, 30, 40, 40}

	b.ReportAllocs()
	b.ResetTimer()
	for b.Loop() {
		pdf := fpdf.New("P", "mm", "A4", "")
		pdf.SetMargins(15, 15, 15)
		pdf.AddPage()

		// Header
		pdf.SetFont("Helvetica", "B", 16)
		pdf.Cell(90, 10, "ACME Corp")
		pdf.SetFont("Helvetica", "", 12)
		pdf.CellFormat(90, 10, "Invoice #001", "", 1, "R", false, 0, "")

		// Info
		pdf.SetFont("Helvetica", "", 10)
		pdf.Cell(90, 7, "Bill To: John Doe")
		pdf.CellFormat(90, 7, "Date: 2026-01-15", "", 1, "R", false, 0, "")
		pdf.Cell(90, 7, "123 Main St")
		pdf.CellFormat(90, 7, "Due: 2026-02-15", "", 1, "R", false, 0, "")
		pdf.Cell(90, 7, "Anytown, ST 12345")
		pdf.Ln(12)

		// Table header
		pdf.SetFont("Helvetica", "B", 10)
		for i, h := range header {
			pdf.CellFormat(colW[i], 8, h, "1", 0, "C", false, 0, "")
		}
		pdf.Ln(-1)

		// Table data
		pdf.SetFont("Helvetica", "", 10)
		for _, r := range rows {
			for j, cell := range r {
				pdf.CellFormat(colW[j], 8, cell, "1", 0, "L", false, 0, "")
			}
			pdf.Ln(-1)
		}
		pdf.Ln(5)

		// Footer
		pdf.Cell(0, 7, "Thank you for your business!")
		pdf.Ln(7)
		pdf.Cell(0, 7, "Payment terms: Net 30")

		var buf bytes.Buffer
		if err := pdf.Output(&buf); err != nil {
			b.Fatal(err)
		}
		_ = buf.Bytes()
	}
}

// ===========================================================================
// signintech/gopdf benchmarks (requires TTF font — uses Go Regular)
// ===========================================================================

// goFont holds pre-loaded TTF data to avoid filesystem I/O in the benchmark loop.
var goFont = goregular.TTF

func BenchmarkGopdf_SinglePage(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for b.Loop() {
		pdf := gopdf.GoPdf{}
		pdf.Start(gopdf.Config{PageSize: *gopdf.PageSizeA4})
		pdf.SetMargins(20, 20, 20, 20)
		if err := pdf.AddTTFFontData("go", goFont); err != nil {
			b.Fatal(err)
		}
		if err := pdf.SetFont("go", "", 12); err != nil {
			b.Fatal(err)
		}
		pdf.AddPage()
		pdf.SetXY(20, 30)
		_ = pdf.Cell(nil, "Hello, World!")
		data := pdf.GetBytesPdf()
		_ = data
	}
}

func BenchmarkGopdf_Table(b *testing.B) {
	header := []string{"Name", "Age", "City", "Score"}
	rows := make([][]string, 10)
	for i := range rows {
		rows[i] = []string{
			fmt.Sprintf("User %d", i+1),
			fmt.Sprintf("%d", 20+i),
			fmt.Sprintf("City %d", i+1),
			fmt.Sprintf("%d", 80+i),
		}
	}

	b.ReportAllocs()
	b.ResetTimer()
	for b.Loop() {
		pdf := gopdf.GoPdf{}
		pdf.Start(gopdf.Config{PageSize: *gopdf.PageSizeA4})
		pdf.SetMargins(15, 15, 15, 15)
		if err := pdf.AddTTFFontData("go", goFont); err != nil {
			b.Fatal(err)
		}
		pdf.AddPage()

		colW := 42.5
		y := 15.0

		// Header
		_ = pdf.SetFont("go", "", 10)
		for i, h := range header {
			pdf.SetXY(15+float64(i)*colW, y)
			_ = pdf.CellWithOption(&gopdf.Rect{W: colW, H: 8}, h, gopdf.CellOption{
				Border: gopdf.AllBorders,
			})
		}
		y += 8

		// Data
		for _, r := range rows {
			for j, cell := range r {
				pdf.SetXY(15+float64(j)*colW, y)
				_ = pdf.CellWithOption(&gopdf.Rect{W: colW, H: 8}, cell, gopdf.CellOption{
					Border: gopdf.AllBorders,
				})
			}
			y += 8
		}

		data := pdf.GetBytesPdf()
		_ = data
	}
}

func BenchmarkGopdf_100Pages(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for b.Loop() {
		pdf := gopdf.GoPdf{}
		pdf.Start(gopdf.Config{PageSize: *gopdf.PageSizeA4})
		pdf.SetMargins(20, 20, 20, 20)
		if err := pdf.AddTTFFontData("go", goFont); err != nil {
			b.Fatal(err)
		}
		if err := pdf.SetFont("go", "", 12); err != nil {
			b.Fatal(err)
		}
		for i := range 100 {
			pdf.AddPage()
			pdf.SetXY(20, 30)
			_ = pdf.Cell(nil, fmt.Sprintf("Page %d", i+1))
		}
		data := pdf.GetBytesPdf()
		_ = data
	}
}

func BenchmarkGopdf_ComplexDocument(b *testing.B) {
	header := []string{"Product", "Qty", "Price", "Total"}
	rows := make([][]string, 10)
	for i := range rows {
		rows[i] = []string{
			fmt.Sprintf("Product %d", i+1),
			fmt.Sprintf("%d", (i+1)*2),
			fmt.Sprintf("$%d.00", 10+i),
			fmt.Sprintf("$%d.00", (10+i)*(i+1)*2),
		}
	}

	b.ReportAllocs()
	b.ResetTimer()
	for b.Loop() {
		pdf := gopdf.GoPdf{}
		pdf.Start(gopdf.Config{PageSize: *gopdf.PageSizeA4})
		pdf.SetMargins(15, 15, 15, 15)
		if err := pdf.AddTTFFontData("go", goFont); err != nil {
			b.Fatal(err)
		}
		pdf.AddPage()

		// Header
		_ = pdf.SetFont("go", "", 16)
		pdf.SetXY(15, 15)
		_ = pdf.Cell(nil, "ACME Corp")
		_ = pdf.SetFont("go", "", 12)
		pdf.SetXY(105, 15)
		_ = pdf.Cell(nil, "Invoice #001")

		// Info
		_ = pdf.SetFont("go", "", 10)
		pdf.SetXY(15, 30)
		_ = pdf.Cell(nil, "Bill To: John Doe")
		pdf.SetXY(105, 30)
		_ = pdf.Cell(nil, "Date: 2026-01-15")
		pdf.SetXY(15, 38)
		_ = pdf.Cell(nil, "123 Main St")
		pdf.SetXY(105, 38)
		_ = pdf.Cell(nil, "Due: 2026-02-15")
		pdf.SetXY(15, 46)
		_ = pdf.Cell(nil, "Anytown, ST 12345")

		// Table
		colW := 42.5
		y := 58.0
		for i, h := range header {
			pdf.SetXY(15+float64(i)*colW, y)
			_ = pdf.CellWithOption(&gopdf.Rect{W: colW, H: 8}, h, gopdf.CellOption{
				Border: gopdf.AllBorders,
			})
		}
		y += 8
		for _, r := range rows {
			for j, cell := range r {
				pdf.SetXY(15+float64(j)*colW, y)
				_ = pdf.CellWithOption(&gopdf.Rect{W: colW, H: 8}, cell, gopdf.CellOption{
					Border: gopdf.AllBorders,
				})
			}
			y += 8
		}

		// Footer
		y += 5
		pdf.SetXY(15, y)
		_ = pdf.Cell(nil, "Thank you for your business!")
		pdf.SetXY(15, y+8)
		_ = pdf.Cell(nil, "Payment terms: Net 30")

		data := pdf.GetBytesPdf()
		_ = data
	}
}
