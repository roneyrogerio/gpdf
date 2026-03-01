package gpdf_test

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
	"os"
	"testing"

	gpdf "github.com/gpdf-dev/gpdf"
	"github.com/gpdf-dev/gpdf/document"
	"github.com/gpdf-dev/gpdf/pdf"
	"github.com/gpdf-dev/gpdf/template"
)

const (
	outputDir = "_output"
	goldenDir = "testdata/golden"
)

var updateGolden = os.Getenv("UPDATE_GOLDEN") != ""

func writePDF(t *testing.T, name string, data []byte) {
	t.Helper()
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		t.Fatalf("Failed to create output dir: %v", err)
	}
	path := outputDir + "/" + name
	if err := os.WriteFile(path, data, 0644); err != nil {
		t.Fatalf("Failed to write %s: %v", path, err)
	}
	t.Logf("Written %s (%d bytes)", path, len(data))
}

func assertValidPDF(t *testing.T, data []byte) {
	t.Helper()
	if len(data) == 0 {
		t.Fatal("Generated PDF is empty")
	}
	if string(data[:5]) != "%PDF-" {
		t.Fatalf("Invalid PDF header: %q", string(data[:5]))
	}
}

// assertMatchesGolden compares data against the golden file in testdata/golden/.
// When UPDATE_GOLDEN=1, it updates the golden file instead of comparing.
func assertMatchesGolden(t *testing.T, filename string, data []byte) {
	t.Helper()
	goldenPath := goldenDir + "/" + filename

	if updateGolden {
		if err := os.MkdirAll(goldenDir, 0755); err != nil {
			t.Fatalf("Failed to create golden dir: %v", err)
		}
		if err := os.WriteFile(goldenPath, data, 0644); err != nil {
			t.Fatalf("Failed to update golden file %s: %v", goldenPath, err)
		}
		t.Logf("Updated golden file %s (%d bytes)", goldenPath, len(data))
		return
	}

	golden, err := os.ReadFile(goldenPath)
	if err != nil {
		t.Fatalf("Failed to read golden file %s (run with UPDATE_GOLDEN=1 to create): %v", goldenPath, err)
	}
	if !bytes.Equal(data, golden) {
		t.Errorf("Output does not match golden file %s (got %d bytes, want %d bytes; run with UPDATE_GOLDEN=1 to update)", goldenPath, len(data), len(golden))
	}
}

// generatePDF is a helper that calls Generate, validates the output, writes the file,
// and compares against the golden file.
func generatePDF(t *testing.T, filename string, doc *template.Document) {
	t.Helper()
	data, err := doc.Generate()
	if err != nil {
		t.Fatalf("Generate failed: %v", err)
	}
	assertValidPDF(t, data)
	writePDF(t, filename, data)
	assertMatchesGolden(t, filename, data)
}

// testImagePNG creates a small test PNG image (colored rectangle).
func testImagePNG(t *testing.T, w, h int, c color.Color) []byte {
	t.Helper()
	img := image.NewRGBA(image.Rect(0, 0, w, h))
	for y := range h {
		for x := range w {
			img.Set(x, y, c)
		}
	}
	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		t.Fatalf("Failed to create test PNG: %v", err)
	}
	return buf.Bytes()
}

// testImageJPEG creates a small test JPEG image (colored rectangle).
func testImageJPEG(t *testing.T, w, h int, c color.Color) []byte {
	t.Helper()
	img := image.NewRGBA(image.Rect(0, 0, w, h))
	for y := range h {
		for x := range w {
			img.Set(x, y, c)
		}
	}
	var buf bytes.Buffer
	if err := jpeg.Encode(&buf, img, &jpeg.Options{Quality: 80}); err != nil {
		t.Fatalf("Failed to create test JPEG: %v", err)
	}
	return buf.Bytes()
}

// ---------------------------------------------------------------------------
// Basic: Hello World
// ---------------------------------------------------------------------------

func TestExample_01_HelloWorld(t *testing.T) {
	doc := template.New(
		template.WithPageSize(document.A4),
		template.WithMargins(document.UniformEdges(document.Mm(20))),
	)

	page := doc.AddPage()
	page.AutoRow(func(r *template.RowBuilder) {
		r.Col(12, func(c *template.ColBuilder) {
			c.Text("Hello, World!", template.FontSize(24), template.Bold())
		})
	})

	generatePDF(t, "01_hello_world.pdf", doc)
}

// ---------------------------------------------------------------------------
// Text: Font size, weight, style, color, alignment
// ---------------------------------------------------------------------------

func TestExample_02_TextStyling(t *testing.T) {
	doc := template.New(
		template.WithPageSize(document.A4),
		template.WithMargins(document.UniformEdges(document.Mm(20))),
	)

	page := doc.AddPage()

	// Title
	page.AutoRow(func(r *template.RowBuilder) {
		r.Col(12, func(c *template.ColBuilder) {
			c.Text("Text Styling Examples", template.FontSize(20), template.Bold())
			c.Spacer(document.Mm(5))
		})
	})

	// Font sizes
	page.AutoRow(func(r *template.RowBuilder) {
		r.Col(12, func(c *template.ColBuilder) {
			c.Text("Font Size 8pt", template.FontSize(8))
			c.Text("Font Size 12pt (default)", template.FontSize(12))
			c.Text("Font Size 18pt", template.FontSize(18))
			c.Text("Font Size 24pt", template.FontSize(24))
			c.Text("Font Size 36pt", template.FontSize(36))
			c.Spacer(document.Mm(5))
		})
	})

	// Font weight and style
	page.AutoRow(func(r *template.RowBuilder) {
		r.Col(12, func(c *template.ColBuilder) {
			c.Text("Normal text")
			c.Text("Bold text", template.Bold())
			c.Text("Italic text", template.Italic())
			c.Text("Bold + Italic text", template.Bold(), template.Italic())
			c.Spacer(document.Mm(5))
		})
	})

	// Text colors
	page.AutoRow(func(r *template.RowBuilder) {
		r.Col(12, func(c *template.ColBuilder) {
			c.Text("Red text", template.TextColor(pdf.Red))
			c.Text("Green text", template.TextColor(pdf.Green))
			c.Text("Blue text", template.TextColor(pdf.Blue))
			c.Text("Custom color (orange)", template.TextColor(pdf.RGB(1.0, 0.5, 0.0)))
			c.Text("Hex color (#336699)", template.TextColor(pdf.RGBHex(0x336699)))
			c.Spacer(document.Mm(5))
		})
	})

	// Background colors
	page.AutoRow(func(r *template.RowBuilder) {
		r.Col(12, func(c *template.ColBuilder) {
			c.Text("Yellow background", template.BgColor(pdf.Yellow))
			c.Text("Cyan background", template.BgColor(pdf.Cyan))
			c.Text("White text on dark background",
				template.TextColor(pdf.White),
				template.BgColor(pdf.RGBHex(0x333333)),
			)
			c.Spacer(document.Mm(5))
		})
	})

	// Text alignment
	page.AutoRow(func(r *template.RowBuilder) {
		r.Col(12, func(c *template.ColBuilder) {
			c.Text("Left aligned (default)", template.AlignLeft())
			c.Text("Center aligned", template.AlignCenter())
			c.Text("Right aligned", template.AlignRight())
		})
	})

	generatePDF(t, "02_text_styling.pdf", doc)
}

// ---------------------------------------------------------------------------
// Layout: 12-column grid system
// ---------------------------------------------------------------------------

func TestExample_03_GridLayout(t *testing.T) {
	doc := template.New(
		template.WithPageSize(document.A4),
		template.WithMargins(document.UniformEdges(document.Mm(15))),
	)

	page := doc.AddPage()

	// Title
	page.AutoRow(func(r *template.RowBuilder) {
		r.Col(12, func(c *template.ColBuilder) {
			c.Text("12-Column Grid Layout", template.FontSize(18), template.Bold())
			c.Spacer(document.Mm(5))
		})
	})

	// Full width (12 columns)
	page.AutoRow(func(r *template.RowBuilder) {
		r.Col(12, func(c *template.ColBuilder) {
			c.Text("Col 12 (full width)", template.BgColor(pdf.RGBHex(0xE3F2FD)))
		})
	})

	page.AutoRow(func(r *template.RowBuilder) {
		r.Col(12, func(c *template.ColBuilder) {
			c.Spacer(document.Mm(3))
		})
	})

	// Two columns (6 + 6)
	page.AutoRow(func(r *template.RowBuilder) {
		r.Col(6, func(c *template.ColBuilder) {
			c.Text("Col 6 (left)", template.BgColor(pdf.RGBHex(0xE8F5E9)))
		})
		r.Col(6, func(c *template.ColBuilder) {
			c.Text("Col 6 (right)", template.BgColor(pdf.RGBHex(0xFFF3E0)))
		})
	})

	page.AutoRow(func(r *template.RowBuilder) {
		r.Col(12, func(c *template.ColBuilder) {
			c.Spacer(document.Mm(3))
		})
	})

	// Three columns (4 + 4 + 4)
	page.AutoRow(func(r *template.RowBuilder) {
		r.Col(4, func(c *template.ColBuilder) {
			c.Text("Col 4", template.BgColor(pdf.RGBHex(0xFCE4EC)))
		})
		r.Col(4, func(c *template.ColBuilder) {
			c.Text("Col 4", template.BgColor(pdf.RGBHex(0xF3E5F5)))
		})
		r.Col(4, func(c *template.ColBuilder) {
			c.Text("Col 4", template.BgColor(pdf.RGBHex(0xE8EAF6)))
		})
	})

	page.AutoRow(func(r *template.RowBuilder) {
		r.Col(12, func(c *template.ColBuilder) {
			c.Spacer(document.Mm(3))
		})
	})

	// Four columns (3 + 3 + 3 + 3)
	page.AutoRow(func(r *template.RowBuilder) {
		r.Col(3, func(c *template.ColBuilder) {
			c.Text("Col 3", template.BgColor(pdf.RGBHex(0xE0F7FA)))
		})
		r.Col(3, func(c *template.ColBuilder) {
			c.Text("Col 3", template.BgColor(pdf.RGBHex(0xE0F2F1)))
		})
		r.Col(3, func(c *template.ColBuilder) {
			c.Text("Col 3", template.BgColor(pdf.RGBHex(0xFFF9C4)))
		})
		r.Col(3, func(c *template.ColBuilder) {
			c.Text("Col 3", template.BgColor(pdf.RGBHex(0xFFECB3)))
		})
	})

	page.AutoRow(func(r *template.RowBuilder) {
		r.Col(12, func(c *template.ColBuilder) {
			c.Spacer(document.Mm(3))
		})
	})

	// Asymmetric layout (3 + 9)
	page.AutoRow(func(r *template.RowBuilder) {
		r.Col(3, func(c *template.ColBuilder) {
			c.Text("Sidebar (3)", template.BgColor(pdf.RGBHex(0xD7CCC8)))
		})
		r.Col(9, func(c *template.ColBuilder) {
			c.Text("Main content (9)", template.BgColor(pdf.RGBHex(0xF5F5F5)))
		})
	})

	page.AutoRow(func(r *template.RowBuilder) {
		r.Col(12, func(c *template.ColBuilder) {
			c.Spacer(document.Mm(3))
		})
	})

	// Asymmetric layout (8 + 4)
	page.AutoRow(func(r *template.RowBuilder) {
		r.Col(8, func(c *template.ColBuilder) {
			c.Text("Article area (8)", template.BgColor(pdf.RGBHex(0xE1F5FE)))
		})
		r.Col(4, func(c *template.ColBuilder) {
			c.Text("Side panel (4)", template.BgColor(pdf.RGBHex(0xFBE9E7)))
		})
	})

	page.AutoRow(func(r *template.RowBuilder) {
		r.Col(12, func(c *template.ColBuilder) {
			c.Spacer(document.Mm(3))
		})
	})

	// Multiple content in columns
	page.AutoRow(func(r *template.RowBuilder) {
		r.Col(6, func(c *template.ColBuilder) {
			c.Text("Left column - line 1")
			c.Text("Left column - line 2")
			c.Text("Left column - line 3")
		})
		r.Col(6, func(c *template.ColBuilder) {
			c.Text("Right column - line 1")
			c.Text("Right column - line 2")
			c.Text("Right column - line 3")
		})
	})

	generatePDF(t, "03_grid_layout.pdf", doc)
}

// ---------------------------------------------------------------------------
// Layout: Fixed-height rows
// ---------------------------------------------------------------------------

func TestExample_04_FixedHeightRow(t *testing.T) {
	doc := template.New(
		template.WithPageSize(document.A4),
		template.WithMargins(document.UniformEdges(document.Mm(20))),
	)

	page := doc.AddPage()

	page.AutoRow(func(r *template.RowBuilder) {
		r.Col(12, func(c *template.ColBuilder) {
			c.Text("Fixed-Height Row Examples", template.FontSize(18), template.Bold())
			c.Spacer(document.Mm(5))
		})
	})

	// Fixed height row: 30mm
	page.Row(document.Mm(30), func(r *template.RowBuilder) {
		r.Col(12, func(c *template.ColBuilder) {
			c.Text("This row is 30mm tall", template.BgColor(pdf.RGBHex(0xE3F2FD)))
		})
	})

	page.AutoRow(func(r *template.RowBuilder) {
		r.Col(12, func(c *template.ColBuilder) {
			c.Spacer(document.Mm(3))
		})
	})

	// Fixed height row: 50mm
	page.Row(document.Mm(50), func(r *template.RowBuilder) {
		r.Col(6, func(c *template.ColBuilder) {
			c.Text("Left: 50mm row", template.BgColor(pdf.RGBHex(0xE8F5E9)))
		})
		r.Col(6, func(c *template.ColBuilder) {
			c.Text("Right: 50mm row", template.BgColor(pdf.RGBHex(0xFFF3E0)))
		})
	})

	page.AutoRow(func(r *template.RowBuilder) {
		r.Col(12, func(c *template.ColBuilder) {
			c.Spacer(document.Mm(3))
		})
	})

	// Auto-height row for comparison
	page.AutoRow(func(r *template.RowBuilder) {
		r.Col(12, func(c *template.ColBuilder) {
			c.Text("This row has auto height (fits content)", template.BgColor(pdf.RGBHex(0xFCE4EC)))
		})
	})

	generatePDF(t, "04_fixed_height_row.pdf", doc)
}

// ---------------------------------------------------------------------------
// Line: Horizontal rules with styling
// ---------------------------------------------------------------------------

func TestExample_05_Line(t *testing.T) {
	doc := template.New(
		template.WithPageSize(document.A4),
		template.WithMargins(document.UniformEdges(document.Mm(20))),
	)

	page := doc.AddPage()

	page.AutoRow(func(r *template.RowBuilder) {
		r.Col(12, func(c *template.ColBuilder) {
			c.Text("Line / Horizontal Rule Examples", template.FontSize(18), template.Bold())
			c.Spacer(document.Mm(5))
		})
	})

	// Default line
	page.AutoRow(func(r *template.RowBuilder) {
		r.Col(12, func(c *template.ColBuilder) {
			c.Text("Default line (gray, 1pt):")
			c.Line()
			c.Spacer(document.Mm(5))
		})
	})

	// Colored lines
	page.AutoRow(func(r *template.RowBuilder) {
		r.Col(12, func(c *template.ColBuilder) {
			c.Text("Red line:")
			c.Line(template.LineColor(pdf.Red))
			c.Spacer(document.Mm(3))
			c.Text("Blue line:")
			c.Line(template.LineColor(pdf.Blue))
			c.Spacer(document.Mm(3))
			c.Text("Green line:")
			c.Line(template.LineColor(pdf.Green))
			c.Spacer(document.Mm(5))
		})
	})

	// Thick lines
	page.AutoRow(func(r *template.RowBuilder) {
		r.Col(12, func(c *template.ColBuilder) {
			c.Text("Thin line (0.5pt):")
			c.Line(template.LineThickness(document.Pt(0.5)))
			c.Spacer(document.Mm(3))
			c.Text("Medium line (2pt):")
			c.Line(template.LineThickness(document.Pt(2)))
			c.Spacer(document.Mm(3))
			c.Text("Thick line (5pt):")
			c.Line(template.LineThickness(document.Pt(5)))
			c.Spacer(document.Mm(5))
		})
	})

	// Combined: color + thickness
	page.AutoRow(func(r *template.RowBuilder) {
		r.Col(12, func(c *template.ColBuilder) {
			c.Text("Thick red line (3pt):")
			c.Line(template.LineColor(pdf.Red), template.LineThickness(document.Pt(3)))
			c.Spacer(document.Mm(3))
			c.Text("Thick blue line (4pt):")
			c.Line(template.LineColor(pdf.Blue), template.LineThickness(document.Pt(4)))
		})
	})

	generatePDF(t, "05_line.pdf", doc)
}

// ---------------------------------------------------------------------------
// Spacer: Vertical spacing
// ---------------------------------------------------------------------------

func TestExample_06_Spacer(t *testing.T) {
	doc := template.New(
		template.WithPageSize(document.A4),
		template.WithMargins(document.UniformEdges(document.Mm(20))),
	)

	page := doc.AddPage()

	page.AutoRow(func(r *template.RowBuilder) {
		r.Col(12, func(c *template.ColBuilder) {
			c.Text("Spacer Examples", template.FontSize(18), template.Bold())
		})
	})

	page.AutoRow(func(r *template.RowBuilder) {
		r.Col(12, func(c *template.ColBuilder) {
			c.Text("Text before 5mm spacer")
			c.Spacer(document.Mm(5))
			c.Text("Text after 5mm spacer")
		})
	})

	page.AutoRow(func(r *template.RowBuilder) {
		r.Col(12, func(c *template.ColBuilder) {
			c.Text("Text before 15mm spacer")
			c.Spacer(document.Mm(15))
			c.Text("Text after 15mm spacer")
		})
	})

	page.AutoRow(func(r *template.RowBuilder) {
		r.Col(12, func(c *template.ColBuilder) {
			c.Text("Text before 30mm spacer")
			c.Spacer(document.Mm(30))
			c.Text("Text after 30mm spacer")
		})
	})

	generatePDF(t, "06_spacer.pdf", doc)
}

// ---------------------------------------------------------------------------
// Table: Basic
// ---------------------------------------------------------------------------

func TestExample_07_TableBasic(t *testing.T) {
	doc := template.New(
		template.WithPageSize(document.A4),
		template.WithMargins(document.UniformEdges(document.Mm(15))),
	)

	page := doc.AddPage()

	page.AutoRow(func(r *template.RowBuilder) {
		r.Col(12, func(c *template.ColBuilder) {
			c.Text("Basic Table", template.FontSize(18), template.Bold())
			c.Spacer(document.Mm(5))
		})
	})

	page.AutoRow(func(r *template.RowBuilder) {
		r.Col(12, func(c *template.ColBuilder) {
			c.Table(
				[]string{"Name", "Age", "City"},
				[][]string{
					{"Alice", "30", "Tokyo"},
					{"Bob", "25", "New York"},
					{"Charlie", "35", "London"},
					{"Diana", "28", "Paris"},
				},
			)
		})
	})

	generatePDF(t, "07_table_basic.pdf", doc)
}

// ---------------------------------------------------------------------------
// Table: Styled (header colors, stripes, column widths)
// ---------------------------------------------------------------------------

func TestExample_08_TableStyled(t *testing.T) {
	doc := template.New(
		template.WithPageSize(document.A4),
		template.WithMargins(document.UniformEdges(document.Mm(15))),
	)

	page := doc.AddPage()

	page.AutoRow(func(r *template.RowBuilder) {
		r.Col(12, func(c *template.ColBuilder) {
			c.Text("Styled Table", template.FontSize(18), template.Bold())
			c.Spacer(document.Mm(5))
		})
	})

	darkBlue := pdf.RGBHex(0x1A237E)
	lightGray := pdf.RGBHex(0xF5F5F5)

	page.AutoRow(func(r *template.RowBuilder) {
		r.Col(12, func(c *template.ColBuilder) {
			c.Table(
				[]string{"Product", "Category", "Qty", "Unit Price", "Total"},
				[][]string{
					{"Laptop Pro 15", "Electronics", "2", "$1,299.00", "$2,598.00"},
					{"Wireless Mouse", "Accessories", "10", "$29.99", "$299.90"},
					{"USB-C Hub", "Accessories", "5", "$49.99", "$249.95"},
					{"Monitor 27\"", "Electronics", "3", "$399.00", "$1,197.00"},
					{"Keyboard", "Accessories", "10", "$79.99", "$799.90"},
					{"Webcam HD", "Electronics", "4", "$89.99", "$359.96"},
				},
				template.ColumnWidths(30, 20, 10, 20, 20),
				template.TableHeaderStyle(
					template.TextColor(pdf.White),
					template.BgColor(darkBlue),
				),
				template.TableStripe(lightGray),
			)
		})
	})

	generatePDF(t, "08_table_styled.pdf", doc)
}

// ---------------------------------------------------------------------------
// Table: Multiple tables in columns
// ---------------------------------------------------------------------------

func TestExample_09_TableInColumns(t *testing.T) {
	doc := template.New(
		template.WithPageSize(document.A4),
		template.WithMargins(document.UniformEdges(document.Mm(15))),
	)

	page := doc.AddPage()

	page.AutoRow(func(r *template.RowBuilder) {
		r.Col(12, func(c *template.ColBuilder) {
			c.Text("Tables in Grid Columns", template.FontSize(18), template.Bold())
			c.Spacer(document.Mm(5))
		})
	})

	page.AutoRow(func(r *template.RowBuilder) {
		r.Col(6, func(c *template.ColBuilder) {
			c.Text("Team A", template.Bold())
			c.Spacer(document.Mm(2))
			c.Table(
				[]string{"Player", "Score"},
				[][]string{
					{"Alice", "95"},
					{"Bob", "87"},
					{"Charlie", "92"},
				},
				template.ColumnWidths(60, 40),
			)
		})
		r.Col(6, func(c *template.ColBuilder) {
			c.Text("Team B", template.Bold())
			c.Spacer(document.Mm(2))
			c.Table(
				[]string{"Player", "Score"},
				[][]string{
					{"Diana", "91"},
					{"Eve", "88"},
					{"Frank", "85"},
				},
				template.ColumnWidths(60, 40),
			)
		})
	})

	generatePDF(t, "09_table_in_columns.pdf", doc)
}

// ---------------------------------------------------------------------------
// Image: PNG and JPEG
// ---------------------------------------------------------------------------

func TestExample_10_Image(t *testing.T) {
	doc := template.New(
		template.WithPageSize(document.A4),
		template.WithMargins(document.UniformEdges(document.Mm(20))),
	)

	page := doc.AddPage()

	page.AutoRow(func(r *template.RowBuilder) {
		r.Col(12, func(c *template.ColBuilder) {
			c.Text("Image Examples", template.FontSize(18), template.Bold())
			c.Spacer(document.Mm(5))
		})
	})

	// Create test images
	pngData := testImagePNG(t, 200, 100, color.RGBA{R: 66, G: 133, B: 244, A: 255})
	jpegData := testImageJPEG(t, 200, 100, color.RGBA{R: 234, G: 67, B: 53, A: 255})

	// PNG image
	page.AutoRow(func(r *template.RowBuilder) {
		r.Col(12, func(c *template.ColBuilder) {
			c.Text("PNG image (blue):")
			c.Spacer(document.Mm(2))
			c.Image(pngData)
			c.Spacer(document.Mm(5))
		})
	})

	// JPEG image
	page.AutoRow(func(r *template.RowBuilder) {
		r.Col(12, func(c *template.ColBuilder) {
			c.Text("JPEG image (red):")
			c.Spacer(document.Mm(2))
			c.Image(jpegData)
			c.Spacer(document.Mm(5))
		})
	})

	// Images in columns
	page.AutoRow(func(r *template.RowBuilder) {
		r.Col(12, func(c *template.ColBuilder) {
			c.Text("Images side by side in grid columns:")
			c.Spacer(document.Mm(2))
		})
	})

	greenImg := testImagePNG(t, 150, 80, color.RGBA{R: 52, G: 168, B: 83, A: 255})
	yellowImg := testImagePNG(t, 150, 80, color.RGBA{R: 251, G: 188, B: 4, A: 255})

	page.AutoRow(func(r *template.RowBuilder) {
		r.Col(6, func(c *template.ColBuilder) {
			c.Text("Green PNG")
			c.Image(greenImg)
		})
		r.Col(6, func(c *template.ColBuilder) {
			c.Text("Yellow PNG")
			c.Image(yellowImg)
		})
	})

	generatePDF(t, "10_image.pdf", doc)
}

// ---------------------------------------------------------------------------
// Image: Fit options
// ---------------------------------------------------------------------------

func TestExample_11_ImageFit(t *testing.T) {
	doc := template.New(
		template.WithPageSize(document.A4),
		template.WithMargins(document.UniformEdges(document.Mm(20))),
	)

	page := doc.AddPage()

	page.AutoRow(func(r *template.RowBuilder) {
		r.Col(12, func(c *template.ColBuilder) {
			c.Text("Image Fit Options", template.FontSize(18), template.Bold())
			c.Spacer(document.Mm(5))
		})
	})

	imgData := testImagePNG(t, 300, 200, color.RGBA{R: 100, G: 149, B: 237, A: 255})

	// FitWidth
	page.AutoRow(func(r *template.RowBuilder) {
		r.Col(12, func(c *template.ColBuilder) {
			c.Text("FitWidth(80mm):")
			c.Spacer(document.Mm(2))
			c.Image(imgData, template.FitWidth(document.Mm(80)))
			c.Spacer(document.Mm(5))
		})
	})

	// FitHeight
	page.AutoRow(func(r *template.RowBuilder) {
		r.Col(12, func(c *template.ColBuilder) {
			c.Text("FitHeight(30mm):")
			c.Spacer(document.Mm(2))
			c.Image(imgData, template.FitHeight(document.Mm(30)))
			c.Spacer(document.Mm(5))
		})
	})

	// Default (no fit options)
	page.AutoRow(func(r *template.RowBuilder) {
		r.Col(12, func(c *template.ColBuilder) {
			c.Text("Default (no fit options):")
			c.Spacer(document.Mm(2))
			c.Image(imgData)
		})
	})

	generatePDF(t, "11_image_fit.pdf", doc)
}

// ---------------------------------------------------------------------------
// Multi-page document
// ---------------------------------------------------------------------------

func TestExample_12_MultiPage(t *testing.T) {
	doc := template.New(
		template.WithPageSize(document.A4),
		template.WithMargins(document.UniformEdges(document.Mm(20))),
	)

	for i := 1; i <= 5; i++ {
		page := doc.AddPage()
		page.AutoRow(func(r *template.RowBuilder) {
			r.Col(12, func(c *template.ColBuilder) {
				c.Text("Multi-Page Document", template.FontSize(20), template.Bold())
				c.Spacer(document.Mm(5))
				c.Line()
				c.Spacer(document.Mm(10))
			})
		})

		// Fill the page with some content
		for j := 1; j <= 10; j++ {
			page.AutoRow(func(r *template.RowBuilder) {
				r.Col(12, func(c *template.ColBuilder) {
					c.Text("Lorem ipsum dolor sit amet, consectetur adipiscing elit. " +
						"Sed do eiusmod tempor incididunt ut labore et dolore magna aliqua.")
				})
			})
		}
	}

	generatePDF(t, "12_multi_page.pdf", doc)
}

// ---------------------------------------------------------------------------
// Header and Footer
// ---------------------------------------------------------------------------

func TestExample_13_HeaderFooter(t *testing.T) {
	doc := template.New(
		template.WithPageSize(document.A4),
		template.WithMargins(document.UniformEdges(document.Mm(20))),
	)

	// Header: company name on left, document title on right
	doc.Header(func(p *template.PageBuilder) {
		p.AutoRow(func(r *template.RowBuilder) {
			r.Col(6, func(c *template.ColBuilder) {
				c.Text("ACME Corporation", template.Bold(), template.FontSize(10))
			})
			r.Col(6, func(c *template.ColBuilder) {
				c.Text("Confidential Report", template.AlignRight(), template.FontSize(10),
					template.TextColor(pdf.Gray(0.5)))
			})
		})
		p.AutoRow(func(r *template.RowBuilder) {
			r.Col(12, func(c *template.ColBuilder) {
				c.Line(template.LineColor(pdf.RGBHex(0x1565C0)), template.LineThickness(document.Pt(2)))
				c.Spacer(document.Mm(5))
			})
		})
	})

	// Footer: centered text with line separator
	doc.Footer(func(p *template.PageBuilder) {
		p.AutoRow(func(r *template.RowBuilder) {
			r.Col(12, func(c *template.ColBuilder) {
				c.Spacer(document.Mm(5))
				c.Line(template.LineColor(pdf.Gray(0.7)))
				c.Spacer(document.Mm(2))
			})
		})
		p.AutoRow(func(r *template.RowBuilder) {
			r.Col(6, func(c *template.ColBuilder) {
				c.Text("Generated by gpdf", template.FontSize(8), template.TextColor(pdf.Gray(0.5)))
			})
			r.Col(6, func(c *template.ColBuilder) {
				c.Text("Confidential", template.FontSize(8),
					template.AlignRight(), template.TextColor(pdf.Gray(0.5)))
			})
		})
	})

	// Three pages of content
	for range 3 {
		page := doc.AddPage()
		page.AutoRow(func(r *template.RowBuilder) {
			r.Col(12, func(c *template.ColBuilder) {
				c.Text("Main Content Area", template.FontSize(16), template.Bold())
				c.Spacer(document.Mm(5))
				c.Text("This page demonstrates header and footer repeated on every page. " +
					"The header contains the company name and document title. " +
					"The footer contains generation info and a confidentiality notice.")
			})
		})
	}

	generatePDF(t, "13_header_footer.pdf", doc)
}

// ---------------------------------------------------------------------------
// Metadata
// ---------------------------------------------------------------------------

func TestExample_14_Metadata(t *testing.T) {
	doc := template.New(
		template.WithPageSize(document.A4),
		template.WithMargins(document.UniformEdges(document.Mm(20))),
		template.WithMetadata(document.DocumentMetadata{
			Title:   "Annual Report 2026",
			Author:  "gpdf Library",
			Subject: "Example of document metadata",
			Creator: "gpdf example_test.go",
		}),
	)

	page := doc.AddPage()
	page.AutoRow(func(r *template.RowBuilder) {
		r.Col(12, func(c *template.ColBuilder) {
			c.Text("Document with Metadata", template.FontSize(20), template.Bold())
			c.Spacer(document.Mm(5))
			c.Text("This PDF has the following metadata set:")
			c.Spacer(document.Mm(3))
			c.Text("Title: Annual Report 2026")
			c.Text("Author: gpdf Library")
			c.Text("Subject: Example of document metadata")
			c.Text("Creator: gpdf example_test.go")
			c.Text("Producer: gpdf (set automatically)")
			c.Spacer(document.Mm(5))
			c.Text("Open the PDF properties in your viewer to verify.", template.Italic())
		})
	})

	generatePDF(t, "14_metadata.pdf", doc)
}

// ---------------------------------------------------------------------------
// Page sizes
// ---------------------------------------------------------------------------

func TestExample_15_PageSizes(t *testing.T) {
	sizes := []struct {
		name string
		size document.Size
		file string
	}{
		{"A4 (210mm x 297mm)", document.A4, "15a_pagesize_a4.pdf"},
		{"A3 (297mm x 420mm)", document.A3, "15b_pagesize_a3.pdf"},
		{"Letter (8.5in x 11in)", document.Letter, "15c_pagesize_letter.pdf"},
		{"Legal (8.5in x 14in)", document.Legal, "15d_pagesize_legal.pdf"},
	}

	for _, s := range sizes {
		t.Run(s.name, func(t *testing.T) {
			doc := template.New(
				template.WithPageSize(s.size),
				template.WithMargins(document.UniformEdges(document.Mm(20))),
			)

			page := doc.AddPage()
			page.AutoRow(func(r *template.RowBuilder) {
				r.Col(12, func(c *template.ColBuilder) {
					c.Text("Page Size: "+s.name, template.FontSize(20), template.Bold())
					c.Spacer(document.Mm(10))
					c.Text("This page demonstrates the " + s.name + " page format.")
				})
			})

			generatePDF(t, s.file, doc)
		})
	}
}

// ---------------------------------------------------------------------------
// Custom margins
// ---------------------------------------------------------------------------

func TestExample_16_Margins(t *testing.T) {
	doc := template.New(
		template.WithPageSize(document.A4),
		template.WithMargins(document.Edges{
			Top:    document.Mm(10),
			Right:  document.Mm(40),
			Bottom: document.Mm(10),
			Left:   document.Mm(40),
		}),
	)

	page := doc.AddPage()
	page.AutoRow(func(r *template.RowBuilder) {
		r.Col(12, func(c *template.ColBuilder) {
			c.Text("Custom Margins", template.FontSize(18), template.Bold())
			c.Spacer(document.Mm(5))
			c.Text("This page has asymmetric margins: 10mm top/bottom, 40mm left/right. " +
				"The wide side margins create a narrower text area, similar to a book layout.")
			c.Spacer(document.Mm(5))
			c.Line()
			c.Spacer(document.Mm(5))
			c.Text("Lorem ipsum dolor sit amet, consectetur adipiscing elit. " +
				"Sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. " +
				"Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris.")
		})
	})

	generatePDF(t, "16_margins.pdf", doc)
}

// ---------------------------------------------------------------------------
// Colors: RGB, Hex, Gray, CMYK
// ---------------------------------------------------------------------------

func TestExample_17_Colors(t *testing.T) {
	doc := template.New(
		template.WithPageSize(document.A4),
		template.WithMargins(document.UniformEdges(document.Mm(20))),
	)

	page := doc.AddPage()

	page.AutoRow(func(r *template.RowBuilder) {
		r.Col(12, func(c *template.ColBuilder) {
			c.Text("Color System Examples", template.FontSize(18), template.Bold())
			c.Spacer(document.Mm(5))
		})
	})

	// Predefined colors
	page.AutoRow(func(r *template.RowBuilder) {
		r.Col(12, func(c *template.ColBuilder) {
			c.Text("Predefined Colors:", template.Bold())
			c.Text("Red", template.TextColor(pdf.Red))
			c.Text("Green", template.TextColor(pdf.Green))
			c.Text("Blue", template.TextColor(pdf.Blue))
			c.Text("Yellow", template.TextColor(pdf.Yellow))
			c.Text("Cyan", template.TextColor(pdf.Cyan))
			c.Text("Magenta", template.TextColor(pdf.Magenta))
			c.Spacer(document.Mm(5))
		})
	})

	// RGB colors (0.0-1.0)
	page.AutoRow(func(r *template.RowBuilder) {
		r.Col(12, func(c *template.ColBuilder) {
			c.Text("RGB Colors (float):", template.Bold())
			c.Text("RGB(1.0, 0.5, 0.0) - Orange", template.TextColor(pdf.RGB(1.0, 0.5, 0.0)))
			c.Text("RGB(0.5, 0.0, 0.5) - Purple", template.TextColor(pdf.RGB(0.5, 0.0, 0.5)))
			c.Text("RGB(0.0, 0.5, 0.5) - Teal", template.TextColor(pdf.RGB(0.0, 0.5, 0.5)))
			c.Spacer(document.Mm(5))
		})
	})

	// Hex colors
	page.AutoRow(func(r *template.RowBuilder) {
		r.Col(12, func(c *template.ColBuilder) {
			c.Text("Hex Colors:", template.Bold())
			c.Text("#FF6B6B - Coral", template.TextColor(pdf.RGBHex(0xFF6B6B)))
			c.Text("#4ECDC4 - Turquoise", template.TextColor(pdf.RGBHex(0x4ECDC4)))
			c.Text("#45B7D1 - Sky Blue", template.TextColor(pdf.RGBHex(0x45B7D1)))
			c.Text("#96CEB4 - Sage", template.TextColor(pdf.RGBHex(0x96CEB4)))
			c.Spacer(document.Mm(5))
		})
	})

	// Grayscale
	page.AutoRow(func(r *template.RowBuilder) {
		r.Col(12, func(c *template.ColBuilder) {
			c.Text("Grayscale:", template.Bold())
			c.Text("Gray(0.0) - Black", template.TextColor(pdf.Gray(0.0)))
			c.Text("Gray(0.3) - Dark gray", template.TextColor(pdf.Gray(0.3)))
			c.Text("Gray(0.5) - Medium gray", template.TextColor(pdf.Gray(0.5)))
			c.Text("Gray(0.7) - Light gray", template.TextColor(pdf.Gray(0.7)))
			c.Spacer(document.Mm(5))
		})
	})

	// Background color swatches
	page.AutoRow(func(r *template.RowBuilder) {
		r.Col(12, func(c *template.ColBuilder) {
			c.Text("Background Color Swatches:", template.Bold())
		})
	})

	page.AutoRow(func(r *template.RowBuilder) {
		r.Col(3, func(c *template.ColBuilder) {
			c.Text(" Red ", template.TextColor(pdf.White), template.BgColor(pdf.Red))
		})
		r.Col(3, func(c *template.ColBuilder) {
			c.Text(" Green ", template.TextColor(pdf.White), template.BgColor(pdf.Green))
		})
		r.Col(3, func(c *template.ColBuilder) {
			c.Text(" Blue ", template.TextColor(pdf.White), template.BgColor(pdf.Blue))
		})
		r.Col(3, func(c *template.ColBuilder) {
			c.Text(" Yellow ", template.BgColor(pdf.Yellow))
		})
	})

	generatePDF(t, "17_colors.pdf", doc)
}

// ---------------------------------------------------------------------------
// Combined: Invoice example
// ---------------------------------------------------------------------------

func TestExample_18_Invoice(t *testing.T) {
	doc := template.New(
		template.WithPageSize(document.A4),
		template.WithMargins(document.UniformEdges(document.Mm(20))),
		template.WithMetadata(document.DocumentMetadata{
			Title:  "Invoice #INV-2026-001",
			Author: "ACME Corporation",
		}),
	)

	page := doc.AddPage()

	// Company header
	page.AutoRow(func(r *template.RowBuilder) {
		r.Col(6, func(c *template.ColBuilder) {
			c.Text("ACME Corporation", template.FontSize(24), template.Bold(),
				template.TextColor(pdf.RGBHex(0x1A237E)))
			c.Text("123 Business Street")
			c.Text("Suite 100")
			c.Text("San Francisco, CA 94105")
		})
		r.Col(6, func(c *template.ColBuilder) {
			c.Text("INVOICE", template.FontSize(28), template.Bold(), template.AlignRight(),
				template.TextColor(pdf.RGBHex(0x1A237E)))
			c.Spacer(document.Mm(3))
			c.Text("#INV-2026-001", template.AlignRight(), template.FontSize(12))
			c.Text("Date: March 1, 2026", template.AlignRight())
			c.Text("Due: March 31, 2026", template.AlignRight())
		})
	})

	page.AutoRow(func(r *template.RowBuilder) {
		r.Col(12, func(c *template.ColBuilder) {
			c.Spacer(document.Mm(5))
			c.Line(template.LineColor(pdf.RGBHex(0x1A237E)), template.LineThickness(document.Pt(2)))
			c.Spacer(document.Mm(5))
		})
	})

	// Bill to
	page.AutoRow(func(r *template.RowBuilder) {
		r.Col(6, func(c *template.ColBuilder) {
			c.Text("Bill To:", template.Bold(), template.TextColor(pdf.Gray(0.4)))
			c.Spacer(document.Mm(2))
			c.Text("John Smith", template.Bold())
			c.Text("Tech Solutions Inc.")
			c.Text("456 Client Avenue")
			c.Text("New York, NY 10001")
		})
		r.Col(6, func(c *template.ColBuilder) {
			c.Text("Payment Info:", template.Bold(), template.TextColor(pdf.Gray(0.4)))
			c.Spacer(document.Mm(2))
			c.Text("Bank: First National Bank")
			c.Text("Account: 1234-5678-9012")
			c.Text("Routing: 021000021")
		})
	})

	page.AutoRow(func(r *template.RowBuilder) {
		r.Col(12, func(c *template.ColBuilder) {
			c.Spacer(document.Mm(10))
		})
	})

	// Items table
	headerBlue := pdf.RGBHex(0x1A237E)
	stripeGray := pdf.RGBHex(0xF5F5F5)

	page.AutoRow(func(r *template.RowBuilder) {
		r.Col(12, func(c *template.ColBuilder) {
			c.Table(
				[]string{"Description", "Qty", "Unit Price", "Amount"},
				[][]string{
					{"Web Development - Frontend", "40 hrs", "$150.00", "$6,000.00"},
					{"Web Development - Backend", "60 hrs", "$150.00", "$9,000.00"},
					{"UI/UX Design", "20 hrs", "$120.00", "$2,400.00"},
					{"Database Design", "15 hrs", "$130.00", "$1,950.00"},
					{"QA Testing", "25 hrs", "$100.00", "$2,500.00"},
					{"Project Management", "10 hrs", "$140.00", "$1,400.00"},
				},
				template.ColumnWidths(40, 15, 20, 25),
				template.TableHeaderStyle(
					template.TextColor(pdf.White),
					template.BgColor(headerBlue),
				),
				template.TableStripe(stripeGray),
			)
		})
	})

	page.AutoRow(func(r *template.RowBuilder) {
		r.Col(12, func(c *template.ColBuilder) {
			c.Spacer(document.Mm(5))
		})
	})

	// Totals
	page.AutoRow(func(r *template.RowBuilder) {
		r.Col(8, func(c *template.ColBuilder) {
			// empty left side
		})
		r.Col(4, func(c *template.ColBuilder) {
			c.Text("Subtotal:    $23,250.00", template.AlignRight())
			c.Text("Tax (10%):    $2,325.00", template.AlignRight())
			c.Spacer(document.Mm(2))
			c.Line(template.LineThickness(document.Pt(1)))
			c.Spacer(document.Mm(2))
			c.Text("Total:       $25,575.00", template.AlignRight(),
				template.Bold(), template.FontSize(14))
		})
	})

	page.AutoRow(func(r *template.RowBuilder) {
		r.Col(12, func(c *template.ColBuilder) {
			c.Spacer(document.Mm(15))
			c.Line(template.LineColor(pdf.Gray(0.8)))
			c.Spacer(document.Mm(3))
			c.Text("Thank you for your business!", template.AlignCenter(),
				template.Italic(), template.TextColor(pdf.Gray(0.5)))
		})
	})

	generatePDF(t, "18_invoice.pdf", doc)
}

// ---------------------------------------------------------------------------
// Combined: Report with all features
// ---------------------------------------------------------------------------

func TestExample_19_Report(t *testing.T) {
	doc := template.New(
		template.WithPageSize(document.A4),
		template.WithMargins(document.UniformEdges(document.Mm(20))),
		template.WithMetadata(document.DocumentMetadata{
			Title:   "Quarterly Report Q1 2026",
			Author:  "ACME Corporation",
			Subject: "Q1 2026 Financial Summary",
		}),
	)

	// Header for all pages
	doc.Header(func(p *template.PageBuilder) {
		p.AutoRow(func(r *template.RowBuilder) {
			r.Col(6, func(c *template.ColBuilder) {
				c.Text("ACME Corp", template.Bold(), template.FontSize(9),
					template.TextColor(pdf.RGBHex(0x1565C0)))
			})
			r.Col(6, func(c *template.ColBuilder) {
				c.Text("Q1 2026 Report", template.AlignRight(), template.FontSize(9),
					template.TextColor(pdf.Gray(0.5)))
			})
		})
		p.AutoRow(func(r *template.RowBuilder) {
			r.Col(12, func(c *template.ColBuilder) {
				c.Line(template.LineColor(pdf.RGBHex(0x1565C0)))
				c.Spacer(document.Mm(3))
			})
		})
	})

	// Footer for all pages
	doc.Footer(func(p *template.PageBuilder) {
		p.AutoRow(func(r *template.RowBuilder) {
			r.Col(12, func(c *template.ColBuilder) {
				c.Spacer(document.Mm(3))
				c.Line(template.LineColor(pdf.Gray(0.8)))
				c.Spacer(document.Mm(2))
			})
		})
		p.AutoRow(func(r *template.RowBuilder) {
			r.Col(12, func(c *template.ColBuilder) {
				c.Text("Confidential - For Internal Use Only",
					template.AlignCenter(), template.FontSize(7), template.TextColor(pdf.Gray(0.5)))
			})
		})
	})

	// --- Page 1: Title & Executive Summary ---
	page1 := doc.AddPage()

	page1.AutoRow(func(r *template.RowBuilder) {
		r.Col(12, func(c *template.ColBuilder) {
			c.Spacer(document.Mm(20))
			c.Text("Quarterly Report", template.FontSize(28), template.Bold(),
				template.AlignCenter(), template.TextColor(pdf.RGBHex(0x1A237E)))
			c.Text("Q1 2026 - Financial Summary", template.FontSize(16),
				template.AlignCenter(), template.TextColor(pdf.Gray(0.4)))
			c.Spacer(document.Mm(15))
			c.Line(template.LineColor(pdf.RGBHex(0x1A237E)), template.LineThickness(document.Pt(2)))
			c.Spacer(document.Mm(10))
		})
	})

	page1.AutoRow(func(r *template.RowBuilder) {
		r.Col(12, func(c *template.ColBuilder) {
			c.Text("Executive Summary", template.FontSize(16), template.Bold())
			c.Spacer(document.Mm(3))
			c.Text("This report presents the financial performance of ACME Corporation " +
				"for the first quarter of 2026. Revenue increased by 15% compared to Q4 2025, " +
				"driven primarily by strong growth in the cloud services division. " +
				"Operating margins improved to 22%, up from 19% in the previous quarter.")
			c.Spacer(document.Mm(5))
			c.Text("Key highlights include the successful launch of three new product lines, " +
				"expansion into the European market, and a 20% reduction in customer churn rate. " +
				"The company remains well-positioned for continued growth throughout 2026.")
		})
	})

	// Key metrics in grid
	page1.AutoRow(func(r *template.RowBuilder) {
		r.Col(12, func(c *template.ColBuilder) {
			c.Spacer(document.Mm(10))
			c.Text("Key Metrics", template.FontSize(14), template.Bold())
			c.Spacer(document.Mm(3))
		})
	})

	page1.AutoRow(func(r *template.RowBuilder) {
		r.Col(3, func(c *template.ColBuilder) {
			c.Text("Revenue", template.TextColor(pdf.Gray(0.5)), template.FontSize(9))
			c.Text("$12.5M", template.FontSize(18), template.Bold(),
				template.TextColor(pdf.RGBHex(0x2E7D32)))
		})
		r.Col(3, func(c *template.ColBuilder) {
			c.Text("Growth", template.TextColor(pdf.Gray(0.5)), template.FontSize(9))
			c.Text("+15%", template.FontSize(18), template.Bold(),
				template.TextColor(pdf.RGBHex(0x2E7D32)))
		})
		r.Col(3, func(c *template.ColBuilder) {
			c.Text("Customers", template.TextColor(pdf.Gray(0.5)), template.FontSize(9))
			c.Text("2,450", template.FontSize(18), template.Bold(),
				template.TextColor(pdf.RGBHex(0x1565C0)))
		})
		r.Col(3, func(c *template.ColBuilder) {
			c.Text("Margin", template.TextColor(pdf.Gray(0.5)), template.FontSize(9))
			c.Text("22%", template.FontSize(18), template.Bold(),
				template.TextColor(pdf.RGBHex(0x1565C0)))
		})
	})

	// --- Page 2: Financial Details ---
	page2 := doc.AddPage()

	page2.AutoRow(func(r *template.RowBuilder) {
		r.Col(12, func(c *template.ColBuilder) {
			c.Text("Revenue Breakdown", template.FontSize(16), template.Bold())
			c.Spacer(document.Mm(5))
		})
	})

	darkHeader := pdf.RGBHex(0x1A237E)
	stripe := pdf.RGBHex(0xF5F5F5)

	page2.AutoRow(func(r *template.RowBuilder) {
		r.Col(12, func(c *template.ColBuilder) {
			c.Table(
				[]string{"Division", "Q1 2026", "Q4 2025", "Change"},
				[][]string{
					{"Cloud Services", "$5,200,000", "$4,100,000", "+26.8%"},
					{"Enterprise Software", "$3,800,000", "$3,500,000", "+8.6%"},
					{"Consulting", "$2,100,000", "$1,900,000", "+10.5%"},
					{"Support & Maintenance", "$1,400,000", "$1,350,000", "+3.7%"},
				},
				template.ColumnWidths(35, 22, 22, 21),
				template.TableHeaderStyle(
					template.TextColor(pdf.White),
					template.BgColor(darkHeader),
				),
				template.TableStripe(stripe),
			)
			c.Spacer(document.Mm(10))
		})
	})

	page2.AutoRow(func(r *template.RowBuilder) {
		r.Col(12, func(c *template.ColBuilder) {
			c.Text("Expense Summary", template.FontSize(16), template.Bold())
			c.Spacer(document.Mm(5))
		})
	})

	page2.AutoRow(func(r *template.RowBuilder) {
		r.Col(12, func(c *template.ColBuilder) {
			c.Table(
				[]string{"Category", "Amount", "% of Revenue"},
				[][]string{
					{"Personnel", "$5,500,000", "44.0%"},
					{"Infrastructure", "$1,800,000", "14.4%"},
					{"Marketing", "$1,200,000", "9.6%"},
					{"R&D", "$950,000", "7.6%"},
					{"General & Admin", "$300,000", "2.4%"},
				},
				template.ColumnWidths(40, 30, 30),
				template.TableHeaderStyle(
					template.TextColor(pdf.White),
					template.BgColor(darkHeader),
				),
				template.TableStripe(stripe),
			)
			c.Spacer(document.Mm(10))
		})
	})

	// Two-column commentary
	page2.AutoRow(func(r *template.RowBuilder) {
		r.Col(6, func(c *template.ColBuilder) {
			c.Text("Highlights", template.Bold(), template.TextColor(pdf.RGBHex(0x2E7D32)))
			c.Spacer(document.Mm(2))
			c.Text("Cloud services revenue grew 26.8%, exceeding projections by 5%. " +
				"New enterprise clients added: 47.")
		})
		r.Col(6, func(c *template.ColBuilder) {
			c.Text("Challenges", template.Bold(), template.TextColor(pdf.RGBHex(0xC62828)))
			c.Spacer(document.Mm(2))
			c.Text("Infrastructure costs rose 12% due to scaling needs. " +
				"Two major client renewals deferred to Q2.")
		})
	})

	generatePDF(t, "19_report.pdf", doc)
}

// ---------------------------------------------------------------------------
// Render to io.Writer
// ---------------------------------------------------------------------------

func TestExample_20_RenderToWriter(t *testing.T) {
	doc := template.New(
		template.WithPageSize(document.A4),
		template.WithMargins(document.UniformEdges(document.Mm(20))),
	)

	page := doc.AddPage()
	page.AutoRow(func(r *template.RowBuilder) {
		r.Col(12, func(c *template.ColBuilder) {
			c.Text("Render to io.Writer", template.FontSize(18), template.Bold())
			c.Spacer(document.Mm(5))
			c.Text("This PDF was rendered using doc.Render(w) instead of doc.Generate().")
		})
	})

	// Use Render instead of Generate
	var buf bytes.Buffer
	if err := doc.Render(&buf); err != nil {
		t.Fatalf("Render failed: %v", err)
	}

	data := buf.Bytes()
	assertValidPDF(t, data)
	writePDF(t, "20_render_to_writer.pdf", data)
	assertMatchesGolden(t, "20_render_to_writer.pdf", data)
}

// ---------------------------------------------------------------------------
// gpdf facade (root package convenience API)
// ---------------------------------------------------------------------------

func TestExample_21_GpdfFacade(t *testing.T) {
	// Use the root gpdf package convenience functions
	// (imported as gpdf_test, so we use template.New which gpdf.NewDocument wraps)
	doc := template.New(
		template.WithPageSize(document.A4),
		template.WithMargins(document.UniformEdges(document.Mm(20))),
		template.WithDefaultFont("", 14),
		template.WithMetadata(document.DocumentMetadata{
			Title:  "Facade Example",
			Author: "gpdf",
		}),
	)

	page := doc.AddPage()
	page.AutoRow(func(r *template.RowBuilder) {
		r.Col(12, func(c *template.ColBuilder) {
			c.Text("gpdf Facade API", template.FontSize(20), template.Bold())
			c.Spacer(document.Mm(5))
			c.Text("This document uses WithDefaultFont to set the base font size to 14pt.")
			c.Text("All text in this document inherits the 14pt default.")
		})
	})

	generatePDF(t, "21_facade.pdf", doc)
}

// ---------------------------------------------------------------------------
// LetterSpacing (WP1)
// ---------------------------------------------------------------------------

func TestExample_22_LetterSpacing(t *testing.T) {
	doc := template.New(
		template.WithPageSize(document.A4),
		template.WithMargins(document.UniformEdges(document.Mm(20))),
	)

	page := doc.AddPage()
	page.AutoRow(func(r *template.RowBuilder) {
		r.Col(12, func(c *template.ColBuilder) {
			c.Text("Letter Spacing Demo", template.FontSize(20), template.Bold())
			c.Spacer(document.Mm(8))

			c.Text("Normal spacing (0pt)")
			c.Spacer(document.Mm(3))

			c.Text("Letter spacing 1pt", template.LetterSpacing(1))
			c.Spacer(document.Mm(3))

			c.Text("Letter spacing 3pt", template.LetterSpacing(3))
			c.Spacer(document.Mm(3))

			c.Text("WIDE HEADER", template.FontSize(16), template.Bold(),
				template.LetterSpacing(5))
			c.Spacer(document.Mm(3))

			c.Text("Tight spacing -0.5pt", template.LetterSpacing(-0.5))
		})
	})

	generatePDF(t, "22_letter_spacing.pdf", doc)
}

// ---------------------------------------------------------------------------
// TextIndent (WP2)
// ---------------------------------------------------------------------------

func TestExample_23_TextIndent(t *testing.T) {
	doc := template.New(
		template.WithPageSize(document.A4),
		template.WithMargins(document.UniformEdges(document.Mm(20))),
	)

	page := doc.AddPage()
	page.AutoRow(func(r *template.RowBuilder) {
		r.Col(12, func(c *template.ColBuilder) {
			c.Text("Text Indent Demo", template.FontSize(20), template.Bold())
			c.Spacer(document.Mm(8))

			c.Text("This paragraph has a 24pt first-line indent. "+
				"The first line starts further to the right, while subsequent "+
				"lines wrap at the normal left margin. This is commonly used "+
				"in book typography to indicate new paragraphs.",
				template.TextIndent(document.Pt(24)))
			c.Spacer(document.Mm(5))

			c.Text("This paragraph uses a larger 48pt indent for a more dramatic "+
				"effect. The indentation makes it easy to distinguish where a "+
				"new paragraph begins without adding extra vertical space.",
				template.TextIndent(document.Pt(48)))
			c.Spacer(document.Mm(5))

			c.Text("No indent on this paragraph for comparison. " +
				"Standard left-aligned text without any first-line indentation " +
				"starts flush with the left margin.")
		})
	})

	generatePDF(t, "23_text_indent.pdf", doc)
}

// ---------------------------------------------------------------------------
// TextDecoration (WP3)
// ---------------------------------------------------------------------------

func TestExample_24_TextDecoration(t *testing.T) {
	doc := template.New(
		template.WithPageSize(document.A4),
		template.WithMargins(document.UniformEdges(document.Mm(20))),
	)

	page := doc.AddPage()
	page.AutoRow(func(r *template.RowBuilder) {
		r.Col(12, func(c *template.ColBuilder) {
			c.Text("Text Decoration Demo", template.FontSize(20), template.Bold())
			c.Spacer(document.Mm(8))

			c.Text("Normal text without decoration")
			c.Spacer(document.Mm(4))

			c.Text("Underlined text for emphasis", template.Underline())
			c.Spacer(document.Mm(4))

			c.Text("Strikethrough text for deletions", template.Strikethrough())
			c.Spacer(document.Mm(4))

			c.Text("Combined underline and strikethrough",
				template.Underline(), template.Strikethrough())
			c.Spacer(document.Mm(4))

			c.Text("Colored underlined text",
				template.Underline(),
				template.TextColor(pdf.RGBHex(0x1565C0)),
				template.FontSize(14))
			c.Spacer(document.Mm(4))

			c.Text("Bold underlined heading",
				template.Bold(), template.Underline(), template.FontSize(16))
		})
	})

	generatePDF(t, "24_text_decoration.pdf", doc)
}

// ---------------------------------------------------------------------------
// VerticalAlign in Table Cells (WP4)
// ---------------------------------------------------------------------------

func TestExample_25_TableVerticalAlign(t *testing.T) {
	doc := template.New(
		template.WithPageSize(document.A4),
		template.WithMargins(document.UniformEdges(document.Mm(20))),
	)

	page := doc.AddPage()
	page.AutoRow(func(r *template.RowBuilder) {
		r.Col(12, func(c *template.ColBuilder) {
			c.Text("Table Vertical Align Demo", template.FontSize(20), template.Bold())
			c.Spacer(document.Mm(8))

			// Default (top) alignment
			c.Text("Default (Top) Alignment:", template.Bold())
			c.Spacer(document.Mm(3))
			c.Table(
				[]string{"Short", "Tall Cell"},
				[][]string{
					{"A", "This cell has\nmuch more content\nthat spans\nmultiple lines"},
					{"B", "Another tall\ncell with\nlong text"},
				},
				template.TableHeaderStyle(
					template.BgColor(pdf.RGBHex(0x1565C0)),
					template.TextColor(pdf.White),
				),
			)
			c.Spacer(document.Mm(8))

			// Middle alignment
			c.Text("Middle Alignment:", template.Bold())
			c.Spacer(document.Mm(3))
			c.Table(
				[]string{"Short", "Tall Cell"},
				[][]string{
					{"A", "This cell has\nmuch more content\nthat spans\nmultiple lines"},
					{"B", "Another tall\ncell with\nlong text"},
				},
				template.TableHeaderStyle(
					template.BgColor(pdf.RGBHex(0x2E7D32)),
					template.TextColor(pdf.White),
				),
				template.TableCellVAlign(document.VAlignMiddle),
			)
			c.Spacer(document.Mm(8))

			// Bottom alignment
			c.Text("Bottom Alignment:", template.Bold())
			c.Spacer(document.Mm(3))
			c.Table(
				[]string{"Short", "Tall Cell"},
				[][]string{
					{"A", "This cell has\nmuch more content\nthat spans\nmultiple lines"},
					{"B", "Another tall\ncell with\nlong text"},
				},
				template.TableHeaderStyle(
					template.BgColor(pdf.RGBHex(0xE65100)),
					template.TextColor(pdf.White),
				),
				template.TableCellVAlign(document.VAlignBottom),
			)
		})
	})

	generatePDF(t, "25_table_vertical_align.pdf", doc)
}

// ---------------------------------------------------------------------------
// PageNumber (WP5)
// ---------------------------------------------------------------------------

func TestExample_26_PageNumber(t *testing.T) {
	doc := template.New(
		template.WithPageSize(document.A4),
		template.WithMargins(document.UniformEdges(document.Mm(20))),
		template.WithMetadata(document.DocumentMetadata{
			Title:  "Page Number Demo",
			Author: "gpdf",
		}),
	)

	// Header with total pages
	doc.Header(func(p *template.PageBuilder) {
		p.AutoRow(func(r *template.RowBuilder) {
			r.Col(6, func(c *template.ColBuilder) {
				c.Text("Page Number Demo", template.Bold(), template.FontSize(10))
			})
			r.Col(6, func(c *template.ColBuilder) {
				c.TotalPages(template.AlignRight(), template.FontSize(9),
					template.TextColor(pdf.Gray(0.5)))
			})
		})
		p.AutoRow(func(r *template.RowBuilder) {
			r.Col(12, func(c *template.ColBuilder) {
				c.Line(template.LineColor(pdf.RGBHex(0x1565C0)))
				c.Spacer(document.Mm(3))
			})
		})
	})

	// Footer with page number
	doc.Footer(func(p *template.PageBuilder) {
		p.AutoRow(func(r *template.RowBuilder) {
			r.Col(12, func(c *template.ColBuilder) {
				c.Spacer(document.Mm(3))
				c.Line(template.LineColor(pdf.Gray(0.7)))
				c.Spacer(document.Mm(2))
			})
		})
		p.AutoRow(func(r *template.RowBuilder) {
			r.Col(6, func(c *template.ColBuilder) {
				c.Text("Generated by gpdf", template.FontSize(8),
					template.TextColor(pdf.Gray(0.5)))
			})
			r.Col(6, func(c *template.ColBuilder) {
				c.PageNumber(template.AlignRight(), template.FontSize(8),
					template.TextColor(pdf.Gray(0.5)))
			})
		})
	})

	// Create 4 pages of content
	for i := range 4 {
		page := doc.AddPage()
		page.AutoRow(func(r *template.RowBuilder) {
			r.Col(12, func(c *template.ColBuilder) {
				titles := []string{
					"Introduction",
					"Background",
					"Analysis",
					"Conclusion",
				}
				c.Text(titles[i], template.FontSize(18), template.Bold())
				c.Spacer(document.Mm(5))
				c.Text("This is the content of the page. The footer displays the " +
					"current page number, and the header shows the total number " +
					"of pages in the document. Both are automatically updated " +
					"after pagination, including on overflow pages.")
			})
		})
	}

	generatePDF(t, "26_page_number.pdf", doc)
}

// ===========================================================================
// Example functions for GoDoc (pkg.go.dev)
// ===========================================================================

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
