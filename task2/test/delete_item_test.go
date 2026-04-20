package test

import (
	"testing"

	"github.com/ozontech/allure-go/pkg/allure"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"
	"github.com/stretchr/testify/require"

	"github.com/zamurabims/QA_avito/task2/helpers"
	"github.com/zamurabims/QA_avito/task2/internal/suiteRun"
	"github.com/zamurabims/QA_avito/task2/test/testdata"
)

type DeleteItemSuite struct {
	suiteRun.BaseSuite
}

func (s *DeleteItemSuite) TestDelete_Success(t provider.T) {
	t.Title("TC-32: Успешное удаление существующего объявления")
	t.Severity(allure.CRITICAL)
	t.Feature("DELETE /api/2/item/:id")
	t.Tags("positive")

	var itemID string
	t.WithNewStep("Создать объявление для удаления", func(ctx provider.StepCtx) {
		created, err := s.ItemAPI.MustCreate(testdata.NewItem().Build())
		require.NoError(ctx, err)
		itemID = created.ID
	})

	t.WithNewStep("DELETE /api/2/item/:id → 200", func(ctx provider.StepCtx) {
		resp, err := s.DeleteAPI.DeleteByID(itemID)
		helpers.RequireOK(t, resp, err)
	})
}

func (s *DeleteItemSuite) TestDelete_NotFound(t provider.T) {
	t.Title("TC-33: Удаление по несуществующему UUID → 404")
	t.Severity(allure.NORMAL)
	t.Feature("DELETE /api/2/item/:id")
	t.Tags("negative")

	resp, err := s.DeleteAPI.DeleteByID(testdata.NonExistentID)
	helpers.RequireNotFound(t, resp, err)
}

func (s *DeleteItemSuite) TestDelete_InvalidID(t provider.T) {
	t.Title("TC-34: Удаление по невалидному ID → 400")
	t.Severity(allure.NORMAL)
	t.Feature("DELETE /api/2/item/:id")
	t.Tags("negative")

	resp, err := s.DeleteAPI.DeleteByID(testdata.InvalidID)
	helpers.RequireBadRequest(t, resp, err)
}

func (s *DeleteItemSuite) TestDelete_E2E_CreateDeleteGet(t provider.T) {
	t.Title("TC-35: E2E — создать → удалить → убедиться что объявление недоступно")
	t.Severity(allure.CRITICAL)
	t.Feature("DELETE /api/2/item/:id")
	t.Tags("e2e")

	var itemID string

	t.WithNewStep("Создать объявление", func(ctx provider.StepCtx) {
		created, err := s.ItemAPI.MustCreate(testdata.NewItem().Build())
		require.NoError(ctx, err)
		itemID = created.ID
		helpers.RequireValidUUID(t, itemID)
	})

	t.WithNewStep("Удалить объявление → 200", func(ctx provider.StepCtx) {
		resp, err := s.DeleteAPI.DeleteByID(itemID)
		helpers.RequireOK(t, resp, err)
	})

	t.WithNewStep("GET удалённого объявления → 404", func(ctx provider.StepCtx) {
		resp, err := s.ItemAPI.GetByID(itemID)
		helpers.RequireNotFound(t, resp, err)
	})
}

func (s *DeleteItemSuite) TestDelete_Idempotent(t provider.T) {
	t.Title("TC-36: Повторное удаление одного объявления → 404")
	t.Severity(allure.NORMAL)
	t.Feature("DELETE /api/2/item/:id")
	t.Tags("idempotency")

	created, err := s.ItemAPI.MustCreate(testdata.NewItem().Build())
	require.NoError(t, err)

	resp1, err := s.DeleteAPI.DeleteByID(created.ID)
	helpers.RequireOK(t, resp1, err)

	resp2, err := s.DeleteAPI.DeleteByID(created.ID)
	require.NoError(t, err)
	helpers.RequireNotFound(t, resp2, err)
}

func TestDeleteItemSuite(t *testing.T) {
	suite.RunSuite(t, new(DeleteItemSuite))
}
