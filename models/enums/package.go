package enums

const (
	PackageBag = iota + 1
	PackageBox
)

var PackageNames = map[int16]string{
	PackageBag: "Bag",
	PackageBox: "Box",
}
