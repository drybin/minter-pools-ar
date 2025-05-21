package usecase

import (
    "context"
    "fmt"
    "math"
    "strconv"
    "strings"
    "time"
    "unicode/utf8"
    
    "github.com/drybin/minter-pools-ar/internal/adapter/webapi"
    "github.com/drybin/minter-pools-ar/internal/domain/model"
    "github.com/drybin/minter-pools-ar/pkg/telegram"
    "github.com/drybin/minter-pools-ar/pkg/wrap"
)

type ISearchWeb interface {
    Process(ctx context.Context) error
    ProcessOther(ctx context.Context) error
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
    attempts := []int{1, 2, 3}
    
    for _, attempt := range attempts {
        processSuccess := false
        
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
                
                msg := generateTgMessage(*res, *commission, attempt)
                _, err = u.TgWebapi.Send(msg)
                if err != nil {
                    return wrap.Errorf("failed to send TG message: %w", err)
                }
                
                processSuccess = true
            }
        }
        
        if !processSuccess {
            break
        }
    }
    
    fmt.Printf("All done %s\n", time.Now())
    
    return nil
}

func generateTgMessage(response model.BuyRawResponse, commission float64, attempt int) string {
    newLine := "\n"
    
    return fmt.Sprintf(
        "Баланс %.2f"+newLine+
            "amountIn %d amountOut %d commission %.2f attempt %d"+newLine,
        response.Balance,
        response.AmountIn,
        response.AmountOut,
        commission,
        attempt,
    )
}

func (u *SearchWeb) ProcessOther(ctx context.Context) error {
    prices := []int{300, 500, 1000, 3000, 5000, 10000, 15000, 20000}
    coins := []string{
        //coin without route to it self
        //"USDTBSC",
        //"DOUBLE",
        //"BTC",
        //"USDTE",
        //"HUB",
        //"TORTUGA",
        //"BNB",
        //"BTCBSC",
        //"BUSDBSC",
        //"TON",
        //"USDCE",
        //"MILE",
        //"TONBSC",
        //"MONSTERHUB",
        //"ETH",
        //"WTF",
        //"AUSDPLUS",
        //"DEXCOIN",
        //"OBSIDIAN",
        //"NEGATIVE",
        //"USDCBSC",
        //"MEMECOIN",
        //"LONG",
        //"POINTS",
        //"ETHBSC",
        //"TWTBSC",
        //"ALPACA",
        //"ADABSC",
        //"USDEQ",
        //"CAKEBSC",
        //"XRPBSC",
        //"LTCBSC",
        //"LINKBSC",
        //"DAIBSC",
        //"BLACKPINK",
        //"SOLBSC",
        //"TORNBSC",
        //"BAKEBSC",
        //"1INCH",
        //"YELLOW",
        //"ETCBSC",
        //"SFPBSC",
        //"PERCHERON",
        //"COMPBSC",
        //"STGBSC",
        //"CRYPTONAC",
        
        "HODL",
        "BEE",
        "ARCONA",
        "BVSD",
        "REDDCOIN",
        "MUSD",
        "EVOLUTION",
        "BIPXBIP",
        "UTLCLUB",
        "CASHBSC",
        "VIZCHAIN",
        "POSITIVE",
        "COFOUNDER",
        "JOOCOIN",
        "SHEKEL",
        "LUNABSC",
        "GOODZONE",
        "NFTBSC",
        "ANKRBSC",
        "MINTERINU",
        "DAIQUIRI",
    }
    
    for _, coin := range coins {
        priceInBipData, err := u.MinterWeb.GetPriceOther(ctx, coin, "BIP", 10000)
        if err != nil {
            return wrap.Errorf("failed to get minter swap info: %w", err)
        }
        
        priceInBipAmountIn, err := strconv.ParseFloat(strings.TrimSpace(priceInBipData.AmountIn), 64)
        if err != nil {
            return wrap.Errorf("failed to parse amountIn as float (when get BIP price): %w", err)
        }
        
        priceInBipAmountOut, err := strconv.ParseFloat(strings.TrimSpace(priceInBipData.AmountOut), 64)
        if err != nil {
            return wrap.Errorf("failed to parse amountIn as float (when get BIP price): %w", err)
        }
        
        priceInBip := priceInBipAmountOut / priceInBipAmountIn
        priceInBipReverse := priceInBipAmountIn / priceInBipAmountOut
        
        for _, price := range prices {
            priceToTryInBipFloat := float64(price) * priceInBipReverse
            priceToTryInBip := int(math.Round(priceToTryInBipFloat))
            
            r, err := u.MinterWeb.GetPriceOther(ctx, coin, coin, priceToTryInBip)
            if err != nil {
                return wrap.Errorf("failed to get minter swap info (when get coin price): %w", err)
            }
            
            if utf8.RuneCountInString(r.AmountIn) == 0 || utf8.RuneCountInString(r.AmountOut) == 0 {
                continue
            }
            
            amountIn, err := strconv.ParseFloat(strings.TrimSpace(r.AmountIn), 64)
            if err != nil {
                return wrap.Errorf("failed to parse amountIn as float (when get coin price): %w", err)
            }
            
            amountOut, err := strconv.ParseFloat(strings.TrimSpace(r.AmountOut), 64)
            if err != nil {
                return wrap.Errorf("failed to parse amountIn as float (when get coin price %s): %w", r.AmountOut, err)
            }
            
            if amountOut < amountIn {
                continue
            }
            
            commission, err := u.MinterWeb.GetCommissionOther(ctx, r, coin, priceToTryInBip)
            if err != nil {
                return wrap.Errorf("failed to get minter swap commission info: %w", err)
            }
            
            profit := amountOut - amountIn
            profitInBip := profit * priceInBip
            
            if profitInBip > *commission {
                //fmt.Println("SUCCESS")
                fmt.Printf("%s  profitInBip: %.2f comission: %.2f ", coin, profitInBip, *commission)
                fmt.Printf("amountIn: %.2f amountOut: %.2f\n", amountIn, amountOut)
                continue
            }
        }
    }
    
    fmt.Printf("All done %s\n", time.Now())
    return nil
}
