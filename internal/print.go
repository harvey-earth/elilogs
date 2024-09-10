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
