package utils

import (
	"encoding/json"

	"github.com/harvey-earth/elilogs/internal/models"
)

// HandleResponse handles esapi responses for everything except search
func HandleResponse(resp []byte) ([]map[string]string, error) {
	var data []map[string]string

	err := json.Unmarshal(resp, &data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// HandleSearchResponse handles responses for searches
func HandleSearchResponse(resp []byte) (models.SearchResponse, error) {
	var data models.SearchResponse

	// err := json.NewDecoder(resp).Decode(&data)
	err := json.Unmarshal(resp, &data)
	if err != nil {
		return data, err
	}
	return data, nil
}
