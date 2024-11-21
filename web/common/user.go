package common

import (
	"crypto"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"nbientry/web"
	"net/http"
	"os/exec"
	"regexp"
	"strings"
	"time"
)

var ValidNames = regexp.MustCompile(`^([0-9a-z\_\.\-]+)@([0-9a-z\.\-\_]+)\.\w+$`)
var _ = web.Migrate(User{})

type User struct {
	Id       int    `gorm:"primaryKey"`
	Email    string `gorm:"unique"`
	FullName string
	Phone    string
	Password string
	Image    string
	JWT      string
	GitHash  string
	Admin    bool
}

func get_git_hash() string {
	cmd := exec.Command("git", "rev-parse", "HEAD")
	out, err := cmd.Output()
	if err != nil {
		panic(err)
	}
	return string(out)
}

var GIT_HASH string = get_git_hash()

func (u User) HashPassword(pass string) string {
	// Generate a number to use as a salt & to make the password more difficult to crack without reading source code
	cycle := uint8(len(u.Email))
	for _, char := range pass {
		cycle <<= 1
		cycle ^= uint8(char)
	}

	if cycle == 0 {
		cycle = uint8(len(u.Email))
	}

	salt := []byte(fmt.Sprintf("%s%d%s", u.Email, cycle, pass))

	// Hash the password [cycle] times
	for i := 0; i < int(cycle); i++ {
		sha256 := crypto.SHA256.New()
		sha256.Write(salt)
		salt = sha256.Sum(nil)
	}

	return hex.EncodeToString(salt)
}

func (u *User) GenerateJWT() error {
	if u.Email == "" || u.Password == "" {
		return fmt.Errorf("missing email or password")
	}

	check := web.GetFirst(User{Email: u.Email, Password: u.Password})

	if check.Email == "" || check.Email != u.Email {
		return fmt.Errorf("bad credentials")
	}

	u.Id = check.Id
	u.GitHash = GIT_HASH
	u.JWT = base64.StdEncoding.EncodeToString(
		[]byte(fmt.Sprintf(`{"name":"%s","timestamp":%d}`, u.Email, time.Now().Unix())),
	)

	web.Save(u)
	return nil
}

func (u *User) Exists() bool {
	if u.Email == "" {
		return false
	}

	u.Email = strings.ToLower(u.Email)

	check := web.GetFirst(User{Email: u.Email})
	if check.Email == "" {
		return false
	}

	u.Email = check.Email
	u.FullName = check.FullName
	return true
}

func (u *User) ValidateEmail() error {
	if u.Exists() {
		return nil
	}

	if len(u.Email) < 3 {
		return fmt.Errorf("email too short")
	}

	if len(u.Email) > 128 {
		return fmt.Errorf("email too long")
	}

	if !ValidNames.MatchString(u.Email) {
		return fmt.Errorf("email invalid")
	}

	return nil
}

func (u *User) ValidatePassword(pass string) error {
	err := u.ValidateEmail()
	if err != nil {
		return err
	}

	u.Password = pass

	if u.Password == "" {
		return fmt.Errorf("password missing")
	}

	if len(u.Password) < 8 {
		return fmt.Errorf("password too short")
	}

	if !u.Exists() {
		return nil
	}

	u.Password = u.HashPassword(pass)
	check := web.GetFirst(User{Email: u.Email, Password: u.Password})
	if check.Email == "" || check.Email != u.Email {
		u.Password = ""
		return fmt.Errorf("bad credentials")
	}

	u.Admin = check.Admin
	u.GitHash = GIT_HASH
	u.JWT = check.JWT
	u.Image = check.Image
	u.Phone = check.Phone
	u.FullName = check.FullName

	return nil
}

func (u User) Cookie() *http.Cookie {
	return &http.Cookie{
		Name:  "jwt",
		Value: u.JWT,
		Path:  "/",
	}
}

func CookieAuth(w http.ResponseWriter, r *http.Request) User {
	cookies := r.Cookies()
	user := User{}

	for _, cookie := range cookies {
		if cookie.Name != "jwt" {
			continue
		}

		user = web.GetFirst(User{JWT: cookie.Value, GitHash: GIT_HASH})
		if user.Email != "" {
			break
		}
	}

	return user
}
