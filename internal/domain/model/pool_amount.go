package model

type SwapPoolInfo struct {

	// amount0
	Amount0 float64 `json:"amount0,omitempty"`

	// amount1
	Amount1 float64 `json:"amount1,omitempty"`

	// id
	ID uint64 `json:"id,omitempty,string"`

	// liquidity
	Liquidity float64 `json:"liquidity,omitempty"`

	// price
	Price float64 `json:"price,omitempty"`
}
