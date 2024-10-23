package executor

import (
	"github.com/Shopify/go-lua"
)

func diceLibrary() []lua.RegistryFunction {
	return []lua.RegistryFunction{
		{
			Name: "roll",
			Function: func(l *lua.State) int {
				numberOfDice := lua.CheckInteger(l, 1)
				sizeOfDice := lua.CheckInteger(l, 2)

				total := rollDie(numberOfDice, sizeOfDice)
				return deepPush(l, total)
			},
		},
	}
}
