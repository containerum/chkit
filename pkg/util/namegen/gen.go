package namegen

import (
	"math/rand"
	"time"
)

var (
	rnd = rand.New(rand.NewSource(time.Now().UnixNano()))
)

func GenerateName() string {
	i := rnd.Intn(len(Colors))
	j := rnd.Intn(len(Physicists))
	return Colors[i] + "-" + Physicists[j]
}
