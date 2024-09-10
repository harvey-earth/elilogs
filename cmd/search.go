package cmd

import (
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/harvey-earth/elilogs/internal"
	"github.com/harvey-earth/elilogs/utils"
)

// searchCmd represents the search command
var searchCmd = &cobra.Command{
	Use:   "search [-i|--index] <index> query",
	Short: "Search within indices for matching docs",
	Long: `This command searches an optional index, (or all indices if left blank) with a given query string. The query string should use Lucene query string syntax. It returns the index(es) and document(s) (with document fields in no particular order) that match the query string.

EXIT STATUS
0 if search is successful and returns results,
1 if there was an error with the request/response,
2 if search is successful but returned no results.`,

	Example: `elilogs search 'query'
elilogs search -i ["index"] 'query'`,

	Run: func(cmd *cobra.Command, args []string) {
		// Store query as var
		query := args[0]

		exitCode := 0

		// Get flags
		indexF, _ := cmd.Flags().GetString("index")
		indexStrings := strings.Split(indexF, ",")

		utils.Info("search called")
		utils.Debug("query: " + query)

		// Connect to cluster
		conn, err := utils.Connect()
		if err != nil {
			utils.Error("error connecting", err)
		}
		utils.Info("check successful")

		searchData, err := internal.Search(conn, indexStrings, query)

		if searchData.Hits.HitsCount.Total == 0 {
			exitCode = 2
		}

		// Print results unless quiet
		if q := viper.GetBool("quiet"); !q {
			internal.PrintSearchResults(searchData)
		}
		if exitCode != 0 {
			os.Exit(exitCode)
		}
	},
}

func init() {
	rootCmd.AddCommand(searchCmd)

	searchCmd.Flags().StringP("index", "i", "", "search within specified indexes")
}
