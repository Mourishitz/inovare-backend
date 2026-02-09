package enums

const (
	Collections = iota + 1
	Nightgown
	LongPajama
	Body
	Babydoll
	Bikini
)

var AllowedModelsNames = map[int16]string{
	Collections: "Conjuntos",
	Nightgown:   "Camisola",
	LongPajama:  "Pijama Longo",
	Body:        "Bodys",
	Babydoll:    "Baby Doll",
	Bikini:      "Biquíni",
}
