package model

type SwapCoin struct {
	ID     int    `json:"id"`
	Symbol string `json:"symbol"`
}

type SwapData struct {
	SwapType  string     `json:"swap_type"`
	AmountIn  string     `json:"amount_in"`
	AmountOut string     `json:"amount_out"`
	Coins     []SwapCoin `json:"coins"`
}

type SwapDetails struct {
	Commission string `json:"commission"`
	SwapFrom   string `json:"swap_from"`
	WillPay    string `json:"will_pay"`
}
