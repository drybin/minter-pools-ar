package webapi

import (
    "context"
    "encoding/json"
    "fmt"
    "math/big"
    
    "github.com/MinterTeam/minter-go-sdk/v2/transaction"
    "github.com/drybin/minter-pools-ar/internal/domain/helpers"
    "github.com/drybin/minter-pools-ar/internal/domain/model"
    "github.com/drybin/minter-pools-ar/pkg/wrap"
    "github.com/go-resty/resty/v2"
)

const calc_url = "https://explorer-api.minter.network/api/v2/pools/coins/BIP/BIP/estimate?type=output&amount=%s"
const comission_url = "https://gate-api.minter.network/api/v2/estimate_coin_buy?coin_to_buy=BIP&value_to_buy=%s&coin_to_sell=BIP&swap_from=poold&coin_commission=BIP"

const calc_url_other = "https://explorer-api.minter.network/api/v2/pools/coins/%s/%s/estimate?type=output&amount=%s"
const comission_url_other = "https://gate-api.minter.network/api/v2/estimate_coin_buy?coin_to_buy=%s&value_to_buy=%s&coin_to_sell=%s&swap_from=poold&coin_commission=BIP"

type MinterWeb struct {
    client *resty.Client
}

func NewMinterWeb(
    client *resty.Client,
) *MinterWeb {
    return &MinterWeb{
        client: client,
    }
}

func (c *MinterWeb) GetPrice(ctx context.Context, value int) (*model.SwapData, error) {
    pip := transaction.BipToPip(big.NewInt(int64(value)))
    
    res, err := c.client.R().Get(
        fmt.Sprintf(calc_url, pip),
    )
    if err != nil {
        return nil, wrap.Errorf("failed to get swap info from minter api: %w", err)
    }
    
    result := model.SwapData{}
    err = json.Unmarshal(res.Body(), &result)
    if err != nil {
        return nil, wrap.Errorf("failed to unmarshal swap info: %w", err)
    }
    
    return &result, nil
}

func (c *MinterWeb) GetPriceOther(ctx context.Context, coin1 string, coin2 string, value int) (*model.SwapData, error) {
    pip := transaction.BipToPip(big.NewInt(int64(value)))
    
    res, err := c.client.R().Get(
        fmt.Sprintf(calc_url_other, coin1, coin2, pip),
    )
    if err != nil {
        return nil, wrap.Errorf("failed to get swap info from minter api: %w", err)
    }
    
    result := model.SwapData{}
    err = json.Unmarshal(res.Body(), &result)
    if err != nil {
        return nil, wrap.Errorf("failed to unmarshal swap info: %w", err)
    }
    
    return &result, nil
}

func (c *MinterWeb) GetCommission(ctx context.Context, swapData *model.SwapData, value int) (*float64, error) {
    pip := transaction.BipToPip(big.NewInt(int64(value)))
    
    url := fmt.Sprintf(comission_url, pip)
    
    for _, coin := range swapData.Coins {
        if coin.Symbol == "BIP" {
            continue
        }
        url += fmt.Sprintf("&route=%d", coin.ID)
    }
    res, err := c.client.R().Get(url)
    if err != nil {
        return nil, wrap.Errorf("failed to get swap commission info from minter api: %w", err)
    }
    
    result := model.SwapDetails{}
    err = json.Unmarshal(res.Body(), &result)
    if err != nil {
        return nil, wrap.Errorf("failed to unmarshal swap commission info: %w", err)
    }
    
    if len(result.Commission) == 0 {
        res := 10000.0
        return &res, nil
    }
    
    commission, _ := helpers.BipFromApiToFloat(result.Commission)
    
    return commission, nil
}

func (c *MinterWeb) GetCommissionOther(ctx context.Context, swapData *model.SwapData, coin string, value int) (*float64, error) {
    pip := transaction.BipToPip(big.NewInt(int64(value)))
    
    //url := fmt.Sprintf(comission_url, pip)
    url := fmt.Sprintf(comission_url_other, coin, pip.String(), coin)
    
    for _, swapCoin := range swapData.Coins {
        if swapCoin.Symbol == coin {
            continue
        }
        url += fmt.Sprintf("&route=%d", swapCoin.ID)
    }
    
    res, err := c.client.R().Get(url)
    if err != nil {
        return nil, wrap.Errorf("failed to get swap commission info from minter api: %w", err)
    }
    
    result := model.SwapDetails{}
    err = json.Unmarshal(res.Body(), &result)
    if err != nil {
        return nil, wrap.Errorf("failed to unmarshal swap commission info: %w", err)
    }
    
    if len(result.Commission) == 0 {
        res := 10000.0
        return &res, nil
    }
    
    commission, _ := helpers.BipFromApiToFloat(result.Commission)
    
    return commission, nil
}
