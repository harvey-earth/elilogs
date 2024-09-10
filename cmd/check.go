package cmd

import (
	"github.com/spf13/cobra"

	"github.com/harvey-earth/elilogs/internal"
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

		c, _ := cmd.Flags().GetBool("cache")
		err := internal.Check(c)
		if err != nil {
			utils.Error("Error checking connection", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(checkCmd)

	checkCmd.Flags().BoolP("cache", "c", false, "save index list to cache")
}
