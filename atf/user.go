/*
 * user.go 
 *
 * History:
 *  1   Jun11 MR Initial version, limited testing
 */

package atf

import (
	"encoding/json"
	"fmt"
	"strings"
)

/*
 * User
 */
type User struct {
	Username  string // username 
	Password  string // password (should always be hashed, use CreateUser()! ) 
	Name      string // full name
    Role      string // user role, limited to (guest, user, admin)
	Email     string // e-mail address
}

// A list of allowed user roles
var AllowedRoles = []string{"admin", "user", "guest"}

// String representation of the User 
func (u *User) String() (s string) {
	s = fmt.Sprintf("%s [%s]: %s %q", u.Username, u.Email, u.Password, u.Role)
	return s
}

// Convert a User into JSON representation
func (u *User) Json() (string, error) {
	b, err := json.Marshal(u)
	if err != nil {
		return "JSON marshal error - User", err
	}
	return string(b[:]), err
}

// Handle roles for users, since this is a bit tricky, we don't allowed direct
// management of roles, but via method that automates procedure and does some
// basic error checking.
func (u *User) SetRole(role string) error {

    role = strings.ToLower(role)

    if role == "administrator" { role = "admin" }

    for _, r := range AllowedRoles {
        if role == r {
            u.Role = role
            return nil
        }
    }
    return fmt.Errorf("User role value not valid.")
}

// Create a new user with username, password and role as mandatory information.
// This one is used to create a non-existing user (in some sort of DB), so role
// is a vital information about the user.
func CreateNewUser(username, password, role string) (*User, error) {
	p := new(Password)
	p.Set(password)
    u := &User{username, p.Get(), "", "user", ""}
    if err := u.SetRole(role); err != nil { return nil, err }
	return u, nil
}

// Create a user with username and password. This one is used to authenticate
// the existing user (role is stored in some DB...)
func CreateUser(username, password string) *User {
	p := new(Password)
	p.Set(password)
    return &User{username, p.Get(), "", "user", ""}
}
