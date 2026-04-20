package suiteRun

import (
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"

	"github.com/zamurabims/QA_avito/task2/internal/api"
	"github.com/zamurabims/QA_avito/task2/internal/client"
	"github.com/zamurabims/QA_avito/task2/internal/config"
)

type BaseSuite struct {
	suite.Suite
	ItemAPI      *api.ItemAPI
	StatisticAPI *api.StatisticAPI
	DeleteAPI    *api.DeleteItemAPI
}

func (s *BaseSuite) BeforeAll(t provider.T) {
	cfg := config.Load()
	httpClient := client.New(cfg.BaseURL)
	s.ItemAPI = api.NewItemAPI(httpClient)
	s.StatisticAPI = api.NewStatisticAPI(httpClient)
	s.DeleteAPI = api.NewDeleteItemAPI(httpClient)
}
