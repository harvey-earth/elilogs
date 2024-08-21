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

	"github.com/harvey-earth/elilogs/utils"
)

// searchCmd represents the search command
var searchCmd = &cobra.Command{
	Use:   "search [-i|--index] <index> query",
	Short: "Search within indices for matching docs",
	Long: `This command searches an optional index, (or all indices if left blank) with a given query string. It returns the index and document (with document fields in no particular order) that match the query string.

EXIT STATUS
0 if search is successful and returns results
1 if search is successful but returned no results
2 if there was an error with the request/response`,

	Example: `elilogs search 'query'
elilogs search -i ["index"] 'query'`,

	Run: func(cmd *cobra.Command, args []string) {
		// Store query as var
		query := args[0]

		// Get flags
		indexF, _ := cmd.Flags().GetString("index")
		indexStrings := strings.Split(indexF, ",")

		utils.Info("search called")
		utils.Info("query: " + query)

		// Connect to cluster
		conn, err := utils.Connect()
		if err != nil {
			utils.Error("error connecting", err)
			os.Exit(2)
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
			os.Exit(2)
		}
		if err != nil {
			utils.Error("error searching:", err)
			os.Exit(2)
		}
		defer searchResp.Body.Close()

		resp, _ := io.ReadAll(searchResp.Body)
		utils.LogRequest(resp)
		searchData, err := utils.HandleSearchResponse(resp)
		if err != nil {
			utils.Error("error handling response:", err)
			os.Exit(2)
		}

		if searchData.Hits.HitsCount.Total == 0 {
			fmt.Println("No results found")
			os.Exit(1)
		} else {
			fmt.Println("index", "\t", "document")
		}
		for i := 0; i < searchData.Hits.HitsCount.Total; i++ {
			fmt.Print(searchData.Hits.HitsMap[i].Index, "\t", "{")
			for k, v := range searchData.Hits.HitsMap[i].Source {
				fmt.Print("\"", k, "\": ", v, ", ")
			}
			fmt.Print("}\n")
		}
	},
}

func init() {
	rootCmd.AddCommand(searchCmd)

	searchCmd.Flags().StringP("index", "i", "", "search within specified indexes")
}
