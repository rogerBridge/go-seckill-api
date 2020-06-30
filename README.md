# 单体redis商品抢购Demo

- 思路

    用户请求的三个参数分别为: 用户Id, 商品Id, 商品数量, 其中, 用户Id默认是经过网关校验的, 这个demo里没有校验用户Id的功能的
    
- 功能清单:
    - [x] 支持多个商品抢购
    - [x] 限制用户可以购买的商品的总数量
    - [x] 自定义请求头部的解析 && 响应头部的解析

- 流程

    用户发送过来的请求, 首先会判断是否可以购买, 细节有: (不能超过限购数量, 不能超过库存数量), 如果满足购买条件, 就会生成订单, key为: `order:[orderId]`, value type为`hash`, value为: `userId int, productId string, orderNum int orderTime string`, 然后给用户相关的订单里面添加list, key为: `user:[userId]:orderNumList`, value type为: `list`, value为: `[orderNum]`, (orderNum的规则自定义, 这里定义的是单个字符的范围是: a-z, A-Z, 0-9, 长度为10的随机字符串), 最后是用户已经购买的商品id:商品数量, key为: `[user:userId:bought]`, value type为: `hash`, value为: `productId: purchaseNum`, 用这个可以快速的知道用户想要购买的某种商品是否已经超出了购买数量;

- 部署方法

    打包成二进制文件, 通过nginx转发, 或者直接使用裸二进制文件
    
- 性能测试结果

    
- 展望未来

    - [] 主从redis超卖问题
    
    主从redis超卖的问题, 等我有空了再写(master用来写入数据, slave用来读出数据, 有可能库存为1的时候, 多个请求独到的库存都是1, 然后master减去了多个1, 这就完蛋了呀...)