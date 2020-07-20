package main

// map of users and their roles - again, never use something like this in prod.
var userRoles = map[string]string{
	"joey@jojo.com":    "admin",
	"average@john.com": "user",
}

// hasAdminPermission is somewhat like "middleware" which
// checks the role of received of user in our mock database.
func hasAdminPermissions(pid string) bool {
	if _, ok := userRoles[pid]; ok {
		return true
	}

	return false
}
