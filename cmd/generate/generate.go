package generate

import (
	"context"
	"fmt"
	"io/ioutil"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/etiennemarais/world-cities-database/pkg/countries"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

func Command(ctx context.Context, logger *zap.Logger) *cobra.Command {
	generateCmd := &cobra.Command{
		Use:   "generate",
		Short: "Generate a resource in the specified format",
	}

	mysqlCmd := &cobra.Command{
		Use:   "mysql",
		Short: "Generate the .sql MySQL import files for easy migration use",
		RunE: func(cmd *cobra.Command, args []string) error {
			logger.Info("Generating mysql output")
			countries := countries.GetAll(logger)

			var bc strings.Builder
			bc.WriteString(getCountryTableStructure())
			bc.WriteString("\n\n")

			var br strings.Builder
			br.WriteString(getRegionTableStructure())
			br.WriteString("\n\n")

			regionPID := 1

			for ci, country := range countries {
				countryPID := ci + 1
				fmt.Fprintf(&bc, "\nINSERT INTO `countries` (`id`, `code`, `name`) VALUES (%d, \"%s\", \"%s\");", countryPID, country.Code, country.Name)

				for _, region := range country.Regions {
					// Continue on if region name is empty
					if region.Name == "" {
						continue
					}
					fmt.Fprintf(&br, "\nINSERT INTO `regions` (`id`, `region_id`, `country_id`, `name`) VALUES (%d, \"%s\", %d, \"%s\");", regionPID, region.ID, countryPID, region.Name)
					regionPID++
				}
			}

			countriesFilename := fmt.Sprintf("dist/%s.sql", getFilename("countries"))
			err := writeFile(bc, countriesFilename, logger)
			if err != nil {
				logger.Error(fmt.Sprintf("Error writing file to %s", countriesFilename), zap.Error(err))
			}

			regionsFilename := fmt.Sprintf("dist/%s.sql", getFilename("regions"))
			err = writeFile(br, regionsFilename, logger)
			if err != nil {
				logger.Error(fmt.Sprintf("Error writing file to %s", regionsFilename), zap.Error(err))
			}

			logger.Info("Generated mysql output: Complete!")
			return nil
		},
	}

	generateCmd.AddCommand(mysqlCmd)
	return generateCmd
}

func getFilename(name string) string {
	src := rand.NewSource(time.Now().UnixNano())
	rnd := rand.New(src)
	return name + "_" + strconv.FormatInt(rnd.Int63(), 10)
}

func writeFile(b strings.Builder, filename string, logger *zap.Logger) error {
	content := []byte(b.String())
	return ioutil.WriteFile(filename, content, 0644)
}

func getCountryTableStructure() string {
	return `CREATE TABLE countries (
id int(11) unsigned NOT NULL AUTO_INCREMENT,
code varchar(2) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '',
name varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '',
PRIMARY KEY (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;`
}

func getRegionTableStructure() string {
	return `CREATE TABLE regions (
id int(11) unsigned NOT NULL AUTO_INCREMENT,
region_id varchar(20) NOT NULL,
country_id int(11) NOT NULL,
name varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '',
PRIMARY KEY (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;`
}
