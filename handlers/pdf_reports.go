package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"generatePDF/models"
	"net/http"
	"os"
	"time"

	"github.com/jung-kurt/gofpdf"
)

func (h *PDFHandler) GeneratePDF(w http.ResponseWriter, r *http.Request) {
	// Parse request body
	var req models.GeneratePDFRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Validate
	if req.ProcessID == "" || req.TransactionID == "" {
		http.Error(w, "ProcessID and TransactionID are required", http.StatusBadRequest)
		return
	}

	// Generate PDF
	// pdfBytes, err := h.convertPdf(&req)
	// if err != nil {
	// 	http.Error(w, err.Error(), http.StatusInternalServerError)
	// 	return
	// }

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

	pdf, err := h.GenerateReport(data)
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
func (s *PDFHandler) GenerateReport(data interface{}) (*gofpdf.Fpdf, error) {
	pdf := gofpdf.New("P", "mm", "A4", "assets/font")
	pdf.SetAutoPageBreak(true, 25)

	pdf.AddFont("BRIDigital", "", "BRIDigitalText-Regular.json")
	pdf.AddFont("BRIDigital", "B", "BRIDigitalText-SemiBold.json")
	pdf.AddFont("BRIDigitalLogo", "B", "BRIDigitalText-SemiBold.json")

	pdf.SetHeaderFunc(
		func() {
			s.addHeader(pdf)
			// s.addWatermark(pdf)
		})
	pdf.AddPageFormat("P", gofpdf.SizeType{Wd: 297, Ht: 420})
	// Add header
	// Add watermark background

	fontSize := 12.0
	fontName := "BRIDigital"
	pageW, _ := pdf.GetPageSize()
	// heightRow := 10.0
	// heightHeader := 15.0
	pdf.SetFont(fontName, "", fontSize)
	// paddingDetailTrx := 20.0

	pdf.SetFont(fontName, "B", 16)

	pdf.AliasNbPages("{nb}")

	pdf.SetFooterFunc(func() {
		pdf.SetY(-19)
		pdf.Ln(2)
		pdf.Line(10, pdf.GetY(), pageW-8, pdf.GetY())
		pdf.Ln(1)
		pdf.SetFont(fontName, "", 11)

		// Simpan posisi awal
		// startX := pdf.GetX()
		startY := pdf.GetY()

		// Text kiri
		pdf.MultiCell(pageW-100, 7, "Terima kasih telah bertransaksi menggunakan Qlola BRI, bila menemui kendala silakan hubungi kami di 1500001 atau qlola@bri.co.id", "", "L", false)

		// Kembali ke posisi awal untuk bagian kanan
		pdf.SetY(startY)
		pdf.SetX(225)

		currentTime := time.Now()
		formattedTime := currentTime.Format("02/01/2006 15:04:05")

		// Text kanan - juga pakai MultiCell untuk alignment yang konsisten
		pdf.MultiCell(pageW/2, 7, fmt.Sprintf("%v - Halaman %d/{nb}", formattedTime, pdf.PageNo()), "", "L", false)
	})

	// Add content
	s.addContent(pdf, data)

	// // Add footer
	// s.addFooter(pdf)

	return pdf, nil
}

func (s *PDFHandler) addWatermark(pdf *gofpdf.Fpdf, paper Paper) {
	pdf.SetAlpha(0.5, "Normal")
	pdf.ImageOptions("./assets/images/newwatermark3.png", paper.RectSetup.X, paper.RectSetup.Y, paper.RectSetup.W, paper.RectSetup.H, false, gofpdf.ImageOptions{
		ReadDpi:   false,
		ImageType: "", // biarkan kosong agar otomatis deteksi dari ekstensi (jpg, png, dll)
	}, 0, "")
	pdf.SetAlpha(1.0, "Normal")
}

func setBackgroundColor(pdf *gofpdf.Fpdf, x, y, width, height float64, r, g, b int) {
	// Simpan warna sebelumnya
	r1, g1, b1 := pdf.GetFillColor()

	// Set warna fill baru
	pdf.SetFillColor(r, g, b)

	// Gambar rectangle dengan warna fill
	pdf.Rect(x, y, width, height, "F")

	// Kembalikan warna semula
	pdf.SetFillColor(r1, g1, b1)
}

func (s *PDFHandler) addHeader(pdf *gofpdf.Fpdf) {
	pageW, _ := pdf.GetPageSize()

	setBackgroundColor(pdf, 0, 0, pageW, 50, 16, 47, 50)
	// Add logos
	qlolaLogo := "./assets/images/qlola.png"
	briLogo := "./assets/images/bri.png"

	opts := gofpdf.ImageOptions{ImageType: "", ReadDpi: false}

	// Left logo (Qlola)
	if _, err := os.Stat(qlolaLogo); err == nil {
		pdf.ImageOptions(qlolaLogo, 10, 10, 40, 15, false, opts, 0, "")
	}

	// Right logo (BRI)
	if _, err := os.Stat(briLogo); err == nil {
		pageWidth, _ := pdf.GetPageSize()
		pdf.ImageOptions(briLogo, pageWidth-50, 10, 40, 15, false, opts, 0, "")
	}

	// Reset position for content
	pdf.SetY(45)
}

func (s *PDFHandler) addContent(pdf *gofpdf.Fpdf, data interface{}) {
	pageW, _ := pdf.GetPageSize()
	x, y := pdf.GetXY()
	fieldMap := s.extractFieldMap(data)

	// Title
	pdf.SetY(y - 8)
	pdf.SetFont("Arial", "B", 16)
	pdf.Cell(0, 10, "Laporan Transaksi")

	// Status
	pdf.SetY(y - 6)
	pdf.SetFillColor(253, 239, 216) // #fdefd8
	pdf.SetTextColor(144, 94, 10)   // #905e0a
	pdf.SetX(pageW - 50)

	pdf.CellFormat(40, 8, fieldMap["status"], "0", 0, "C", true, 0, "")
	pdf.Ln(15)

	// Reset text color
	pdf.SetTextColor(0, 0, 0)

	pdf.SetX(x)
	pdf.Ln(0)

	// Transaction ID
	pdf.SetY(y)
	pdf.SetFont("Arial", "", 12)
	pdf.Cell(0, 8, "Transaction ID: "+fieldMap["transactionId"])
	pdf.SetY(y)
	pdf.Ln(16)

	// Account Information - 3 columns
	s.addAccountInfo(pdf, fieldMap)

	// Line separator
	pdf.Line(10, pdf.GetY()+8, pageW-10, pdf.GetY()+8)
	pdf.Ln(10)

	// Transaction Details
	s.addTransactionDetails(pdf, fieldMap)

	// Line separator
	pdf.Line(10, pdf.GetY(), pageW-8, pdf.GetY())
	pdf.Ln(4)

	// Transaction Information
	s.addTransactionInfo(pdf, fieldMap)

	// Line separator
	pdf.Line(10, pdf.GetY(), pageW-8, pdf.GetY())
	pdf.Ln(4)

	s.addDetailProduct(pdf, fieldMap)

	// Line separator
	pdf.Line(10, pdf.GetY(), pageW-8, pdf.GetY())
	pdf.Ln(4)

	s.addLogActivity(pdf, fieldMap)

}

// Helper function untuk cek perlu page break
func (s *PDFHandler) needPageBreak(pdf *gofpdf.Fpdf, rowHeight float64) bool {
	return pdf.GetY()+rowHeight > 250 // 297 - 25 - buffer
}

// Helper function untuk pastikan ada space untuk footer
func (s *PDFHandler) ensureFooterSpace(pdf *gofpdf.Fpdf) {
	currentY := pdf.GetY()
	if currentY < 260 { // Jika masih ada space
		pdf.SetY(260) // Set posisi ke 260mm dari atas
	}
}

func (s *PDFHandler) addLogActivity(pdf *gofpdf.Fpdf, fieldMap map[string]string) {
	// Title
	pdf.SetX(10)
	pdf.SetFont("Arial", "B", 12)
	pdf.Cell(0, 10, "Log Aktivitas")
	pdf.Ln(12)

	// Log activity table (simplified)
	logs := []struct {
		no          string
		date        string
		name        string
		role        string
		action      string
		description string
	}{
		{"1", "30/01/2025 10:00:00", "Andre", "Maker", "Send for Approval", ""},
		{"2", "30/01/2025 10:00:00", "Naufal", "Checker", "Aprrove", ""},
		{"3", "30/01/2025 10:00:00", "Satya", "Checker", "Approve", ""},
		{"4", "30/01/2025 10:00:00", "Arya", "Cheker", "Approve", ""},
		{"5", "30/01/2025 10:00:00", "Billy", "Cheker", "Approve", ""},
		{"6", "30/01/2025 10:00:00", "Tomo", "Checker", "Approve", ""},
		{"7", "30/01/2025 10:00:00", "Fajar", "Signer", "Approve", ""},
		{"8", "30/01/2025 10:00:00", "Alfa", "Releaser", "Approve", ""},
		{"1", "30/01/2025 10:00:00", "Andre", "Maker", "Send for Approval", ""},
		{"2", "30/01/2025 10:00:00", "Naufal", "Checker", "Aprrove", ""},
		{"3", "30/01/2025 10:00:00", "Satya", "Checker", "Approve", ""},
		{"4", "30/01/2025 10:00:00", "Arya", "Cheker", "Approve", ""},
		{"5", "30/01/2025 10:00:00", "Billy", "Cheker", "Approve", ""},
		{"6", "30/01/2025 10:00:00", "Tomo", "Checker", "Approve", ""},
		{"7", "30/01/2025 10:00:00", "Fajar", "Signer", "Approve", ""},
		{"8", "30/01/2025 10:00:00", "Alfa", "Releaser", "Approve", ""},
		{"1", "30/01/2025 10:00:00", "Andre", "Maker", "Send for Approval", ""},
		{"2", "30/01/2025 10:00:00", "Naufal", "Checker", "Aprrove", ""},
		{"3", "30/01/2025 10:00:00", "Satya", "Checker", "Approve", ""},
		{"4", "30/01/2025 10:00:00", "Arya", "Cheker", "Approve", ""},
		{"5", "30/01/2025 10:00:00", "Billy", "Cheker", "Approve", ""},
		{"6", "30/01/2025 10:00:00", "Tomo", "Checker", "Approve", ""},
		{"7", "30/01/2025 10:00:00", "Fajar", "Signer", "Approve", ""},
		{"8", "30/01/2025 10:00:00", "Alfa", "Releaser", "Approve", ""},
		{"1", "30/01/2025 10:00:00", "Andre", "Maker", "Send for Approval", ""},
		{"2", "30/01/2025 10:00:00", "Naufal", "Checker", "Aprrove", ""},
		{"3", "30/01/2025 10:00:00", "Satya", "Checker", "Approve", ""},
		{"4", "30/01/2025 10:00:00", "Arya", "Cheker", "Approve", ""},
		{"5", "30/01/2025 10:00:00", "Billy", "Cheker", "Approve", ""},
		{"6", "30/01/2025 10:00:00", "Tomo", "Checker", "Approve", ""},
		{"7", "30/01/2025 10:00:00", "Fajar", "Signer", "Approve", ""},
		{"8", "30/01/2025 10:00:00", "Alfa", "Releaser", "Approve", ""},
		{"1", "30/01/2025 10:00:00", "Andre", "Maker", "Send for Approval", ""},
		{"2", "30/01/2025 10:00:00", "Naufal", "Checker", "Aprrove", ""},
		{"3", "30/01/2025 10:00:00", "Satya", "Checker", "Approve", ""},
		{"4", "30/01/2025 10:00:00", "Arya", "Cheker", "Approve", ""},
		{"5", "30/01/2025 10:00:00", "Billy", "Cheker", "Approve", ""},
		{"6", "30/01/2025 10:00:00", "Tomo", "Checker", "Approve", ""},
		{"7", "30/01/2025 10:00:00", "Fajar", "Signer", "Approve", ""},
		{"8", "30/01/2025 10:00:00", "Alfa", "Releaser", "Approve", ""},
		{"1", "30/01/2025 10:00:00", "Andre", "Maker", "Send for Approval", ""},
		{"2", "30/01/2025 10:00:00", "Naufal", "Checker", "Aprrove", ""},
		{"3", "30/01/2025 10:00:00", "Satya", "Checker", "Approve", ""},
		{"4", "30/01/2025 10:00:00", "Arya", "Cheker", "Approve", ""},
		{"5", "30/01/2025 10:00:00", "Billy", "Cheker", "Approve", ""},
		{"6", "30/01/2025 10:00:00", "Tomo", "Checker", "Approve", ""},
		{"7", "30/01/2025 10:00:00", "Fajar", "Signer", "Approve", ""},
		{"8", "30/01/2025 10:00:00", "Alfa", "Releaser", "Approve", ""},
		{"1", "30/01/2025 10:00:00", "Andre", "Maker", "Send for Approval", ""},
		{"2", "30/01/2025 10:00:00", "Naufal", "Checker", "Aprrove", ""},
		{"3", "30/01/2025 10:00:00", "Satya", "Checker", "Approve", ""},
		{"4", "30/01/2025 10:00:00", "Arya", "Cheker", "Approve", ""},
		{"5", "30/01/2025 10:00:00", "Billy", "Cheker", "Approve", ""},
		{"6", "30/01/2025 10:00:00", "Tomo", "Checker", "Approve", ""},
		{"7", "30/01/2025 10:00:00", "Fajar", "Signer", "Approve", ""},
		{"8", "30/01/2025 10:00:00", "Alfa", "Releaser", "Approve", ""},
		{"1", "30/01/2025 10:00:00", "Andre", "Maker", "Send for Approval", ""},
		{"2", "30/01/2025 10:00:00", "Naufal", "Checker", "Aprrove", ""},
		{"3", "30/01/2025 10:00:00", "Satya", "Checker", "Approve", ""},
		{"4", "30/01/2025 10:00:00", "Arya", "Cheker", "Approve", ""},
		{"5", "30/01/2025 10:00:00", "Billy", "Cheker", "Approve", ""},
		{"6", "30/01/2025 10:00:00", "Tomo", "Checker", "Approve", ""},
		{"7", "30/01/2025 10:00:00", "Fajar", "Signer", "Approve", ""},
		{"8", "30/01/2025 10:00:00", "Alfa", "Releaser", "Approve", ""},
		{"1", "30/01/2025 10:00:00", "Andre", "Maker", "Send for Approval", ""},
		{"2", "30/01/2025 10:00:00", "Naufal", "Checker", "Aprrove", ""},
		{"3", "30/01/2025 10:00:00", "Satya", "Checker", "Approve", ""},
		{"4", "30/01/2025 10:00:00", "Arya", "Cheker", "Approve", ""},
		{"5", "30/01/2025 10:00:00", "Billy", "Cheker", "Approve", ""},
		{"6", "30/01/2025 10:00:00", "Tomo", "Checker", "Approve", ""},
		{"7", "30/01/2025 10:00:00", "Fajar", "Signer", "Approve", ""},
		{"8", "30/01/2025 10:00:00", "Alfa", "Releaser", "Approve", ""},
		{"1", "30/01/2025 10:00:00", "Andre", "Maker", "Send for Approval", ""},
		{"2", "30/01/2025 10:00:00", "Naufal", "Checker", "Aprrove", ""},
		{"3", "30/01/2025 10:00:00", "Satya", "Checker", "Approve", ""},
		{"4", "30/01/2025 10:00:00", "Arya", "Cheker", "Approve", ""},
		{"5", "30/01/2025 10:00:00", "Billy", "Cheker", "Approve", ""},
		{"6", "30/01/2025 10:00:00", "Tomo", "Checker", "Approve", ""},
		{"7", "30/01/2025 10:00:00", "Fajar", "Signer", "Approve", ""},
		{"8", "30/01/2025 10:00:00", "Alfa", "Releaser", "Approve", ""},
		{"1", "30/01/2025 10:00:00", "Andre", "Maker", "Send for Approval", ""},
		{"2", "30/01/2025 10:00:00", "Naufal", "Checker", "Aprrove", ""},
		{"3", "30/01/2025 10:00:00", "Satya", "Checker", "Approve", ""},
		{"4", "30/01/2025 10:00:00", "Arya", "Cheker", "Approve", ""},
		{"5", "30/01/2025 10:00:00", "Billy", "Cheker", "Approve", ""},
		{"6", "30/01/2025 10:00:00", "Tomo", "Checker", "Approve", ""},
		{"7", "30/01/2025 10:00:00", "Fajar", "Signer", "Approve", ""},
		{"8", "30/01/2025 10:00:00", "Alfa", "Releaser", "Approve", ""},
		{"1", "30/01/2025 10:00:00", "Andre", "Maker", "Send for Approval", ""},
		{"2", "30/01/2025 10:00:00", "Naufal", "Checker", "Aprrove", ""},
		{"3", "30/01/2025 10:00:00", "Satya", "Checker", "Approve", ""},
		{"4", "30/01/2025 10:00:00", "Arya", "Cheker", "Approve", ""},
		{"5", "30/01/2025 10:00:00", "Billy", "Cheker", "Approve", ""},
		{"6", "30/01/2025 10:00:00", "Tomo", "Checker", "Approve", ""},
		{"7", "30/01/2025 10:00:00", "Fajar", "Signer", "Approve", ""},
		{"8", "30/01/2025 10:00:00", "Alfa", "Releaser", "Approve", ""},
		{"1", "30/01/2025 10:00:00", "Andre", "Maker", "Send for Approval", ""},
		{"2", "30/01/2025 10:00:00", "Naufal", "Checker", "Aprrove", ""},
		{"3", "30/01/2025 10:00:00", "Satya", "Checker", "Approve", ""},
		{"4", "30/01/2025 10:00:00", "Arya", "Cheker", "Approve", ""},
		{"5", "30/01/2025 10:00:00", "Billy", "Cheker", "Approve", ""},
		{"6", "30/01/2025 10:00:00", "Tomo", "Checker", "Approve", ""},
		{"7", "30/01/2025 10:00:00", "Fajar", "Signer", "Approve", ""},
		{"8", "30/01/2025 10:00:00", "Alfa", "Releaser", "Approve", ""},
		{"1", "30/01/2025 10:00:00", "Andre", "Maker", "Send for Approval", ""},
		{"2", "30/01/2025 10:00:00", "Naufal", "Checker", "Aprrove", ""},
		{"3", "30/01/2025 10:00:00", "Satya", "Checker", "Approve", ""},
		{"4", "30/01/2025 10:00:00", "Arya", "Cheker", "Approve", ""},
		{"5", "30/01/2025 10:00:00", "Billy", "Cheker", "Approve", ""},
		{"6", "30/01/2025 10:00:00", "Tomo", "Checker", "Approve", ""},
		{"7", "30/01/2025 10:00:00", "Fajar", "Signer", "Approve", ""},
		{"8", "30/01/2025 10:00:00", "Alfa", "Releaser", "Approve", ""},
		{"1", "30/01/2025 10:00:00", "Andre", "Maker", "Send for Approval", ""},
		{"2", "30/01/2025 10:00:00", "Naufal", "Checker", "Aprrove", ""},
		{"3", "30/01/2025 10:00:00", "Satya", "Checker", "Approve", ""},
		{"4", "30/01/2025 10:00:00", "Arya", "Cheker", "Approve", ""},
		{"5", "30/01/2025 10:00:00", "Billy", "Cheker", "Approve", ""},
		{"6", "30/01/2025 10:00:00", "Tomo", "Checker", "Approve", ""},
		{"7", "30/01/2025 10:00:00", "Fajar", "Signer", "Approve", ""},
		{"8", "30/01/2025 10:00:00", "Alfa", "Releaser", "Approve", ""},
		{"1", "30/01/2025 10:00:00", "Andre", "Maker", "Send for Approval", ""},
		{"2", "30/01/2025 10:00:00", "Naufal", "Checker", "Aprrove", ""},
		{"3", "30/01/2025 10:00:00", "Satya", "Checker", "Approve", ""},
		{"4", "30/01/2025 10:00:00", "Arya", "Cheker", "Approve", ""},
		{"5", "30/01/2025 10:00:00", "Billy", "Cheker", "Approve", ""},
		{"6", "30/01/2025 10:00:00", "Tomo", "Checker", "Approve", ""},
		{"7", "30/01/2025 10:00:00", "Fajar", "Signer", "Approve", ""},
		{"8", "30/01/2025 10:00:00", "Alfa", "Releaser", "Approve", ""},
	}

	// Table header
	pdf.SetFont("Arial", "B", 10)
	pdf.CellFormat(20, 8, "No", "1", 0, "C", false, 0, "")
	pdf.CellFormat(50, 8, "Tanggal", "1", 0, "C", false, 0, "")
	pdf.CellFormat(30, 8, "Nama", "1", 0, "C", false, 0, "")
	pdf.CellFormat(30, 8, "Peran", "1", 0, "C", false, 0, "")
	pdf.CellFormat(70, 8, "Aksi", "1", 0, "C", false, 0, "")
	pdf.CellFormat(80, 8, "Deskripsi", "1", 0, "C", false, 0, "")
	pdf.Ln(8)

	// Table rows
	pdf.SetFont("Arial", "", 10)
	for _, log := range logs {
		pdf.CellFormat(20, 8, log.no, "1", 0, "C", false, 0, "")
		pdf.CellFormat(50, 8, log.date, "1", 0, "C", false, 0, "")
		pdf.CellFormat(30, 8, log.name, "1", 0, "C", false, 0, "")
		pdf.CellFormat(30, 8, log.role, "1", 0, "C", false, 0, "")
		pdf.CellFormat(70, 8, log.action, "1", 0, "C", false, 0, "")
		pdf.CellFormat(80, 8, log.description, "1", 0, "C", false, 0, "")
		pdf.Ln(8)
	}

	pdf.Ln(20)
}

func (s *PDFHandler) addDetailProduct(pdf *gofpdf.Fpdf, fieldMap map[string]string) {
	// Title
	pdf.SetX(10)
	pdf.SetFont("Arial", "B", 12)
	pdf.Cell(0, 10, "Detail Produk")
	pdf.Ln(12)

	// Product details table (simplified)
	products := []struct {
		material    string
		description string
		sendDate    string
		trip        string
		transporter string
		uom         string
		totalqty    string
	}{
		{"A040900083", "", "30/01/2024", "1", "", "BB6", "4"},
		{"A040900083", "", "30/01/2024", "1", "", "BB6", "4"},
		{"A040900083", "", "30/01/2024", "1", "", "BB6", "4"},
		{"A040900083", "", "30/01/2024", "1", "", "BB6", "4"},
		{"A040900083", "", "30/01/2024", "1", "", "BB6", "4"},
		{"A040900083", "", "30/01/2024", "1", "", "BB6", "4"},
		{"A040900083", "", "30/01/2024", "1", "", "BB6", "4"},
		{"A040900083", "", "30/01/2024", "1", "", "BB6", "4"},
		{"A040900083", "", "30/01/2024", "1", "", "BB6", "4"},
		{"A040900083", "", "30/01/2024", "1", "", "BB6", "4"},
		{"A040900083", "", "30/01/2024", "1", "", "BB6", "4"},
		{"A040900083", "", "30/01/2024", "1", "", "BB6", "4"},
		{"A040900083", "", "30/01/2024", "1", "", "BB6", "4"},
		{"A040900083", "", "30/01/2024", "1", "", "BB6", "4"},
		{"A040900083", "", "30/01/2024", "1", "", "BB6", "4"},
		{"A040900083", "", "30/01/2024", "1", "", "BB6", "4"},
		{"A040900083", "", "30/01/2024", "1", "", "BB6", "4"},
		{"A040900083", "", "30/01/2024", "1", "", "BB6", "4"},
		{"A040900083", "", "30/01/2024", "1", "", "BB6", "4"},
		{"A040900083", "", "30/01/2024", "1", "", "BB6", "4"},
		{"A040900083", "", "30/01/2024", "1", "", "BB6", "4"},
		{"A040900083", "", "30/01/2024", "1", "", "BB6", "4"},
		{"A040900083", "", "30/01/2024", "1", "", "BB6", "4"},
		{"A040900083", "", "30/01/2024", "1", "", "BB6", "4"},
		{"A040900083", "", "30/01/2024", "1", "", "BB6", "4"},
		{"A040900083", "", "30/01/2024", "1", "", "BB6", "4"},
		{"A040900083", "", "30/01/2024", "1", "", "BB6", "4"},
		{"A040900083", "", "30/01/2024", "1", "", "BB6", "4"},
		{"A040900083", "", "30/01/2024", "1", "", "BB6", "4"},
		{"A040900083", "", "30/01/2024", "1", "", "BB6", "4"},
		{"A040900083", "", "30/01/2024", "1", "", "BB6", "4"},
		{"A040900083", "", "30/01/2024", "1", "", "BB6", "4"},
		{"A040900083", "", "30/01/2024", "1", "", "BB6", "4"},
		{"A040900083", "", "30/01/2024", "1", "", "BB6", "4"},
		{"A040900083", "", "30/01/2024", "1", "", "BB6", "4"},
	}

	// Table header
	pdf.SetFont("Arial", "B", 10)
	pdf.CellFormat(30, 8, "Material", "1", 0, "C", false, 0, "")
	pdf.CellFormat(25, 8, "Perjalanan", "1", 0, "C", false, 0, "")
	pdf.CellFormat(30, 8, "Total Kuantitas", "1", 0, "C", false, 0, "")
	pdf.CellFormat(20, 8, "UOM", "1", 0, "C", false, 0, "")
	pdf.CellFormat(55, 8, "Pengangkutan", "1", 0, "C", false, 0, "")
	pdf.CellFormat(40, 8, "Tanggal Kirim", "1", 0, "C", false, 0, "")
	pdf.CellFormat(80, 8, "Keterangan", "1", 0, "C", false, 0, "")
	pdf.Ln(8)

	// Table rows
	pdf.SetFont("Arial", "", 10)
	for _, prod := range products {
		pdf.CellFormat(30, 8, prod.material, "1", 0, "C", false, 0, "")
		pdf.CellFormat(25, 8, prod.trip, "1", 0, "C", false, 0, "")
		pdf.CellFormat(30, 8, prod.totalqty, "1", 0, "C", false, 0, "")
		pdf.CellFormat(20, 8, prod.uom, "1", 0, "C", false, 0, "")
		pdf.CellFormat(55, 8, prod.transporter, "1", 0, "C", false, 0, "")
		pdf.CellFormat(40, 8, prod.sendDate, "1", 0, "C", false, 0, "")
		pdf.CellFormat(80, 8, prod.description, "1", 0, "C", false, 0, "")
		pdf.Ln(8)
	}

	pdf.CellFormat(55, 8, "Total", "1", 0, "C", false, 0, "")
	pdf.CellFormat(225, 8, "128", "1", 0, "C", false, 0, "")

	pdf.Ln(20)
}

func (s *PDFHandler) addAccountInfo(pdf *gofpdf.Fpdf, fieldMap map[string]string) {
	currentY := pdf.GetY()

	// Column 1: Sumber Dana
	pdf.SetX(10)
	pdf.SetFont("Arial", "", 10)
	pdf.SetTextColor(107, 114, 128) // #6b7280
	pdf.Cell(0, 6, "Sumber Dana")
	pdf.Ln(6)

	// Set text color to black for the first part
	pdf.SetTextColor(0, 0, 0)
	imagePathIconBRI := "assets" + "/images/icon.png"
	imageIconBRIWidth := 14.0
	imageIconBRIHeight := 14.0
	x, y := pdf.GetXY()
	pdf.ImageOptions(imagePathIconBRI, x+2, y+2, imageIconBRIWidth, imageIconBRIHeight, false, gofpdf.ImageOptions{ImageType: "PNG"}, 0, "")
	pdf.CellFormat(0, 6, "", "", 0, "L", false, 0, "")

	pdf.SetX(10 + imageIconBRIWidth + 4)
	pdf.SetTextColor(0, 0, 0)
	pdf.SetFont("Arial", "", 10)
	pdf.Cell(0, 8, fmt.Sprintf("%v - %v | %v...", fieldMap["debitAccountName"][0:5], fieldMap["debitAccountName"], fieldMap["accountNumber"][0:8]))
	pdf.Ln(6)

	pdf.SetX(10 + imageIconBRIWidth + 4)
	pdf.SetFont("Arial", "", 10)
	pdf.SetTextColor(107, 114, 128)

	// Tulis bagian pertama
	pdf.CellFormat(0, 6, "BRI", "", 0, "L", false, 0, "")

	briWidth := pdf.GetStringWidth("BRI")

	// Tambahkan icon
	iconX := 10 + imageIconBRIWidth + 4 + briWidth + 3
	pdf.Image("assets/images/ellips.png", iconX, pdf.GetY()+2.5, 1, 1, false, "", 0, "")

	// Lanjutkan dengan account number
	pdf.SetX(iconX + 1 + 1) // iconWidth (1) + spacing (1)
	pdf.Cell(0, 6, fmt.Sprintf("%v", fieldMap["accountNumber"]))

	pdf.Ln(6)

	iconsX := 10 + imageIconBRIWidth + 4 + 1
	pdf.Image("assets/images/indonesia.png", iconsX, pdf.GetY()+1, 6, 4, false, "", 0, "")
	pdf.SetX(iconsX + 6 + 2) // iconWidth (1) + spacing (1)
	pdf.Cell(0, 6, "IDR")
	pdf.Ln(10)

	// Column 2: Tujuan Transaksi (position manually)
	pdf.SetXY(100, currentY)
	pdf.SetTextColor(107, 114, 128)
	pdf.Cell(0, 6, "Tujuan Transaksi")
	pdf.Ln(6)

	// Set text color to black for the first part
	pdf.SetTextColor(0, 0, 0)
	imagePathIconPertamina := "assets" + "/images/pertamina.png"
	imageIconPertaminaWidth := 14.0
	imageIconPertaminaHeight := 14.0
	pdf.SetXY(100, currentY)
	x2, y2 := pdf.GetXY()
	pdf.ImageOptions(imagePathIconPertamina, x2+2, y2+8, imageIconPertaminaWidth, imageIconPertaminaHeight, false, gofpdf.ImageOptions{ImageType: "PNG"}, 0, "")
	pdf.CellFormat(0, 6, "", "", 0, "L", false, 0, "")

	// pdf.SetX(100 + imageIconPertaminaWidth + 2)
	// pdf.SetFont("Arial", "B", 12)
	// pdf.Cell(0, 8, "DO Pertamina")
	pdf.Ln(6)

	pdf.SetX(100 + imageIconPertaminaWidth + 6)
	pdf.SetTextColor(0, 0, 0)
	pdf.SetFont("Arial", "", 10)
	pdf.Cell(0, 8, "DO Pertamina")
	pdf.Ln(8)

	pdf.SetX(100 + imageIconPertaminaWidth + 6)
	pdf.SetFont("Arial", "", 10)
	pdf.SetTextColor(107, 114, 128)
	pdf.Cell(0, 6, fieldMap["moduleType"])
	pdf.Ln(10)

	// Column 3: Total Pembayaran
	pdf.SetXY(200, currentY)
	pdf.SetTextColor(107, 114, 128)
	pdf.Cell(0, 6, "Total Pembayaran")
	pdf.Ln(6)

	pdf.SetX(200)
	pdf.SetTextColor(0, 0, 0)
	pdf.SetFont("Arial", "B", 14)
	pdf.Cell(0, 10, fieldMap["totalPayment"])
	pdf.Ln(16)
}

func (s *PDFHandler) addTransactionDetails(pdf *gofpdf.Fpdf, fieldMap map[string]string) {
	currentY := pdf.GetY() + 3
	pageW, _ := pdf.GetPageSize()

	// Title
	pdf.SetX(10)
	pdf.SetFont("Arial", "B", 12)
	pdf.Cell(0, 10, "Rincian Transaksi")
	pdf.Ln(12)

	// Details list
	details := []struct {
		label string
		value string
	}{
		{"Nilai Bersih", fieldMap["netValue"]},
		{"Pajak PPN", fieldMap["ppnTax"]},
		{"Pajak PBBKB", fieldMap["pbbkbTax"]},
		{"Pajak PPH", fieldMap["pphTax"]},
		{"Nilai Kotor", fieldMap["grossValue"]},
		{"Nilai Debit atau Kredit", fieldMap["debitCreditValue"]},
		{"Total Nominal", fieldMap["totalAmount"]},
	}

	// Position details on the right side
	for _, detail := range details {
		pdf.SetXY(100, currentY)
		pdf.SetFont("Arial", "", 10)
		pdf.SetTextColor(107, 114, 128)
		pdf.Cell(60, 6, detail.label)

		pdf.SetXY(pageW-38, currentY)
		pdf.SetTextColor(0, 0, 0)
		pdf.CellFormat(30, 6, detail.value, "", 0, "R", false, 0, "")

		currentY += 8
	}

	// Dotted line
	pdf.Line(100, currentY+5, pageW-8, currentY+5)
	currentY += 10

	// Total Payment
	pdf.SetXY(100, currentY)
	pdf.SetFont("Arial", "B", 10)
	pdf.SetTextColor(107, 114, 128)
	pdf.Cell(60, 6, "Total Pembayaran")

	pdf.SetXY(160, currentY)
	pdf.SetTextColor(0, 0, 0)
	pdf.Cell(30, 6, fieldMap["totalPayment"])

	pdf.SetY(currentY + 15)
}

func (s *PDFHandler) addTransactionInfo(pdf *gofpdf.Fpdf, fieldMap map[string]string) {
	currentY := pdf.GetY()

	// Title
	pdf.SetX(10)
	pdf.SetFont("Arial", "B", 12)
	pdf.Cell(0, 10, "Informasi Transaksi")
	pdf.Ln(3)

	// Two column layout
	leftColumn := []struct {
		label string
		value string
	}{
		{"ID Transaksi", fieldMap["transactionId"]},
		{"Tujuan Pengiriman", fieldMap["soldToName"]},
		{"Pembeli", fieldMap["buyer"]},
		{"Depo", fieldMap["depoName"]},
	}

	rightColumn := []struct {
		label string
		value string
	}{
		{"Organisasi Penjualan", fieldMap["salesOrganization"]},
		{"Grup Produk", fieldMap["productGroup"]},
		{"Channel Distribusi", fieldMap["distributionChannel"]},
		{"Pembayar", fieldMap["payer"]},
	}

	// Left column
	for _, info := range leftColumn {
		pdf.SetX(100)
		pdf.SetFont("Arial", "B", 10)
		pdf.Cell(0, 6, info.label)
		pdf.Ln(6)

		pdf.SetX(100)
		pdf.SetFont("Arial", "", 10)
		pdf.Cell(0, 6, info.value)
		pdf.Ln(10)
	}

	// Right column
	rightColumnY := currentY + 3
	for _, info := range rightColumn {
		pdf.SetXY(200, rightColumnY)
		pdf.SetFont("Arial", "B", 10)
		pdf.Cell(0, 6, info.label)

		pdf.SetXY(200, rightColumnY+6)
		pdf.SetFont("Arial", "", 10)
		pdf.Cell(0, 6, info.value)

		rightColumnY += 16
	}

	pdf.SetY(rightColumnY + 10)
}

func (s *PDFHandler) addFooter(pdf *gofpdf.Fpdf) {
	pageW, _ := pdf.GetPageSize()
	// Dotted line
	pdf.Line(10, pdf.GetY()+8, pageW-10, pdf.GetY()+8)
	pdf.Ln(10)

	// Footer text
	pdf.SetFont("Arial", "", 10)
	pdf.SetTextColor(107, 114, 128) // #6b7280
	pdf.MultiCell(0, 5, "Terima kasih telah bertransaksi menggunakan Qlola BRI, bila menemui kendala silakan hubungi kami di 1500001 atau qlola@bri.co.id", "", "L", false)
}

func (s *PDFHandler) extractFieldMap(data interface{}) map[string]string {
	fieldMap := make(map[string]string)

	// Default values
	defaultFields := map[string]string{
		"transactionId":       "N/A",
		"status":              "Completed",
		"debitAccountName":    "John Doe",
		"accountNumber":       "1234567890",
		"moduleType":          "Pertamina",
		"totalPayment":        "Rp 10.000.000",
		"netValue":            "Rp 8.000.000",
		"ppnTax":              "Rp 800.000",
		"pbbkbTax":            "Rp 600.000",
		"pphTax":              "Rp 400.000",
		"grossValue":          "Rp 9.800.000",
		"debitCreditValue":    "Rp 9.800.000",
		"totalAmount":         "Rp 10.000.000",
		"soldToName":          "PT. Pertamina Jakarta",
		"buyer":               "PT. Customer Indonesia",
		"depoName":            "Depo Jakarta Selatan",
		"salesOrganization":   "SO001",
		"productGroup":        "Bahan Bakar",
		"distributionChannel": "Direct Sales",
		"payer":               "PT. Customer Indonesia",
	}

	// Copy defaults
	for k, v := range defaultFields {
		fieldMap[k] = v
	}

	// Override with actual data
	if m, ok := data.(map[string]interface{}); ok {
		for k, v := range m {
			if strVal, ok := v.(string); ok {
				fieldMap[k] = strVal
			} else {
				fieldMap[k] = fmt.Sprintf("%v", v)
			}
		}
	}

	return fieldMap
}
