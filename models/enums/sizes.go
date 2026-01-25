package enums

const (
	SizePP = iota + 1
	SizeP
	SizeM
	SizeG
	SizeGG
	SizeXG
)

var SizeNames = map[int16]string{
	SizePP: "PP",
	SizeP:  "P",
	SizeM:  "M",
	SizeG:  "G",
	SizeGG: "GG",
	SizeXG: "XG",
}
