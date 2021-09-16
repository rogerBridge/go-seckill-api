package shop_orm_test

import (
	"go-seckill/internal/db/shop_orm"
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestUser(t *testing.T) {
	// // create admin
	// admin := randomUser("admin")
	// require.NoError(t, admin.CreateAdmin(conn))
	// create user
	user := randomUser("user")
	require.NoError(t, user.CreateUser(conn))
	// query user by email
	require.Equal(t, user, user.QueryUserByEmail(user.Email))
	// update user
	user.Address = randomString(16)
	user.Birthday = time.Now().UnixNano()/1e6
	require.NoError(t, user.UpdateUserInfo(conn))
	require.Equal(t, user.Address, user.QueryUserByEmail(user.Email).Address)
	require.Equal(t, user.Email, user.QueryUserByEmail(user.Email).Email)
	// delete user
	require.NoError(t, user.DeleteUserByUserEmail(user.Email))
	require.NotEqual(t, user.QueryUserByEmail(user.Email), user)
}

func randomUser(group string) *shop_orm.User {
	rand.Seed(time.Now().UnixNano())
	username := randomString(8)
	birthday := time.Now().UnixNano() / 1e6
	user := &shop_orm.User{
		SelfDefine: shop_orm.SelfDefine{Version: "v0.0.0"},
		Username:   username,
		Password:   "12345678",
		Group:      group,
		Address:    randomString(16),
		Email:      username + "@gmail.com",
		Birthday:   birthday,
	}
	return user
}
