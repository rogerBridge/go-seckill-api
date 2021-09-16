package shop_orm_test

import (
	"go-seckill/internal/db/shop_orm"
	"math/rand"
	"testing"
	"time"

	"github.com/segmentio/ksuid"
	"github.com/stretchr/testify/require"
)

func TestOrder(t *testing.T) {
	// create order
	order := randomOrder()
	require.NoError(t, order.CreateOrder(conn))
	// query order
	require.Equal(t, order, order.QueryOrderByOrderNumber(order.OrderNumber))
	// update order
	order.ProductID = 1099 // updated: order.ProductID+100
	order.PurchaseNum = 2
	require.NoError(t, order.UpdateOrder(conn))
	require.Equal(t, order.ProductID, order.QueryOrderByOrderNumber(order.OrderNumber).ProductID)
	require.Equal(t, order.PurchaseNum, order.QueryOrderByOrderNumber(order.OrderNumber).PurchaseNum)
	// delete order
	require.NoError(t, order.DeleteOrderByOrderNumber(conn))
	require.NotEqual(t, order, order.QueryOrderByOrderNumber(order.OrderNumber))
}

func randomOrder() *shop_orm.Order {
	rand.Seed(time.Now().UnixNano())
	order := &shop_orm.Order{
		SelfDefine:  shop_orm.SelfDefine{Version: "v0.0.0"},
		OrderNumber: ksuid.New().String(),
		Username:    randomString(8),
		ProductID:   999,
		PurchaseNum: 1,
		Price:       1999,
		Status:      randomOrderStatus(),
	}
	return order
}

func randomString(length int) string {
	baseString := "abcdefghigklmnopqrstuvwxyz0123456789"
	wantBytes := make([]byte, length)
	for i := 0; i < length; i++ {
		wantBytes[i] = baseString[rand.Intn(len(baseString))]
	}
	return string(wantBytes)
}

func randomOrderStatus() string {
	statusSlice := []string{
		"cancel", "process", "finished",
	}
	return statusSlice[rand.Intn(len(statusSlice))]
}
