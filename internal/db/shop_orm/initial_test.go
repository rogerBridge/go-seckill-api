package shop_orm_test

import (
	"go-seckill/internal/db/shop_orm"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestInitDB(t *testing.T) {
	require.NoError(t, shop_orm.InitialMysql())
}