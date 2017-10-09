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
	seededRand *rand.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))
)

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
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

func RandomPercent() int {
	return seededRand.Intn(100)
}