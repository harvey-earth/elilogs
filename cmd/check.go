package cmd

import (
	"context"
	"errors"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/elastic/go-elasticsearch/v8/esapi"
	"github.com/spf13/cobra"

	"github.com/harvey-earth/elilogs/utils"
)

// checkCmd represents the check command
var checkCmd = &cobra.Command{
	Use:   "check",
	Short: "Tests connection to Elasticsearch",
	Long: `This command tests the connection to Elasticsearch. This is only useful for initial configuration. You do not need to call this before other commands, as they each perform their own connection test.
	
EXIT STATUS
0 if check is sucessful,
1 if check is not successful.`,
	Run: func(cmd *cobra.Command, args []string) {
		utils.Info("check called")

		// connect to ES instance
		conn, err := utils.Connect()
		if err != nil {
			utils.Error("error connecting", err)
		}
		utils.Info("check successful")

		if c, _ := cmd.Flags().GetBool("cache"); c {
			utils.Info("getting cache")

			// Make request for indexes and error check
			indexResp, err := esapi.CatIndicesRequest{Format: "json"}.Do(context.Background(), conn)
			if err != nil {
				utils.Error("error requesting indexes:", err)
			}
			if indexResp.StatusCode != http.StatusOK {
				r, _ := io.ReadAll(indexResp.Body)
				utils.Error("error requesting indexes:", errors.New(string(r)))
			}

			resp, _ := io.ReadAll(indexResp.Body)
			utils.LogRequest(resp)
			indexData, err := utils.HandleResponse(resp)
			if err != nil {
				utils.Error("error unmarshalling response:", err)
			}

			// Get user home directory and write to .cache/elilogs.txt
			// TODO: Add error handling if .cache doesn't exist
			homePath, _ := os.UserHomeDir()
			filePath, _ := filepath.Abs(homePath + "/.cache/elilogs.txt")
			f, err := os.Create(filePath)
			defer f.Close()

			for i := 0; i < len(indexData); i++ {
				f.WriteString(indexData[i]["index"])
				f.WriteString("\n")
			}
			utils.Info("created cache")
		}
	},
}

func init() {
	rootCmd.AddCommand(checkCmd)

	checkCmd.Flags().BoolP("cache", "c", false, "save index list to cache")
}
