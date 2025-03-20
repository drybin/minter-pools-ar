package helpers

import (
	"fmt"

	"github.com/drybin/minter-pools-ar/internal/domain/model"
)

func PathAsString(data model.Path) string {

	result := ""
	for _, item := range data.Path {
		result += item.Pair.Coin0Symbol
		result += "-"
		result += item.Pair.Coin1Symbol
		result += " "
	}

	return result
}

func PathAsStringWithLiquidity(data model.Path) string {
	result := ""

	prevCoinAmount := data.MinCoinAmount
	prevCoinName := data.Coin.Name

	for _, item := range data.Path {
		amountToBuy := 0.0
		poolAmount := 0.0
		if item.Pair.Coin1Symbol == prevCoinName {
			amountToBuy = prevCoinAmount / item.Pair.Price
			prevCoinAmount = amountToBuy
			prevCoinName = item.Pair.Coin0Symbol
			poolAmount = item.Pool.Amount0
		} else {
			amountToBuy = prevCoinAmount / item.Pool.Price
			prevCoinAmount = amountToBuy
			prevCoinName = item.Pair.Coin1Symbol
			poolAmount = item.Pool.Amount1
		}

		result += item.Pair.Coin0Symbol
		result += "-"
		result += item.Pair.Coin1Symbol
		result += fmt.Sprintf("(buy %.6f %s (liq %.4f))", amountToBuy, prevCoinName, poolAmount)
		if poolAmount < amountToBuy {
			result += "ALARM_LOW_LIQUIDITY"
		}
		result += "\n"
	}

	return result
}
