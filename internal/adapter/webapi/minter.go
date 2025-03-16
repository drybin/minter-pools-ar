package webapi

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"os"
	"strconv"

	"github.com/MinterTeam/minter-go-sdk/v2/api/http_client"
	"github.com/MinterTeam/minter-go-sdk/v2/transaction"
	"github.com/MinterTeam/minter-go-sdk/v2/wallet"
	"github.com/drybin/minter-pools-ar/internal/domain/helpers"
	"github.com/drybin/minter-pools-ar/internal/domain/model"
	"github.com/drybin/minter-pools-ar/pkg/wrap"
)

type MinterWebapi struct {
	client     *http_client.Client
	passPhrase string
}

func NewMinterWebapi(
	client *http_client.Client,
	passPhrase string,
) *MinterWebapi {
	return &MinterWebapi{
		client:     client,
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
