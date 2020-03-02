package utils

import (
	cryptoRand "crypto/rand"
	"io"
	"math/rand"
	"strings"
	"time"
)

const alphabet string = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
var table = []byte{'1', '2', '3', '4', '5', '6', '7', '8', '9', '0'}

var r1 *rand.Rand

func GenerateRandomName() string {
	sb := strings.Builder{}

	for  i := 0; i < 8; i++ {
		j := r1.Intn(len(alphabet))
		sb.WriteByte(alphabet[j])
	} 

	return sb.String()
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