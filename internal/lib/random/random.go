package random

import (
	"math/rand/v2"
	"strings"
)

func CreateRandomString(length int) string {
	symbols := "QWERTYUIOPASDFGHJKLZXCVBNMqwertyuiopasdfghjklzxcvbnm0123456789"
	n := len(symbols)
	var s strings.Builder
	for range length {
		idx := rand.IntN(n)
		s.WriteByte(symbols[idx])
	}
	return s.String()
}
