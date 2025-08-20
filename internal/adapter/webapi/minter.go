package webapi

import (
    "context"
    "fmt"
    "log"
    "math"
    "math/big"
    "os"
    "strconv"
    "strings"
    
    "github.com/MinterTeam/minter-go-sdk/v2/api/http_client"
    "github.com/MinterTeam/minter-go-sdk/v2/api/http_client/models"
    "github.com/MinterTeam/minter-go-sdk/v2/transaction"
    "github.com/MinterTeam/minter-go-sdk/v2/wallet"
    "github.com/drybin/minter-pools-ar/internal/domain/helpers"
    "github.com/drybin/minter-pools-ar/internal/domain/model"
    "github.com/drybin/minter-pools-ar/pkg/wrap"
)

type MinterWebapi struct {
    client     *http_client.Client
    clientGate *http_client.Client
    passPhrase string
}

func NewMinterWebapi(
    client *http_client.Client,
    clientGate *http_client.Client,
    passPhrase string,
) *MinterWebapi {
    return &MinterWebapi{
        client:     client,
        clientGate: clientGate,
        passPhrase: passPhrase,
    }
}

func (c *MinterWebapi) GetBalance(ctx context.Context, address string) (*float64, error) {
    res, err := c.client.Address(address)
    if err != nil {
        return nil, wrap.Errorf("failed to get balance: %w", err)
    }
    
    balance, err := helpers.BipFromApiToFloat(res.BipValue)
    if err != nil {
        return nil, wrap.Errorf("failed to convert balance to float: %w", err)
    }
    
    return balance, nil
}

func (c *MinterWebapi) GetSwapPoolInfo(ctx context.Context, pair model.ChainikCoin) (
    *model.SwapPoolInfo,
    error,
) {
    res, err := c.client.SwapPool(uint64(pair.Coin0Id), uint64(pair.Coin1Id), 0)
    if err != nil {
        return nil, wrap.Errorf("failed to get swap pool info: %w", err)
    }
    
    amount0, err := helpers.BipFromApiToFloat(res.Amount0)
    if err != nil {
        return nil, wrap.Errorf("failed to convert amount0 to float: %w", err)
    }
    
    amount1, err := helpers.BipFromApiToFloat(res.Amount1)
    if err != nil {
        return nil, wrap.Errorf("failed to convert amount1 to float: %w", err)
    }
    
    liquidity, err := helpers.BipFromApiToFloat(res.Liquidity)
    if err != nil {
        return nil, wrap.Errorf("failed to convert liquidity to float: %w", err)
    }
    
    price, err := strconv.ParseFloat(res.Price, 64)
    if err != nil {
        return nil, wrap.Errorf("failed to parse price as float: %w", err)
    }
    
    return &model.SwapPoolInfo{
        Amount0:   *amount0,
        Amount1:   *amount1,
        ID:        res.ID,
        Liquidity: *liquidity,
        Price:     price,
    }, nil
}

func (c *MinterWebapi) Buy(ctx context.Context, path model.Path) error {
    w, _ := wallet.Create(c.passPhrase, "")
    nonce, _ := c.client.Nonce(w.Address)
    
    amountToSpend := transaction.BipToPip(transaction.BipToPip(big.NewInt(int64(path.MinCoinAmount + 1.0))))
    minReceive := transaction.BipToPip(transaction.BipToPip(big.NewInt(int64(path.MinCoinAmount + 20.0))))
    
    data := transaction.NewBuySwapPoolData().SetValueToBuy(minReceive).SetMaximumValueToSell(amountToSpend)
    for _, coinId := range path.CoinIds {
        data.AddCoin(uint64(coinId))
    }
    
    // Формируем BuySwapPool транзакцию через несколько пулов
    tx, err := transaction.NewBuilder(transaction.MainNetChainID).NewTransaction(data)
    
    if err != nil {
        log.Fatalf("Ошибка создания транзакции: %v", err)
    }
    
    sign, _ := tx.SetNonce(nonce).SetGasPrice(220).SetGasCoin(0).Sign(w.PrivateKey)
    encode, _ := sign.Encode()
    //fmt.Printf("encode: %v\n", encode)
    hash, _ := sign.Hash()
    fmt.Printf("hash: %v\n", hash)
    
    //res, err := c.client.WithDebug(true).SendTransaction(encode)
    res, err := c.client.SendTransaction(encode)
    if err != nil {
        _, m, _ := c.client.ErrorBody(err)
        
        if m.Error.Code != "703" {
            os.Exit(1)
        }
        fmt.Println("Error 703")
        return err
    }
    if res.Code != 0 {
        panic(res.Log)
    }
    fmt.Printf("sendData: %v\n", res)
    return nil
    
    //time.Sleep(5 * time.Second)
    //response, err := c.client.Transaction(hash)
    //if err != nil {
    //	log.Fatalf("Ошибка создания транзакции: %w", err)
    //}
    //fmt.Printf("response: %v\n", response)
    //_, _ = c.client.Marshal(response)
    //sendData := new(models.SendData)
    //_ = response.Data.UnmarshalTo(sendData)
    //_, _ = c.client.Marshal(sendData)
    //fmt.Printf("sendData: %v\n", sendData)
    //os.Exit(1)
    //return nil
}

func (c *MinterWebapi) BuyRaw(ctx context.Context, swapData model.SwapData) (*model.BuyRawResponse, error) {
    w, _ := wallet.Create(c.passPhrase, "")
    nonce, _ := c.client.Nonce(w.Address)
    
    amountInFloat, err := strconv.ParseFloat(strings.TrimSpace(swapData.AmountIn), 64)
    if err != nil {
        log.Fatalf("failed to convert amountIn to float: %v", err)
    }
    
    amountIn := int64(math.Ceil(amountInFloat))
    
    amountOutFloat, err := strconv.ParseFloat(strings.TrimSpace(swapData.AmountOut), 64)
    if err != nil {
        log.Fatalf("failed to convert amountIn to float: %v", err)
    }
    
    amountOut := int64(math.Ceil(amountOutFloat))
    
    fmt.Printf("amountIn: %v\n", amountInFloat)
    fmt.Printf("amountIn: %v\n", amountIn)
    
    fmt.Printf("amountOut: %v\n", amountOutFloat)
    fmt.Printf("amountOut: %v\n", amountOut)
    
    amountInMinter := transaction.BipToPip(big.NewInt(amountIn))
    amountOutMinter := transaction.BipToPip(big.NewInt(amountOut))
    
    data := transaction.NewBuySwapPoolData().SetValueToBuy(amountOutMinter).SetMaximumValueToSell(amountInMinter)
    for _, coin := range swapData.Coins {
        data.AddCoin(uint64(coin.ID))
    }
    
    // Формируем BuySwapPool транзакцию через несколько пулов
    tx, err := transaction.NewBuilder(transaction.MainNetChainID).NewTransaction(data)
    
    if err != nil {
        return nil, wrap.Errorf("Ошибка создания транзакции: %w", err)
    }
    
    sign, _ := tx.SetNonce(nonce).SetGasPrice(1).SetGasCoin(0).Sign(w.PrivateKey)
    encode, _ := sign.Encode()
    //fmt.Printf("encode: %v\n", encode)
    hash, _ := sign.Hash()
    fmt.Printf("hash: %v\n", hash)
    
    res, err := c.clientGate.WithDebug(true).SendTransaction(encode)
    //res, err := c.client.SendTransaction(encode)
    if err != nil {
        respCode, m, errBody := c.clientGate.ErrorBody(err)
        if respCode == 0 {
            fmt.Println("TRANSACTION WARNING")
        } else {
            
            fmt.Println("TRANSACTION ERROR")
            
            fmt.Printf("error=%v\n", err)
            fmt.Printf("errorBody=%v\n", errBody)
            
            fmt.Printf("respCode=%v\n", respCode)
            if res != nil {
                fmt.Printf("res=%v\n", res)
            }
            
            if m != nil {
                fmt.Printf("m=%v\n", m)
                fmt.Printf("errorCode=%v\n", m.Error.Code)
            }
            
            return nil, wrap.Errorf("Ошибка проведения транзакции: %w", err)
        }
    }
    
    if res != nil && res.Code != 0 {
        return nil, wrap.Errorf("Код транзакции: %d", res.Code)
    }
    
    balance, err := c.GetBalance(ctx, w.Address)
    if err != nil {
        return nil, wrap.Errorf("failed to get balance: %w", err)
    }
    
    result := model.BuyRawResponse{
        AmountIn:        amountIn,
        AmountOut:       amountOut,
        TransactionHash: hash,
        Balance:         *balance,
    }
    
    return &result, nil
}

func (c *MinterWebapi) BuyRawFloat(ctx context.Context, swapData model.SwapData) (*model.BuyRawResponse, error) {
    fmt.Println("TRANSACTION ATTEMPT")
    w, _ := wallet.Create(c.passPhrase, "")
    nonce, err := c.client.Nonce(w.Address)
    if err != nil {
        log.Fatalf("failed to get nonce: %v", err)
    }
    
    fmt.Printf("swapData: %v\n", swapData)
    amountInFloat, err := strconv.ParseFloat(strings.TrimSpace(swapData.AmountIn), 64)
    if err != nil {
        log.Fatalf("failed to convert amountIn to float: %v", err)
    }
    
    amountOutFloat, err := strconv.ParseFloat(strings.TrimSpace(swapData.AmountOut), 64)
    if err != nil {
        log.Fatalf("failed to convert amountIn to float: %v", err)
    }
    
    amountInMinter := transaction.FloatBipToPip(amountInFloat)
    amountOutMinter := transaction.FloatBipToPip(amountOutFloat)
    
    data := transaction.NewBuySwapPoolData().SetValueToBuy(amountOutMinter).SetMaximumValueToSell(amountInMinter)
    for _, coin := range swapData.Coins {
        data.AddCoin(uint64(coin.ID))
    }
    
    tx, err := transaction.NewBuilder(transaction.MainNetChainID).NewTransaction(data)
    
    if err != nil {
        return nil, wrap.Errorf("Failed to create transaction: %w", err)
    }
    
    sign, _ := tx.SetNonce(nonce).SetGasCoin(uint64(swapData.Coins[0].ID)).Sign(w.PrivateKey)
    encode, _ := sign.Encode()
    hash, _ := sign.Hash()
    
    res, err := c.clientGate.WithDebug(true).SendTransaction(encode)
    if err != nil {
        respCode, m, errBody := c.clientGate.ErrorBody(err)
        
        fmt.Println("TRANSACTION ERROR")
        fmt.Printf("error=%v\n", err)
        if m != nil {
            fmt.Printf("m=%v\n", m)
            fmt.Printf("errorBody=%v\n", errBody)
            needVal := c.tryToParseAmountError(m)
            if needVal == nil {
                return nil, wrap.Errorf("Failed to make transaction: %w", err)
            }
            
            result := model.BuyRawResponse{
                AmountInFloat: *needVal,
            }
            
            return &result, wrap.Errorf("Failed to make transaction: %w", err)
        }
        
        if respCode == 0 {
            fmt.Println("TRANSACTION WARNING")
        } else {
            fmt.Printf("error=%v\n", err)
            fmt.Printf("errorBody=%v\n", errBody)
            
            fmt.Printf("respCode=%v\n", respCode)
        }
        
        return nil, wrap.Errorf("Failed to make transaction: %w", err)
    }
    
    if res.Code != 0 {
        return nil, wrap.Errorf("Transaction code: %d", res.Code)
    }
    
    balance, err := c.GetBalance(ctx, w.Address)
    if err != nil {
        return nil, wrap.Errorf("failed to get balance: %w", err)
    }
    
    result := model.BuyRawResponse{
        AmountInFloat:   amountInFloat,
        AmountOutFloat:  amountOutFloat,
        TransactionHash: hash,
        Balance:         *balance,
    }
    
    return &result, nil
}

func (c *MinterWebapi) tryToParseAmountError(m *models.ErrorBody) *float64 {
    val, ok := m.Error.Data["needed_spend_value"]
    
    if ok {
        spendVal, err := helpers.BipFromApiToFloat(val)
        if err != nil {
            return nil
        }
        
        return spendVal
    }
    
    return nil
}
