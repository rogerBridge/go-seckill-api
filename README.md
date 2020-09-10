# 单体redis商品抢购Demo

- 思路

    用户请求的三个参数分别为: 用户Id, 商品Id, 商品数量, 其中, 用户Id默认是经过网关校验的, 这里不做验证.
    
- 特征:

    - [x] 支持多个商品库存初始化, 同时抢购
    - [x] 限制用户是否可购买, 可购买的商品的总数量(userFilter函数内容可以自定义)
    - [x] 支持订单取消
    - [x] 订单号使用ksuid, 类似于uuid
    - [x] 生成订单后内容传输给队列, 如rabbitmq, 写入Mysql数据库
    - [ ] 完善前端+用户ID(合法性, cookie或者token)校验

- 应用情景

    1. 卖出商品a 100件, 限制每个人只能购买同款商品1件;
    2. 卖出商品a 100件, 商品b 100件, 商品a限制每个人只能购买2件, 商品b限制每个人只能购买1件;
    3. 多件商品, 设置特定时间段, 例如: a商品限制在xxxx.xx.xx xx:xx:xx ~ yyyy:yy.yy yy:yy:yy时间段内进行购买,
    每个人限制购买2件, 商品b购买特定时间段, 每个人限制购买5件, 商品c不做限制;
    
    
- 流程
    

- 部署方法
    - Docker部署:
    1. 实验版:
        - 两个Redis实例, 分别位于目录: redisDocker/runRedis.sh, redisDocker/orderInfoRedis/runRedis.sh;
        - 一个MySQL实例, 位于目录: mysql/runMysql.sh
        - 一个webapp实例, 位于: buils.sh
        - 一键部署: `bash sharpRun.sh`
        - 建议: 运行在本地, 搭配Nginx反向代理服务器, 配合TLS使用

- 性能测试

    - 测试场景
    
        系统: Ubuntu 20.04 LTS
    
        CPU: Intel i5 8250U (4H8T) (runtime.GOMAXPROCS设置为1个线程, 但通过htop发现, 8个线程全部被占用了, 我也不知道咋回事,  This call will go away when the scheduler improves, 意思是升级了scheduler所以这个go sway了?)
    
        go version: go1.14.4 linux/amd64
        
        内存占用(反复跑pressure_test后停在了这个地方, peak value 400MB):
    ![pressure_test_memory](./img/pressure.png)
    1. 用户20000名, 请求: /buy, 购买商品, 商品ID: "10000", 购买数量: 1, 库存数量: 200件, 测试结果如下:
    连续 5 次测试:
    
    ![1.1](./img/1.1.png)
    
    2. 20000名用户, 其中10000名用户抢购商品"10000", 库存为200, 10000名用户抢购商品"10001", 数量为200, 测试结果如下:
    连续 5 次测试:
    
    ![2.1](./img/2.1.png)