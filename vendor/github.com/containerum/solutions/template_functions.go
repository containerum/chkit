package solutions

import (
	"crypto/sha256"
	"encoding/hex"
	"math/rand"
	"text/template"
	"time"
)

func init() {
	rand.Seed(time.Now().Unix())
}

const letterBytes = "abcdefghijklmnopqrstuvwxyz0123456789"

var templateFunctions = template.FuncMap{
	"rand_string": randString,
}

func randString(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

func NamespaceSelector(namespace string) string {
	nsHash := sha256.Sum256([]byte(namespace))
	return hex.EncodeToString(nsHash[:])[:32]
}
