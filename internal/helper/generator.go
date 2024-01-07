package helper

import (
	"math/rand"
	"strings"
	"time"
)

// BaseLetter source to generate the random string
type BaseLetter string

// BaseLetter constants
const (
	AlphabetCaps  BaseLetter = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	AlphabetLower BaseLetter = "abcdefghijklmnopqrstuvwxyz"
	Numeric       BaseLetter = "0123456789"
	AlphaNumeric  BaseLetter = AlphabetLower + AlphabetCaps + Numeric
)

const (
	letterIdxBits = 6
	letterIdxMask = 1<<letterIdxBits - 1
	letterIdxMax  = 63 / letterIdxBits
)

var (
	src = rand.NewSource(time.Now().UnixNano())
)

// GenerateRandomString Generate random alphanumeric character adapted from
// https://stackoverflow.com/questions/22892120/how-to-generate-a-random-string-of-a-fixed-length-in-go
func GenerateRandomString(base BaseLetter, n int) string {
	sb := strings.Builder{}
	sb.Grow(n)
	for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(base) {
			sb.WriteByte(base[idx])
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return sb.String()
}
