package executor

import (
	"github.com/Shopify/go-lua"
	"math/rand"
)

func diceLibrary() []lua.RegistryFunction {
	return []lua.RegistryFunction{
		{
			Name:     "roll",
			Function: rollDie,
		},
	}
}

func rollDie(l *lua.State) int {
	numberOfDice := lua.CheckInteger(l, 1)
	sizeOfDice := lua.CheckInteger(l, 2)

	var total int
	for range numberOfDice {
		total += rand.Intn(sizeOfDice)
	}
	return deepPush(l, total)
}
