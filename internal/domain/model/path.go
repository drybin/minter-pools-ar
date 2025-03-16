package model

type Path struct {
	Coin          Coin
	Path          []Pair
	CoinIds       []int
	MinCoinAmount float64
	Profit        float64
}
