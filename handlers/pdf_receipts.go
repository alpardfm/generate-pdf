package handlers

import (
	"bytes"
	"fmt"
	"net/http"
	"time"

	"github.com/jung-kurt/gofpdf"
)

/* ---------------------- struct.go ---------------------- */
type Paper struct {
	PaperSize            PaperSize
	MarginSetup          MarginSetup
	RectSetup            RectSetup
	TransformSetup       TransformSetup
	LineHt               float64
	TotalPaymentFont     FontSize
	ValueFont            FontSize
	FooterSetup          FooterSetup
	WLogo1               LogoSetup
	WLogo2               LogoSetup
	TransactionTextSetup TransactionTextSetup
	BottomSetup          BottomSetup
	ValueCellSetup       CellSetup
	HeaderSetup          HeaderSetup
	IconSetup            IconSetup
}

type PaperSize struct {
	Width  float64
	Height float64
}

type FontSize struct {
	ValueFontSize  float64
	HeaderFontSize float64
}

type MarginSetup struct {
	XMargin float64
	YMargin float64
}

type RectSetup struct {
	X      float64
	Y      float64
	W      float64
	H      float64
	InnerX float64
	InnerY float64
	InnerW float64
	InnerH float64
}

type TransformSetup struct {
	X struct {
		A float64
		B float64
	}
	Y struct {
		A float64
		B float64
	}
	TextX struct {
		A float64
		B float64
	}
	TextY struct {
		A float64
		B float64
	}
	Angle float64
	I     struct {
		Min float64
		Max float64
	}
	J struct {
		Min float64
		Max float64
	}
}

type LogoSetup struct {
	X float64
	Y float64
	W float64
	H float64
}

type FooterSetup struct {
	Y           float64
	RectHeight  float64
	WordSpacing float64
	FontSize    float64
}

type TransactionTextSetup struct {
	FontSize   float64
	UpperSpace float64
	LowerSpace float64
}

type BottomSetup struct {
	BottomLimit      float64
	BottomLimitMinus float64
	FontSize         float64
}

type CellSetup struct {
	W1         float64
	W2         float64
	WMultiCell float64
	H1         float64
	H2         float64
	HMultiCell float64
	Ln1        float64
	Ln2        float64
}

type HeaderSetup struct {
	Space1   float64
	Space2   float64
	W        float64
	H        float64
	X        float64
	Y        float64
	FontSize float64
}

type IconSetup struct {
	X float64
	Y float64
	W float64
	H float64
}

/* -------------------------------------------------------------------------- */

type PDFHandler struct {
	fontPath    string
	assetsPath  string
	templateDir string
}

func NewPDFHandler() *PDFHandler {
	return &PDFHandler{
		fontPath:    "assets/font",
		assetsPath:  "assets",
		templateDir: "assets/template",
	}
}

func (h *PDFHandler) GenerateDownloadReceiptPDF(w http.ResponseWriter, r *http.Request) {
	data := []Block{}

	data = append(data, Block{
		Fields: []Field{
			{
				Key:   "Tanggal Transaksi",
				Value: "Tipe Transaksi",
			},
			{
				Key:   "05/10/2025 09:10:33",
				Value: "Product Allocation",
			},
		},
		Type:      BLOCK_TYPE_ROWS,
		ShowTitle: false,

		// Override style per column to make first column 30% width and second //
		// column 70% width
		ColumnStyleOverrides: []ColumnStyleOverride{
			{
				ColIndex: 0,
				Width:    0.35,
			},
			{
				ColIndex: 1,
				Width:    0.65,
			},
		},
	})

	data = append(data, Block{
		Title: "Informasi Transaksi",
		Fields: []Field{
			{
				Key:   "Status SO",
				Value: "Release",
			},
			{
				Key:   "Nomor SO",
				Value: "4000012322",
			},
			{
				Key:   "ID Aplikasi",
				Value: "250214024675",
			},
			{
				Key:   "Nomor Perjanjian Penjadwalan",
				Value: "352424410",
			},
			{
				Key:   "Sumber Dana",
				Value: "Khris - Khrisna Joh",
			},
			{
				Key:   "",
				Value: "PT. BANK RAKYAT INDONESIA (PERSERO) TBK - BRINIDJA",
			},
			{
				Key:   "",
				Value: "1001******890",
			},
			{
				Key:   "Organisasi Penjualan",
				Value: "007-C&T LPG Retail",
			},
			{
				Key:   "Grup Produk",
				Value: "001-LPG/BBG",
			},
			{
				Key:   "Tujuan Pengiriman",
				Value: "100123 - PT. Raya Jaya Mulya",
			},
			{
				Key:   "Pembeli",
				Value: "966787 - PT. Makmur Sentosa",
			},
			{
				Key:   "Depo",
				Value: "2150 - SPBE Wanantara D. Satria",
			},
			{
				Key:   "Pembayar",
				Value: "966787 - PT. Makmur Sentosa",
			},
		},
		Type:      BLOCK_TYPE_ROWS,
		ShowTitle: true,

		// Override style per column to make first column 30% width and second //
		// column 70% width
		ColumnStyleOverrides: []ColumnStyleOverride{
			{
				ColIndex: 0,
				Width:    0.35,
			},
			{
				ColIndex: 1,
				Width:    0.65,
			},
		},
	})

	data = append(data, Block{
		Title:     "Rincian Transaksi",
		Type:      BLOCK_TYPE_ROWS,
		ShowTitle: true,
		Fields: []Field{
			{
				Key:   "Nilai Bersih",
				Value: "IDR 16.000.000,00",
			},
			{
				Key:   "Pajak PPN",
				Value: "IDR 2.086.000,00",
			},
			{
				Key:   "Pajak PBBKB",
				Value: "IDR 0,00",
			},
			{
				Key:   "Pajak PPH",
				Value: "IDR 0,00",
			},
			{
				Key:   "Nilai Kotor",
				Value: "IDR 18.068.000,00",
			},
			{
				Key:   "Nilai Debit atau Kredit",
				Value: "IDR 1.000.000,00",
			},
			{
				Key:   "Biaya Admin",
				Value: "IDR 2.000,00",
			},
			{
				Key:     "Total Pembayaran",
				Value:   "IDR 18.071.000,00",
				IsTotal: true,
			},
		},
		StartFrom: 0.35,
	})

	data = append(data, Block{
		Title:     "Detail Produk",
		ShowTitle: true,
		Type:      BLOCK_TYPE_TABLE,
		TableData: TableData{
			Rows: [][]string{
				{"Material", "Deskripsi Material", "Trip", "Qty", "UoM", "Transporter", "Tanggal Kirim"},
				{"A040900001", "LPG BR1 3KG", "1", "1000", "KG", "PT. Rahayu Sentosa", "05/10/2025"},
				{"A040900002", "LPG BR1 3KG", "1", "1000", "KG", "PT. Rahayu Sentosa", "05/10/2025"},
				{"A040900003", "LPG BR1 3KG", "1", "1000", "KG", "PT. Rahayu Sentosa", "05/10/2025"},
				{"A040900004", "LPG BR1 3KG", "1", "1000", "KG", "PT. Rahayu Sentosa", "05/10/2025"},
				{"A040900005", "LPG BR1 3KG", "1", "1000", "KG", "PT. Rahayu Sentosa", "05/10/2025"},
				{"Total", "", "", "5000", "", "", ""},
			},
			ColSize:     []float64{0.15, 0.20, 0.10, 0.10, 0.10, 0.15, 0.20},
			LastRowBold: true,
		},
	})

	var buf bytes.Buffer
	pdf, err := h.GenerateReceipt("DOP1234567890", "success", data)
	if err != nil {
		fmt.Println("Error generating report:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = pdf.Output(&buf)
	if err != nil {
		fmt.Println("Error outputting PDF:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Set header untuk PDF preview
	w.Header().Set("Content-Type", "application/pdf")
	w.Header().Set("Content-Length", fmt.Sprintf("%d", len(buf.Bytes())))

	// Optional: beri nama file untuk context
	w.Header().Set("Content-Disposition", "inline; filename=\"preview.pdf\"")

	// Write PDF binary langsung ke response
	w.Write(buf.Bytes())
}

/*
 * BlockType adalah tipe dari blok yang akan dirender di dalam PDF.
 */

const BLOCK_TYPE_ROWS = "rows"
const BLOCK_TYPE_TABLE = "table"

/*
 * Another specific struct for block type. This time for column style override.
 * You can change for example font size, boldness, text align per column. Just
 * specify the column index (0-based). For example: If you want to change the
 * style of the second column to have bold text and center alignment, you can
 * create a ColumnStyleOverride like this:
 *
 * ColumnStyleOverride{
 *     ColIndex:   1,
 *     FontSize:   12,
 *     FontWeight: "bold",
 *     TextAlign:  "CM",
 *     Margin:     2,
 * }
 */
type ColumnStyleOverride struct {
	ColIndex   int
	Width      float64
	FontSize   float64
	FontWeight string
	TextAlign  string
	Margin     float64
}

type RowStyleOverride struct {
	RowIndex   int
	FontSize   float64
	FontWeight string
	TextAlign  string
	Margin     float64
}

type Field struct {
	Key     string
	Value   interface{}
	IsTotal bool
}

type TableData struct {
	Headers []string
	Rows    [][]string
	ColSize []float64 // Column size in percentage (0.0 - 1.0). 1.0 = 100% width
	// of table
	LastRowBold bool
}

/*
 * Kalau kita cermati sebetulnya data setiap bagan itu mirip, hanya saja
 * stylenya yang berbeda, ada yang table dan ada yang berbentuk list. Namun
 * secara data, mereka semua memiliki key, value dan title. Namun di sini
 * kita akan membuat field "type" untuk membedakan apakah itu grid atau list
 * atau rows.
 */
type Block struct {
	Title                string
	ShowTitle            bool
	Fields               []Field
	TableData            TableData
	Type                 string
	ColumnStyleOverrides []ColumnStyleOverride

	/*
	 * Special case for a row that needs to start from specific X position. To
	 * make consistent unit, we're using float64 between 0.0 to 1.0, where 0.0
	 * is the very left and 1.0 is the very right.
	 */
	StartFrom float64
}

/*
 * GenerateReceipt generates a PDF receipt for a given transaction.
 *
 * @param transactionId The ID of the transaction.
 * @param status. TODO: use proper enum or type for status.
 * @param Data[] the data to be included in the receipt.
 */
func (s *PDFHandler) GenerateReceipt(transactionId string, status string, data []Block) (*gofpdf.Fpdf, error) {
	var paper Paper = GetPaperA4()

	pdf := gofpdf.New("P", "mm", "A4", "assets/font")

	pdf.AddUTF8Font("BRIDigital-Light", "", "BRIDigitalText-Light.ttf")
	pdf.AddFont("BRIDigital", "", "BRIDigitalText-Regular.json")
	pdf.AddFont("BRIDigital", "B", "BRIDigitalText-SemiBold.json")
	pdf.AddFont("BRIDigitalLogo", "B", "BRIDigitalText-SemiBold.json")

	pdf.SetMargins(paper.MarginSetup.XMargin, paper.MarginSetup.YMargin, paper.MarginSetup.XMargin)
	pdf.SetAutoPageBreak(true, 10)

	pdf.SetHeaderFunc(func() {
		pdf.ImageOptions("./assets/images/receipt-header.png", 0, 0, 210, 0, false, gofpdf.ImageOptions{
			ReadDpi:   false,
			ImageType: "",
		}, 0, "")

		drawBackgroundRounded(pdf, paper)

		s.addWatermark(pdf, paper)
	})

	pdf.AliasNbPages("{nb}")

	pdf.SetFooterFunc(func() {
		currentTime := time.Now()
		formattedTime := currentTime.Format("02/01/2006 15:04:05")

		pdf.SetY(paper.FooterSetup.Y)
		pdf.SetFont("BRIDigital", "", paper.FooterSetup.FontSize)
		pdf.SetTextColor(128, 128, 128)
		pdf.CellFormat(0, 10, fmt.Sprintf("%s - Halaman %d/{nb}", formattedTime, pdf.PageNo()), "", 0, "R", false, 0, "")

		pdf.Ln(5)
		pdf.SetFillColor(16, 47, 50)
		pdf.Rect(0, paper.PaperSize.Height-paper.FooterSetup.RectHeight, paper.PaperSize.Width, paper.FooterSetup.RectHeight, "F")
		pdf.Ln(2)
	})

	pdf.AddPage()
	drawTransactionDetails(pdf, paper, transactionId)
	drawContentsReceipt(pdf, paper, data)

	pdf.CellFormat(paper.RectSetup.InnerX, 12, "Dokumen ini merupakan bukti transaksi yang sah dan dicetak otomatis oleh sistem.", "", 0, "L", false, 0, "")

	pdf.SetFont("BRIDigital", "B", 9)

	fontSize, _ := pdf.GetFontSize()
	footerContainerHeight := fontSize * 2

	pdf.SetY((paper.RectSetup.Y + paper.RectSetup.H) - footerContainerHeight - 4)

	pdf.SetAlpha(0.16, "Normal")
	pdf.SetFillColor(61, 134, 143)
	pdf.RoundedRect(paper.RectSetup.InnerX, pdf.GetY(), paper.RectSetup.InnerW, footerContainerHeight, 3, "1234", "F")
	pdf.SetAlpha(1, "Normal")

	pdf.SetTextColor(107, 104, 128)

	pdf.SetX(paper.RectSetup.InnerX)
	pdf.MultiCell(paper.RectSetup.InnerW, fontSize, "Dokumen ini merupakan bukti transaksi yang sah dan dicetak otomatis oleh sistem. Terima kasih telah bertransaksi menggunakan Qlola BRI, bila menemui kendala silakan hubungi kami di 500001 atau qlola@bri.co.id", "", "C", false)

	return pdf, nil
}

func drawBackgroundRounded(pdf *gofpdf.Fpdf, paper Paper) {
	borderSize := 0.5

	pdf.SetFillColor(224, 224, 224)
	pdf.RoundedRect(paper.RectSetup.X-borderSize, paper.RectSetup.Y-borderSize, paper.RectSetup.W+(2*borderSize), paper.RectSetup.H+(2*borderSize), 3, "1234", "F")

	pdf.SetFillColor(255, 255, 255)
	pdf.RoundedRect(paper.RectSetup.X, paper.RectSetup.Y, paper.RectSetup.W, paper.RectSetup.H, 3, "1234", "F")
}

func drawTransactionDetails(pdf *gofpdf.Fpdf, paper Paper, transactionId string) {
	successIconPath := "assets/images/Icon-2.png"
	pdf.ImageOptions(successIconPath, paper.IconSetup.X, paper.IconSetup.Y, paper.IconSetup.W, paper.IconSetup.H, false, gofpdf.ImageOptions{ImageType: "PNG"}, 0, "")

	lineHeight := 5.0
	containerStartX := paper.RectSetup.InnerX + paper.IconSetup.W + (paper.MarginSetup.XMargin / 2)
	containerBottomY := paper.RectSetup.InnerY + paper.IconSetup.H + paper.MarginSetup.YMargin

	pdf.SetFont("BRIDigital", "B", 13)
	pdf.SetTextColor(24, 24, 24)
	pdf.Text(containerStartX, paper.RectSetup.InnerY+6, "Transaksi Sukses")
	pdf.SetFont("BRIDigital-Light", "", 10)
	pdf.Text(containerStartX, paper.RectSetup.InnerY+7+lineHeight, transactionId)

	pdf.SetFont("BRIDigital", "B", 10)
	strWidth := pdf.GetStringWidth("DO Pertamina")
	pdf.Text((paper.RectSetup.InnerW+paper.RectSetup.InnerX)-strWidth, paper.RectSetup.InnerY+10, "DO Pertamina")

	drawDashedLine(pdf, paper.RectSetup.InnerX, containerBottomY-2, (paper.RectSetup.InnerW + paper.RectSetup.InnerX), containerBottomY-2)
}

func drawDashedLine(pdf *gofpdf.Fpdf, x1, y1, x2, y2 float64) {
	pdf.SetDrawColor(224, 224, 224)
	pdf.SetLineWidth(0.5)

	// Pattern dash sederhana: [panjang_dash, panjang_gap]
	pdf.SetDashPattern([]float64{2, 2}, 0)
	pdf.Line(x1, y1, x2, y2)
	pdf.SetDashPattern([]float64{}, 0) // Reset ke solid
}

func drawContentsReceipt(pdf *gofpdf.Fpdf, paper Paper, data []Block) {
	pdf.SetY(paper.RectSetup.InnerY + paper.IconSetup.H + paper.MarginSetup.YMargin)

	// Loop through each block and render based on its type
	for _, block := range data {
		if block.ShowTitle {
			pdf.SetX(paper.RectSetup.InnerX)
			startX, startY := pdf.GetX(), pdf.GetY()

			// 1. Draw semi-transparent background rectangle
			pdf.SetAlpha(0.16, "Normal")
			pdf.SetFillColor(61, 136, 143)
			pdf.Rect(startX, startY, paper.RectSetup.InnerW, paper.HeaderSetup.H, "F")
			pdf.SetAlpha(1, "Normal")

			// 2. Draw text on top (fully opaque)
			pdf.SetFont("BRIDigital", "B", 10)
			pdf.SetTextColor(0, 0, 0)

			// 3. Calculate text positioning (left-aligned with some padding)
			textX := startX + 1 // 2mm padding from left
			textY := startY

			pdf.SetXY(textX, textY)
			pdf.CellFormat(paper.RectSetup.InnerW-8, paper.HeaderSetup.H, block.Title,
				"", 0, "L", false, 0, "") // fill = false

			// 3. Move to next position after the block
			pdf.SetXY(startX, startY+paper.HeaderSetup.H)
		}

		switch block.Type {
		case BLOCK_TYPE_ROWS:
			drawBlockRows(pdf, paper, block)
		case BLOCK_TYPE_TABLE:
			drawBlockTable(pdf, paper, block)
		}
	}
}

func drawBlockRows(pdf *gofpdf.Fpdf, paper Paper, block Block) {
	if block.StartFrom > 0 {
		pdf.SetX(paper.RectSetup.InnerX + (paper.RectSetup.InnerW * block.StartFrom))
	} else {
		pdf.SetX(paper.RectSetup.InnerX)
	}

	columnPersentage := []int{50, 50}

	if len(block.ColumnStyleOverrides) > 0 {
		for _, colStyle := range block.ColumnStyleOverrides {
			switch colStyle.ColIndex {
			case 0:
				columnPersentage[0] = int(colStyle.Width * 100)
			case 1:
				columnPersentage[1] = int(colStyle.Width * 100)
			}
		}
	}

	// Calculate column widths based on percentage
	totalWidth := paper.RectSetup.InnerW - (paper.RectSetup.InnerW * block.StartFrom)
	colWidths := []float64{
		(totalWidth * float64(columnPersentage[0])) / 100,
		(totalWidth * float64(columnPersentage[1])) / 100,
	}

	/*
	 * Loop through each field and render key-value pairs in two columns.
	 * Pair always have 2 columns. So it's okay to hard code the column
	 * override variable.
	 */
	for _, field := range block.Fields {
		startX, startY := pdf.GetX(), pdf.GetY()

		if field.IsTotal {
			pdf.Ln(1)
			drawDashedLine(pdf, startX, pdf.GetY(), startX+colWidths[0]+colWidths[1], pdf.GetY())

			pdf.SetFont("BRIDigital", "B", 10)
			pdf.SetTextColor(0, 0, 0)
			pdf.SetXY(startX, startY)
			pdf.CellFormat(colWidths[0], 10, fmt.Sprintf("%v", field.Key), "", 0, "L", false, 0, "")
			// Draw Value
			pdf.SetXY(startX+colWidths[0], startY+1)
			pdf.CellFormat(colWidths[1], 10, fmt.Sprintf("%v", field.Value), "", 0, "R", false, 0, "")
			pdf.Ln(7)

			continue
		}

		// Now we can draw the key and value in their respective columns. Loop
		// through each column.
		pdf.SetFont("BRIDigital-Light", "", 10)
		pdf.SetTextColor(107, 104, 128)

		// Draw Key
		pdf.SetXY(startX, startY)
		pdf.CellFormat(colWidths[0], 7, fmt.Sprintf("%v", field.Key), "", 0, "L", false, 0, "")
		// Draw Value
		pdf.SetXY(startX+colWidths[0], startY)
		pdf.SetTextColor(0, 0, 0)
		if block.Title == "Rincian Transaksi" {
			pdf.CellFormat(colWidths[1], 7, fmt.Sprintf("%v", field.Value), "", 0, "R", false, 0, "")
		} else {
			pdf.CellFormat(colWidths[1], 7, fmt.Sprintf("%v", field.Value), "", 0, "L", false, 0, "")
		}
		pdf.Ln(7)

		// Set X back to left margin for next row
		if block.StartFrom > 0 {
			pdf.SetX(paper.RectSetup.InnerX + (paper.RectSetup.InnerW * block.StartFrom))
		} else {
			pdf.SetX(paper.RectSetup.InnerX)
		}
	}

	// Set Y to current Y + gap for the end of the block
	pdf.SetY(pdf.GetY() + 4)
}

/*
 * Draw table with headers and rows where the top left, top right,
 * bottom left, bottom right corners are rounded. Make sure that if
 * the table is too long, it will break into multiple pages while
 * and recreate the header on the new page.
 */
func drawBlockTable(pdf *gofpdf.Fpdf, paper Paper, block Block) {
	pdf.SetY(pdf.GetY() + 4)

	pdf.SetFont("BRIDigital", "B", 10)
	pdf.SetTextColor(0, 0, 0)
	borderRadius := 3.0

	// Draw header for first page
	drawTableHeader(pdf, paper, block, false)

	for rowIdx, tableRow := range block.TableData.Rows {
		// Skip header row (index 0) since we already drew it
		if rowIdx == 0 {
			continue
		}

		// Calculate row height
		rowHeight := calculateRowHeight(pdf, tableRow, block, paper)

		// Check if we need a new page. 16 is the header height
		// plus some padding.
		rectHeight := paper.RectSetup.InnerH + 16

		if pdf.GetY() > rectHeight {
			pdf.AddPage()
			pdf.SetY(paper.RectSetup.InnerY)
			drawTableHeader(pdf, paper, block, true)
		}

		// Draw the row
		drawTableRow(pdf, paper, block, tableRow, rowIdx, len(block.TableData.Rows), rowHeight, borderRadius)
	}
}

// Helper function to calculate row height
func calculateRowHeight(pdf *gofpdf.Fpdf, row []string, block Block, paper Paper) float64 {
	_, lineHt := pdf.GetFontSize()
	maxHeight := lineHt + 4 // Minimum height

	for colIdx, txt := range row {
		width := paper.RectSetup.InnerW * block.TableData.ColSize[colIdx]
		lines := pdf.SplitLines([]byte(txt), width)
		cellHeight := float64(len(lines))*(lineHt+2.0) + 4.0
		if cellHeight > maxHeight {
			maxHeight = cellHeight
		}
	}

	return maxHeight
}

// Helper function to draw a single table row
func drawTableRow(pdf *gofpdf.Fpdf, paper Paper, block Block, tableRow []string, rowIdx, totalRows int, height, borderRadius float64) {
	pdf.SetX(paper.RectSetup.InnerX)

	startX, startY := pdf.GetX(), pdf.GetY()
	x := paper.RectSetup.InnerX

	for colIdx, txt := range tableRow {
		width := paper.RectSetup.InnerW * block.TableData.ColSize[colIdx]

		// Determine border style based on position
		if rowIdx == totalRows-1 { // Last row
			if colIdx == 0 {
				pdf.RoundedRect(x, startY, width, height, borderRadius, "4", "D")
			} else if colIdx == len(tableRow)-1 {
				pdf.RoundedRect(x, startY, width, height, borderRadius, "3", "D")
			} else {
				pdf.Rect(x, startY, width, height, "D")
			}
		} else { // Middle rows
			pdf.Rect(x, startY, width, height, "D")
		}

		padding := 2.0

		innerX := x + padding
		innerW := width - padding*2
		pdf.SetXY(innerX, pdf.GetY()+padding)

		// Draw text
		_, unitSize := pdf.GetFontSize()

		if block.TableData.LastRowBold && rowIdx == totalRows-1 {
			pdf.SetFont("BRIDigital", "B", 10)
		} else {
			pdf.SetFont("BRIDigital", "", 10)
		}

		pdf.MultiCell(innerW, unitSize+2.0, txt, "", "L", false)
		x += width
		pdf.SetXY(x, startY)
	}

	// Move to next row position
	pdf.SetXY(startX, startY+height)
}

// Updated header function with page break awareness
func drawTableHeader(pdf *gofpdf.Fpdf, paper Paper, block Block, isNewTab bool) {
	// Set position after adding new page
	if isNewTab {
		pdf.SetY(paper.RectSetup.InnerY)
	}

	pdf.SetFont("BRIDigital", "B", 10)
	pdf.SetTextColor(0, 0, 0)
	pdf.SetFillColor(249, 249, 249)
	borderRadius := 3.0

	pdf.SetX(paper.RectSetup.InnerX)
	startX, startY := pdf.GetX(), pdf.GetY()
	x := paper.RectSetup.InnerX

	// Calculate header height (find the tallest cell)
	headerHeight := 10.0
	_, lineHt := pdf.GetFontSize()
	for col, txt := range block.TableData.Rows[0] {
		width := paper.RectSetup.InnerW * block.TableData.ColSize[col]
		lines := pdf.SplitLines([]byte(txt), width)
		cellHeight := float64(len(lines))*(lineHt+2.0) + 4.0
		if cellHeight > headerHeight {
			headerHeight = cellHeight
		}
	}

	for col, txt := range block.TableData.Rows[0] {
		width := paper.RectSetup.InnerW * block.TableData.ColSize[col]

		// Draw rounded borders for header
		if col == 0 {
			pdf.RoundedRect(x, startY, width, headerHeight, borderRadius, "1", "FD")
		} else if col == len(block.TableData.Rows[0])-1 {
			pdf.RoundedRect(x, startY, width, headerHeight, borderRadius, "2", "FD")
		} else {
			pdf.RoundedRect(x, startY, width, headerHeight, borderRadius, "TB", "FD")
		}

		padding := 2.0

		innerX := x + padding
		innerW := width - padding*2
		pdf.SetXY(innerX, pdf.GetY())

		// Draw header text
		pdf.MultiCell(innerW, lineHt+2.0, txt, "", "AL", false)
		x += width
		pdf.SetXY(x, startY)
	}

	// Move to position for first data row
	pdf.SetXY(startX, startY+headerHeight)
}

func GetPaperA4() Paper {
	var paper = Paper{
		PaperSize: PaperSize{
			Width:  210,
			Height: 297,
		},
		MarginSetup: MarginSetup{
			XMargin: 7.9375,
			YMargin: 7.9375,
		},
		TransformSetup: TransformSetup{
			X: struct {
				A float64
				B float64
			}{A: 7.9375, B: 30},
			Y: struct {
				A float64
				B float64
			}{A: 15, B: 10.7},
			TextX: struct {
				A float64
				B float64
			}{A: 7.9375, B: 30},
			TextY: struct {
				A float64
				B float64
			}{A: 15, B: 10.7},
			Angle: 30,
			I: struct {
				Min float64
				Max float64
			}{
				Min: 0.04,
				Max: 10,
			},
			J: struct {
				Min float64
				Max float64
			}{
				Min: 0.9,
				Max: 26.75,
			},
		},
		LineHt: 5.5,
		TotalPaymentFont: FontSize{
			ValueFontSize:  17,
			HeaderFontSize: 17,
		},
		ValueFont: FontSize{
			ValueFontSize:  15,
			HeaderFontSize: 15,
		},
		FooterSetup: FooterSetup{
			Y:           -16.5,
			RectHeight:  4.5, // Green Rect
			WordSpacing: 1,
			FontSize:    10,
		},
		TransactionTextSetup: TransactionTextSetup{
			FontSize:   15,
			UpperSpace: 10,
			LowerSpace: 25.4,
		},
		BottomSetup: BottomSetup{
			BottomLimit:      45,
			BottomLimitMinus: 40,
			FontSize:         15,
		},
		ValueCellSetup: CellSetup{
			W1:         10,
			W2:         1,
			WMultiCell: 0,
			H2:         5.5,
			H1:         5.5,
			HMultiCell: 5.5,
			Ln1:        2,
			Ln2:        2,
		},
		HeaderSetup: HeaderSetup{
			Space1:   1.5,
			Space2:   13.2,
			W:        0,
			H:        7,
			X:        30.1,
			Y:        5,
			FontSize: 12,
		},
	}

	yAfterHeaderImage := 15.9375
	paper.RectSetup = RectSetup{
		X: paper.MarginSetup.XMargin,
		Y: paper.MarginSetup.YMargin + yAfterHeaderImage,
		W: paper.PaperSize.Width - (paper.MarginSetup.XMargin * 2),
		H: paper.PaperSize.Height - (paper.FooterSetup.RectHeight + paper.MarginSetup.YMargin + yAfterHeaderImage + paper.MarginSetup.YMargin*2),

		/* Start drawing from here */
		InnerX: paper.MarginSetup.XMargin + (paper.MarginSetup.XMargin / 2),
		InnerY: paper.MarginSetup.YMargin + yAfterHeaderImage + (paper.MarginSetup.YMargin / 2),
		InnerW: paper.PaperSize.Width - (paper.MarginSetup.XMargin * 2) - (paper.MarginSetup.XMargin / 2 * 2),
		InnerH: paper.PaperSize.Height - (paper.FooterSetup.RectHeight + paper.MarginSetup.YMargin + yAfterHeaderImage + paper.MarginSetup.YMargin*2),
	}

	/* Success Icon */
	paper.IconSetup = IconSetup{
		X: paper.RectSetup.InnerX,
		Y: paper.RectSetup.InnerY,
		W: 14,
		H: 14,
	}

	return paper
}
