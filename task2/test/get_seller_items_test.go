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

type GetSellerItemsSuite struct {
	suiteRun.BaseSuite
	sellerID int
}

func (s *GetSellerItemsSuite) BeforeAll(t provider.T) {
	s.BaseSuite.BeforeAll(t)
	s.sellerID = testdata.RandomSellerID()

	_, err := s.ItemAPI.MustCreate(testdata.NewItem().WithSellerID(s.sellerID).Build())
	require.NoError(t, err, "BeforeAll: не удалось создать первое объявление")

	_, err = s.ItemAPI.MustCreate(
		testdata.NewItem().WithSellerID(s.sellerID).WithName("Second Item").WithPrice(500).Build(),
	)
	require.NoError(t, err, "BeforeAll: не удалось создать второе объявление")
}

func (s *GetSellerItemsSuite) TestGetBySellerID_Success(t provider.T) {
	t.Title("TC-14: Успешное получение всех объявлений продавца")
	t.Severity(allure.CRITICAL)
	t.Feature("GET /api/1/:sellerID/item")
	t.Tags("positive")

	var items []models.Item
	t.WithNewStep("GET /api/1/:sellerID/item → 200", func(ctx provider.StepCtx) {
		resp, err := s.ItemAPI.GetBySellerID(s.sellerID)
		helpers.RequireOK(t, resp, err)
		require.NoError(ctx, resp.Decode(&items))
	})

	t.WithNewStep("Список содержит минимум 2 объявления", func(ctx provider.StepCtx) {
		require.GreaterOrEqual(ctx, len(items), 2)
	})

	t.WithNewStep("Все объявления принадлежат нашему продавцу", func(ctx provider.StepCtx) {
		for _, item := range items {
			require.Equal(ctx, s.sellerID, item.SellerID)
		}
	})
}

func (s *GetSellerItemsSuite) TestGetBySellerID_ItemStructure(t provider.T) {
	t.Title("TC-15: Каждый элемент списка содержит все обязательные поля")
	t.Severity(allure.NORMAL)
	t.Feature("GET /api/1/:sellerID/item")
	t.Tags("positive")

	resp, err := s.ItemAPI.GetBySellerID(s.sellerID)
	helpers.RequireOK(t, resp, err)

	var items []models.Item
	require.NoError(t, resp.Decode(&items))

	for _, item := range items {
		helpers.RequireValidItem(t, item)
		helpers.RequireValidStatistics(t, item.Statistics)
	}
}

func (s *GetSellerItemsSuite) TestGetBySellerID_Idempotent(t provider.T) {
	t.Title("TC-16: GET /:sellerID/item идемпотентен")
	t.Severity(allure.NORMAL)
	t.Feature("GET /api/1/:sellerID/item")
	t.Tags("positive", "idempotency")

	resp1, err := s.ItemAPI.GetBySellerID(s.sellerID)
	helpers.RequireOK(t, resp1, err)
	resp2, err := s.ItemAPI.GetBySellerID(s.sellerID)
	helpers.RequireOK(t, resp2, err)

	require.Equal(t, string(resp1.Body), string(resp2.Body))
}

func (s *GetSellerItemsSuite) TestGetBySellerID_NoItems(t provider.T) {
	t.Title("TC-17: Продавец без объявлений возвращает пустой список")
	t.Severity(allure.NORMAL)
	t.Feature("GET /api/1/:sellerID/item")
	t.Tags("positive", "corner-case")

	resp, err := s.ItemAPI.GetBySellerID(testdata.RandomSellerID())
	require.NoError(t, err)

	switch resp.StatusCode {
	case http.StatusOK:
		var items []models.Item
		require.NoError(t, resp.Decode(&items))
		require.Empty(t, items)
	case http.StatusNotFound:
	default:
		t.Errorf("unexpected status: %d, body: %s", resp.StatusCode, resp.Body)
	}
}

func (s *GetSellerItemsSuite) TestGetBySellerID_StringID(t provider.T) {
	t.Title("TC-18: Строковый sellerID (не число) → 400")
	t.Severity(allure.NORMAL)
	t.Feature("GET /api/1/:sellerID/item")
	t.Tags("negative")

	resp, err := s.ItemAPI.GetBySellerIDRaw("abc")
	helpers.RequireBadRequest(t, resp, err)
}

func TestGetSellerItemsSuite(t *testing.T) {
	suite.RunSuite(t, new(GetSellerItemsSuite))
}
