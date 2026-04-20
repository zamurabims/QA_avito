package api

import (
	"fmt"
	"net/http"

	"github.com/zamurabims/QA_avito/task2/internal/client"
	"github.com/zamurabims/QA_avito/task2/internal/models"
)

type ItemAPI struct {
	client *client.HTTPClient
}

func NewItemAPI(c *client.HTTPClient) *ItemAPI {
	return &ItemAPI{client: c}
}

func (a *ItemAPI) Create(req models.CreateItemRequest) (*client.Response, error) {
	return a.client.Do(http.MethodPost, "/api/1/item", req)
}

func (a *ItemAPI) CreateRaw(body map[string]interface{}) (*client.Response, error) {
	return a.client.Do(http.MethodPost, "/api/1/item", body)
}

func (a *ItemAPI) GetByID(id string) (*client.Response, error) {
	return a.client.Do(http.MethodGet, fmt.Sprintf("/api/1/item/%s", id), nil)
}

func (a *ItemAPI) GetBySellerID(sellerID int) (*client.Response, error) {
	return a.client.Do(http.MethodGet, fmt.Sprintf("/api/1/%d/item", sellerID), nil)
}

func (a *ItemAPI) GetBySellerIDRaw(sellerID string) (*client.Response, error) {
	return a.client.Do(http.MethodGet, fmt.Sprintf("/api/1/%s/item", sellerID), nil)
}

func (a *ItemAPI) MustCreate(req models.CreateItemRequest) (models.CreateItemResponse, error) {
	resp, err := a.Create(req)
	if err != nil {
		return models.CreateItemResponse{}, fmt.Errorf("create item: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		return models.CreateItemResponse{}, fmt.Errorf(
			"expected 200, got %d: %s", resp.StatusCode, resp.Body,
		)
	}
	var result models.CreateItemResponse
	if err := resp.Decode(&result); err != nil {
		return models.CreateItemResponse{}, fmt.Errorf("decode response: %w", err)
	}
	return result, nil
}
