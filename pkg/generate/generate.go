package generate

import "crypto/rand"

func RandString(l int) string {
	alphaNum := "0123456789ABCDEabcde"

	var bytes = make([]byte, l)
	rand.Read(bytes)
	for i, b := range bytes {
		bytes[i] = alphaNum[b%byte(len(alphaNum))]
	}
	return string(bytes)
}
