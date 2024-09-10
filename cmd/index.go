package cmd

import (
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/harvey-earth/elilogs/internal"
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

		indexData, exitCode, err := internal.ListIndex(conn, indexStrings)
		if err != nil {
			utils.Error("error listing index: %w", err)
		}

		// Print report if not quiet
		if q := viper.GetBool("quiet"); !q {
			internal.PrintListIndexResults(indexData)
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
