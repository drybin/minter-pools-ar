package command

import (
	"context"

	"github.com/drybin/minter-pools-ar/internal/app/cli/usecase"
	"github.com/urfave/cli/v2"
)

func NewSearchWebCommand(service usecase.ISearchWeb) *cli.Command {
	return &cli.Command{
		Name:  "search-web",
		Usage: "search web command",
		Flags: []cli.Flag{},
		Action: func(c *cli.Context) error {
			return service.Process(context.Background())
		},
	}
}
