package enums

const (
	RomanticStyle = iota + 1
	SensualStyle
	ElegantStyle
	FunnyStyle
	ClassicStyle
	ModernStyle
)

var StyleNames = map[int]string{
	RomanticStyle: "Romântico",
	SensualStyle:  "Sensual",
	ElegantStyle:  "Elegante",
	FunnyStyle:    "Divertido",
	ClassicStyle:  "Classico",
	ModernStyle:   "Moderno",
}
