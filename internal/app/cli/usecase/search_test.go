package usecase

import (
	"reflect"
	"testing"

	"github.com/drybin/minter-pools-ar/internal/adapter/webapi"
	"github.com/drybin/minter-pools-ar/internal/domain/model"
)

func TestSearch_getMinCoinAmount(t *testing.T) {
	type fields struct {
		ChainikApi   *webapi.ChainikWebapi
		MinterWebapi *webapi.MinterWebapi
	}
	type args struct {
		coinName string
		profit   float64
		fee      float64
		pairs    []model.ChainikCoin
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    float64
		wantErr bool
	}{
		{
			name: "BIP2 100.01",
			fields: fields{
				ChainikApi:   nil,
				MinterWebapi: nil,
			},
			args: args{
				coinName: "BIP",
				profit:   100.01,
				fee:      200.0,
				pairs:    []model.ChainikCoin{},
			},
			want:    0,
			wantErr: false,
		},
		{
			name: "BIP 100.5",
			fields: fields{
				ChainikApi:   nil,
				MinterWebapi: nil,
			},
			args: args{
				coinName: "BIP",
				profit:   100.5,
				fee:      200.0,
				pairs:    []model.ChainikCoin{},
			},
			want:    44101,
			wantErr: false,
		},
		{
			name: "BIP 101",
			fields: fields{
				ChainikApi:   nil,
				MinterWebapi: nil,
			},
			args: args{
				coinName: "BIP",
				profit:   101.0,
				fee:      200.0,
				pairs:    []model.ChainikCoin{},
			},
			want:    22050,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &Search{
				ChainikApi:   tt.fields.ChainikApi,
				MinterWebapi: tt.fields.MinterWebapi,
			}
			got, err := u.getMinCoinAmount(tt.args.coinName, tt.args.profit, tt.args.fee, tt.args.pairs)
			if (err != nil) != tt.wantErr {
				t.Errorf("getMinCoinAmount() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(*got, tt.want) {
				t.Errorf("getMinCoinAmount() got = %v, want %v", *got, tt.want)
			}
		})
	}
}
