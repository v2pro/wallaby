package infra

import (
	"math/rand"
	"strings"
	"time"
)

// Charsets
const (
	Uppercase    string = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	Lowercase           = "abcdefghijklmnopqrstuvwxyz"
	Alphabetic          = Uppercase + Lowercase
	Numeric             = "0123456789"
	Alphanumeric        = Alphabetic + Numeric
	Symbols             = "`" + `~!@#$%^&*()-_+={}[]|\;:"<>,./?`
	Hex                 = Numeric + "abcdef"
)

var (
	seed = rand.New(rand.NewSource(time.Now().UnixNano()))
)

// RandomString returns a random string
// length: str length
// charsets:
//   - uppercase，
//   - lowercase，
//   - include uppercase & lowercase，
//   - numbers，
//   - include uppercase & lowercase & numbers，
//   - symbols，
//   - hex
func RandomString(length uint8, charsets ...string) string {
	charset := strings.Join(charsets, "")
	if charset == "" {
		charset = Alphanumeric
	}
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seed.Intn(len(charset))]
	}
	return string(b)
}

// RandomPercent get random value from 0 ~ 99
func RandomPercent() int {
	return seed.Intn(100)
}
