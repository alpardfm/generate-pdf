package receipt_generator

import (
	"errors"
	"log"
	"strings"
	"time"

	"github.com/jung-kurt/gofpdf"
)

func Receipt(transactionId, transactionName string, isFirstHeaderVisible bool, listModel []ListModelData, paperSize string) (gofpdf.Pdf, error) {
	var borderRadius = 3.0
	var paper Paper
	log.Println(paper)

	if strings.ToUpper(paperSize) == "A5" {
		paper = GetPaperA5()
	} else if strings.ToUpper(paperSize) == "A4" {
		paper = GetPaperA4()
	} else {
		log.Println("Paper Size not Supported")
		return nil, errors.New("Paper Size not Supported")
	}

	pdf := gofpdf.New("P", "mm", paperSize, "assets/font")

	pdf.AddFont("BRIDigital", "", "BRIDigitalText-Regular.json")
	pdf.AddFont("BRIDigital", "B", "BRIDigitalText-SemiBold.json")
	//pdf.AddFont("LatoLight", "", "Lato-Light.json")
	pdf.AddFont("BRIDigitalLogo", "B", "BRIDigitalText-SemiBold.json")

	//var marginH = paper.MarginSetup.TMargin
	var lineHt = paper.LineHt

	location, err := time.LoadLocation("Asia/Jakarta")
	if err != nil {
		log.Println(err)
	}

	exportedAt := time.Now().In(location).Format("02/01/2006, 15:04:05 (GMT Z07)")
	exportedAt = strings.Replace(exportedAt, "+0", "+", -1)

	pagew, pageh := pdf.GetPageSize()

	pdf.SetMargins(paper.MarginSetup.LMargin*3, paper.MarginSetup.TMargin*3, paper.MarginSetup.RMargin*3)
	lf, tp, _, _ := pdf.GetMargins()

	pdf.SetHeaderFuncMode(func() {
		pdf.SetMargins(paper.MarginSetup.LMargin*3, paper.MarginSetup.TMargin*3, paper.MarginSetup.RMargin*3)
		x, y := pdf.GetXY()
		r, g, b := pdf.GetFillColor()
		log.Printf("RGB: %d%d%d", r, g, b)
		pdf.SetFillColor(235, 235, 235)
		pdf.Rect(0, 0, pagew, pageh, "F")
		pdf.SetFillColor(r, g, b)
		pdf.SetXY(x, y)
		pdf.SetFillColor(255, 255, 255)
		////pdf.Rect(12.5, 12.5, 135.3, 197.3, "F")
		pdf.RoundedRect(paper.RectSetup.X, paper.RectSetup.Y, paper.RectSetup.W, paper.RectSetup.H, 3, "1234", "F")
		pdf.SetFillColor(r, g, b)
		pdf.SetXY(x, y)
		pdf.SetFont("BRIDigital", "B", 13)
		pdf.SetAlpha(0.15, "Normal")

		pdf.SetTextColor(235, 235, 235)

		for i := paper.TransformSetup.I.Min; i < paper.TransformSetup.I.Max; i++ {
			for j := paper.TransformSetup.J.Min; j < paper.TransformSetup.J.Max; j++ {
				pdf.TransformBegin()
				pdf.TransformRotate(paper.TransformSetup.Angle, paper.TransformSetup.X.A+(i*paper.TransformSetup.X.B), paper.TransformSetup.Y.A+(j*paper.TransformSetup.Y.B))
				pdf.Text(paper.TransformSetup.TextX.A+(i*paper.TransformSetup.TextX.B), paper.TransformSetup.TextY.A+(j*paper.TransformSetup.TextY.B), "Qlola Cash Management")
				pdf.TransformRotate(0, paper.TransformSetup.X.A+(i*paper.TransformSetup.X.B), paper.TransformSetup.Y.A+(j*paper.TransformSetup.Y.B))
				pdf.Text(paper.TransformSetup.TextX.A+(i*paper.TransformSetup.TextX.B), paper.TransformSetup.TextY.A+(j*paper.TransformSetup.TextY.B), "Qlola Cash Management")
				pdf.TransformEnd()
			}
		}

		pdf.SetAlpha(1, "Normal")
	}, true)

	pdf.AddPage()
	pdf.SetMargins(paper.MarginSetup.LMargin*3, paper.MarginSetup.TMargin*3, paper.MarginSetup.RMargin*3)

	pdf.SetTextColor(0, 82, 156)
	pdf.SetFont("BRIDigital", "B", 20)

	pdf.Image("assets/qlola_cash_management.png", paper.WLogo1.X, paper.WLogo1.Y, paper.WLogo1.W, paper.WLogo1.H, paper.WLogo1.Flow, "", 0, "")
	pdf.Image("assets/bri.png", pagew-paper.WLogo2.X, paper.WLogo2.Y, paper.WLogo2.W, paper.WLogo2.H, paper.WLogo2.Flow, "", 0, "")

	//html := pdf.HTMLBasicNew()
	pdf.SetFooterFunc(func() {
		pdf.SetY(paper.FooterSetup.Y)
		pdf.SetFont("BRIDigital", "", paper.FooterSetup.FontSize)
		pdf.SetWordSpacing(paper.FooterSetup.WordSpacing)
		pdf.CellFormat(0, lineHt, "This Receipt is printed by BRI Qlola Cash Management System", "0", 0, "CM", false, 0, "")
		pdf.SetY(pdf.GetY() + 4)
		pdf.CellFormat(0, lineHt, "PT Bank Rakyat Indonesia (Persero) Tbk", "0", 0, "CM", false, 0, "")
		pdf.SetY(pdf.GetY() + 4)
		pdf.CellFormat(0, lineHt, exportedAt, "0", 0, "CM", false, 0, "")
	})

	pdf.SetY(tp + paper.MarginSetup.TMargin*3)
	pdf.SetDrawColor(224, 224, 224)
	pdf.SetTextColor(66, 66, 66)

	pdf.SetY(pdf.GetY() - paper.TransactionTextSetup.UpperSpace)
	pdf.SetFont("BRIDigital", "B", paper.TransactionTextSetup.FontSize)

	pdf.CellFormat(0, lineHt, transactionId, "0", 0, "LM", false, 0, "")
	pdf.CellFormat(1, lineHt, "", "0", 0, "CM", false, 0, "")
	pdf.CellFormat(0, lineHt, strings.ToUpper(transactionName), "0", 0, "RM", false, 0, "")

	y := paper.TransactionTextSetup.LowerSpace
	if len(listModel) > 0 {
		pdf.SetY(tp + y)
		for i, datum := range listModel {
			spaceLeft := pageh - (pdf.GetY() + lineHt)
			if spaceLeft < paper.BottomSetup.BottomLimit {
				pdf.SetY(pageh - paper.BottomSetup.BottomLimitMinus)
				pdf.SetFont("BRIDigital", "B", paper.BottomSetup.FontSize)
				pdf.CellFormat(0, lineHt, "...", "0", 0, "CM", false, 0, "")
				pdf.AddPage()
			}
			//log.Printf("pdf.GetX: %f", pdf.GetX())
			//log.Printf("pdf.GetY: %f", pdf.GetY())
			//xHeader, y = pdf.GetXY()
			if !datum.IsTable {
				if len(datum.ModelData) > 0 {
					if !isFirstHeaderVisible {
						if strings.Trim(datum.HeaderData, " ") == "" {
							datum.HeaderData = "-"
						}
						for _, modelDatum := range datum.ModelData {
							spaceLeft = pageh - (pdf.GetY() + lineHt)
							if spaceLeft < paper.BottomSetup.BottomLimit {
								pdf.SetY(pageh - paper.BottomSetup.BottomLimitMinus)
								pdf.SetFont("BRIDigital", "B", paper.BottomSetup.FontSize-2.0)
								pdf.CellFormat(0, lineHt, "...", "0", 0, "CM", false, 0, "")
								pdf.AddPage()
							}
							if modelDatum.IsTotalPayment {
								pdf.Ln(paper.ValueCellSetup.Ln1)
								pdf.SetFont("BRIDigital", "B", paper.TotalPaymentFont.HeaderFontSize)
								pdf.CellFormat(paper.ValueCellSetup.W1, paper.ValueCellSetup.H1, modelDatum.Key, "0", 0, "LM", false, 0, "")
								pdf.CellFormat(paper.ValueCellSetup.W2, paper.ValueCellSetup.H2, "", "0", 0, "CM", false, 0, "")
								pdf.SetFont("BRIDigital", "B", paper.TotalPaymentFont.ValueFontSize)
								pdf.MultiCell(paper.ValueCellSetup.WMultiCell, paper.ValueCellSetup.HMultiCell, modelDatum.Value, "0", "RM", false)
								pdf.SetFont("BRIDigital", "", paper.TotalPaymentFont.HeaderFontSize)
								pdf.Ln(paper.ValueCellSetup.Ln2)
							} else {
								pdf.SetFont("BRIDigital", "", paper.ValueFont.HeaderFontSize)
								pdf.CellFormat(paper.ValueCellSetup.W1, paper.ValueCellSetup.H1, modelDatum.Key, "0", 0, "LM", false, 0, "")
								pdf.CellFormat(paper.ValueCellSetup.W2, paper.ValueCellSetup.H2, "", "0", 0, "CM", false, 0, "")
								pdf.SetFont("BRIDigital", "B", paper.ValueFont.ValueFontSize)
								pdf.MultiCell(paper.ValueCellSetup.WMultiCell, paper.ValueCellSetup.HMultiCell, modelDatum.Value, "0", "RM", false)
								pdf.SetFont("BRIDigital", "", paper.ValueFont.HeaderFontSize)
								pdf.Ln(paper.ValueCellSetup.Ln2)
							}
						}
						if i < len(listModel)-1 {
							pdf.SetY(pdf.GetY() + paper.HeaderSetup.Space1)
							pdf.SetFillColor(0, 82, 156)
							pdf.SetTextColor(255, 255, 255)
							pdf.SetFont("BRIDigital", "B", paper.HeaderSetup.FontSize)
							pdf.CellFormat(paper.HeaderSetup.W, paper.HeaderSetup.H, "",
								"", 0, "C", true, 0, "")
							pdf.Text(paper.HeaderSetup.X, pdf.GetY()+paper.HeaderSetup.Y, listModel[i+1].HeaderData)
							pdf.SetFillColor(255, 255, 255)
							pdf.SetX(lf)
						}
						pdf.SetTextColor(66, 66, 66)
						y = pdf.GetY()
						pdf.SetY(tp + y - paper.HeaderSetup.Space2)
					} else {
						if strings.Trim(datum.HeaderData, " ") == "" {
							datum.HeaderData = "-"
						}
						pdf.SetY(pdf.GetY() + paper.HeaderSetup.Space1)
						pdf.SetFillColor(0, 82, 156)
						pdf.SetTextColor(255, 255, 255)
						pdf.SetFont("BRIDigital", "B", paper.HeaderSetup.FontSize)
						pdf.CellFormat(paper.HeaderSetup.W, paper.HeaderSetup.H, "",
							"", 0, "C", true, 0, "")
						pdf.Text(paper.HeaderSetup.X, pdf.GetY()+paper.HeaderSetup.Y, datum.HeaderData)
						pdf.SetFillColor(255, 255, 255)
						pdf.SetX(lf)
						pdf.SetTextColor(66, 66, 66)

						y = pdf.GetY()
						pdf.SetY(tp + y)

						for _, modelDatum := range datum.ModelData {
							if modelDatum.IsTotalPayment {
								pdf.Ln(paper.ValueCellSetup.Ln1)
								pdf.SetFont("BRIDigital", "B", paper.TotalPaymentFont.HeaderFontSize)
								pdf.CellFormat(paper.ValueCellSetup.W1, paper.ValueCellSetup.H1, modelDatum.Key, "0", 0, "LM", false, 0, "")
								pdf.CellFormat(paper.ValueCellSetup.W2, paper.ValueCellSetup.H2, "", "0", 0, "CM", false, 0, "")
								pdf.SetFont("BRIDigital", "B", paper.TotalPaymentFont.ValueFontSize)
								pdf.MultiCell(paper.ValueCellSetup.WMultiCell, paper.ValueCellSetup.HMultiCell, modelDatum.Value, "0", "RM", false)
								pdf.SetFont("BRIDigital", "", paper.TotalPaymentFont.HeaderFontSize)
								pdf.Ln(paper.ValueCellSetup.Ln2)
							} else {
								pdf.SetFont("BRIDigital", "", paper.TotalPaymentFont.HeaderFontSize)
								pdf.CellFormat(paper.ValueCellSetup.W1, paper.ValueCellSetup.H1, modelDatum.Key, "0", 0, "LM", false, 0, "")
								pdf.CellFormat(paper.ValueCellSetup.W2, paper.ValueCellSetup.H2, "", "0", 0, "CM", false, 0, "")
								pdf.SetFont("BRIDigital", "B", paper.TotalPaymentFont.ValueFontSize)
								pdf.MultiCell(paper.ValueCellSetup.WMultiCell, paper.ValueCellSetup.HMultiCell, modelDatum.Value, "0", "RM", false)
								pdf.SetFont("BRIDigital", "", paper.TotalPaymentFont.HeaderFontSize)
								pdf.Ln(paper.ValueCellSetup.Ln2)
							}
						}
					}
				}
			} else {
				if len(datum.TableData) > 0 {
					header := true // Flag to identify the header row
					footer := false
					for tableIterator, tableRow := range datum.TableData {
						curx, y := pdf.GetXY()
						x := curx

						height := 0.
						_, lineHt := pdf.GetFontSize()

						if pdf.GetY()+height > pageh-paper.BottomSetup.BottomLimit {
							pdf.AddPage()
							y = pdf.GetY()
						}

						for i, txt := range tableRow {
							lines := pdf.SplitLines([]byte(txt), float64(datum.ColumnWidth[i]))
							h := float64(len(lines))*lineHt + 4.*float64(len(lines))
							if h > height {
								height = h
							}
						}
						for i, txt := range tableRow {
							width := float64(datum.ColumnWidth[i])

							if header {
								pdf.SetFillColor(249, 249, 249)

								if i == 0 {
									pdf.RoundedRect(x, y, width, height, borderRadius, "1", "FD")
								} else if i == len(tableRow)-1 {
									pdf.RoundedRect(x, y, width, height, borderRadius, "2", "FD")
								} else {
									pdf.RoundedRect(x, y, width, height, borderRadius, "", "FD")
								}
							} else if footer {
								if i == 0 {
									pdf.RoundedRect(x, y, width, height, borderRadius, "4", "D")
								} else if i == len(tableRow)-1 {
									pdf.RoundedRect(x, y, width, height, borderRadius, "3", "D")
								} else {
									pdf.RoundedRect(x, y, width, height, borderRadius, "", "D")
								}
							} else {
								pdf.SetFillColor(255, 255, 255)
								pdf.Rect(x, y, width, height, "D")
							}

							pdf.MultiCell(width, lineHt+2., txt, "", "AC", false)
							x += width
							pdf.SetXY(x, y)
						}
						pdf.SetXY(curx, y+height)

						if header {
							header = false // Reset the flag after the first row
						}

						if tableIterator == len(datum.TableData)-2 {
							footer = true
						}
					}

					pdf.SetTextColor(66, 66, 66)
					y = pdf.GetY()
					pdf.SetY(tp + y - paper.HeaderSetup.Space2)

					// Create line break
					pdf.Ln(paper.HeaderSetup.Space1 * 2)
				}
			}
		}
	}

	return pdf, nil
}
