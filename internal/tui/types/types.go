package types

type Page int

const (
	MenuPage Page = iota
	MergePage
	SplitPage
	EncryptPage
	DecryptPage
	Img2PdfPage
	ExtractImgsPage
	Doc2PdfPage
)
