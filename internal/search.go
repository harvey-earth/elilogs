package internal

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"

	"github.com/harvey-earth/elilogs/internal/models"
	"github.com/harvey-earth/elilogs/utils"
)

func Search(conn *elasticsearch.Client, indexes []string, query string) (searchData models.SearchResponse, exitCode int, err error) {
	var searchResp *esapi.Response
	exitCode = 0

	if len(indexes) > 0 {
		searchResp, err = esapi.SearchRequest{Index: indexes, Query: string(query)}.Do(context.Background(), conn)
	} else {
		searchResp, err = esapi.SearchRequest{Query: string(query)}.Do(context.Background(), conn)
	}
	if searchResp.StatusCode != http.StatusOK {
		r, _ := io.ReadAll(searchResp.Body)
		utils.Error("error searching:", errors.New(string(r)))
		err = fmt.Errorf("error searching: %w", errors.New(string(r)))
		return searchData, 1, err
	}
	if err != nil {
		err = fmt.Errorf("error searching: %w", err)
		return searchData, 1, err
	}
	defer searchResp.Body.Close()

	resp, _ := io.ReadAll(searchResp.Body)
	utils.LogRequest(resp)
	searchData, err = utils.HandleSearchResponse(resp)
	if err != nil {
		err = fmt.Errorf("error handling response: %w", err)
		return searchData, 1, err
	}
	if searchData.Hits.HitsCount.Total == 0 {
		exitCode = 2
	}

	return searchData, exitCode, nil
}
