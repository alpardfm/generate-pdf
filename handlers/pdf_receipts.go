package handlers

import (
	"bytes"
	"fmt"
	"net/http"
	"time"

	"github.com/jung-kurt/gofpdf"
)

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
	data := map[string]interface{}{
		"transactionId":       "DOP1234567890",
		"status":              "Sukses",
		"debitAccountName":    "Khrisna Joh",
		"accountNumber":       "00123456789",
		"moduleType":          "Product Allocation",
		"totalPayment":        "IDR 17.070.000,00",
		"netValue":            "IDR 16.000.000,00",
		"ppnTax":              "IDR 2.086.000,00",
		"pbbkbTax":            "IDR 0,00",
		"pphTax":              "IDR 0,00",
		"grossValue":          "IDR 18.068.000,00",
		"debitCreditValue":    "IDR 1.000.000,00",
		"totalAmount":         "IDR 17.068.000,00",
		"soldToName":          "PT. Pertamina Tbk",
		"buyer":               "PT. Perusahaan Customer",
		"depoName":            "Depo Jakarta Utara",
		"salesOrganization":   "SO002",
		"productGroup":        "Bahan Bakar Minyak",
		"distributionChannel": "Corporate Sales",
		"payer":               "PT. Perusahaan Customer",
	}

	var buf bytes.Buffer

	pdf, err := h.GenerateReceipt(data)
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

func (s *PDFHandler) GenerateReceipt(data interface{}) (*gofpdf.Fpdf, error) {
	pdf := gofpdf.New("P", "mm", "A4", "assets/font")
	pdf.SetAutoPageBreak(true, 45)

	pdf.AddFont("BRIDigital", "", "BRIDigitalText-Regular.json")
	pdf.AddFont("BRIDigital", "B", "BRIDigitalText-SemiBold.json")
	pdf.AddFont("BRIDigitalLogo", "B", "BRIDigitalText-SemiBold.json")

	pdf.SetHeaderFunc(
		func() {
			s.addHeader(pdf)
			drawBox(pdf, 10, 35, 297-20, 420-55)
		})
	pdf.AddPageFormat("P", gofpdf.SizeType{Wd: 297, Ht: 420})
	// Add header
	// Add watermark background

	fontSize := 12.0
	fontName := "BRIDigital"
	// heightRow := 10.0
	// heightHeader := 15.0
	pageW, pageHeight := pdf.GetPageSize()
	pdf.SetFont(fontName, "", fontSize)
	// paddingDetailTrx := 20.0

	pdf.SetFont(fontName, "B", 16)

	pdf.AliasNbPages("{nb}")

	pdf.SetFooterFunc(func() {
		pdf.SetY(-19)
		pdf.SetX(225)

		currentTime := time.Now()
		formattedTime := currentTime.Format("02/01/2006 15:04:05")

		// Text kanan - juga pakai MultiCell untuk alignment yang konsisten
		pdf.MultiCell(pageW/2, 7, fmt.Sprintf("%v - Halaman %d/{nb}", formattedTime, pdf.PageNo()), "", "L", false)
	})

	// Add content
	s.addContentReceipt(pdf, data)

	pdf.SetX(15)
	pdf.CellFormat(pageW-30, 12, "Dokumen ini merupakan bukti transaksi yang sah dan dicetak otomatis oleh sistem.", "0", 0, "L", false, 0, "")

	// pdf.Ln(10)

	pdf.SetAutoPageBreak(false, 0)

	drawDashedLine(pdf, 15, pageHeight-45, pageW-16, pageHeight-45)

	// pdf.Ln(10)

	pdf.SetFillColor(235, 244, 245) // Warna background
	pdf.SetTextColor(0, 0, 0)
	pdf.SetXY(15, pageHeight-40)
	pdf.CellFormat(pageW-30, 15, "Terima kasih telah bertransaksi menggunakan Qlola BRI, bila menemui kendala silakan hubungi kami di 500001 atau qlola@bri.co.id", "0", 0, "C", true, 0, "")

	// // Add footer
	// s.addFooter(pdf)

	return pdf, nil
}

func drawBox(pdf *gofpdf.Fpdf, x, y, width, height float64) {
	radius := 5.0 // Radius untuk rounded corners

	// 1. Background PUTIH dengan rounded corners
	pdf.SetFillColor(255, 255, 255)

	// Rectangle utama (tanpa area corners)
	pdf.Rect(x+radius, y, width-(2*radius), height, "F") // Middle horizontal
	pdf.Rect(x, y+radius, width, height-(2*radius), "F") // Middle vertical

	// 2. Rounded corners (pakai ellipse yang bersih)
	pdf.Ellipse(x+radius, y+radius, radius, radius, 0, "F")              // Top-left
	pdf.Ellipse(x+width-radius, y+radius, radius, radius, 0, "F")        // Top-right
	pdf.Ellipse(x+width-radius, y+height-radius, radius, radius, 0, "F") // Bottom-right
	pdf.Ellipse(x+radius, y+height-radius, radius, radius, 0, "F")       // Bottom-left

	// 3. Border dengan rounded corners (pakai line yang bersih)
	pdf.SetDrawColor(0, 0, 0)

	// Garis lurus saja, tanpa corner functions yang bikin segitiga
	pdf.Line(x+radius, y, x+width-radius, y)               // Top
	pdf.Line(x+width, y+radius, x+width, y+height-radius)  // Right
	pdf.Line(x+radius, y+height, x+width-radius, y+height) // Bottom
	pdf.Line(x, y+radius, x, y+height-radius)              // Left

	// 4. Watermark
	pdf.ImageOptions("./assets/images/newwatermark3.png", x, y, width, height, false, gofpdf.ImageOptions{
		ReadDpi:   false,
		ImageType: "",
	}, 0, "")
}

// Function untuk membuat rounded rectangle
func drawRoundedRect(pdf *gofpdf.Fpdf, x, y, w, h, r float64, style string) {
	// Style: "F" untuk fill, "D" untuk draw, "FD" untuk both

	// Jika radius terlalu besar, adjust
	if r > w/2 {
		r = w / 2
	}
	if r > h/2 {
		r = h / 2
	}

	// Draw the four corners and sides
	pdf.Curve(x+r, y, x, y, x, y+r, style)             // Top-left corner
	pdf.Curve(x+w-r, y, x+w, y, x+w, y+r, style)       // Top-right corner
	pdf.Curve(x+w, y+h-r, x+w, y+h, x+w-r, y+h, style) // Bottom-right corner
	pdf.Curve(x+r, y+h, x, y+h, x, y+h-r, style)       // Bottom-left corner

	// Fill the center rectangle if needed
	if style == "F" || style == "FD" {
		pdf.Rect(x+r, y, w-2*r, h, "F") // Top and bottom rectangles
		pdf.Rect(x, y+r, w, h-2*r, "F") // Middle rectangle
	}

	// Draw the straight sides if needed
	if style == "D" || style == "FD" {
		pdf.Line(x+r, y, x+w-r, y)     // Top side
		pdf.Line(x+w, y+r, x+w, y+h-r) // Right side
		pdf.Line(x+r, y+h, x+w-r, y+h) // Bottom side
		pdf.Line(x, y+r, x, y+h-r)     // Left side
	}
}

func drawDashedLine(pdf *gofpdf.Fpdf, x1, y1, x2, y2 float64) {
	pdf.SetDrawColor(0, 0, 0)
	pdf.SetLineWidth(0.5)

	// Pattern dash sederhana: [panjang_dash, panjang_gap]
	pdf.SetDashPattern([]float64{3, 2}, 0) // Dash 3mm, Gap 2mm
	pdf.Line(x1, y1, x2, y2)
	pdf.SetDashPattern([]float64{}, 0) // Reset ke solid
}

func (s *PDFHandler) addContentReceipt(pdf *gofpdf.Fpdf, data interface{}) {
	// Set text color to black for the first part
	getX, getY := pdf.GetXY()
	pdf.SetXY(getX+2, getY-7)
	pdf.SetTextColor(0, 0, 0)
	imagePathIconSuccess := "assets" + "/images/icon-2.png"
	imageIconSucessWidth := 15.0
	imageIconSucessHeight := 15.0
	x, y := pdf.GetXY()
	pdf.ImageOptions(imagePathIconSuccess, x+2, y+2, imageIconSucessWidth, imageIconSucessHeight, false, gofpdf.ImageOptions{ImageType: "PNG"}, 0, "")
	pdf.CellFormat(0, 6, "", "", 0, "L", false, 0, "")

	pdf.SetX(10 + imageIconSucessWidth + 8)
	pdf.SetTextColor(0, 0, 0)
	pdf.SetFont("Arial", "", 14)
	pdf.Cell(0, 8, fmt.Sprintf("%v", "Transaksi Sukses"))
	pdf.Ln(10)

	pdf.SetX(10 + imageIconSucessWidth + 8)
	pdf.SetFont("Arial", "", 14)
	pdf.SetTextColor(107, 114, 128)

	// Tulis bagian pertama
	pdf.CellFormat(0, 6, "DOP1234567890", "", 0, "L", false, 0, "")

	pageW, _ := pdf.GetPageSize()
	pdf.SetX(pageW - 50)
	pdf.CellFormat(0, 6, "DO Pertamina", "", 0, "L", false, 0, "")

	pdf.Ln(6)

	drawDashedLine(pdf, getX+6, getY+20, pageW-16, getY+20)

	// Details list
	details := []struct {
		label string
		value string
	}{
		{"Tanggal Transaksi", "05/11/2025 09:10:33"},
		{"Tipe Transaksi", "Product Allocation"},
	}

	// Position details on the right side
	for _, detail := range details {
		pdf.SetXY(15, getY+30)
		pdf.SetFont("Arial", "", 10)
		pdf.SetTextColor(107, 114, 128)
		pdf.Cell(60, 6, detail.label)

		pdf.SetXY(100, getY+30)
		pdf.SetTextColor(0, 0, 0)
		pdf.CellFormat(30, 6, detail.value, "", 0, "L", false, 0, "")

		getY += 8
	}

	pdf.Ln(10)
	pdf.SetX(15)
	pdf.SetFillColor(235, 244, 245) // Warna background
	pdf.SetTextColor(0, 0, 0)       // Text putih agar kontras

	pdf.CellFormat(pageW-30, 12, "Informasi Transaksi", "0", 0, "L", true, 0, "")

	pdf.Ln(10)

	// Details list
	transactionInformation := []struct {
		label string
		value string
	}{
		{"Status SO", "Release"},
		{"Nomor SO", "4000012322"},
		{"ID Aplikasi", "250214024675"},
		{"Nomor Perjanjian Penjadwalan", "352424411"},
		{"Sumber Dana", "Khris - Khrisna Joh - PT. BANK RAKYAT INDONESIA (PERSERO) TBK - BRINIDJA 1001******890"},
		{"Organisasi Penjualan", "007-C&T LPG Retail"},
		{"Grup Produk", "001-LPG/BBG"},
		{"Tujuan Pengiriman ", "100123 - PT. Raya Jaya Mulya"},
		{"Pembeli", "966787 - PT. Makmur Sentosa"},
		{"Depo", "2150 - SPBE Wanantara D. Satria"},
		{"Pembayar", "966787 - PT. Makmur Sentosa"},
	}

	// Position details on the right side
	for _, transaction := range transactionInformation {
		pdf.SetXY(15, getY+50)
		pdf.SetFont("Arial", "", 10)
		pdf.SetTextColor(107, 114, 128)
		pdf.Cell(60, 6, transaction.label)

		pdf.SetXY(100, getY+50)
		pdf.SetTextColor(0, 0, 0)
		pdf.CellFormat(30, 6, transaction.value, "", 0, "L", false, 0, "")

		getY += 8
	}

	pdf.Ln(10)
	pdf.SetX(15)
	pdf.SetFillColor(235, 244, 245) // Warna background
	pdf.SetTextColor(0, 0, 0)       // Text putih agar kontras

	pdf.CellFormat(pageW-30, 12, "Rincian Transaksi", "0", 0, "L", true, 0, "")

	pdf.Ln(10)

	// Details list
	detailTransaction := []struct {
		label string
		value string
	}{
		{"Nilai Bersih", "IDR 16.000.000,00"},
		{"Pajak PPN", "IDR 2.086.000,00"},
		{"Pajak PBBKB", "IDR 0,00"},
		{"Pajak PPH", "IDR 0,00"},
		{"Nilai Kotor", "IDR 18.068.000,00"},
		{"Nilai Debet/Kredit", "IDR 1.000.000,00"},
		{"Biaya Admin", "IDR 2.000,00"},
	}

	// Position details on the right side
	for _, detail := range detailTransaction {
		pdf.SetXY(100, getY+70)
		pdf.SetFont("Arial", "", 10)
		pdf.SetTextColor(107, 114, 128)
		pdf.Cell(60, 6, detail.label)

		pdf.SetXY(pageW-50, getY+70)
		pdf.SetTextColor(0, 0, 0)
		pdf.CellFormat(30, 6, detail.value, "", 0, "R", false, 0, "")

		getY += 8
	}

	// Dotted line
	drawDashedLine(pdf, 100, getY+75, pageW-16, getY+75)
	getY += 10

	// Total Payment
	pdf.SetXY(100, getY+70)
	pdf.SetFont("Arial", "B", 10)
	pdf.SetTextColor(107, 114, 128)
	pdf.Cell(60, 6, "Total Pembayaran")

	pdf.SetXY(pageW-50, getY+70)
	pdf.SetTextColor(0, 0, 0)
	pdf.Cell(30, 6, "IDR 17.070.000,00")

	pdf.Ln(10)
	pdf.SetX(15)
	pdf.SetFillColor(235, 244, 245) // Warna background
	pdf.SetTextColor(0, 0, 0)       // Text putih agar kontras

	pdf.CellFormat(pageW-30, 12, "Detail Produk", "0", 0, "L", true, 0, "")

	pdf.Ln(15)

	products := []struct {
		material    string
		description string
		trip        string
		qty         string
		uom         string
		transporter string
		sendDate    string
	}{
		{"A040900002", "ADV 0001 DRUM 220 KG", "51", "12", "B03", "13bBAHBD738847Y", "23/12/2025"},
		{"A040900002", "ADV 0001 DRUM 220 KG", "51", "12", "B03", "13bBAHBD738847Y", "23/12/2025"},
		{"A040900002", "ADV 0001 DRUM 220 KG", "51", "12", "B03", "13bBAHBD738847Y", "23/12/2025"},
		{"A040900002", "ADV 0001 DRUM 220 KG", "51", "12", "B03", "13bBAHBD738847Y", "23/12/2025"},
		// {"A040900002", "ADV 0001 DRUM 220 KG", "51", "12", "B03", "13bBAHBD738847Y", "23/12/2025"},
		// {"A040900002", "ADV 0001 DRUM 220 KG", "51", "12", "B03", "13bBAHBD738847Y", "23/12/2025"},
		// {"A040900002", "ADV 0001 DRUM 220 KG", "51", "12", "B03", "13bBAHBD738847Y", "23/12/2025"},
		// {"A040900002", "ADV 0001 DRUM 220 KG", "51", "12", "B03", "13bBAHBD738847Y", "23/12/2025"},
		// {"A040900002", "ADV 0001 DRUM 220 KG", "51", "12", "B03", "13bBAHBD738847Y", "23/12/2025"},
	}

	// Table header
	pdf.SetFont("Arial", "", 10)
	pdf.SetX(15)
	pdf.SetFillColor(240, 243, 243)

	pdf.CellFormat(40, 8, "Material", "1", 0, "C", true, 0, "")
	pdf.CellFormat(80, 8, "Deskripsi Material", "1", 0, "C", true, 0, "")
	pdf.CellFormat(20, 8, "Trip", "1", 0, "C", true, 0, "")
	pdf.CellFormat(20, 8, "Qty", "1", 0, "C", true, 0, "")
	pdf.CellFormat(20, 8, "UOM", "1", 0, "C", true, 0, "")
	pdf.CellFormat(40, 8, "Transporter", "1", 0, "C", true, 0, "")
	pdf.CellFormat(45, 8, "Tgl Kirim", "1", 0, "C", true, 0, "")
	pdf.Ln(8)

	// Table rows
	for _, prod := range products {
		pdf.SetX(15)
		pdf.CellFormat(40, 8, prod.material, "1", 0, "C", false, 0, "")
		pdf.CellFormat(80, 8, prod.description, "1", 0, "C", false, 0, "")
		pdf.CellFormat(20, 8, prod.trip, "1", 0, "C", false, 0, "")
		pdf.CellFormat(20, 8, prod.qty, "1", 0, "C", false, 0, "")
		pdf.CellFormat(20, 8, prod.uom, "1", 0, "C", false, 0, "")
		pdf.CellFormat(40, 8, prod.transporter, "1", 0, "C", false, 0, "")
		pdf.CellFormat(45, 8, prod.sendDate, "1", 0, "C", false, 0, "")
		pdf.Ln(8)
	}

	pdf.SetX(15)
	pdf.CellFormat(40, 8, "Total", "1", 0, "C", false, 0, "")
	pdf.CellFormat(225, 8, "128", "1", 0, "C", false, 0, "")

	pdf.Ln(10)

}
