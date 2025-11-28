package handlers

import (
	"bytes"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/jung-kurt/gofpdf"
)

type TransactionInfo struct {
	TransactionId    string
	Status           string
	DebitAccount     string
	DebitAccountName string
	ModuleType       string
	TotalPayment     string
	Message          string
}

func (h *PDFHandler) GeneratePDF(w http.ResponseWriter, req *http.Request) {
	trxInfo := TransactionInfo{
		TransactionId:    "DOP1234567890",
		Status:           "Gagal",
		DebitAccount:     "00123456789",
		DebitAccountName: "Khrisna Joh",
		ModuleType:       "Product Allocation",
		TotalPayment:     "IDR 17.070.000,00",
		Message:          "You are a fucking disgrace",
	}

	data := []Block{}

	data = append(data, Block{
		Title:         "Rincian Transaksi",
		Type:          BLOCK_TYPE_ROWS,
		ShowTitle:     true,
		IsTitleInside: true,
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
		Title:         "Informasi Transaksi",
		Type:          BLOCK_TYPE_COLS,
		ShowTitle:     true,
		IsTitleInside: true,
		Fields: []Field{
			{
				Key:   "Nomor SO",
				Value: "SO123456789",
			},
			{
				Key:   "Status SO",
				Value: "Release",
			},
			{
				Key:   "ID Aplikasi",
				Value: "DOP1234567890",
			},
			{
				Key:   "Nomor Perjanjian Penjadwalan",
				Value: "NP1234567890",
			},
			{
				Key:   "Tujuan Pengiriman",
				Value: "100123 - PT Klara Jaya",
			},
			{
				Key:   "Pembeli",
				Value: "123123 - PT. Contoh Pembeli",
			},
			{
				Key:   "Depo",
				Value: "2150 - SPBU Ibu Kota Negara",
			},
			{
				Key:   "Organisasi Penjualan",
				Value: "007 - C&T LPG Retail",
			},
			{
				Key:   "Grup Produk",
				Value: "LPG",
			},
			{
				Key:   "Channel Distribusi",
				Value: "LPG Retail",
			},
			{
				Key:   "Pembayar",
				Value: "1000013 - Customer SP",
			},
		},
		StartFrom: 0.35,
	})

	data = append(data, Block{
		Title:          "Detail Produk",
		ShowTitle:      true,
		Type:           BLOCK_TYPE_TABLE,
		HideBackground: true,
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

	data = append(data, Block{
		Title:          "Log Aktivitas",
		ShowTitle:      true,
		Type:           BLOCK_TYPE_TABLE,
		HideBackground: true,
		TableData: TableData{
			Rows: [][]string{
				{"No", "Tanggal", "Nama", "Peran", "Aksi", "Deskripsi"},
				{"1", "30/01/2025 10:00:00", "Andre", "Maker", "Send for Approval", ""},
				{"2", "30/01/2025 10:00:00", "Naufal", "Checker", "Approve", ""},
				{"3", "30/01/2025 10:00:00", "Satya", "Checker", "Approve", ""},
				{"4", "30/01/2025 10:00:00", "Arya", "Signer", "Approve", ""},
			},
			ColSize:     []float64{0.07, 0.25, 0.13, 0.15, 0.2, 0.2},
			LastRowBold: true,
		},
	})

	var buf bytes.Buffer

	pdf, err := h.GenerateReport(trxInfo, data)
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

func (s *PDFHandler) GenerateReport(trxInfo TransactionInfo, data []Block) (*gofpdf.Fpdf, error) {
	var paper Paper = GetPaperA4()

	pdf := gofpdf.New("P", "mm", "A4", "assets/font")

	pdf.AddUTF8Font("BRIDigital-Light", "", "BRIDigitalText-Light.ttf")
	pdf.AddFont("BRIDigital", "", "BRIDigitalText-Regular.json")
	pdf.AddFont("BRIDigital", "B", "BRIDigitalText-SemiBold.json")
	pdf.AddFont("BRIDigitalLogo", "B", "BRIDigitalText-SemiBold.json")

	pdf.SetMargins(paper.MarginSetup.XMargin, paper.MarginSetup.YMargin, paper.MarginSetup.XMargin)
	pdf.SetAutoPageBreak(true, 10)

	pdf.SetHeaderFunc(func() {
		pdf.ImageOptions("./assets/images/report-header.png", 0, 0, 210, 0, false, gofpdf.ImageOptions{
			ReadDpi:   false,
			ImageType: "",
		}, 0, "")

		if pdf.PageCount() > 1 {
			pdf.SetFont("BRIDigital", "", 10)
			_, lineHt := pdf.GetFontSize()
			pdf.Text(paper.RectSetup.InnerX, paper.RectSetup.InnerY, fmt.Sprintf("Laporan Transaksi - %s", trxInfo.TransactionId))
			pdf.SetY(paper.RectSetup.InnerY + lineHt + 8.0)
		}

		s.addWatermark(pdf, paper)
	})

	pdf.SetFooterFunc(func() {
		footerContainerHeight := paper.FooterSetup.RectHeight
		spaceX := 16.0
		colSize := []float64{0.7, 0.3}
		pdf.AliasNbPages("{nb}")

		x := paper.RectSetup.X + 2.0
		y := paper.RectSetup.H + paper.RectSetup.Y - footerContainerHeight

		// Left side
		pdf.SetXY(x, y)
		pdf.SetFont("BRIDigital", "", 10)

		_, lineHeight := pdf.GetFontSize()
		lineHeight *= 2
		pdf.MultiCell(paper.RectSetup.W*colSize[0]-spaceX, lineHeight, "Terima kasih telah bertransaksi menggunakan Qlola BRI, bila menemui kendala silakan hubungi kami di 1500001 atau qlola@bri.co.id", "", "L", false)

		// Right side
		currentTime := time.Now()
		formattedTime := currentTime.Format("02/01/2006 15:04:05")
		pdf.SetXY(x+paper.RectSetup.W*colSize[0]+spaceX, y+(lineHeight-1))
		pdf.CellFormat(paper.RectSetup.W*colSize[1]-spaceX, footerContainerHeight, fmt.Sprintf("%v - Halaman %d/{nb}", formattedTime, pdf.PageNo()), "", 0, "R", false, 0, "")

		pdf.Ln(5)
		pdf.SetFillColor(16, 47, 50)
		pdf.Rect(0, paper.PaperSize.Height-paper.FooterSetup.RectHeight, paper.PaperSize.Width, paper.FooterSetup.RectHeight, "F")
		pdf.Ln(2)
	})

	pdf.AddPage()

	lastY := drawReportHeader(pdf, paper, trxInfo)

	if trxInfo.Message != "" && trxInfo.Status != "Sukses" {
		pdf.SetX(paper.RectSetup.InnerX)
		pdf.SetFont("BRIDigital", "", 11)

		_, lineHt := pdf.GetFontSize()
		containerYStart := lastY
		containerHeight := lineHt*2.0 + 3.0

		// x mark icon image
		pdf.ImageOptions("./assets/images/icon-x.png", paper.RectSetup.InnerX+2.0, containerYStart+2.0, 6.0, 6.0, false, gofpdf.ImageOptions{
			ReadDpi:   false,
			ImageType: "",
		}, 0, "")

		pdf.SetAlpha(0.16, "Normal")
		pdf.SetFillColor(198, 40, 40) // #c62828
		pdf.RoundedRect(paper.RectSetup.InnerX, containerYStart, paper.RectSetup.InnerW, containerHeight, 2.0, "1234", "F")
		pdf.SetAlpha(1.0, "Normal")
		pdf.SetTextColor(0, 0, 0) // white

		pdf.Text(paper.RectSetup.InnerX+10, containerYStart+(lineHt+3.0), trxInfo.Message)

		pdf.SetTextColor(0, 0, 0) // reset to black
		pdf.SetY(pdf.GetY() + containerHeight + (lineHt + 12.0))
	} else {
		pdf.SetY(pdf.GetY() + 12.0)
	}

	addThreeColumnInfo(pdf, paper, trxInfo)
	drawContentsReport(pdf, paper, data)

	return pdf, nil
}

func drawReportHeader(pdf *gofpdf.Fpdf, paper Paper, trxInfo TransactionInfo) float64 {
	pdf.SetXY(paper.RectSetup.InnerX, paper.RectSetup.InnerY)

	pdf.SetFont("BRIDigital", "B", 12)
	_, lineHt := pdf.GetFontSize()

	pdf.SetTextColor(24, 24, 24)
	pdf.Text(pdf.GetX(), pdf.GetY(), "Laporan Transaksi")
	pdf.SetFont("BRIDigital-Light", "", 12)
	pdf.Text(pdf.GetX(), pdf.GetY()+lineHt+2.0, trxInfo.TransactionId)

	/*
	 * The badge total height is first column line height - 2mm padding top and
	 * bottom and put it centered with first line text. The row height is lineHt
	 * + 4mm (2mm top and 2mm bottom)
	 */
	badgeHeight := (lineHt * 2.0) + 2.0

	pdf.SetFont("BRIDigital", "", 12)
	strWidth := pdf.GetStringWidth(trxInfo.Status) + 6.0 // padding 4mm kiri kanan
	badgePositionX := pdf.GetX() + paper.RectSetup.InnerW - strWidth

	switch trxInfo.Status {
	case "Sukses":
		pdf.SetFillColor(210, 233, 218) // #d2e9da
		pdf.RoundedRect(badgePositionX, pdf.GetY()-3.0, strWidth, badgeHeight, 2, "1234", "F")

		pdf.SetTextColor(6, 100, 40) // #066428
	case "Gagal":
		pdf.SetAlpha(0.16, "Normal")
		pdf.SetFillColor(205, 13, 19) // #cd0d13
		pdf.RoundedRect(badgePositionX, pdf.GetY()-3.0, strWidth, badgeHeight, 2, "1234", "F")
		pdf.SetAlpha(1, "Normal")

		pdf.SetTextColor(121, 11, 15) // #790b0f
	default:
		pdf.SetFillColor(255, 243, 205) // #fff3cd
	}

	pdf.Text(badgePositionX+3.0, pdf.GetY()+lineHt-1.0, trxInfo.Status)

	// Reset text color
	pdf.SetTextColor(0, 0, 0)

	// Return last Y position
	return pdf.GetY() + badgeHeight + 2.0
}

func addThreeColumnInfo(pdf *gofpdf.Fpdf, paper Paper, trxInfo TransactionInfo) {
	startY := pdf.GetY() + 8.0

	drawSourceOfFunds(pdf, paper, trxInfo, startY)
	drawTransactionType(pdf, paper, trxInfo, startY)
	drawTotalPayment(pdf, paper, trxInfo, startY)
}

func drawContentsReport(pdf *gofpdf.Fpdf, paper Paper, data []Block) {
	pdf.SetY(pdf.GetY() + 28.0)

	// Loop through each block and render based on its type
	rectHeight := paper.RectSetup.InnerH

	for i, block := range data {
		pdf.SetX(paper.RectSetup.InnerX)
		pdf.SetY(pdf.GetY())

		// Last attempt
		if block.Title == "Detail Produk" {
			pdf.AddPage()
		}

		if block.ShowTitle {
			startX, startY := paper.RectSetup.InnerX, pdf.GetY()

			if !block.IsTitleInside && !block.HideBackground {
				// 1. Draw semi-transparent background rectangle
				pdf.SetAlpha(0.16, "Normal")
				pdf.SetFillColor(61, 136, 143)
				pdf.Rect(startX, startY, paper.RectSetup.InnerW, paper.HeaderSetup.H, "F")
				pdf.SetAlpha(1, "Normal")
			}

			// 2. Draw text on top (fully opaque)
			pdf.SetFont("BRIDigital", "B", 11)
			pdf.SetTextColor(0, 0, 0)

			if block.IsTitleInside {
				startY += 5.0
			}

			pdf.Text(startX, startY, block.Title)
		}

		switch block.Type {
		case BLOCK_TYPE_ROWS:
			drawBlockRows(pdf, paper, block)
		case BLOCK_TYPE_TABLE:
			drawBlockTable(pdf, paper, block)
		case BLOCK_TYPE_COLS:
			drawBlockCols(pdf, paper, block)
		}

		endMargin := 4.0
		if i < len(data)-1 && pdf.GetY()+endMargin < rectHeight {
			pdf.SetY(pdf.GetY() + endMargin)
			pdf.Line(paper.RectSetup.InnerX, pdf.GetY(), paper.RectSetup.InnerX+paper.RectSetup.InnerW, pdf.GetY())
			pdf.SetY(pdf.GetY() + endMargin*2.0)
		}
	}
}

func drawSourceOfFunds(pdf *gofpdf.Fpdf, paper Paper, data TransactionInfo, startY float64) {
	x := paper.RectSetup.InnerX
	y := startY
	space := 3.0
	offsetCell := 0.0
	marginTop := 5.0
	imageSize := 12.0
	pdf.SetXY(x, y)

	pdf.SetFont("BRIDigital", "B", 12)
	pdf.Text(x+offsetCell, y, "Sumber Dana")

	// BRI Logo
	briLogo := "./assets/images/Icon.png"
	opts := gofpdf.ImageOptions{ImageType: "", ReadDpi: false}
	if _, err := os.Stat(briLogo); err == nil {
		pdf.ImageOptions(briLogo, x, y+marginTop, imageSize, 0, false, opts, 0, "")
	}

	pdf.SetFont("BRIDigital", "B", 11)
	_, fontSize := pdf.GetFontSize()
	ascent := fontSize * 0.7 // kira-kira ascender

	name := fmt.Sprintf("%s | %s", data.DebitAccount, data.DebitAccountName)
	strLen := pdf.GetStringWidth(name)
	rightWidth := (paper.RectSetup.InnerW / 3) - (imageSize + 4.0) // logo width + padding

	if strLen > rightWidth {
		for strLen > rightWidth {
			name = name[:len(name)-1]
			strLen = pdf.GetStringWidth(name + "...")
		}

		name += "..."
	}

	pdf.Text(x+imageSize+space, y+marginTop+ascent, name)

	pdf.SetFont("BRIDigital", "", 11)
	pdf.Text(x+imageSize+space, y+marginTop+ascent+fontSize+2.0, fmt.Sprintf("%s | %s", "BRI", data.DebitAccount))

	// Indonesian Flag
	indonesiaFlag := "./assets/images/Indonesia.png"
	flagSize := 8.0
	if _, err := os.Stat(indonesiaFlag); err == nil {
		pdf.ImageOptions(indonesiaFlag, x+imageSize+space, y+marginTop+ascent+fontSize+4.0, flagSize, 0, false, opts, 0, "")
	}

	pdf.Text(x+imageSize+space+flagSize+2.0, y+marginTop+ascent+(fontSize*2)+4.0, "IDN")
}

func drawTransactionType(pdf *gofpdf.Fpdf, paper Paper, data TransactionInfo, startY float64) {
	x := paper.RectSetup.InnerX + (paper.RectSetup.InnerW / 3)
	y := startY
	space := 3.0
	marginTop := 5.0
	imageSize := 12.0
	pdf.SetXY(x, y)

	pdf.SetFont("BRIDigital", "B", 12)
	pdf.Text(x+1.0, y, "Jenis Transaksi")

	// Transaction Type Logo
	pertaminaLogo := "./assets/images/pertamina.png"
	opts := gofpdf.ImageOptions{ImageType: "", ReadDpi: false}
	if _, err := os.Stat(pertaminaLogo); err == nil {
		pdf.ImageOptions(pertaminaLogo, x, y+marginTop, imageSize, 0, false, opts, 0, "")
	}

	pdf.SetFont("BRIDigital", "B", 11)
	_, fontSize := pdf.GetFontSize()
	ascent := fontSize * 0.7 // kira-kira ascender

	pdf.Text(x+imageSize+space, y+marginTop+ascent, "DO Pertamina")

	pdf.SetFont("BRIDigital", "", 11)
	pdf.Text(x+imageSize+space, y+marginTop+ascent+fontSize+2.0, data.ModuleType)
}

func drawTotalPayment(pdf *gofpdf.Fpdf, paper Paper, data TransactionInfo, startY float64) {
	x := paper.RectSetup.InnerX + 2*(paper.RectSetup.InnerW/3)
	y := startY
	pdf.SetXY(x, y)

	pdf.SetFont("BRIDigital", "B", 12)
	pdf.Text(x, y, "Total Pembayaran")
	pdf.SetFont("BRIDigital", "B", 16)
	pdf.Text(x, y+10.0, data.TotalPayment)
}

func (s *PDFHandler) addWatermark(pdf *gofpdf.Fpdf, paper Paper) {
	pdf.SetAlpha(0.5, "Normal")
	pdf.ImageOptions("./assets/images/newwatermark3.png", paper.RectSetup.X, paper.RectSetup.Y, paper.RectSetup.W, paper.RectSetup.H, false, gofpdf.ImageOptions{
		ReadDpi:   false,
		ImageType: "", // biarkan kosong agar otomatis deteksi dari ekstensi (jpg, png, dll)
	}, 0, "")
	pdf.SetAlpha(1.0, "Normal")
}
