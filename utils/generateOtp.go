package utils

import (
	"math/rand"
	"strconv"
	"time"
)

func Generate6DigitOtp() string {
	// Create a new random number generator with a specific seed
	source := rand.NewSource(time.Now().UnixNano())
	r := rand.New(source)

	randInt := r.Intn(900000) + 100000
	return strconv.Itoa(randInt)
}
