package contextkey

type ContextKey string

const UserIDCtxKey ContextKey = "userID"

type Role string

const (
	AdminRole    Role = "moderator"
	EmployeeRole Role = "employee"
)
