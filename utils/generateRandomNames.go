package utils

import (
	"math/rand"
	"strings"
	"time"
)

const alphabet string = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

var r1 *rand.Rand

func GenerateRandomName() string {
	sb := strings.Builder{}

	for  i := 0; i < 8; i++ {
		j := r1.Intn(len(alphabet))
		sb.WriteByte(alphabet[j])
	} 

	return sb.String()
}

func init() {
	s1 := rand.NewSource(time.Now().UnixNano())
	r1 = rand.New(s1)
}