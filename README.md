# 单体redis商品抢购Demo

- 思路
从请求的body中解析出用户的id, 将成功的用户id(就是发送请求的时候: 库存商品数量 > -1)存放到redis自带的hash表中, 
暂定表名:
```go
type OrderForm struct{
    ProductName string
    ProductId   string
    ProductNum     int
    OrderTime   string
}
```
库存表名:
```go
type StoreForm struct{
    ProductName string
    ProductId   string
    StoreNum       int
}
```

- 啥时候做主从redis读写Demo?

    主从redis超卖的问题, 等我想明白了再写(master用来写入数据, slave用来读出数据, 有可能库存为1的时候, 多个请求独到的库存都是1, 然后master减去了多个1, 这就完蛋了呀...)
