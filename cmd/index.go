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

// indexCmd represents the list index command
var indexCmd = &cobra.Command{
	Use:   "index",
	Short: "List index information",
	Long: `This command lists information about Elasticsearch indexes.
	
EXIT STATUS
0 if all indexes have green health status,
1 if there was an error with the request/response,
2 if any returned indexes do not have green health status.`,
	Run: func(cmd *cobra.Command, args []string) {
		exitCode := 0
		// Get index args and split into slices on ","
		indexF, _ := cmd.Flags().GetString("index")
		indexStrings := strings.Split(indexF, ",")

		utils.Info("index called")

		// Connect to cluster
		conn, err := utils.Connect()
		if err != nil {
			utils.Error("error connecting", err)
		}
		utils.Info("check successful")

		// Hold response as variable
		var indexResp *esapi.Response

		// If there are strings searching for a specific index use that
		if len(indexStrings) > 0 {
			indexResp, err = esapi.CatIndicesRequest{Index: indexStrings, Format: "json"}.Do(context.Background(), conn)
		} else {
			indexResp, err = esapi.CatIndicesRequest{Format: "json"}.Do(context.Background(), conn)
		}
		if indexResp.StatusCode != http.StatusOK {
			r, _ := io.ReadAll(indexResp.Body)
			utils.Error("error getting indexes:", errors.New(string(r)))
		}
		if err != nil {
			utils.Error("error getting indexes:", err)
		}
		defer indexResp.Body.Close()

		resp, _ := io.ReadAll(indexResp.Body)
		utils.LogRequest(resp)
		indexData, err := utils.HandleResponse(resp)
		if err != nil && len(resp) != 0 {
			utils.Error("error unmarshalling response:", err)
		}

		// Print report if not quiet
		if q := viper.GetBool("quiet"); !q {
			fmt.Printf("%-15s %-10s %-10s\n", "index", "status", "health")
			if len(resp) == 0 {
				fmt.Println("No matching indexes found")
				exitCode = 2
			}
			for i := 0; i < len(indexData); i++ {
				if indexData[i]["health"] != "green" {
					exitCode = 2
				}
				fmt.Printf("%-15.15s %-10.10s %-10.10s\n", indexData[i]["index"], indexData[i]["status"], indexData[i]["health"])
			}
		}
		if exitCode != 0 {
			os.Exit(exitCode)
		}
	},
}

func init() {
	listCmd.AddCommand(indexCmd)

	indexCmd.Flags().StringP("index", "i", "", "get information about specific index")
}
