package enums

const (
	RoleUser int16 = iota + 1
	RoleAdmin
)

var RoleNames = map[int16]string{
	RoleUser:  "User",
	RoleAdmin: "Admin",
}
