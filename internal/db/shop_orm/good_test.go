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
	rand.Seed(time.Now().UnixNano())
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

func TestCreateGood(t *testing.T) {
	good := CreateRandomGood()
	require.NoError(t, nil, good.CreateGood(conn))
	// require.Equal(t, 10, num)
}
