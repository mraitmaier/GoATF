/*
 * password.go 
 *
 * History:
 *  0.1   Jun11 MR Initial version, limited testing
 */

package atf

import ("fmt"
        "crypto/md5"
        "rand"
        )

/*
 * Password - a type defining a user password
 *
 * This struct has only one field. We hide this field into struct, so it 
 * cannot be accessible from outside. We define only operations on password:
 * Get, Set and Cmp (compare) to arbitrary string. From outside we deal with
 * strings, while internally everything is done using bytes.
 * Also, when password is set, it is hashed and stored that way. All operations
 * on passwords are performed using hashed passwords.
 */
type Password struct {
    pwd []byte
}

/*
 * Password.Set - set a password
 *
 * Password will, of course, be hashed.
 */
func (p *Password) Set(passwd string) { p.pwd = p.hashPwd(passwd) }

/*
 * Password.Get - return a stored hashed password
 */
func (p *Password) Get() string {
    return fmt.Sprintf("%x", (p.pwd))
}

/*
 * Password.Cmp - compare arbitrary password to the one stored 
 */
func (p *Password) Cmp(passwd string) bool {
    hashed := p.removeSalt(p.hashPwd(passwd))
    stored := p.removeSalt(p.pwd)
    // we compare two []byte: lengths must be the same and every value  of the
    // buffer must the same
    if len(hashed) != len(stored) { return false }
    for cnt, val := range hashed { if val != stored[cnt] { return false } }
    return true
}

/*
 * Password.addSalt - add a 'salt' to password hashing function
 *
 * Salt is a series of bytes that are appended to hashed password. The reason 
 * is, of course, that password is even harder to crack. 
 * Note that this operation is defined as a method; there's a simple reason: we
 * don't want this operation to be used elsewhere in the code. It is bound to
 * Password type and that's it.
 */
const saltSize int = 8
func (p *Password) addSalt() []byte {
    var b [saltSize]byte
    for cnt := 0; cnt < saltSize; cnt++ { b[cnt] = byte(rand.Intn(255)) }
    return b[:]
}

/*
 * Password.removeSalt - returns the hashed password without the 'salt' part
 *
 * Note that this operation is defined as a method; there's a simple reason: we
 * don't want this operation to be used elsewhere in the code. It is bound to
 * Password type and that's it.
 */
func (p *Password) removeSalt(saltedPwd []byte) []byte {
    return saltedPwd[saltSize:]
}

/*
 * Password.hashPwd - hash a string password using MD5
 *
 * Note that this operation is defined as a method; there's a simple reason: we
 * don't want this operation to be used elsewhere in the code. It is bound to
 * Password type and that's it.
 */
const numOfHashIteration int = 1017
func (p *Password) hashPwd(passwd string) []byte {
    h := md5.New()
    // we need to convert a string into slice of bytes
    b := []byte(passwd)
    // let's iterate the MD5 hash more than a 1000-times
    for cnt := 0; cnt < numOfHashIteration; cnt++ {
        h.Write(b)
        b = h.Sum()
    }
    return append(p.addSalt(), b...)
}

