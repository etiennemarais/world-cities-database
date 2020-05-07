package list

import (
	"context"

	"github.com/etiennemarais/world-cities-database/pkg/countries"
	"github.com/k0kubun/pp"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

func Command(ctx context.Context, logger *zap.Logger) *cobra.Command {
	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List the specified resource",
	}

	listCountries := &cobra.Command{
		Use:   "countries",
		Short: "List the world cities country list",
		RunE: func(cmd *cobra.Command, args []string) error {
			logger.Info("Listing countries")

			countries := countries.GetAll(logger)
			pp.Println(countries)

			return nil
		},
	}

	listCmd.AddCommand(listCountries)
	return listCmd
}
