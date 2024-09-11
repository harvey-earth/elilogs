package cmd

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/elastic/go-elasticsearch/v8/esapi"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/harvey-earth/elilogs/internal"
	"github.com/harvey-earth/elilogs/utils"
)

// clusterCmd represents the list cluster command
var clusterCmd = &cobra.Command{
	Use:   "cluster",
	Short: "Lists cluster information",
	Long: `This command lists information about the Elasticsearch cluster. You can view information about nodes, health, and more.

EXIT STATUS
0 if successful (and cluster health is green if health flag selected),
1 if there was an error with the request/response,
2 if the cluster health is not green and health flag selected.`,
	Run: func(cmd *cobra.Command, args []string) {
		utils.Info("cluster called")

		exitCode := 0
		// Get flag settings; if all are false then make "all" true
		allF, _ := cmd.Flags().GetBool("all")
		healthF, _ := cmd.Flags().GetBool("health")
		nodesF, _ := cmd.Flags().GetBool("nodes")
		pendingF, _ := cmd.Flags().GetBool("pending")
		snapshotF, _ := cmd.Flags().GetBool("snapshot")
		if !allF && !healthF && !nodesF && !pendingF && !snapshotF {
			allF = true
		}

		// Set variables to hold data
		var healthData []map[string]string
		var nodeData []map[string]string
		var pendingData []map[string]string
		var snapData []map[string]string

		// Connect to cluster
		conn, err := utils.Connect()
		if err != nil {
			utils.Error("error connecting", err)
		}
		utils.Info("check successful")

		// Get cluster health information if all or health
		if allF || healthF {
			// Make call for cat health, log, and error check
			healthResp, err := esapi.CatHealthRequest{Format: "json"}.Do(context.Background(), conn)
			if healthResp.StatusCode != http.StatusOK {
				r, _ := io.ReadAll(healthResp.Body)
				utils.Error("error getting health:", errors.New(string(r)))
			}
			if err != nil {
				utils.Error("error getting health:", err)
			}
			defer healthResp.Body.Close()

			// Get response body
			resp, _ := io.ReadAll(healthResp.Body)
			utils.LogRequest(resp)

			healthData, err = utils.HandleResponse(resp)
			if err != nil {
				utils.Error("error unmarshalling response:", err)
			}

			/// Check cluster status
			for i := 0; i < len(healthData); i++ {
				if healthData[i]["status"] != "green" {
					exitCode = 2
				}
			}
		}

		// Get information about nodes if all or nodes
		if allF || nodesF {
			nodeResp, err := esapi.CatNodesRequest{Format: "json"}.Do(context.Background(), conn)
			if nodeResp.StatusCode != http.StatusOK {
				r, _ := io.ReadAll(nodeResp.Body)
				utils.Error("error getting nodes:", errors.New(string(r)))
			}
			if err != nil {
				utils.Error("error getting nodes:", err)
			}
			defer nodeResp.Body.Close()

			resp, _ := io.ReadAll(nodeResp.Body)
			utils.LogRequest(resp)

			nodeData, err = utils.HandleResponse(resp)
			if err != nil {
				utils.Error("error unmarshalling response", err)
			}
		}

		// Get information about pending tasks if all or pending
		if allF || pendingF {
			pendingResp, err := esapi.CatPendingTasksRequest{Format: "json"}.Do(context.Background(), conn)
			if pendingResp.StatusCode != http.StatusOK {
				r, _ := io.ReadAll(pendingResp.Body)
				utils.Error("error getting pending tasks:", errors.New(string(r)))
			}
			if err != nil {
				utils.Error("error getting pending tasks:", err)
			}
			defer pendingResp.Body.Close()

			resp, _ := io.ReadAll(pendingResp.Body)
			utils.LogRequest(resp)

			pendingData, err = utils.HandleResponse(resp)
			if err != nil {
				utils.Error("error unmarshalling response:", err)
			}
		}

		// Get information about snapshots if all or snapshot 
		if allF || snapshotF {
			snapResp, err := esapi.CatSnapshotsRequest{Format: "json"}.Do(context.Background(), conn)
			if snapResp.StatusCode != http.StatusOK {
				r, _ := io.ReadAll(snapResp.Body)
				utils.Error("error getting snapshots:", errors.New(string(r)))
			}
			if err != nil {
				utils.Error("error getting snapshots:", err)
			}
			defer snapResp.Body.Close()

			resp, _ := io.ReadAll(snapResp.Body)
			utils.LogRequest(resp)

			snapData, err = utils.HandleResponse(resp)
			if err != nil {
				utils.Error("error unmarshalling response:", err)
			}

		}

		// Print data
		if q := viper.GetBool("quiet"); !q {
			if allF || healthF {
				// Print health report
				internal.PrintHealthInformation(healthData)

				if allF || nodesF || pendingF || snapshotF {
					fmt.Println("")
				}
			}
			if allF || nodesF {
				// Print node report
				internal.PrintNodeInformation(nodeData)

				if allF || pendingF || snapshotF {
					fmt.Println("")
				}
			}
			if allF || pendingF {
				// Print pending tasks report
				internal.PrintPendingTasks(pendingData)

				if allF || snapshotF {
					fmt.Println("")
				}
			}
			if allF || snapshotF {
				// Print snapshot report
				internal.PrintSnapshots(snapData)
			}
		}
		if exitCode != 0 {
			os.Exit(exitCode)
		}
	},
}

func init() {
	listCmd.AddCommand(clusterCmd)

	clusterCmd.Flags().BoolP("all", "a", false, "get all information")
	clusterCmd.Flags().BoolP("health", "l", false, "get health information")
	clusterCmd.Flags().BoolP("nodes", "n", false, "get information for nodes")
	clusterCmd.Flags().BoolP("pending", "p", false, "get information for pending tasks")
	clusterCmd.Flags().BoolP("snapshot", "s", false, "get information for snapshots")
}
