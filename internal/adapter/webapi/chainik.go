package webapi

import (
	"context"
	"encoding/json"

	"github.com/drybin/minter-pools-ar/internal/domain/model"
	"github.com/drybin/minter-pools-ar/pkg/wrap"
	"github.com/go-resty/resty/v2"
)

const Chainik_list_url = "https://chainik.io/api/dt/pools/list?draw=2&columns[0][data]=function&columns[0][name]=rank&columns[0][searchable]=true&columns[0][orderable]=true&columns[0][search][value]=&columns[0][search][regex]=false&columns[1][data]=function&columns[1][name]=pool&columns[1][searchable]=true&columns[1][orderable]=false&columns[1][search][value]=&columns[1][search][regex]=false&columns[2][data]=function&columns[2][name]=price&columns[2][searchable]=true&columns[2][orderable]=true&columns[2][search][value]=&columns[2][search][regex]=false&columns[3][data]=function&columns[3][name]=change24h&columns[3][searchable]=true&columns[3][orderable]=true&columns[3][search][value]=&columns[3][search][regex]=false&columns[4][data]=function&columns[4][name]=liquidity&columns[4][searchable]=true&columns[4][orderable]=true&columns[4][search][value]=&columns[4][search][regex]=false&columns[5][data]=function&columns[5][name]=volume24h&columns[5][searchable]=true&columns[5][orderable]=true&columns[5][search][value]=&columns[5][search][regex]=false&columns[6][data]=function&columns[6][name]=transactions24h&columns[6][searchable]=true&columns[6][orderable]=true&columns[6][search][value]=&columns[6][search][regex]=false&columns[7][data]=function&columns[7][name]=fees24h&columns[7][searchable]=true&columns[7][orderable]=true&columns[7][search][value]=&columns[7][search][regex]=false&columns[8][data]=function&columns[8][name]=yield_rate&columns[8][searchable]=true&columns[8][orderable]=false&columns[8][search][value]=&columns[8][search][regex]=false&columns[9][data]=function&columns[9][name]=priceTrend&columns[9][searchable]=true&columns[9][orderable]=true&columns[9][search][value]=&columns[9][search][regex]=false&order[0][column]=0&order[0][dir]=asc&start=0&length=2000&search[value]=mode:trusted&search[regex]=false&_=1653909579425"
const Chainik_fees_url = "https://chainik.io/api/dt/stats/fees?length=100"

type ChainikWebapi struct {
	client *resty.Client
}

func NewChainikWebapi(
	client *resty.Client,
) *ChainikWebapi {
	return &ChainikWebapi{
		client: client,
	}
}

type ChainikFees struct {
	Draw            string `json:"draw"`
	RecordsTotal    int    `json:"recordsTotal"`
	RecordsFiltered int    `json:"recordsFiltered"`
	Data            []struct {
		Code  string  `json:"code"`
		Name  string  `json:"name"`
		Value float64 `json:"value"`
		Hex   string  `json:"hex"`
		Cost  string  `json:"cost"`
	} `json:"data"`
	Extra []interface{} `json:"extra"`
}

func (c *ChainikWebapi) GetList(ctx context.Context) (*model.ChainikList, error) {
	res, err := c.client.R().Get(Chainik_list_url)
	if err != nil {
		return nil, wrap.Errorf("failed to get list from Chainik: %w", err)
	}

	list := model.ChainikList{}
	err = json.Unmarshal(res.Body(), &list)
	if err != nil {
		return nil, wrap.Errorf("failed to unmarshal chainik list: %w", err)
	}

	return &list, nil
}

func (c *ChainikWebapi) GetMinterFee(ctx context.Context) (*float64, error) {
	res, err := c.client.R().Get(Chainik_fees_url)
	if err != nil {
		return nil, wrap.Errorf("failed to get fees from Chainik: %w", err)
	}

	data := ChainikFees{}
	err = json.Unmarshal(res.Body(), &data)
	if err != nil {
		return nil, wrap.Errorf("failed to unmarshal chainik fees: %w", err)
	}

	for _, info := range data.Data {
		if info.Code == "sell_bancor" {
			return &info.Value, nil
		}

	}

	return nil, wrap.Errorf("failed to failed to find 'sell_bancor': %w", err)
}
