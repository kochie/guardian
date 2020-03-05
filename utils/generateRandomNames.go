package utils

import (
	cryptoRand "crypto/rand"
	"io"
	"math/rand"
	"strings"
	"time"
)

const alphabet string = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const numbers string = "1234567890"
const specialCharacters = "!@#$%^&*"

var table = []byte{'1', '2', '3', '4', '5', '6', '7', '8', '9', '0'}

var r1 *rand.Rand

func GenerateRandomName() string {
	sb := strings.Builder{}

	for i := 0; i < 8; i++ {
		j := r1.Intn(len(alphabet))
		sb.WriteByte(alphabet[j])
	}

	sb.WriteRune('-')

	for i := 0; i < 2; i++ {
		j := r1.Intn(len(numbers))
		sb.WriteByte(numbers[j])
	}

	return sb.String()
}

func GeneratePassword() string {
	//sb := strings.Builder{}
	chars := make([]rune, 8*3)

	for i := 0; i < 8; i++ {
		j := r1.Intn(len(alphabet))
		chars[i] = rune(alphabet[j])
	}

	for i := 0; i < 8; i++ {
		j := r1.Intn(len(numbers))
		chars[i+8] = rune(numbers[j])
	}

	for i := 0; i < 8; i++ {
		j := r1.Intn(len(specialCharacters))
		chars[i+16] = rune(specialCharacters[j])
	}

	for i := range chars {
		j := rand.Intn(i + 1)
		chars[i], chars[j] = chars[j], chars[i]
	}

	return string(chars)
}

func GenerateDigitCode(max int) string {
	b := make([]byte, max)
	n, err := io.ReadAtLeast(cryptoRand.Reader, b, max)
	if n != max {
		panic(err)
	}
	for i := 0; i < len(b); i++ {
		b[i] = table[int(b[i])%len(table)]
	}
	return string(b)
}

func init() {
	s1 := rand.NewSource(time.Now().UnixNano())
	r1 = rand.New(s1)
}
