package userstatus

// UserStatus user status
type UserStatus string

const (
	Active   UserStatus = "active"
	Inactive UserStatus = "inactive"
	Blocked  UserStatus = "blocked"
)
