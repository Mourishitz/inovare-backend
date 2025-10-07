package models

const (
	RoleUser int16 = iota + 1
	RoleAdmin
)

var roleNames = map[int16]string{
	RoleUser:  "User",
	RoleAdmin: "Admin",
}

func GetRoleName(role int16) string {
	if name, exists := roleNames[role]; exists {
		return name
	}
	return "Unknown"
}

func IsValidRole(role int16) bool {
	_, exists := roleNames[role]
	return exists
}

func GetAllRoles() map[int16]string {
	return roleNames
}
