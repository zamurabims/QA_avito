package testdata

import (
	"math/rand"

	"github.com/zamurabims/QA_avito/task2/internal/models"
)

const (
	NonExistentID = "00000000-0000-0000-0000-000000000000"
	InvalidID     = "not-a-valid-id"

	sellerIDMin = 111111
	sellerIDMax = 999999
)

func RandomSellerID() int {
	return rand.Intn(sellerIDMax-sellerIDMin+1) + sellerIDMin
}

func DefaultStatistics() models.Statistics {
	return models.Statistics{Likes: 5, ViewCount: 100, Contacts: 10}
}

type ItemBuilder struct {
	req models.CreateItemRequest
}

func NewItem() *ItemBuilder {
	return &ItemBuilder{
		req: models.CreateItemRequest{
			SellerID:   RandomSellerID(),
			Name:       "Test Item",
			Price:      1000,
			Statistics: DefaultStatistics(),
		},
	}
}

func (b *ItemBuilder) WithSellerID(id int) *ItemBuilder {
	b.req.SellerID = id
	return b
}

func (b *ItemBuilder) WithName(name string) *ItemBuilder {
	b.req.Name = name
	return b
}

func (b *ItemBuilder) WithPrice(price int) *ItemBuilder {
	b.req.Price = price
	return b
}

func (b *ItemBuilder) WithStatistics(s models.Statistics) *ItemBuilder {
	b.req.Statistics = s
	return b
}

func (b *ItemBuilder) Build() models.CreateItemRequest {
	return b.req
}
