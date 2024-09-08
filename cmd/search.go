package cmd

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/elastic/go-elasticsearch/v8/esapi"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

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
		utils.Info("query: " + query)

		// Connect to cluster
		conn, err := utils.Connect()
		if err != nil {
			utils.Error("error connecting", err)
		}
		utils.Info("check successful")

		// Hold response as variable
		var searchResp *esapi.Response

		// Run query with esapi while checking for specific indexes
		if len(indexStrings) > 0 {
			searchResp, err = esapi.SearchRequest{Index: indexStrings, Query: string(query)}.Do(context.Background(), conn)
		} else {
			searchResp, err = esapi.SearchRequest{Query: "test"}.Do(context.Background(), conn)
		}
		if searchResp.StatusCode != http.StatusOK {
			r, _ := io.ReadAll(searchResp.Body)
			utils.Error("error searching:", errors.New(string(r)))
		}
		if err != nil {
			utils.Error("error searching:", err)
		}
		defer searchResp.Body.Close()

		resp, _ := io.ReadAll(searchResp.Body)
		utils.LogRequest(resp)
		searchData, err := utils.HandleSearchResponse(resp)
		if err != nil {
			utils.Error("error handling response:", err)
		}

		if searchData.Hits.HitsCount.Total == 0 {
			exitCode = 2
		}

		// Print results unless quiet
		if q := viper.GetBool("quiet"); !q {
			if exitCode == 2 {
				fmt.Println("No results found")
			} else {
				fmt.Println("index", "\t", "document")

				for i := 0; i < searchData.Hits.HitsCount.Total; i++ {
					fmt.Print(searchData.Hits.HitsMap[i].Index, "\t", "{")
					for k, v := range searchData.Hits.HitsMap[i].Source {
						fmt.Print("\"", k, "\": ", v, ", ")
					}
					fmt.Print("}\n")
				}
			}
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
