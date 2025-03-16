package model

type ChainikList struct {
	Draw            string        `json:"draw"`
	RecordsTotal    int           `json:"recordsTotal"`
	RecordsFiltered int           `json:"recordsFiltered"`
	Data            []ChainikCoin `json:"data"`
	Extra           []interface{} `json:"extra"`
}
