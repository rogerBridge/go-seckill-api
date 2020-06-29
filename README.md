# 单体redis商品抢购Demo

- 思路
用户请求的三个参数分别为: 用户Id, 商品Id, 商品数量
首先, 根据商品Id和商品数量判断是否可以购买, 如果可以, 生成订单信息(hash), key为:order:orderId, 例如: `order:fj34fjw`, 并且将信息rpush到: `user:userId:orderList` 这个list里面, 方便再次查找
```go
userId:orderList []string
// 使用rpush往orderList里面添加信息, 类似于:
// userId:orderList = []string{"uuid12345", "uuid12354"}
```

```go
order:orderId struct{
    UserId      string
    ProductId   string
    ProductName string
    OrderNum       int
    OrderTime   string
}
```

库存的数据结构表示:
```go
type StoreProductId struct{
    ProductId   string
    ProductName string
    StoreNum       int
}
```

- 啥时候做主从redis读写Demo?

    主从redis超卖的问题, 等我有空了再写(master用来写入数据, slave用来读出数据, 有可能库存为1的时候, 多个请求独到的库存都是1, 然后master减去了多个1, 这就完蛋了呀...)
