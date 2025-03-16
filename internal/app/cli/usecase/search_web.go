package usecase

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/drybin/minter-pools-ar/internal/adapter/webapi"
	"github.com/drybin/minter-pools-ar/pkg/wrap"
)

type ISearchWeb interface {
	Process(ctx context.Context) error
}

type SearchWeb struct {
	MinterWeb *webapi.MinterWeb
}

func NewSearchWebUsecase(minterWeb *webapi.MinterWeb) *SearchWeb {
	return &SearchWeb{
		MinterWeb: minterWeb,
	}
}

func (u *SearchWeb) Process(ctx context.Context) error {
	prices := []int{300, 500, 1000, 3000, 5000, 10000}
	for _, price := range prices {
		r, err := u.MinterWeb.GetPrice(ctx, price)
		if err != nil {
			return wrap.Errorf("failed to get minter swap info: %w", err)
		}

		commission, err := u.MinterWeb.GetCommission(ctx, r, price)
		if err != nil {
			return wrap.Errorf("failed to get minter swap commission info: %w", err)
		}

		amountIn, err := strconv.ParseFloat(strings.TrimSpace(r.AmountIn), 64)
		if err != nil {
			return wrap.Errorf("failed to parse amountIn as float: %w", err)
		}

		amountOut, err := strconv.ParseFloat(strings.TrimSpace(r.AmountOut), 64)
		if err != nil {
			return wrap.Errorf("failed to parse amountIn as float: %w", err)
		}

		result := amountIn + *commission
		if result < amountOut {
			fmt.Println("SUCCESS")
			fmt.Printf("Processing price %d\n", price)
			fmt.Printf("r: %v\n", r)
			fmt.Printf("result: %f\n", result)
			fmt.Printf("com: %f\n", *commission)
		}
	}

	fmt.Println("All done")
	return nil
}
