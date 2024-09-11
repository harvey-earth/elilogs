package internal

import (
	"fmt"

	"github.com/harvey-earth/elilogs/internal/models"
)

// PrintSearchResults prints the results of the search command
func PrintSearchResults(searchData models.SearchResponse) {
	if searchData.Hits.HitsCount.Total == 0 {
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

// PrintListIndexResults prints results returned by the ListIndex function
func PrintListIndexResults(indexData []map[string]string) {
	fmt.Printf("%-15s %-10s %-10s\n", "index", "status", "health")
	if len(indexData) == 0 {
		fmt.Println("No matching indexes found")
	}
	for i := 0; i < len(indexData); i++ {
		fmt.Printf("%-15.15s %-10.10s %-10.10s\n", indexData[i]["index"], indexData[i]["status"], indexData[i]["health"])
	}
}

func PrintHealthInformation(healthData []map[string]string) {
	fmt.Println("HEALTH")
	fmt.Printf("%-20s %-10s\n", "cluster name", "status")
	for i := 0; i < len(healthData); i++ {
		fmt.Printf("%-20.20s %-10s\n", healthData[i]["cluster"], healthData[i]["status"])
	}
}

func PrintNodeInformation(nodeData []map[string]string) {
	fmt.Println("NODES")
	fmt.Printf("%-20.20s %-10.10s\n", "node name", "ip")
	for i := 0; i < len(nodeData); i++ {
		fmt.Printf("%-20.20s %-10.10s", nodeData[i]["name"], nodeData[i]["ip"])
		if nodeData[i]["master"] == "*" {
			fmt.Println(" - MASTER")
		} else {
			fmt.Println("")
		}
	}
}

func PrintPendingTasks(pendingData []map[string]string) {
	fmt.Println("PENDING TASKS")
	fmt.Printf("%-20.20s %-10.10s\n", "source", "priority")
	for i := 0; i < len(pendingData); i++ {
		fmt.Printf("%-20.20s %-10.10s\n", pendingData[i]["source"], pendingData[i]["priority"])
	}
}

func PrintSnapshots(snapData []map[string]string) {
	fmt.Println("SNAPSHOTS")
	fmt.Printf("%-20.20s %-10.10s\n", "id", "status")
	for i := 0; i < len(snapData); i++ {
		fmt.Printf("%-20.20s %-10.10s\n", snapData[i]["id"], snapData[i]["status"])
	}
}
