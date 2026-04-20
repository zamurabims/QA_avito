package api

import (
	"fmt"
	"net/http"

	"github.com/zamurabims/QA_avito/task2/internal/client"
)

type DeleteItemAPI struct {
	client *client.HTTPClient
}

func NewDeleteItemAPI(c *client.HTTPClient) *DeleteItemAPI {
	return &DeleteItemAPI{client: c}
}

func (a *DeleteItemAPI) DeleteByID(id string) (*client.Response, error) {
	return a.client.Do(http.MethodDelete, fmt.Sprintf("/api/2/item/%s", id), nil)
}
