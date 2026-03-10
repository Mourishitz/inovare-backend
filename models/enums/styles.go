package enums

const (
	RomanticStyle = iota + 1
	SensualStyle
	ElegantStyle
)

var StyleNames = map[int]string{
	RomanticStyle: "Romântico",
	SensualStyle:  "Sensual",
	ElegantStyle:  "Elegante",
}
