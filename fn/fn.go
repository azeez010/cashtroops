package fn

import (
	"errors"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"math/rand"
	"regexp"
	"strings"
	"time"
)

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
var ErrInvalidEmail = errors.New("invalid email address")
var userRegexp = regexp.MustCompile("^[a-zA-Z0-9!#$%&'*+/=?^_`{|}~.-]+$")
var hostRegexp = regexp.MustCompile("^[^\\s]+\\.[^\\s]+$")

func init() { rand.Seed(time.Now().UnixNano()) }

// GenerateRandomString returns a randomly generated string
func GenerateRandomString(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func GenRandomCode() string {
	code := rand.Intn(999999)
	return fmt.Sprintf("%06d", code)
}

func ValidateEmail(email string) error {
	email = strings.TrimSpace(email)
	if len(email) < 6 || len(email) > 254 {
		return ErrInvalidEmail
	}

	at := strings.LastIndex(email, "@")
	if at <= 0 || at > len(email)-3 {
		return ErrInvalidEmail
	}

	user := email[:at]
	host := email[at+1:]

	if len(user) > 64 {
		return ErrInvalidEmail
	}

	if !userRegexp.MatchString(user) || !hostRegexp.MatchString(host) {
		return ErrInvalidEmail
	}

	return nil
}

func ValidatePassword(password string) error {
	numbers := `[0-9]{1}`
	lowerCaseLetters := `[a-z]{1}`
	upperCaseLetters := `[A-Z]{1}`
	symbols := `[!@#~$%^&*()+|_]{1}`
	if len(strings.TrimSpace(password)) < 8 {
		return errors.New("password is too weak. password length must be 8 or more")
	}
	if matched, err := regexp.MatchString(numbers, password); !matched || err != nil {
		return errors.New("password must contain at least a number")
	}
	if matched, err := regexp.MatchString(lowerCaseLetters, password); !matched || err != nil {
		return errors.New("password must contain at least a lowercase letter")
	}
	if matched, err := regexp.MatchString(upperCaseLetters, password); !matched || err != nil {
		return errors.New("password must contain at least an uppercase letter")
	}
	if matched, err := regexp.MatchString(symbols, password); !matched || err != nil {
		return errors.New("password must contain at least a symbol")
	}
	return nil
}

func HashPassword(password string) string {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return ""
	}
	return string(hashed)
}

func VerifyHashPassword(hashed, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashed), []byte(password))
	if err != nil {
		return false
	}
	return true
}
