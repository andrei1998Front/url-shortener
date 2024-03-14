package random

import (
	"math/rand"
	"time"
)

func NewRandomString(size int) (string, error) {
	if size < 0 {
		return "", ErrNegativeSize
	}

	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))

	smbls := []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz1234567890")

	rndList := make([]rune, size)
	for i := range rndList {
		rndList[i] = smbls[rnd.Intn(len(smbls))]
	}

	return string(rndList), nil
}
