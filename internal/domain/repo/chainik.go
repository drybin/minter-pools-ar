package repo

import (
	"context"
)

type IChainikRepository interface {
	GetList(ctx context.Context) error
}
