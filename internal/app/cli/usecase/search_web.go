package usecase

import (
    "context"
    "fmt"
    "strconv"
    "strings"
    
    "github.com/drybin/minter-pools-ar/internal/adapter/webapi"
    "github.com/drybin/minter-pools-ar/internal/domain/model"
    "github.com/drybin/minter-pools-ar/pkg/telegram"
    "github.com/drybin/minter-pools-ar/pkg/wrap"
)

type ISearchWeb interface {
    Process(ctx context.Context) error
}

type SearchWeb struct {
    MinterWeb    *webapi.MinterWeb
    MinterWebapi *webapi.MinterWebapi
    TgWebapi     *telegram.TelegramWebapi
}

func NewSearchWebUsecase(
    minterWeb *webapi.MinterWeb,
    minterWebapi *webapi.MinterWebapi,
    tgWebapi *telegram.TelegramWebapi,
) *SearchWeb {
    return &SearchWeb{
        MinterWeb:    minterWeb,
        MinterWebapi: minterWebapi,
        TgWebapi:     tgWebapi,
    }
}

func (u *SearchWeb) Process(ctx context.Context) error {
    prices := []int{300, 500, 1000, 3000, 5000, 10000, 15000, 20000}
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
            res, err := u.MinterWebapi.BuyRaw(
                ctx,
                *r,
            )
            
            if err != nil {
                return wrap.Errorf("failed to process exchange: %w", err)
            }
            
            msg := generateTgMessage(*res, *commission)
            _, err = u.TgWebapi.Send(msg)
            if err != nil {
                return wrap.Errorf("failed to send TG message: %w", err)
            }
        }
    }
    
    fmt.Println("All done")
    return nil
}

func generateTgMessage(response model.BuyRawResponse, commission float64) string {
    newLine := "\n"
    
    return fmt.Sprintf(
        "Баланс %.2f"+newLine+
            "amountIn %d amountOut %d commission %.2f"+newLine,
        response.Balance,
        response.AmountIn,
        response.AmountOut,
        commission,
    )
}
