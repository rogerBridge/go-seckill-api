# 单体redis商品抢购Demo

- 思路

    用户请求的三个参数分别为: 用户Id, 商品Id, 商品数量, 其中, 用户Id默认是经过网关校验的, 这里不做验证.
    
- 特征:

    - [x] 支持多个商品库存初始化, 同时抢购
    - [x] 限制用户是否可购买, 可购买的商品的总数量(userFilter函数内容可以自定义)
    - [x] 支持订单取消
    - [x] 订单号使用ksuid, 类似于uuid
    - [ ] 生成订单后内容传输给队列, 如rabbitmq, 写入Mysql数据库
    - [ ] 完善前端+用户ID(合法性, cookie或者token)校验

- 应用情景

    1. 卖出商品a 100件, 限制每个人只能购买同款商品1件;
    2. 卖出商品a 100件, 商品b 100件, 商品a限制每个人只能购买2件, 商品b限制每个人只能购买1件;
    3. 卖出商品a x件(x>0 && type(x)=="int"), 限购aa件(aa<=x). 商品b y件(y>0 && type(y)=="int"), 限购bb件(bb<=y), 商品c z件(z>0&& type(z)=="int"), 限购cc件(cc<=z) ......, 
    
    
- 流程
    

- 部署方法

    - 打包成二进制文件, 通过nginx转发, 或者直接使用裸二进制文件
        1. 部署redis, 端口号: 6379, AUTH: "hello"
        2. cd redis_play && go build -o redis_play *.go && ./redis_play
        3. 运行redis_play, 默认情况下, redis只能通过localhost连接, 为了安全嘛~, 通过http来调用redis就好了呀, 然后把端口暴露给公网.
        4. 可以使用Postman测试功能, 自己写或者用我的位于pressure_test目录下的压测脚本测试并发性能, 
           目前只覆盖了两种场景(1. 10000名用户同时抢购productId: 10000的商品(库存:200), 2. 10000名用户同时抢购2种商品, productId:10000, productId:10001, per 200件库存)
    - 做成docker
        - [x] 完成

- 性能测试

    - 测试场景
    
        系统: Ubuntu 20.04 LTS
    
        CPU: Intel i5 8250U (4H8T) (runtime.GOMAXPROCS设置为1个线程, 但通过htop发现, 8个线程全部被占用了, 我也不知道咋回事,  This call will go away when the scheduler improves, 意思是升级了scheduler所以这个go sway了?)
    
        go version: go1.14.4 linux/amd64
        
        内存占用(反复跑pressure_test后停在了这个地方, peak value 400MB):
    ![pressure_test_memory](./img/pressure.png)
    1. 用户10000名, 同时请求: /buy, 购买商品, 商品ID: "10000", 购买数量: 1库存数量: 200件, 测试结果如下:
    
    第一次测试:
    
    ![1.1](./img/1.1.png)
    
    检查redis-db0, 发现没有超卖, 订单生成正常, 用户购物车数值正常
    
    第二次测试:
    
    ![1.2](./img/1.2.png)
    
    第三次测试:
    
    ![1.3](./img/1.3.png)
    
    2. 10000名用户, 其中5000名用户抢购商品"10000", 库存为200, 5000名用户抢购商品"10001", 数量为200, 测试结果如下:
    
    第一次测试:
    
    ![2.1](./img/2.1.png)
    检查redis-db0, 发现没有超卖, 订单生成正常, 用户购物车数值正常
  
    第二次测试:
    
    ![2.2](./img/2.2.png)
    
    第三次测试:
    
    ![2.3](./img/2.3.png)