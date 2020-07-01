# 单体redis商品抢购Demo

- 思路

    用户请求的三个参数分别为: 用户Id, 商品Id, 商品数量, 其中, 用户Id默认是经过网关校验的, 这里不做验证.
    
- 功能清单:
    - [x] 支持多个商品库存初始化, 同时抢购
    - [x] 限制用户是否可购买, 可购买的商品的总数量(userFilter函数内容可以自定义)
    - [x] 自定义请求头部 && 响应头部
    - [ ] 生成订单后内容传输给队列, 如rabbitmq, 写入Mysql数据库
    - [ ] 完善前端+用户ID(合法性, cookie或者token)校验

- 流程

    用户发送过来的请求, 首先会判断是否可以购买, 细节有: (不能超过限购数量, 不能超过库存数量), 如果满足购买条件, 就会生成订单, key为: `user:[userId]:order:[orderId]`, value type为`hash`, value为: `userId int, productId string, orderNum int orderTime string`, 然后给用户相关的订单里面添加list, key为: `user:[userId]:orderNumList`, value type为: `list`, value为: `[orderNum]`, (orderNum的规则自定义, 这里定义的是单个字符的范围是: a-z, A-Z, 0-9, 长度为10的随机字符串), 最后是用户已经购买的商品id:商品数量, key为: `[user:userId:bought]`, value type为: `hash`, value为: `productId: purchaseNum`, 用这个可以快速的知道用户想要购买的某种商品是否已经超出了购买数量;

- 部署方法

    打包成二进制文件, 通过nginx转发, 或者直接使用裸二进制文件

- 部署流程(后续还是搞一个docker吧~)

    1. 部署redis, 端口号: 6379, AUTH: "hello"
    2. cd redis_play && go build -o redis_play *.go && ./redis_play
    3. 运行redis_play, 如果使用的是VPS, 注意打开公网端口号: 6379
    4. 可以使用Postman测试, 或者用我的压测脚本, 位于pressure_test目录下,
       目前只覆盖了两种场景
    
- 性能测试

    - 测试场景
    1. 用户20000名, 同时请求: /buy, 购买商品, 商品ID: "10000", 购买数量: 1库存数量: 200件, 测试结果如下:
    
    第一次测试:
    ```text
    2020/07/01 17:08:33 每秒事务处理量: 3229.36, 20000个客户端请求总时间段: 6.1932s
    2020/07/01 17:08:33 无效请求数量: 0
    2020/07/01 17:08:33 在0~1秒内服务器就有返回的请求数量是: 1718
    2020/07/01 17:08:33 在1~2秒内服务器就有返回的请求数量是: 1835
    2020/07/01 17:08:33 在2~3秒内服务器就有返回的请求数量是: 4988
    2020/07/01 17:08:33 在3~4秒内服务器就有返回的请求数量是: 822
    2020/07/01 17:08:33 在4~5秒内服务器就有返回的请求数量是: 4960
    2020/07/01 17:08:33 在大于5秒内服务器就有返回的请求数量是: 5677

    ```
    检查redis-db0, 发现没有超卖, 订单生成正常, 用户购物车数值正常
    
    第二次测试:
    ```text
    2020/07/01 17:10:05 每秒事务处理量: 2737.18, 20000个客户端请求总时间段: 7.3068s
    2020/07/01 17:10:05 无效请求数量: 0
    2020/07/01 17:10:05 在0~1秒内服务器就有返回的请求数量是: 442
    2020/07/01 17:10:05 在1~2秒内服务器就有返回的请求数量是: 476
    2020/07/01 17:10:05 在2~3秒内服务器就有返回的请求数量是: 1726
    2020/07/01 17:10:05 在3~4秒内服务器就有返回的请求数量是: 1059
    2020/07/01 17:10:05 在4~5秒内服务器就有返回的请求数量是: 4584
    2020/07/01 17:10:05 在大于5秒内服务器就有返回的请求数量是: 11713

    ```
    检查redis-db0, 发现没有超卖, 订单生成正常, 用户购物车数值正常
    第三次测试:
    ```text
    2020/07/01 17:13:00 每秒事务处理量: 3422.53, 20000个客户端请求总时间段: 5.8436s
    2020/07/01 17:13:00 无效请求数量: 0
    2020/07/01 17:13:00 在0~1秒内服务器就有返回的请求数量是: 1193
    2020/07/01 17:13:00 在1~2秒内服务器就有返回的请求数量是: 1144
    2020/07/01 17:13:00 在2~3秒内服务器就有返回的请求数量是: 3800
    2020/07/01 17:13:00 在3~4秒内服务器就有返回的请求数量是: 2964
    2020/07/01 17:13:00 在4~5秒内服务器就有返回的请求数量是: 6628
    2020/07/01 17:13:00 在大于5秒内服务器就有返回的请求数量是: 4271

    ```
    检查redis-db0, 发现没有超卖, 订单生成正常, 用户购物车数值正常
    
    2. 20000名用户, 其中10000名用户抢购商品"10000", 库存为200, 10000名用户抢购商品"10001", 数量为200, 测试结果如下:
    
    第一次测试:
    ```text
    2020/07/01 17:22:24 每秒事务处理量: 2849.93, 20000个客户端请求总时间段: 7.0177s
    2020/07/01 17:22:24 无效请求数量: 0
    2020/07/01 17:22:24 在0~1秒内服务器就有返回的请求数量是: 638
    2020/07/01 17:22:24 在1~2秒内服务器就有返回的请求数量是: 1
    2020/07/01 17:22:24 在2~3秒内服务器就有返回的请求数量是: 258
    2020/07/01 17:22:24 在3~4秒内服务器就有返回的请求数量是: 1371
    2020/07/01 17:22:24 在4~5秒内服务器就有返回的请求数量是: 1800
    2020/07/01 17:22:24 在大于5秒内服务器就有返回的请求数量是: 15932
    ```
    检查redis-db0, 发现没有超卖, 订单生成正常, 用户购物车数值正常
  
    第二次测试:
    ```text
    2020/07/01 17:24:53 每秒事务处理量: 3287.84, 20000个客户端请求总时间段: 6.0830s
    2020/07/01 17:24:53 无效请求数量: 0
    2020/07/01 17:24:53 在0~1秒内服务器就有返回的请求数量是: 511
    2020/07/01 17:24:53 在1~2秒内服务器就有返回的请求数量是: 298
    2020/07/01 17:24:53 在2~3秒内服务器就有返回的请求数量是: 4013
    2020/07/01 17:24:53 在3~4秒内服务器就有返回的请求数量是: 3032
    2020/07/01 17:24:53 在4~5秒内服务器就有返回的请求数量是: 6913
    2020/07/01 17:24:53 在大于5秒内服务器就有返回的请求数量是: 5233
    ```
    检查redis-db0, 发现没有超卖, 订单生成正常, 用户购物车数值正常
  
    第三次测试:
    ```text
    2020/07/01 17:25:50 每秒事务处理量: 3287.01, 20000个客户端请求总时间段: 6.0846s
    2020/07/01 17:25:50 无效请求数量: 0
    2020/07/01 17:25:50 在0~1秒内服务器就有返回的请求数量是: 789
    2020/07/01 17:25:50 在1~2秒内服务器就有返回的请求数量是: 1729
    2020/07/01 17:25:50 在2~3秒内服务器就有返回的请求数量是: 3171
    2020/07/01 17:25:50 在3~4秒内服务器就有返回的请求数量是: 4344
    2020/07/01 17:25:50 在4~5秒内服务器就有返回的请求数量是: 7227
    2020/07/01 17:25:50 在大于5秒内服务器就有返回的请求数量是: 2740
    ```
    检查redis-db0, 发现没有超卖, 订单生成正常, 用户购物车数值正常
    
- 展望未来

    - [] 主从redis超卖问题
    
    主从redis超卖的问题, 等我有空了再写(master用来写入数据, slave用来读出数据, 有可能库存为1的时候, 多个请求独到的库存都是1, 然后master减去了多个1, 这就完蛋了呀...)
