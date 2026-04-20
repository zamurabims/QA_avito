package test

import (
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

type GetStatisticsSuite struct {
	suiteRun.BaseSuite
	createdID  string
	savedStats models.Statistics
}

func (s *GetStatisticsSuite) BeforeAll(t provider.T) {
	s.BaseSuite.BeforeAll(t)
	s.savedStats = models.Statistics{Likes: 8, ViewCount: 200, Contacts: 15}

	created, err := s.ItemAPI.MustCreate(
		testdata.NewItem().WithStatistics(s.savedStats).Build(),
	)
	require.NoError(t, err, "BeforeAll: не удалось создать тестовое объявление")
	s.createdID = created.ID
}

func (s *GetStatisticsSuite) TestGetStatistics_Success(t provider.T) {
	t.Title("TC-25: Успешное получение статистики по существующему ID")
	t.Severity(allure.CRITICAL)
	t.Feature("GET /api/1/statistic/:id")
	t.Tags("positive")

	var stats []models.Statistics
	t.WithNewStep("GET /api/1/statistic/:id → 200", func(ctx provider.StepCtx) {
		resp, err := s.StatisticAPI.GetByID(s.createdID)
		helpers.RequireOK(t, resp, err)
		require.NoError(ctx, resp.Decode(&stats))
		require.NotEmpty(ctx, stats)
	})

	t.WithNewStep("Ответ содержит поля likes, viewCount, contacts", func(ctx provider.StepCtx) {
		helpers.RequireValidStatistics(t, stats[0])
	})
}

func (s *GetStatisticsSuite) TestGetStatistics_ValuesMatchCreation(t provider.T) {
	t.Title("TC-26: Значения статистики совпадают с переданными при создании")
	t.Severity(allure.NORMAL)
	t.Feature("GET /api/1/statistic/:id")
	t.Tags("positive")

	resp, err := s.StatisticAPI.GetByID(s.createdID)
	helpers.RequireOK(t, resp, err)

	var stats []models.Statistics
	require.NoError(t, resp.Decode(&stats))
	require.NotEmpty(t, stats)

	got := stats[0]
	require.Equal(t, s.savedStats.Likes, got.Likes)
	require.Equal(t, s.savedStats.ViewCount, got.ViewCount)
	require.Equal(t, s.savedStats.Contacts, got.Contacts)
}

func (s *GetStatisticsSuite) TestGetStatistics_Idempotent(t provider.T) {
	t.Title("TC-27: GET /statistic/:id идемпотентен")
	t.Severity(allure.NORMAL)
	t.Feature("GET /api/1/statistic/:id")
	t.Tags("positive", "idempotency")

	resp1, err := s.StatisticAPI.GetByID(s.createdID)
	helpers.RequireOK(t, resp1, err)
	resp2, err := s.StatisticAPI.GetByID(s.createdID)
	helpers.RequireOK(t, resp2, err)

	require.Equal(t, string(resp1.Body), string(resp2.Body))
}

func (s *GetStatisticsSuite) TestGetStatistics_ZeroValues(t provider.T) {
	t.Title("TC-28: Нулевые значения статистики возвращаются корректно")
	t.Severity(allure.MINOR)
	t.Feature("GET /api/1/statistic/:id")
	t.Tags("positive", "boundary")

	zeroStats := models.Statistics{Likes: 0, ViewCount: 0, Contacts: 0}
	created, err := s.ItemAPI.MustCreate(testdata.NewItem().WithStatistics(zeroStats).Build())
	require.NoError(t, err)

	resp, err := s.StatisticAPI.GetByID(created.ID)
	helpers.RequireOK(t, resp, err)

	var stats []models.Statistics
	require.NoError(t, resp.Decode(&stats))
	require.NotEmpty(t, stats)

	got := stats[0]
	require.Equal(t, 0, got.Likes)
	require.Equal(t, 0, got.ViewCount)
	require.Equal(t, 0, got.Contacts)
}

func (s *GetStatisticsSuite) TestGetStatistics_NotFound(t provider.T) {
	t.Title("TC-29: Статистика по несуществующему UUID → 404")
	t.Severity(allure.NORMAL)
	t.Feature("GET /api/1/statistic/:id")
	t.Tags("negative")

	resp, err := s.StatisticAPI.GetByID(testdata.NonExistentID)
	helpers.RequireNotFound(t, resp, err)
}

func (s *GetStatisticsSuite) TestGetStatistics_InvalidID(t provider.T) {
	t.Title("TC-30: Статистика по невалидному ID → 400")
	t.Severity(allure.NORMAL)
	t.Feature("GET /api/1/statistic/:id")
	t.Tags("negative")

	resp, err := s.StatisticAPI.GetByID(testdata.InvalidID)
	helpers.RequireBadRequest(t, resp, err)
}

func (s *GetStatisticsSuite) TestE2E_CreateThenGetStatistics(t provider.T) {
	t.Title("TC-31: E2E — создать объявление и убедиться что статистика совпадает")
	t.Severity(allure.CRITICAL)
	t.Feature("GET /api/1/statistic/:id")
	t.Tags("e2e")

	expected := models.Statistics{Likes: 3, ViewCount: 50, Contacts: 5}
	var itemID string

	t.WithNewStep("Создать объявление с заданной статистикой", func(ctx provider.StepCtx) {
		created, err := s.ItemAPI.MustCreate(testdata.NewItem().WithStatistics(expected).Build())
		require.NoError(ctx, err)
		itemID = created.ID
		helpers.RequireValidUUID(t, itemID)
	})

	t.WithNewStep("Получить статистику и сверить значения", func(ctx provider.StepCtx) {
		resp, err := s.StatisticAPI.GetByID(itemID)
		helpers.RequireOK(t, resp, err)

		var stats []models.Statistics
		require.NoError(ctx, resp.Decode(&stats))
		require.NotEmpty(ctx, stats)

		got := stats[0]
		require.Equal(ctx, expected.Likes, got.Likes)
		require.Equal(ctx, expected.ViewCount, got.ViewCount)
		require.Equal(ctx, expected.Contacts, got.Contacts)
	})
}

func TestGetStatisticsSuite(t *testing.T) {
	suite.RunSuite(t, new(GetStatisticsSuite))
}
