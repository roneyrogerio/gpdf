package gpdf_test

import (
	"testing"

	"github.com/gpdf-dev/gpdf/document"
	"github.com/gpdf-dev/gpdf/template"
)

func BenchmarkGenerateSinglePage(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		doc := template.New(
			template.WithPageSize(document.A4),
		)
		page := doc.AddPage()
		page.AutoRow(func(r *template.RowBuilder) {
			r.Col(12, func(c *template.ColBuilder) {
				c.Text("Hello, World!", template.FontSize(24))
			})
		})
		_, err := doc.Generate()
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkGenerateWithTable(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		doc := template.New(
			template.WithPageSize(document.A4),
		)
		page := doc.AddPage()
		page.AutoRow(func(r *template.RowBuilder) {
			r.Col(12, func(c *template.ColBuilder) {
				c.Table(
					[]string{"A", "B", "C", "D"},
					[][]string{
						{"1", "2", "3", "4"},
						{"5", "6", "7", "8"},
						{"9", "10", "11", "12"},
					},
				)
			})
		})
		_, err := doc.Generate()
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkGenerate100Pages(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		doc := template.New(
			template.WithPageSize(document.A4),
		)
		for j := 0; j < 100; j++ {
			page := doc.AddPage()
			page.AutoRow(func(r *template.RowBuilder) {
				r.Col(12, func(c *template.ColBuilder) {
					c.Text("Page content", template.FontSize(12))
				})
			})
		}
		_, err := doc.Generate()
		if err != nil {
			b.Fatal(err)
		}
	}
}
