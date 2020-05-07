package generate

func Command(ctx context.Context, logger *zap.Logger) *cobra.Command {
	generateCmd := &cobra.Command{
		Use:   "generate",
		Short: "Generate a resource in the specified format",
	}

	mysqlCmd := &cobra.Command{
		Use:   "mysql",
		Short: "Generate the .sql MySQL import files for easy migration use",
		RunE: func(cmd *cobra.Command, args []string) error {
			logger.Info("Listing countries")
			countries := countries.GetAll(logger)
			

			return nil
		},
	}

	generateCmd.AddCommand(mysqlCmd)
	return generateCmd
}