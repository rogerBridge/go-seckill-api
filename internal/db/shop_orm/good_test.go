package shop_orm_test

import (
	"go-seckill/internal/db/shop_orm"
	"math/rand"
	"testing"
	"time"

	"github.com/segmentio/ksuid"
	"github.com/stretchr/testify/require"
)

func CreateRandomGood() *shop_orm.Good {
	rand.Seed(time.Now().UTC().Unix())
	good := &shop_orm.Good{
		ProductCategory: "test purpose",
		ProductName:     ksuid.New().String(),
		Inventory:       randInt(100, 200),
		Price:           randInt(3000, 6000),
	}
	return good
}

// generate random int number, range: [min, max]
func randInt(min int, max int) int {
	return min + rand.Intn(max-min+1)
}

func TestGood(t *testing.T) {
	// create good
	good := CreateRandomGood()
	require.NoError(t, good.CreateGood(conn))
	// query good
	require.Equal(t, good, good.QueryGoodsByProductCategoryAndProductName(conn)[0])
	// update good
	good.Inventory = randInt(100, 200)
	good.Price = randInt(3000, 6000)
	require.NoError(t, good.UpdateGood(conn))
	require.Equal(t, good.Inventory, good.QueryGoodsByProductCategoryAndProductName(conn)[0].Inventory)
	require.Equal(t, good.Price, good.QueryGoodsByProductCategoryAndProductName(conn)[0].Price)
	// delete good
	require.NoError(t, good.DeleteGood(conn))
	require.Equal(t, len(good.QueryGoodsByProductCategoryAndProductName(conn)), 0)
}

