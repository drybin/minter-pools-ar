package helpers

import (
	"fmt"
	"math/big"
	"strconv"

	"github.com/MinterTeam/minter-go-sdk/v2/transaction"
)

func BipFromApiToFloat(value string) (*float64, error) {
	n := new(big.Int)
	n, ok := n.SetString(transaction.StringToBigInt(value).String(), 10)
	if !ok {
		fmt.Println("Failed to convert string to big.Int")
		return nil, nil
	}

	// Конвертируем в float64 и делим на 1000
	floatVal, _ := strconv.ParseFloat(n.String(), 64)
	result := floatVal / 1000000000000000000

	return &result, nil
}
