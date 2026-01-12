package preferred

const (
	Laced int16 = iota + 1
	Seamless
	Both
)

var ModelNames = map[int16]string{
	Laced:    "Com renda",
	Seamless: "Lisos",
	Both:     "Ambos",
}
