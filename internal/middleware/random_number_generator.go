package middleware

import (
	"math/rand"
	"time"
)

func GenerateRandomNumber(min int, max int) int {
	rand.Seed(time.Now().Unix())
	randomNum := min + rand.Intn(max-min+1)
	return randomNum
}
