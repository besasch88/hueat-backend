package hueat_auth

/*
contextAuthenticatedUser represents a key where the authenticated user information
are stored inside the context of the request.
*/
const contextAuthenticatedUser = "authenticatedUser"

/*
AuthenticatedUser represents an authenticated user in the webapp application.
All the information stored here are retrieved by the
JWT in the Authentication header of the request.
*/
type AuthenticatedUser struct {
	ID          string
	Username    string
	Permissions []string
}

/*
List of permissions we can leverage to evaluate if an authenticated user can perform a specific operation
before performing the API logic.
*/
const (
	READ_PRINTER         = "read-printer"
	WRITE_PRINTER        = "write-printer"
	READ_MENU            = "read-menu"
	WRITE_MENU           = "write-menu"
	READ_MY_TABLES       = "read-my-tables"
	WRITE_MY_TABLES      = "write-my-tables"
	READ_OTHER_TABLES    = "read-other-tables"
	WRITE_OTHER_TABLES   = "write-other-tables"
	READ_STATISTICS      = "read-statistics"
	DELETE_STATISTICS    = "delete-statistics"
	UPDATE_REFRESH_TOKEN = "update-refresh-token"
)
