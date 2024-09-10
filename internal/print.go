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
