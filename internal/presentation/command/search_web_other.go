package command

import (
	"context"

	"github.com/drybin/minter-pools-ar/internal/app/cli/usecase"
	"github.com/urfave/cli/v2"
)

func NewSearchWebOtherCommand(service usecase.ISearchWeb) *cli.Command {
	return &cli.Command{
		Name:  "search-web-other",
		Usage: "search web other coins command",
		Flags: []cli.Flag{},
		Action: func(c *cli.Context) error {
			return service.ProcessOther(context.Background())
		},
	}
}
