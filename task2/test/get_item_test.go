package test

import (
	"net/http"
	"testing"

	"github.com/ozontech/allure-go/pkg/allure"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"
	"github.com/stretchr/testify/require"

	"github.com/zamurabims/QA_avito/task2/helpers"
	"github.com/zamurabims/QA_avito/task2/internal/models"
	"github.com/zamurabims/QA_avito/task2/internal/suiteRun"
	"github.com/zamurabims/QA_avito/task2/test/testdata"
)

type GetItemSuite struct {
	suiteRun.BaseSuite
	created models.CreateItemResponse
}

func (s *GetItemSuite) BeforeAll(t provider.T) {
	s.BaseSuite.BeforeAll(t)

	var err error
	s.created, err = s.ItemAPI.MustCreate(testdata.NewItem().Build())
	require.NoError(t, err, "BeforeAll: не удалось создать тестовое объявление")
}

func (s *GetItemSuite) TestGetByID_Success(t provider.T) {
	t.Title("TC-19: Успешное получение существующего объявления")
	t.Severity(allure.CRITICAL)
	t.Feature("GET /api/1/item/:id")
	t.Tags("positive", "e2e")

	var items []models.Item
	t.WithNewStep("GET /api/1/item/:id → 200", func(ctx provider.StepCtx) {
		resp, err := s.ItemAPI.GetByID(s.created.ID)
		helpers.RequireOK(t, resp, err)
		require.NoError(ctx, resp.Decode(&items))
		require.NotEmpty(ctx, items)
	})

	t.WithNewStep("Данные совпадают с созданным объявлением", func(ctx provider.StepCtx) {
		item := items[0]
		helpers.RequireValidItem(t, item)
		require.Equal(ctx, s.created.ID, item.ID)
		require.Equal(ctx, s.created.Name, item.Name)
		require.Equal(ctx, s.created.Price, item.Price)
		require.Equal(ctx, s.created.SellerID, item.SellerID)
	})

	t.WithNewStep("Статистика содержит все поля", func(ctx provider.StepCtx) {
		helpers.RequireValidStatistics(t, items[0].Statistics)
	})
}

func (s *GetItemSuite) TestGetByID_Idempotent(t provider.T) {
	t.Title("TC-20: GET /item/:id идемпотентен")
	t.Severity(allure.NORMAL)
	t.Feature("GET /api/1/item/:id")
	t.Tags("positive", "idempotency")

	resp1, err := s.ItemAPI.GetByID(s.created.ID)
	helpers.RequireOK(t, resp1, err)
	resp2, err := s.ItemAPI.GetByID(s.created.ID)
	helpers.RequireOK(t, resp2, err)

	require.Equal(t, string(resp1.Body), string(resp2.Body))
}

func (s *GetItemSuite) TestGetByID_ContentType(t provider.T) {
	t.Title("TC-21: Content-Type ответа — application/json")
	t.Severity(allure.MINOR)
	t.Feature("GET /api/1/item/:id")
	t.Tags("positive", "non-functional")

	resp, err := s.ItemAPI.GetByID(s.created.ID)
	helpers.RequireOK(t, resp, err)
	helpers.RequireContentTypeJSON(t, resp)
}

func (s *GetItemSuite) TestGetByID_NotFound(t provider.T) {
	t.Title("TC-22: Получение по несуществующему UUID → 404")
	t.Severity(allure.NORMAL)
	t.Feature("GET /api/1/item/:id")
	t.Tags("negative")

	resp, err := s.ItemAPI.GetByID(testdata.NonExistentID)
	helpers.RequireNotFound(t, resp, err)
}

func (s *GetItemSuite) TestGetByID_InvalidID(t provider.T) {
	t.Title("TC-23: Получение по невалидному ID (не UUID) → 400")
	t.Severity(allure.NORMAL)
	t.Feature("GET /api/1/item/:id")
	t.Tags("negative")

	resp, err := s.ItemAPI.GetByID(testdata.InvalidID)
	helpers.RequireBadRequest(t, resp, err)
}

func (s *GetItemSuite) TestGetByID_EmptyID(t provider.T) {
	t.Title("TC-24: Пустой ID в пути — сервер не должен вернуть 200")
	t.Severity(allure.MINOR)
	t.Feature("GET /api/1/item/:id")
	t.Tags("negative")

	resp, err := s.ItemAPI.GetByID("")
	require.NoError(t, err)
	require.Contains(t,
		[]int{http.StatusBadRequest, http.StatusNotFound, http.StatusMethodNotAllowed},
		resp.StatusCode,
	)
}

func TestGetItemSuite(t *testing.T) {
	suite.RunSuite(t, new(GetItemSuite))
}
