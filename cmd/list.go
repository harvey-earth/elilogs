package cmd

import (
	"github.com/spf13/cobra"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list {cluster|index}",
	Short: "Lists information about an index or cluster",
	Long:  `This command lists information about an Elasticsearch index or cluster. This will retreive health, shard, and cluster information.`,
}

func init() {
	rootCmd.AddCommand(listCmd)
}
