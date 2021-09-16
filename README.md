# 商品抢购Demo

- 属性
  - [x] 限制商品单个用户可购买数量, 可购买时间段
  - [x] 支持订单取消
  - [x] 生成订单后内容传输给rabbitmq, 写入Mysql数据库
  - [x] 用户ID(JWT)校验, 权限校验(URL匹配)
  - [x] 一键docker-compose部署
  - [ ] React SPA
  - [ ] 单元测试 

- 应用情景
  
    卖出商品a n件, 限制单个用户购买 m件, 限制购买时间段为: t1~t2
    
- 结构图
  - Docker容器结构
  
  ![docker-structure](https://i.ibb.co/gW6NN4F/docker-structure.png)

- 部署方法
  - docker-compose部署:
    ```bash
    # 测试过的docker-compose版本为: 1.29.2
    cd deploy && docker-compose up -d
    ```

- 性能测试
  - 测试场景
    
    系统: Ubuntu 20.04 LTS
    
    go version: go1.16.5
  
    配置: Intel 8250U 10W功率
  
    ```text
    Start test
    2021/06/25 13:27:28 客户端总共发送请求: 10000个, 客户端角度的没有被服务器处理的请求数量:0
    2021/06/25 13:27:28 在0~1秒内服务器就有返回的请求数量是: 6360
    2021/06/25 13:27:28 在1~2秒内服务器就有返回的请求数量是: 3640
    2021/06/25 13:27:28 在2~3秒内服务器就有返回的请求数量是: 0
    2021/06/25 13:27:28 在3~4秒内服务器就有返回的请求数量是: 0
    2021/06/25 13:27:28 在4~5秒内服务器就有返回的请求数量是: 0
    2021/06/25 13:27:28 大于5秒服务器返回的请求数量是: 0
    2021/06/25 13:27:28 最大响应时间: 1620.0376ms, 最小响应时间: 18.7147ms, 平均响应时间: 906.2274ms, TPS: 6173
    2021/06/25 13:27:28 0~1s 内处理的请求数量: 6360, 占总体请求数量的63.600%
    INFO[2021-06-25T13:27:28+08:00]/home/leo/Source/goLearn/go-seckill/pressuretest/main.go:60 main.test() errChan info:  []                             app=go-seckill component=pressuretest
    ```
    
    ```text
    Start test 
    2021/06/25 13:29:23 客户端总共发送请求: 10000个, 客户端角度的没有被服务器处理的请求数量:0
    2021/06/25 13:29:23 在0~1秒内服务器就有返回的请求数量是: 5009
    2021/06/25 13:29:23 在1~2秒内服务器就有返回的请求数量是: 4991
    2021/06/25 13:29:23 在2~3秒内服务器就有返回的请求数量是: 0
    2021/06/25 13:29:23 在3~4秒内服务器就有返回的请求数量是: 0
    2021/06/25 13:29:23 在4~5秒内服务器就有返回的请求数量是: 0
    2021/06/25 13:29:23 大于5秒服务器返回的请求数量是: 0
    2021/06/25 13:29:23 最大响应时间: 1499.8092ms, 最小响应时间: 10.3933ms, 平均响应时间: 895.7024ms, TPS: 6668
    2021/06/25 13:29:23 0~1s 内处理的请求数量: 5009, 占总体请求数量的50.090%
    INFO[2021-06-25T13:29:23+08:00]/home/leo/Source/goLearn/go-seckill/pressuretest/main.go:60 main.test() errChan info:  []                             app=go-seckill component=pressuretest
    ```
  
    ```text
    Start test 
    2021/06/25 13:31:18 客户端总共发送请求: 10000个, 客户端角度的没有被服务器处理的请求数量:0
    2021/06/25 13:31:18 在0~1秒内服务器就有返回的请求数量是: 9989
    2021/06/25 13:31:18 在1~2秒内服务器就有返回的请求数量是: 11
    2021/06/25 13:31:18 在2~3秒内服务器就有返回的请求数量是: 0
    2021/06/25 13:31:18 在3~4秒内服务器就有返回的请求数量是: 0
    2021/06/25 13:31:18 在4~5秒内服务器就有返回的请求数量是: 0
    2021/06/25 13:31:18 大于5秒服务器返回的请求数量是: 0
    2021/06/25 13:31:18 最大响应时间: 1050.9305ms, 最小响应时间: 10.9800ms, 平均响应时间: 605.1730ms, TPS: 9515
    2021/06/25 13:31:18 0~1s 内处理的请求数量: 9989, 占总体请求数量的99.890%
    INFO[2021-06-25T13:31:18+08:00]/home/leo/Source/goLearn/go-seckill/pressuretest/main.go:60 main.test() errChan info:  []                             app=go-seckill component=pressuretest
    ```

    ```text
    容器资源使用:
    CONTAINER ID   NAME                                 CPU %     MEM USAGE / LIMIT     MEM %     NET I/O           BLOCK I/O        PIDS
    c045a2a02864   docker-compose_go-seckill_1          0.06%     536.2MiB / 15.38GiB   3.41%     159MB / 119MB     1.08MB / 0B      25
    a1f21b1ed696   docker-compose_rabbitmq-receiver_1   0.00%     7.094MiB / 15.38GiB   0.05%     706kB / 374kB     291kB / 0B       10
    e3a413cfcfa7   docker-compose_goodRedis_1           0.39%     15.31MiB / 15.38GiB   0.10%     18.4MB / 9.29MB   618kB / 20.5kB   5
    5df24df33b7e   docker-compose_db_1                  0.58%     342.1MiB / 15.38GiB   2.17%     8.67MB / 36.4MB   758kB / 30.7MB   41
    6d3d03d25e56   docker-compose_tokenRedis_1          1.08%     58.03MiB / 15.38GiB   0.37%     32.8MB / 36.6MB   213kB / 6.37MB   5
    f8f9f12d5ab0   docker-compose_orderRedis_1          0.33%     13.21MiB / 15.38GiB   0.08%     2.28MB / 1.74MB   8.19kB / 131kB   5
    dc6ece37bbf6   docker-compose_rabbitmq-server_1     0.89%     103.3MiB / 15.38GiB   0.66%     619kB / 340kB     1.14MB / 885kB   37
    ```
    