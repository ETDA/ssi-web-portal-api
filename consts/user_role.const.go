package consts

type UserRole string

const (
	UserRoleAdmin  UserRole = "ADMIN"
	UserRoleMod    UserRole = "MOD"
	UserRoleMember UserRole = "MEMBER"
)

var UserRoles = []string{string(UserRoleAdmin), string(UserRoleMod), string(UserRoleMember)}
