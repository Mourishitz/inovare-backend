package preferred

const (
	PaddedBra = iota + 1
	NonPaddedBra
	BothBra
)

var BraNames = map[int]string{
	PaddedBra:    "Com Bojo",
	NonPaddedBra: "Sem Bojo",
	BothBra:      "Ambos",
}
