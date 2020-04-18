package pingen

import (
	"crypto/rand"
	"log"
	"math"
	"math/big"
	"strconv"
)

// GeneratePIN generates random PIN code of provided digits count.
// If count is 0 or no arguments provided 4 digit PIN will be generated.
func GeneratePIN(digitCount ...int) int {
	var lim *big.Int
	digits := 4 // default
	if len(digitCount) == 0 || digitCount[0] == 0 {
		lim = big.NewInt(int64(getLimit(digits)))
	} else {
		lim = big.NewInt(int64(getLimit(digitCount[0])))
		digits = digitCount[0]
	}
	var (
		n   *big.Int
		err error
		pin int
	)
	for {
		n, err = rand.Int(rand.Reader, lim)
		if err != nil {
			log.Fatal(err)
		}
		pin = int(n.Int64())
		if s := strconv.Itoa(pin); len(s) == digits {
			break
		}
	}
	return pin
}

// getLimit returns max number argument for rand.Int func.
func getLimit(n int) float64 {
	if n == 1 {
		return 9
	}
	return math.Pow10(n-1)*9 + getLimit(n-1)
}
