package es

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
)

type Client struct {
	es         *elasticsearch.Client
	searchSize int
}

func NewClient(address string, searchSize int) (*Client, error) {
	cfg := elasticsearch.Config{
		Addresses: []string{address},
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}

	esClient, err := elasticsearch.NewClient(cfg)
	if err != nil {
		return nil, fmt.Errorf("es.NewClient: error creating ES client: %w", err)
	}

	res, err := esClient.Info()
	if err != nil {
		return nil, fmt.Errorf("es.NewClient: error pinging ES: %w", err)
	}

	defer func() { _ = res.Body.Close() }()

	if res.IsError() {
		return nil, fmt.Errorf("es.NewClient: ES returned error status: %s", res.String())
	}

	return &Client{es: esClient, searchSize: searchSize}, nil
}

func (c *Client) CreateIndex(ctx context.Context, indexName string) error {
	mapping := `{
        "mappings": {
            "properties": {
                "id": { "type": "integer" },
                "text": { "type": "text" }
            }
        }
    }`

	req := esapi.IndicesCreateRequest{
		Index: indexName,
		Body:  strings.NewReader(mapping),
	}

	res, err := req.Do(ctx, c.es)
	if err != nil {
		return fmt.Errorf("es.CreateIndex: error sending create index request: %w", err)
	}

	defer func() { _ = res.Body.Close() }()

	if res.StatusCode == http.StatusBadRequest {
		return nil
	}

	if res.IsError() {
		return fmt.Errorf("es.CreateIndex: error creating index: %s", res.String())
	}

	return nil
}

func (c *Client) IndexDocument(ctx context.Context, indexName string, id int, text string) error {
	doc := map[string]any{
		"id":   id,
		"text": text,
	}

	docBytes, err := json.Marshal(doc)
	if err != nil {
		return fmt.Errorf("es.IndexDocument: error marshaling document: %w", err)
	}

	req := esapi.IndexRequest{
		Index:      indexName,
		DocumentID: strconv.Itoa(id),
		Body:       strings.NewReader(string(docBytes)),
		Refresh:    "true",
	}

	res, err := req.Do(ctx, c.es)
	if err != nil {
		return fmt.Errorf("es.IndexDocument: error indexing document: %w", err)
	}

	defer func() { _ = res.Body.Close() }()

	if res.IsError() {
		return fmt.Errorf("es.IndexDocument: error indexing document response: %s", res.String())
	}

	return nil
}

func (c *Client) SearchDocuments(
	ctx context.Context,
	indexName string,
	searchText string,
) ([]int, error) {
	query := map[string]any{
		"query": map[string]any{
			"match": map[string]any{
				"text": searchText,
			},
		},
		"_source": false,
		"fields":  []string{"id"},
		"size":    c.searchSize,
	}

	queryBytes, err := json.Marshal(query)
	if err != nil {
		return nil, fmt.Errorf("es.SearchDocuments: error marshaling query: %w", err)
	}

	req := esapi.SearchRequest{
		Index: []string{indexName},
		Body:  strings.NewReader(string(queryBytes)),
	}

	res, err := req.Do(ctx, c.es)
	if err != nil {
		return nil, fmt.Errorf("es.SearchDocuments: error sending search request: %w", err)
	}

	defer func() { _ = res.Body.Close() }()

	if res.IsError() {
		return nil, fmt.Errorf("es.SearchDocuments: error search response: %s", res.String())
	}

	var result map[string]any

	err = json.NewDecoder(res.Body).Decode(&result)
	if err != nil {
		return nil, fmt.Errorf("es.SearchDocuments: error parsing search response: %w", err)
	}

	var ids []int

	hits, ok := result["hits"].(map[string]any)["hits"].([]any)
	if !ok {
		return ids, nil
	}

	for _, hit := range hits {
		hitMap, ok := hit.(map[string]any)
		if !ok {
			return nil, errors.New("es.SearchDocuments: invalid hit format")
		}

		esID, ok := hitMap["_id"].(string)
		if !ok {
			return nil, errors.New("es.SearchDocuments: missing or invalid _id in hit")
		}

		id, err := strconv.Atoi(esID)
		if err != nil {
			return nil, fmt.Errorf("es.SearchDocuments: error converting id: %w", err)
		}

		ids = append(ids, id)
	}

	return ids, nil
}

func (c *Client) DeleteDocument(ctx context.Context, indexName string, id int) error {
	req := esapi.DeleteRequest{
		Index:      indexName,
		DocumentID: strconv.Itoa(id),
		Refresh:    "true",
	}

	res, err := req.Do(ctx, c.es)
	if err != nil {
		return fmt.Errorf("es.DeleteDocument: error sending delete request: %w", err)
	}

	defer func() { _ = res.Body.Close() }()

	if res.IsError() {
		return fmt.Errorf("es.DeleteDocument: error delete response: %s", res.String())
	}

	return nil
}

func (c *Client) Close(ctx context.Context) error {
	return c.es.Close(ctx)
}
