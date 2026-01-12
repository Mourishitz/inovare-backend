package utils

func GetName[T comparable](m map[T]string, key T) string {
	if name, ok := m[key]; ok {
		return name
	}
	return "Unknown"
}

func IsValid[T comparable](m map[T]string, key T) bool {
	_, ok := m[key]
	return ok
}

func GetAll[T comparable](m map[T]string) map[T]string {
	return m
}

