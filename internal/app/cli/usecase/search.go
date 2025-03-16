package usecase

import (
	"context"
	"fmt"
	"log"
	"math"
	"os"
	"time"

	"github.com/drybin/minter-pools-ar/internal/adapter/webapi"
	"github.com/drybin/minter-pools-ar/internal/domain/helpers"
	"github.com/drybin/minter-pools-ar/internal/domain/model"
	"github.com/drybin/minter-pools-ar/pkg/wrap"

	"github.com/jedib0t/go-pretty/v6/table"
)

type ISearch interface {
	Process(ctx context.Context) error
}

type Search struct {
	ChainikApi   *webapi.ChainikWebapi
	MinterWebapi *webapi.MinterWebapi
}

func NewSearchUsecase(ChainikApi *webapi.ChainikWebapi, MinterWebapi *webapi.MinterWebapi) *Search {
	return &Search{
		ChainikApi:   ChainikApi,
		MinterWebapi: MinterWebapi,
	}
}

func (u *Search) Process(ctx context.Context) error {
	log.Println("Process Hello World!!!")
	balance, err := u.MinterWebapi.GetBalance(ctx, "Mxc3e9e6bb8ee040439e94fc3ba8296d3a679b49b5")
	if err != nil {
		return wrap.Errorf("failed to get balance: %w", err)
	}

	fmt.Printf("Баланс %.2f\n", *balance)

	list, err := u.ChainikApi.GetList(ctx)
	if err != nil {
		return wrap.Errorf("failed to get coin list: %w", err)
	}

	start := time.Now()
	allProfitPaths, err := u.getAllProfitPaths(ctx, list)
	if err != nil {
		return wrap.Errorf("failed to get all profit paths: %w", err)
	}
	allProfitPathsCount := len(*allProfitPaths)
	allProfitPathsSorted := helpers.SortArrayByCoinAndLiquidity(*allProfitPaths)
	allProfitPaths = &allProfitPathsSorted
	duration := time.Since(start)

	fmt.Printf("Найдено профитных пар %d (за %s)\n", allProfitPathsCount, duration)

	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"#", "Coin", "Profit", "Min Coin Amount", "Path"})
	for index, path := range *allProfitPaths {
		t.AppendRows([]table.Row{
			{index, path.Coin.Name, path.Profit, path.MinCoinAmount, helpers.PathAsString(path), helpers.PathAsStringWithLiquidity(path)},
		})
		t.AppendSeparator()
	}
	t.AppendFooter(table.Row{"", "", "Total", allProfitPathsCount})
	//t.Render()

	/**
	   $want = $amount+1;
	         $str = (string) $want;
	         //$swapPoolTx = new MinterSellSwapPoolTx($coinIds, $amount, $amount);
	         $swapPoolTx = new MinterSellSwapPoolTx($coinIds, $amount, $str);
	         $tx   = new MinterTx($nonce, $swapPoolTx);
	         $transaction = $tx->sign($wallet->getPrivateKey());
	         $comission = $this->minterClient->estimateTxCommission($transaction);

	         echo 'Comission: ' . $comission . PHP_EOL;
	         echo 'Transaction: ' . $transaction . PHP_EOL;
	  //       die;
	         $transactionResult = $this->minterClient->send($transaction);

	         var_dump($transactionResult);die;
	*/

	//first := *allProfitPaths
	for _, path := range *allProfitPaths {
		_ = u.MinterWebapi.Buy(ctx, path)
	}

	fmt.Println("All done")
	return nil
}

func (u *Search) getAllProfitPaths(ctx context.Context, list *model.ChainikList) (*[]model.Path, error) {
	fmt.Printf("Выбрали пулов с chainik %d\n", list.RecordsTotal)
	uniqCoins := getUniqCoins(list)
	fmt.Printf("Уникальных монет %d\n", len(uniqCoins))

	pairs := list.Data
	pairs = helpers.RemoveLowLiquidityPairFromArray(pairs)
	helpers.SortByLiquidity(pairs)

	allProfitPaths := []model.Path{}

	fee, err := u.ChainikApi.GetMinterFee(ctx)
	if err != nil {
		return nil, wrap.Errorf("failed to get minter fee: %w", err)
	}

	for _, coin := range uniqCoins {
		if coin.Name != "BIP" {
			continue
		}
		pairsWithCoin := getAllPairWithCoin(pairs, coin)

		final := []model.ChainikCoin{}

		for _, pairWithCoin := range pairsWithCoin {
			path := []model.ChainikCoin{}
			path = append(path, pairWithCoin)

			cloned := make([]model.ChainikCoin, len(pairs))
			copy(cloned, pairs)
			cloned = helpers.RemovePairFromArray(cloned, pairWithCoin)

			if pairWithCoin.Coin0Symbol == coin.Name {
				final = processPairs(cloned, path, coin, model.Coin{Name: pairWithCoin.Coin1Symbol})
			}
			if pairWithCoin.Coin1Symbol == coin.Name {
				final = processPairs(cloned, path, coin, model.Coin{Name: pairWithCoin.Coin0Symbol})
			}

			if len(final) > 1 {
				calc := 100.0
				currentCoin := coin
				coinIds := []int{}
				if final[0].Coin0Symbol == currentCoin.Name {
					coinIds = append(coinIds, final[0].Coin0Id)
				} else {
					coinIds = append(coinIds, final[0].Coin1Id)
				}

				for _, pair := range final {
					if pair.Coin0Symbol == currentCoin.Name {
						calc = calc * pair.Price
						currentCoin = model.Coin{Name: pair.Coin1Symbol}
						coinIds = append(coinIds, pair.Coin1Id)
					} else {
						calc = (1.0 / pair.Price) * calc
						currentCoin = model.Coin{Name: pair.Coin0Symbol}
						coinIds = append(coinIds, pair.Coin0Id)
					}
				}

				if calc > 100 && len(final) < 5 {
					resultWithAmount, err := u.decorateWithAmount(ctx, final)
					if err != nil {
						return nil, wrap.Errorf("failed to get pair pool info: %w", err)
					}

					minVal, _ := u.getMinCoinAmount(coin.Name, calc, *fee, pairsWithCoin)

					path := model.Path{
						Coin:          coin,
						Path:          *resultWithAmount,
						CoinIds:       coinIds,
						MinCoinAmount: *minVal,
						Profit:        calc,
					}

					allProfitPaths = append(allProfitPaths, path)
				}
			}
		}
	}

	return &allProfitPaths, nil
}

func processPairs(
	data []model.ChainikCoin,
	path []model.ChainikCoin,
	firstCoin model.Coin,
	currentCoin model.Coin,
) []model.ChainikCoin {
	pairsWithCoin := getAllPairWithCoin(data, currentCoin)

	if len(pairsWithCoin) == 0 {
		return []model.ChainikCoin{}
	}

	for _, pairWithCoin := range pairsWithCoin {
		if len(path) > 5 {
			return []model.ChainikCoin{}
		}

		if len(path) > 1 && (pairWithCoin.Coin0Symbol == firstCoin.Name || pairWithCoin.Coin1Symbol == firstCoin.Name) {
			path = append(path, pairWithCoin)
			return path
		}

		if pairWithCoin.Coin0Symbol == currentCoin.Name {
			cloned := make([]model.ChainikCoin, len(data))
			copy(cloned, data)
			cloned = helpers.RemovePairFromArray(cloned, pairWithCoin)

			path = append(path, pairWithCoin)

			return processPairs(cloned, path, firstCoin, model.Coin{Name: pairWithCoin.Coin1Symbol})
		}

		if pairWithCoin.Coin1Symbol == currentCoin.Name {
			cloned := make([]model.ChainikCoin, len(data))
			copy(cloned, data)
			cloned = helpers.RemovePairFromArray(cloned, pairWithCoin)

			path = append(path, pairWithCoin)

			return processPairs(cloned, path, firstCoin, model.Coin{Name: pairWithCoin.Coin0Symbol})
		}
	}

	return path
}

//func (u *Search) getBipUsdPrice(ctx context.Context) (*float64, error) {
//	res, err := u.MinterWebapi.GetSwapPoolInfo(
//		ctx,
//		model.ChainikCoin{Coin0Id: 1837, Coin1Id: 1},
//	)
//	if err != nil {
//		return nil, wrap.Errorf("failed to get pair pool info: %w", err)
//	}
//
//	return nil, nil
//}

func (u *Search) getMinCoinAmount(
	coinName string,
	profit float64,
	fee float64,
	pairs []model.ChainikCoin,
) (*float64, error) {
	found := 0.0

	bipWant := 20.0
	coef := profit / 100.0
	if coef < 1 {
		coef = 1.0 + coef
	}

	for value := 1; value < 90000; value++ {

		valFloat := float64(value)
		valueWithCoeff := valFloat * coef
		valueWithFee := valFloat + fee
		diff := valueWithCoeff - valueWithFee
		if math.Round(diff) > bipWant {
			// math.Round нужно для округления миллионых сотых
			found = valFloat
			break
		}
	}

	if found < 100 {
		found += 100
	}

	if coinName == "BIP" {
		return &found, nil
	}

	for _, pair := range pairs {
		if pair.Coin0Symbol == coinName && pair.Coin1Symbol == "BIP" {
			result := pair.Price * found
			return &result, nil
		}
	}

	return &found, nil
}

func (u *Search) decorateWithAmount(ctx context.Context, pairs []model.ChainikCoin) (*[]model.Pair, error) {
	res := []model.Pair{}

	for _, pair := range pairs {
		poolInfo, err := u.MinterWebapi.GetSwapPoolInfo(ctx, pair)
		if err != nil {
			return nil, wrap.Errorf("failed to get pair pool info: %w", err)
		}
		res = append(res, model.Pair{
			Pair: pair,
			Pool: *poolInfo,
		})
	}

	return &res, nil
}

func getAllPairWithCoin(items []model.ChainikCoin, coin model.Coin) []model.ChainikCoin {
	result := []model.ChainikCoin{}

	for _, item := range items {
		if item.Coin0Symbol != coin.Name && item.Coin1Symbol != coin.Name {
			continue
		}
		result = append(result, item)
	}

	return result
}

func getUniqCoins(list *model.ChainikList) []model.Coin {
	coins := []string{}
	for _, item := range list.Data {
		coins = append(coins, item.Coin0Symbol)
		coins = append(coins, item.Coin1Symbol)
	}

	coins = helpers.UniqueStrArray(coins)
	result := []model.Coin{}

	for _, coin := range coins {
		result = append(result, model.Coin{
			Name: coin,
		})
	}

	return result
}
