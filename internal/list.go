package internal

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"

	"github.com/harvey-earth/elilogs/utils"
)

// ListIndex takes a connection and slice of strings and runs a CatIndicesRequest to return data, an exit code to bubble up, and an error
func ListIndex(conn *elasticsearch.Client, indexes []string) (indexData []map[string]string, exitCode int, err error) {
	var indexResp *esapi.Response
	exitCode = 0

	// If there are strings searching for a specific index use that
	if len(indexes) > 0 {
		indexResp, err = esapi.CatIndicesRequest{Index: indexes, Format: "json"}.Do(context.Background(), conn)
	} else {
		indexResp, err = esapi.CatIndicesRequest{Format: "json"}.Do(context.Background(), conn)
	}
	if indexResp.StatusCode == http.StatusNotFound {
		return make([]map[string]string, 0), 2, nil
	}
	if indexResp.StatusCode != http.StatusOK {
		r, _ := io.ReadAll(indexResp.Body)
		err = fmt.Errorf("error getting indexes: %w", errors.New(string(r)))
		return nil, 1, err
	}
	if err != nil {
		err = fmt.Errorf("error getting indexes: %w", err)
		return nil, 1, err
	}
	defer indexResp.Body.Close()

	resp, _ := io.ReadAll(indexResp.Body)
	utils.LogRequest(resp)
	indexData, err = utils.HandleResponse(resp)
	if err != nil && len(resp) != 0 {
		err = fmt.Errorf("error unmarshalling response: %w", err)
		return nil, 1, err
	}

	if len(indexData) == 0 {
		exitCode = 2
	}
	for i := 0; i < len(indexData); i++ {
		if indexData[i]["health"] != "green" {
			exitCode = 2
		}
	}
	return
}
