package receipt_generator

type ModelData struct {
	Key            string
	Value          string
	LenChar        int
	Height         int
	CutDot         bool
	IsTotalPayment bool
}

type ListModelData struct {
	HeaderData  string
	ModelData   []ModelData
	IsTable     bool
	TableData   [][]string
	ColumnWidth []uint32
}

type Paper struct {
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
}

type FontSize struct {
	ValueFontSize  float64
	HeaderFontSize float64
}

type MarginSetup struct {
	LMargin float64
	TMargin float64
	RMargin float64
}

type RectSetup struct {
	X float64
	Y float64
	W float64
	H float64
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
	X    float64
	Y    float64
	W    float64
	H    float64
	Flow bool
}

type FooterSetup struct {
	Y           float64
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
