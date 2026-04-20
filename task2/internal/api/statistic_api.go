package api

import (
	"fmt"
	"net/http"

	"github.com/zamurabims/QA_avito/task2/internal/client"
)

type StatisticAPI struct {
	client *client.HTTPClient
}

func NewStatisticAPI(c *client.HTTPClient) *StatisticAPI {
	return &StatisticAPI{client: c}
}

func (a *StatisticAPI) GetByID(id string) (*client.Response, error) {
	return a.client.Do(http.MethodGet, fmt.Sprintf("/api/1/statistic/%s", id), nil)
}

func (a *StatisticAPI) GetByIDV2(id string) (*client.Response, error) {
	return a.client.Do(http.MethodGet, fmt.Sprintf("/api/2/statistic/%s", id), nil)
}
