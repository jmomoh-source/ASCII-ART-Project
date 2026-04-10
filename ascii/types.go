package ascii

const (
	CharHeight = 8
	DefaultWidth = 80
	DefaultAlignment = "left"
	DefaultBanner = "standard"
)

const (
	AlignLeft    = "left"
	AlignRight   = "right"
	AlignCenter  = "center"
	AlignJustify = "justify"
)

func IsValidAlignment(align string) bool {
	return align == AlignLeft || align == AlignRight || align == AlignCenter || align == AlignJustify
}
