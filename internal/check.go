package internal

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"

	"github.com/harvey-earth/elilogs/utils"
)

// Check attempts to connect to the elasticsearch instance
func Check(cache bool) error {
	// connect to ES instance
	conn, err := utils.Connect()
	if err != nil {
		err = fmt.Errorf("error connecting: %w", err)
		return err
	}
	utils.Info("check successful")

	if cache {
		err = makeCache(conn)
		if err != nil {
			err = fmt.Errorf("error creating cache: %w", err)
			return err
		}
	}
	return nil
}

func makeCache(conn *elasticsearch.Client) error {
	utils.Info("getting cache")

	// Make request for indexes and error check
	indexResp, err := esapi.CatIndicesRequest{Format: "json"}.Do(context.Background(), conn)
	if err != nil {
		err = fmt.Errorf("error requesting indexes: %w", err)
		return err
	}
	if indexResp.StatusCode != http.StatusOK {
		r, _ := io.ReadAll(indexResp.Body)
		err = fmt.Errorf("error requesting indexes: %w", errors.New(string(r)))
		return err
	}

	resp, _ := io.ReadAll(indexResp.Body)
	utils.LogRequest(resp)
	indexData, err := utils.HandleResponse(resp)
	if err != nil {
		err = fmt.Errorf("error handling response: %w", err)
		return err
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
	return nil
}
