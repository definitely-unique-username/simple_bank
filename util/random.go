package util

import (
	"fmt"
	"math/rand"
	"strings"
)

const alphabet = "abcdefghijklmnopqrstuvwxyz"

func RandInt(min, max int64) int64 {
	return min + rand.Int63n(max-min+1)
}

func RandString(n int) string {
	var sb strings.Builder

	k := len(alphabet)

	for range n {
		c := alphabet[rand.Intn(k)]
		sb.WriteByte(c)
	}

	return sb.String()
}

func RandOwner() string {
	return RandString(6)
}

func RandMoney() int64 {
	return RandInt(0, 1000)
}

func RandCurrency() string {
	currencies := []string{USD, EUR, CAD}
	n := len(currencies)

	return currencies[rand.Intn(n)]
}

func RandEmail() string {
	return fmt.Sprintf("%v@gmail.com", RandString(5))
}
