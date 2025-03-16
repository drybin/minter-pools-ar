package model

type ChainikCoin struct {
	Id              int         `json:"id"`
	Rank            int         `json:"rank"`
	PoolSymbol      string      `json:"pool_symbol"`
	PoolIcon        interface{} `json:"pool_icon"`
	Coin0Symbol     string      `json:"coin0_symbol"`
	Coin0Icon       interface{} `json:"coin0_icon"`
	Coin0Id         int         `json:"coin0_id"`
	Coin1Symbol     string      `json:"coin1_symbol"`
	Coin1Icon       string      `json:"coin1_icon"`
	Coin1Id         int         `json:"coin1_id"`
	Price           float64     `json:"price"`
	PriceUsd        float64     `json:"price_usd"`
	PriceChange24H  float64     `json:"priceChange24h"`
	Volume24HUsd    float64     `json:"volume_24h_usd"`
	Transactions24H int         `json:"transactions24h"`
	Fee24HUsd       float64     `json:"fee_24h_usd"`
	Liquidity       float64     `json:"liquidity"`
	YieldRate       float64     `json:"yield_rate"`
	Trends          []float64   `json:"trends"`
}
