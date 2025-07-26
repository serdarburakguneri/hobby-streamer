package constants

const (
	RoleAdmin = "admin"
	RoleUser  = "user"
)

var AllowedRoles = map[string]struct{}{
	RoleAdmin: {},
	RoleUser:  {},
}

func IsValidRole(r string) bool {
	_, ok := AllowedRoles[r]
	return ok
}
