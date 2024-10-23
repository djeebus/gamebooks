package executor

import (
	"math/rand"
)

func rollDie(number, size int) int {
	var total int
	for range number {
		total += rand.Intn(size) + 1
	}
	return total
}
