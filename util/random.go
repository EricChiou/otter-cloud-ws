package util

import (
	"math/rand"
	"time"
)

// Random int
func Random(times int, rang int) []int {
	rand.Seed(time.Now().UnixNano())
	if times < 1 {
		times = 1
	}
	if rang < 1 {
		rang = 100
	}
	var result []int
	for i := 0; i < times; i++ {
		result = append(result, rand.Intn(rang))
	}

	return result
}
