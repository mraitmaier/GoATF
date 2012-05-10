/*
 * user.go 
 *
 * History:
 *  0.1   Jun11 MR Initial version, limited testing
 */

package atf

import (
	"fmt"
	"os"
	"strings"
	"json"
)

/*
 * Role - type (enum) defining user roles
 */
type Role int

const (
	UnknownRole Role = iota
	GuestRole
	UserRole
	TesterRole
	DeveloperRole
	ManagerRole
	AdminRole
)

/*
 * Role.String - a method returning string representation of the Role value
 */
func (r Role) String() string {
	switch r {
	case AdminRole:
		return "Admin"
	case UserRole:
		return "User"
	case TesterRole:
		return "Tester"
	case DeveloperRole:
		return "Developer"
	case ManagerRole:
		return "Manager"
	case GuestRole:
		return "Guest"
	}
	return "Unknown Role"
}

/*
 * RoleValue - a function that converts the string value into Role value
 */
func RoleValue(s string) Role {
	switch strings.ToLower(s) {
	case "admin", "administrator":
		return AdminRole
	case "manager":
		return ManagerRole
	case "developer":
		return DeveloperRole
	case "tester":
		return TesterRole
	case "user":
		return UserRole
	case "guest":
		return GuestRole
	}
	return UnknownRole
}

/*
 * User
 */
type User struct {
	Username  string // username 
	*Password        // password (is hidden so it can be hashed always) 
	Hint      string // password hint
	Name      string // full name
	Role      Role   // user role (guest, user, admin)
	Email     string // e-mail address
}

/*
 * User.String - return string representation of the User
 */
func (u *User) String() (s string) {
	s = fmt.Sprintf("%s [%s]: %s %q", u.Username, u.Email,
		u.Password.Get(), u.Role.String())
	return s
}

/*
 * User.Json - convert a user into JSON representation
 */
func (u *User) Json() (string, os.Error) {
	b, err := json.Marshal(u)
	if err != nil {
		return "", err
	}
	return string(b[:]), err
}

/*
 * CreateUser - create a user with username and password 
 *
 * Username and password are the mandatory parameters to create a new user.
 */
func CreateUser(username, password string) *User {
	p := new(Password)
	p.Set(password)
	return &User{username, p, "", "", UnknownRole, ""}
}
