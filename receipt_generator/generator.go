package receipt_generator

//-- this function will print receipt to A5 paper --//

func GetPaperA5() Paper {
	var paper Paper

	paper = Paper{
		MarginSetup: MarginSetup{
			LMargin: 6.35,
			TMargin: 6.35,
			RMargin: 6.35,
		},
		RectSetup: RectSetup{
			X: 12.7,
			Y: 12.7,
			W: 148 - 25.4,
			H: 210 - 25.4,
		},
		TransformSetup: TransformSetup{
			X: struct {
				A float64
				B float64
			}{A: 6.35, B: 30},
			Y: struct {
				A float64
				B float64
			}{A: 15, B: 10.7},
			TextX: struct {
				A float64
				B float64
			}{A: 6.35, B: 30},
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
				Max: 5,
			},
			J: struct {
				Min float64
				Max float64
			}{
				Min: 0.9,
				Max: 18.5,
			},
		},
		WLogo1: LogoSetup{
			X:    20.05,
			Y:    19.05,
			W:    50,
			H:    5.5,
			Flow: false,
		},
		WLogo2: LogoSetup{
			X:    35.2,
			Y:    19.05,
			W:    15,
			H:    5.5,
			Flow: false,
		},
		LineHt: 5.5,
		TotalPaymentFont: FontSize{
			ValueFontSize:  12,
			HeaderFontSize: 12,
		},
		ValueFont: FontSize{
			ValueFontSize:  10,
			HeaderFontSize: 10,
		},
		FooterSetup: FooterSetup{
			Y:           -27.5,
			WordSpacing: 1,
			FontSize:    7,
		},
		TransactionTextSetup: TransactionTextSetup{
			FontSize:   13,
			UpperSpace: 8,
			LowerSpace: 20,
		},
		BottomSetup: BottomSetup{
			BottomLimit:      40,
			BottomLimitMinus: 35,
			FontSize:         15,
		},
		ValueCellSetup: CellSetup{
			W1:         10,
			W2:         1,
			WMultiCell: 99.3,
			H2:         5.5,
			H1:         5.5,
			HMultiCell: 5.5,
			Ln1:        1.5,
			Ln2:        0.1,
		},
		HeaderSetup: HeaderSetup{
			Space1:   1,
			Space2:   12.5,
			W:        0,
			H:        5.5,
			X:        22.4,
			Y:        3.5,
			FontSize: 7,
		},
	}

	return paper
}

//-- this function will print receipt to A4 paper --//

func GetPaperA4() Paper {
	var paper Paper

	paper = Paper{
		MarginSetup: MarginSetup{
			LMargin: 7.9375,
			TMargin: 7.9375,
			RMargin: 7.9375,
		},
		RectSetup: RectSetup{
			X: 15.875,
			Y: 15.875,
			W: 210 - 31.75,
			H: 297 - 31.75,
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
		WLogo1: LogoSetup{
			X:    24.8125,
			Y:    23.8125,
			W:    52,
			H:    5.5,
			Flow: false,
		},
		WLogo2: LogoSetup{
			X:    40.8625,
			Y:    23.8125,
			W:    15,
			H:    5.5,
			Flow: false,
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
			Y:           -32.5,
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

	return paper
}
