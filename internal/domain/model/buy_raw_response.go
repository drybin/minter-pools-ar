package model

type BuyRawResponse struct {
	AmountIn        int64
	AmountOut       int64
	AmountInFloat   float64
	AmountOutFloat  float64
	TransactionHash string
	Balance         float64
}
