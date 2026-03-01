# gpdf

[![Go Reference](https://pkg.go.dev/badge/github.com/gpdf-dev/gpdf.svg)](https://pkg.go.dev/github.com/gpdf-dev/gpdf)
[![CI](https://github.com/gpdf-dev/gpdf/actions/workflows/check-code.yml/badge.svg)](https://github.com/gpdf-dev/gpdf/actions/workflows/check-code.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/gpdf-dev/gpdf)](https://goreportcard.com/report/github.com/gpdf-dev/gpdf)
[![Go Version](https://img.shields.io/badge/Go-%3E%3D1.22-blue)](https://go.dev/)

[English](README.md) | [日本語](README_ja.md) | [中文](README_zh.md) | **한국어** | [Español](README_es.md) | [Português](README_pt.md)

순수 Go로 구현된 외부 의존성 없는 PDF 생성 라이브러리. 레이어드 아키텍처와 선언적 빌더 API를 제공합니다.

## 특징

- **외부 의존성 제로** — Go 표준 라이브러리만 사용
- **레이어드 아키텍처** — 저수준 PDF 프리미티브, 문서 모델, 고수준 템플릿 API
- **12컬럼 그리드 시스템** — Bootstrap 스타일의 반응형 레이아웃
- **TrueType 폰트 지원** — 커스텀 폰트 임베딩 및 서브셋팅
- **CJK 지원** — 첫날부터 한중일 텍스트 완벽 지원
- **테이블** — 헤더, 컬럼 너비, 줄무늬 행
- **머리글 및 바닥글** — 모든 페이지에서 일관된 표시
- **다양한 단위** — pt, mm, cm, in, em, %
- **색상 공간** — RGB, 그레이스케일, CMYK
- **이미지** — JPEG 및 PNG 임베딩 (맞춤 옵션 지원)
- **문서 메타데이터** — 제목, 저자, 주제, 작성자

## 아키텍처

```
┌─────────────────────────────────────┐
│  gpdf (진입점)                       │
├─────────────────────────────────────┤
│  template  — 빌더 API, 그리드        │  레이어 3
├─────────────────────────────────────┤
│  document  — 노드, 스타일, 레이아웃    │  레이어 2
├─────────────────────────────────────┤
│  pdf       — Writer, 폰트, 스트림    │  레이어 1
└─────────────────────────────────────┘
```

## 요구사항

- Go 1.22 이상

## 설치

```bash
go get github.com/gpdf-dev/gpdf
```

## 빠른 시작

```go
package main

import (
	"os"

	"github.com/gpdf-dev/gpdf"
	"github.com/gpdf-dev/gpdf/document"
	"github.com/gpdf-dev/gpdf/template"
)

func main() {
	doc := gpdf.NewDocument(
		gpdf.WithPageSize(gpdf.A4),
		gpdf.WithMargins(document.UniformEdges(document.Mm(20))),
	)

	page := doc.AddPage()
	page.AutoRow(func(r *template.RowBuilder) {
		r.Col(12, func(c *template.ColBuilder) {
			c.Text("Hello, World!", template.FontSize(24), template.Bold())
		})
	})

	data, _ := doc.Generate()
	os.WriteFile("hello.pdf", data, 0644)
}
```

## 사용 예시

### 텍스트 스타일링

글꼴 크기, 두께, 스타일, 색상, 배경색, 정렬:

```go
page.AutoRow(func(r *template.RowBuilder) {
	r.Col(12, func(c *template.ColBuilder) {
		c.Text("큰 굵은 제목", template.FontSize(24), template.Bold())
		c.Text("이탤릭 텍스트", template.Italic())
		c.Text("굵게 + 이탤릭", template.Bold(), template.Italic())
		c.Text("빨간 텍스트", template.TextColor(pdf.Red))
		c.Text("커스텀 색상", template.TextColor(pdf.RGBHex(0x336699)))
		c.Text("배경색 포함", template.BgColor(pdf.Yellow))
		c.Text("가운데 정렬", template.AlignCenter())
		c.Text("오른쪽 정렬", template.AlignRight())
	})
})
```

### 12컬럼 그리드 레이아웃

Bootstrap 스타일의 12컬럼 그리드로 레이아웃 구성:

```go
// 2등분 컬럼
page.AutoRow(func(r *template.RowBuilder) {
	r.Col(6, func(c *template.ColBuilder) {
		c.Text("왼쪽 절반")
	})
	r.Col(6, func(c *template.ColBuilder) {
		c.Text("오른쪽 절반")
	})
})

// 사이드바 + 메인 콘텐츠
page.AutoRow(func(r *template.RowBuilder) {
	r.Col(3, func(c *template.ColBuilder) {
		c.Text("사이드바")
	})
	r.Col(9, func(c *template.ColBuilder) {
		c.Text("메인 콘텐츠")
	})
})
```

### 고정 높이 행

`Row()`로 높이를 지정하거나, `AutoRow()`로 콘텐츠에 맞춰 자동 조절:

```go
// 고정 높이: 30mm
page.Row(document.Mm(30), func(r *template.RowBuilder) {
	r.Col(12, func(c *template.ColBuilder) {
		c.Text("이 행의 높이는 30mm입니다")
	})
})
```

### 테이블

기본 테이블:

```go
c.Table(
	[]string{"상품명", "수량", "가격"},
	[][]string{
		{"위젯", "10", "₩5,000"},
		{"가젯", "3", "₩12,000"},
	},
)
```

스타일 적용 테이블 (헤더 색상, 컬럼 너비, 줄무늬 행):

```go
c.Table(
	[]string{"제품", "카테고리", "수량", "단가", "합계"},
	[][]string{
		{"노트북 Pro 15", "전자기기", "2", "₩1,299,000", "₩2,598,000"},
		{"무선 마우스", "주변기기", "10", "₩29,900", "₩299,000"},
	},
	template.ColumnWidths(30, 20, 10, 20, 20),
	template.TableHeaderStyle(
		template.TextColor(pdf.White),
		template.BgColor(pdf.RGBHex(0x1A237E)),
	),
	template.TableStripe(pdf.RGBHex(0xF5F5F5)),
)
```

### 이미지

JPEG 및 PNG 이미지 임베딩 (맞춤 옵션 지원):

```go
c.Image(imgData)                                      // 기본 크기
c.Image(imgData, template.FitWidth(document.Mm(80)))   // 너비에 맞춤
c.Image(imgData, template.FitHeight(document.Mm(30)))  // 높이에 맞춤
```

### 선 및 간격

```go
c.Line()                                           // 기본 (회색, 1pt)
c.Line(template.LineColor(pdf.Red))                 // 색상 지정
c.Line(template.LineThickness(document.Pt(3)))      // 굵은 선
c.Spacer(document.Mm(5))                            // 5mm 세로 간격
```

### 머리글 및 바닥글

모든 페이지에 반복 표시되는 머리글과 바닥글:

```go
doc.Header(func(p *template.PageBuilder) {
	p.AutoRow(func(r *template.RowBuilder) {
		r.Col(6, func(c *template.ColBuilder) {
			c.Text("ACME 주식회사", template.Bold(), template.FontSize(10))
		})
		r.Col(6, func(c *template.ColBuilder) {
			c.Text("대외비", template.AlignRight(), template.FontSize(10),
				template.TextColor(pdf.Gray(0.5)))
		})
	})
})

doc.Footer(func(p *template.PageBuilder) {
	p.AutoRow(func(r *template.RowBuilder) {
		r.Col(12, func(c *template.ColBuilder) {
			c.Text("gpdf로 생성", template.AlignCenter(),
				template.FontSize(8), template.TextColor(pdf.Gray(0.5)))
		})
	})
})
```

### 문서 메타데이터

```go
doc := gpdf.NewDocument(
	gpdf.WithPageSize(gpdf.A4),
	gpdf.WithMetadata(document.DocumentMetadata{
		Title:   "연간 보고서 2026",
		Author:  "gpdf Library",
		Subject: "문서 메타데이터 예시",
		Creator: "My Application",
	}),
)
```

### 페이지 크기 및 여백

```go
// 사용 가능한 페이지 크기
document.A4      // 210mm x 297mm
document.A3      // 297mm x 420mm
document.Letter  // 8.5in x 11in
document.Legal   // 8.5in x 14in

// 균일 여백
template.WithMargins(document.UniformEdges(document.Mm(20)))

// 비대칭 여백
template.WithMargins(document.Edges{
	Top:    document.Mm(10),
	Right:  document.Mm(40),
	Bottom: document.Mm(10),
	Left:   document.Mm(40),
})
```

### 출력 옵션

```go
// Generate는 []byte를 반환
data, err := doc.Generate()

// Render는 임의의 io.Writer에 기록
var buf bytes.Buffer
err := doc.Render(&buf)

// 파일에 직접 기록
f, _ := os.Create("output.pdf")
defer f.Close()
doc.Render(f)
```

## API 레퍼런스

### 문서 옵션

| 함수 | 설명 |
|---|---|
| `WithPageSize(size)` | 페이지 크기 설정 (A4, A3, Letter, Legal) |
| `WithMargins(edges)` | 페이지 여백 설정 |
| `WithFont(family, data)` | TrueType 폰트 등록 |
| `WithDefaultFont(family, size)` | 기본 폰트 설정 |
| `WithMetadata(meta)` | 문서 메타데이터 설정 |

### 컬럼 콘텐츠

| 메서드 | 설명 |
|---|---|
| `c.Text(text, opts...)` | 스타일 옵션을 포함한 텍스트 추가 |
| `c.Table(header, rows, opts...)` | 테이블 추가 |
| `c.Image(data, opts...)` | 이미지 추가 (JPEG/PNG) |
| `c.Line(opts...)` | 수평선 추가 |
| `c.Spacer(height)` | 수직 공간 추가 |

### 텍스트 옵션

| 옵션 | 설명 |
|---|---|
| `template.FontSize(size)` | 글꼴 크기를 포인트 단위로 설정 |
| `template.Bold()` | 굵게 |
| `template.Italic()` | 이탤릭 |
| `template.FontFamily(name)` | 등록된 폰트 사용 |
| `template.TextColor(color)` | 텍스트 색상 설정 |
| `template.BgColor(color)` | 배경 색상 설정 |
| `template.AlignLeft()` | 왼쪽 정렬 (기본값) |
| `template.AlignCenter()` | 가운데 정렬 |
| `template.AlignRight()` | 오른쪽 정렬 |

### 테이블 옵션

| 옵션 | 설명 |
|---|---|
| `template.ColumnWidths(w...)` | 컬럼 너비를 백분율로 설정 |
| `template.TableHeaderStyle(opts...)` | 헤더 행 스타일 설정 |
| `template.TableStripe(color)` | 교차 행 색상 설정 |

### 이미지 옵션

| 옵션 | 설명 |
|---|---|
| `template.FitWidth(value)` | 너비에 맞춰 스케일 (비율 유지) |
| `template.FitHeight(value)` | 높이에 맞춰 스케일 (비율 유지) |

### 선 옵션

| 옵션 | 설명 |
|---|---|
| `template.LineColor(color)` | 선 색상 설정 |
| `template.LineThickness(value)` | 선 두께 설정 |

### 단위

```go
document.Pt(72)    // 포인트 (1/72 인치)
document.Mm(10)    // 밀리미터
document.Cm(2.5)   // 센티미터
document.In(1)     // 인치
document.Em(1.5)   // 폰트 크기 기준 상대값
document.Pct(50)   // 퍼센트
```

### 색상

```go
pdf.RGB(0.2, 0.4, 0.8)   // RGB (0.0–1.0)
pdf.RGBHex(0xFF5733)      // 16진수 RGB
pdf.Gray(0.5)             // 그레이스케일
pdf.CMYK(0, 0.5, 1, 0)   // CMYK

// 사전 정의 색상
pdf.Black, pdf.White, pdf.Red, pdf.Green, pdf.Blue
pdf.Yellow, pdf.Cyan, pdf.Magenta
```

## 벤치마크

[go-pdf/fpdf](https://github.com/go-pdf/fpdf), [signintech/gopdf](https://github.com/signintech/gopdf), [maroto v2](https://github.com/johnfercher/maroto)와 비교.
5회 실행 중앙값, 각 100회 반복. Apple M1, Go 1.25.

**실행 시간** (낮을수록 좋음):

| 벤치마크 | gpdf | fpdf | gopdf | maroto v2 |
|---|--:|--:|--:|--:|
| 단일 페이지 | **13 µs** | 132 µs | 423 µs | 237 µs |
| 테이블 (4x10) | **108 µs** | 241 µs | 835 µs | 8.6 ms |
| 100페이지 | **683 µs** | 11.7 ms | 8.6 ms | 19.8 ms |
| 복합 문서 | **133 µs** | 254 µs | 997 µs | 10.4 ms |

**메모리 사용량** (낮을수록 좋음):

| 벤치마크 | gpdf | fpdf | gopdf | maroto v2 |
|---|--:|--:|--:|--:|
| 단일 페이지 | **16 KB** | 1.2 MB | 1.8 MB | 61 KB |
| 테이블 (4x10) | **209 KB** | 1.3 MB | 1.9 MB | 1.6 MB |
| 100페이지 | **909 KB** | 121 MB | 83 MB | 4.0 MB |
| 복합 문서 | **246 KB** | 1.3 MB | 2.0 MB | 2.0 MB |

### gpdf가 빠른 이유

- **단일 페이지** — 빌드→레이아웃→렌더링 단일 패스 파이프라인으로 중간 데이터 구조가 없음. 전체적으로 구체적인 구조체 타입 사용(`interface{}` 박싱 없음)으로 최소한의 힙 할당으로 문서 트리를 구축.
- **테이블** — 셀 내용을 재사용 가능한 `strings.Builder` 버퍼를 통해 PDF 콘텐츠 스트림 명령으로 직접 기록. 셀별 객체 래핑이나 반복적인 폰트 조회가 없으며, 폰트는 문서당 한 번만 해석.
- **100페이지** — 레이아웃이 O(n)으로 선형 확장. 오버플로우 페이지네이션이 슬라이스 참조로 나머지 노드를 전달(딥 카피 없음). 폰트는 한 번만 파싱되어 모든 페이지에서 공유.
- **복합 문서** — 재측정 없는 단일 패스 레이아웃이 위의 모든 장점을 통합. 폰트 서브셋팅은 실제 사용된 글리프만 임베딩하고, Flate 압축이 기본 적용되어 메모리와 출력 크기 모두 작게 유지.

벤치마크 실행:

```bash
cd _benchmark && go test -bench=. -benchmem -count=5
```

## 라이선스

MIT
