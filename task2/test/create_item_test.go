package test

import (
	"strings"
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

type CreateItemSuite struct {
	suiteRun.BaseSuite
}

func (s *CreateItemSuite) TestCreate_Success(t provider.T) {
	t.Title("TC-01: Успешное создание объявления с валидными данными")
	t.Severity(allure.CRITICAL)
	t.Feature("POST /api/1/item")
	t.Tags("positive")

	req := testdata.NewItem().Build()

	var resp models.CreateItemResponse
	t.WithNewStep("POST /api/1/item → 200", func(ctx provider.StepCtx) {
		created, err := s.ItemAPI.MustCreate(req)
		require.NoError(ctx, err)
		resp = created
	})

	t.WithNewStep("Ответ содержит валидный UUID", func(ctx provider.StepCtx) {
		helpers.RequireValidUUID(t, resp.ID)
	})

	t.WithNewStep("Поля ответа соответствуют полям запроса", func(ctx provider.StepCtx) {
		require.Equal(ctx, req.Name, resp.Name)
		require.Equal(ctx, req.Price, resp.Price)
		require.Equal(ctx, req.SellerID, resp.SellerID)
	})

	t.WithNewStep("Статистика в ответе совпадает с переданной", func(ctx provider.StepCtx) {
		require.Equal(ctx, req.Statistics.Likes, resp.Statistics.Likes)
		require.Equal(ctx, req.Statistics.ViewCount, resp.Statistics.ViewCount)
		require.Equal(ctx, req.Statistics.Contacts, resp.Statistics.Contacts)
	})

	t.WithNewStep("Ответ содержит поле createdAt", func(ctx provider.StepCtx) {
		require.NotEmpty(ctx, resp.CreatedAt)
	})
}

func (s *CreateItemSuite) TestCreate_UniqueIDs(t provider.T) {
	t.Title("TC-02: Каждое созданное объявление получает уникальный ID")
	t.Severity(allure.CRITICAL)
	t.Feature("POST /api/1/item")
	t.Tags("positive", "corner-case")

	req := testdata.NewItem().Build()
	created1, err := s.ItemAPI.MustCreate(req)
	require.NoError(t, err)
	created2, err := s.ItemAPI.MustCreate(req)
	require.NoError(t, err)

	require.NotEqual(t, created1.ID, created2.ID)
}

func (s *CreateItemSuite) TestCreate_ZeroValues(t provider.T) {
	t.Title("TC-03: Создание с нулевыми значениями допустимых полей (price, statistics)")
	t.Severity(allure.NORMAL)
	t.Feature("POST /api/1/item")
	t.Tags("positive", "boundary")

	zeroStats := models.Statistics{Likes: 0, ViewCount: 0, Contacts: 0}
	req := testdata.NewItem().WithPrice(0).WithStatistics(zeroStats).Build()
	created, err := s.ItemAPI.MustCreate(req)
	require.NoError(t, err)

	require.Equal(t, 0, created.Price)
	require.Equal(t, 0, created.Statistics.Likes)
	require.Equal(t, 0, created.Statistics.ViewCount)
	require.Equal(t, 0, created.Statistics.Contacts)
}

func (s *CreateItemSuite) TestCreate_MaxLengthName(t provider.T) {
	t.Title("TC-04: Создание объявления с name длиной 255 символов")
	t.Severity(allure.MINOR)
	t.Feature("POST /api/1/item")
	t.Tags("positive", "boundary")

	req := testdata.NewItem().WithName(strings.Repeat("A", 255)).Build()
	created, err := s.ItemAPI.MustCreate(req)
	require.NoError(t, err)
	helpers.RequireValidUUID(t, created.ID)
}

func (s *CreateItemSuite) TestCreate_MultipleItemsPerSeller(t provider.T) {
	t.Title("TC-05: Несколько объявлений с одним sellerID создаются независимо")
	t.Severity(allure.NORMAL)
	t.Feature("POST /api/1/item")
	t.Tags("positive")

	sellerID := testdata.RandomSellerID()
	created1, err := s.ItemAPI.MustCreate(testdata.NewItem().WithSellerID(sellerID).Build())
	require.NoError(t, err)
	created2, err := s.ItemAPI.MustCreate(
		testdata.NewItem().WithSellerID(sellerID).WithName("Second Item").Build(),
	)
	require.NoError(t, err)

	require.NotEqual(t, created1.ID, created2.ID)
	require.Equal(t, sellerID, created1.SellerID)
	require.Equal(t, sellerID, created2.SellerID)
}

func (s *CreateItemSuite) TestCreate_MissingName(t provider.T) {
	t.Title("TC-06: Создание без обязательного поля name → 400")
	t.Severity(allure.NORMAL)
	t.Feature("POST /api/1/item")
	t.Tags("negative")

	body := map[string]interface{}{
		"sellerID":   testdata.RandomSellerID(),
		"price":      1000,
		"statistics": testdata.DefaultStatistics(),
	}
	resp, err := s.ItemAPI.CreateRaw(body)
	helpers.RequireBadRequest(t, resp, err)
}

func (s *CreateItemSuite) TestCreate_MissingPrice(t provider.T) {
	t.Title("TC-07: Создание без обязательного поля price → 400")
	t.Severity(allure.NORMAL)
	t.Feature("POST /api/1/item")
	t.Tags("negative")

	body := map[string]interface{}{
		"sellerID":   testdata.RandomSellerID(),
		"name":       "Test Item",
		"statistics": testdata.DefaultStatistics(),
	}
	resp, err := s.ItemAPI.CreateRaw(body)
	helpers.RequireBadRequest(t, resp, err)
}

func (s *CreateItemSuite) TestCreate_MissingSellerID(t provider.T) {
	t.Title("TC-08: Создание без обязательного поля sellerID → 400")
	t.Severity(allure.NORMAL)
	t.Feature("POST /api/1/item")
	t.Tags("negative")

	body := map[string]interface{}{
		"name":       "Test Item",
		"price":      1000,
		"statistics": testdata.DefaultStatistics(),
	}
	resp, err := s.ItemAPI.CreateRaw(body)
	helpers.RequireBadRequest(t, resp, err)
}

func (s *CreateItemSuite) TestCreate_NegativePrice(t provider.T) {
	t.Title("TC-09: Создание с отрицательной ценой → 400")
	t.Severity(allure.NORMAL)
	t.Feature("POST /api/1/item")
	t.Tags("negative")

	req := testdata.NewItem().WithPrice(-1).Build()
	resp, err := s.ItemAPI.Create(req)
	helpers.RequireBadRequest(t, resp, err)
}

func (s *CreateItemSuite) TestCreate_EmptyBody(t provider.T) {
	t.Title("TC-10: Создание с пустым телом запроса → 400")
	t.Severity(allure.NORMAL)
	t.Feature("POST /api/1/item")
	t.Tags("negative")

	resp, err := s.ItemAPI.CreateRaw(map[string]interface{}{})
	helpers.RequireBadRequest(t, resp, err)
}

func (s *CreateItemSuite) TestCreate_StringSellerID(t provider.T) {
	t.Title("TC-11: Создание с нечисловым sellerID (строка) → 400")
	t.Severity(allure.NORMAL)
	t.Feature("POST /api/1/item")
	t.Tags("negative")

	body := map[string]interface{}{
		"sellerID": "abc", "name": "Test", "price": 1000,
		"statistics": testdata.DefaultStatistics(),
	}
	resp, err := s.ItemAPI.CreateRaw(body)
	helpers.RequireBadRequest(t, resp, err)
}

func (s *CreateItemSuite) TestCreate_StringPrice(t provider.T) {
	t.Title("TC-12: Создание с нечисловым price (строка) → 400")
	t.Severity(allure.NORMAL)
	t.Feature("POST /api/1/item")
	t.Tags("negative")

	body := map[string]interface{}{
		"sellerID": testdata.RandomSellerID(), "name": "Test", "price": "бесплатно",
		"statistics": testdata.DefaultStatistics(),
	}
	resp, err := s.ItemAPI.CreateRaw(body)
	helpers.RequireBadRequest(t, resp, err)
}

func (s *CreateItemSuite) TestCreate_EmptyName(t provider.T) {
	t.Title("TC-13: Создание с пустым name → 400")
	t.Severity(allure.NORMAL)
	t.Feature("POST /api/1/item")
	t.Tags("negative")

	req := testdata.NewItem().WithName("").Build()
	resp, err := s.ItemAPI.Create(req)
	helpers.RequireBadRequest(t, resp, err)
}

func TestCreateItemSuite(t *testing.T) {
	suite.RunSuite(t, new(CreateItemSuite))
}

func (s *CreateItemSuite) TestCreate_NegativeStatistics(t provider.T) {
	t.Title("TC-14: Создание с отрицательными значениями statistics → 400")
	t.Severity(allure.NORMAL)
	t.Feature("POST /api/1/item")
	t.Tags("negative")

	req := testdata.NewItem().WithStatistics(models.Statistics{
		Likes:     -5,
		ViewCount: -10,
		Contacts:  -1,
	}).Build()
	resp, err := s.ItemAPI.Create(req)
	helpers.RequireBadRequest(t, resp, err)
}

func (s *CreateItemSuite) TestCreate_NegativeSellerID(t provider.T) {
	t.Title("TC-15: Создание с отрицательным sellerID → 400")
	t.Severity(allure.NORMAL)
	t.Feature("POST /api/1/item")
	t.Tags("negative")

	req := testdata.NewItem().WithSellerID(-123456).Build()
	resp, err := s.ItemAPI.Create(req)
	helpers.RequireBadRequest(t, resp, err)
}

func (s *CreateItemSuite) TestCreate_WrongFieldName(t provider.T) {
	t.Title("TC-16: Создание с недокументированным именем поля sellerId → 400")
	t.Severity(allure.NORMAL)
	t.Feature("POST /api/1/item")
	t.Tags("negative")

	body := map[string]interface{}{
		"sellerId":   testdata.RandomSellerID(),
		"name":       "Test Item",
		"price":      1000,
		"statistics": testdata.DefaultStatistics(),
	}
	resp, err := s.ItemAPI.CreateRaw(body)
	helpers.RequireBadRequest(t, resp, err)
}
