package helpers

import (
	"sort"

	"github.com/drybin/minter-pools-ar/internal/domain/model"
)

func UniqueStrArray(arr []string) []string {
	uniqueMap := make(map[string]bool)
	var uniqueSlice []string

	for _, num := range arr {
		if !uniqueMap[num] {
			uniqueSlice = append(uniqueSlice, num)
			uniqueMap[num] = true
		}
	}

	return uniqueSlice
}

func SortByLiquidity(items []model.ChainikCoin) {
	sort.Slice(items, func(i, j int) bool {
		return items[i].Liquidity > items[j].Liquidity
	})
}

func SortPathByProfit(items []model.Path) {
	sort.Slice(items, func(i, j int) bool {
		return items[i].Profit > items[j].Profit
	})
}

func RemovePairFromArray(data []model.ChainikCoin, pairToRemove model.ChainikCoin) []model.ChainikCoin {
	result := []model.ChainikCoin{}
	for _, pair := range data {
		if pair.Id == pairToRemove.Id {
			continue
		}
		result = append(result, pair)
	}

	return result
}

func RemoveLowLiquidityPairFromArray(data []model.ChainikCoin) []model.ChainikCoin {
	result := []model.ChainikCoin{}
	for _, pair := range data {
		if pair.Liquidity < 2 {
			continue
		}
		result = append(result, pair)
	}

	return result
}

func SortArrayByCoinAndLiquidity(data []model.Path) []model.Path {
	uniqCoin := []string{}
	result := []model.Path{}

	for _, item := range data {
		uniqCoin = append(uniqCoin, item.Coin.Name)
	}

	uniqCoin = UniqueStrArray(uniqCoin)
	for _, coin := range uniqCoin {
		allPathWithCoin := []model.Path{}
		for _, item := range data {
			if item.Coin.Name == coin {
				allPathWithCoin = append(allPathWithCoin, item)
			}
		}
		SortPathByProfit(allPathWithCoin)
		result = append(result, allPathWithCoin...)
	}

	return result
}
