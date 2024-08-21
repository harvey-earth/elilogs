package models

import ()

// SearchResponse holds the full response to be broken down
type SearchResponse struct {
	Took     int                `json:"took"`
	TimedOut bool               `json:"timed_out"`
	Shards   map[string]int     `json:"_shards"`
	Hits     SearchResponseHits `json:"hits"`
}

// SearchResponseHits represents the main "hits" portion of the response
type SearchResponseHits struct {
	HitsCount SearchResponseCount          `json:"total"`
	HitsMap   []SearchResponseHitsInstance `json:"hits"`
}

// SearchResponseCount represents "hits":{"total:{"value"}"} of the response
type SearchResponseCount struct {
	Total int `json:"value"`
}

// SearchResponseHitsInstance represents "hits":{"hits"} responses
type SearchResponseHitsInstance struct {
	Index  string         `json:"_index"`
	Source map[string]any `json:"_source"`
}
