package enums

const (
	WhiteColor = iota + 1
	BlackColor
	RedColor
	LilacColor
	PurpleColor
	MarsalaColor
	RoseColor
	BlueColor
	YellowColor
	GreenColor
	GrayColor
	TerracottaColor
	BeigeColor
	RubyColor
)

var ColorNames = map[int16]string{
	WhiteColor:      "Branco",
	BlackColor:      "Preto",
	RedColor:        "Vermelho",
	LilacColor:      "Lilás",
	PurpleColor:     "Roxo",
	MarsalaColor:    "Marsala",
	RoseColor:       "Rosê (Rosa Claro)",
	BlueColor:       "Azul",
	YellowColor:     "Amarelo",
	GreenColor:      "Verde",
	GrayColor:       "Cinza",
	TerracottaColor: "Terracota",
	BeigeColor:      "Bege",
	RubyColor:       "Rubi",
}
