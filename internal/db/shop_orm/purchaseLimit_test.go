package shop_orm_test

import (
	"go-seckill/internal/db/shop_orm"
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestPurchaseLimit(t *testing.T) {
	p := randomPurchaseLimit()
	// create p
	require.NoError(t, p.CreatePurchaseLimit(conn))
	require.Equal(t, p, p.QueryPurchaseLimitByProductID())
	// query p
	require.Equal(t, p, p.QueryPurchaseLimitByProductID())
	// update p
	p.LimitNum = randInt(5, 10)
	p.StartPurchaseTimeStamp = int(time.Now().Unix()) + randInt(6000, 12000)
	p.StopPurchaseTimeStamp = int(time.Now().Unix()) + randInt(6000, 12000)
	require.NoError(t, p.UpdatePurchaseLimit(conn))
	require.Equal(t, p.LimitNum, p.QueryPurchaseLimitByProductID().LimitNum)
	require.Equal(t, p.StartPurchaseTimeStamp, p.QueryPurchaseLimitByProductID().StartPurchaseTimeStamp)
	require.Equal(t, p.StopPurchaseTimeStamp, p.QueryPurchaseLimitByProductID().StopPurchaseTimeStamp)
	// delete p
	require.NoError(t, p.DeletePurchaseLimit(conn))
	require.NotEqual(t, p, p.QueryPurchaseLimitByProductID())
}

func randomPurchaseLimit() *shop_orm.PurchaseLimit {
	rand.Seed(time.Now().UnixNano())
	return &shop_orm.PurchaseLimit{
		ProductID:              randInt(999, 1999),
		LimitNum:               randInt(5, 10),
		StartPurchaseTimeStamp: int(time.Now().Unix()) + randInt(60, 600),
		StopPurchaseTimeStamp:  int(time.Now().Unix()) + randInt(6000, 12000),
	}
}
