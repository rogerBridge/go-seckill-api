# 商品抢购 Demo

- 属性

  - [x] 限制商品单个用户可购买数量, 可购买时间段
  - [x] 支持订单取消
  - [x] 生成订单后内容传输给 rabbitmq, 写入 Mysql 数据库
  - [x] 用户 ID(JWT)校验, 权限校验(URL 匹配)
  - [x] 一键 docker-compose 部署
  - [ ] React SPA
  - [ ] [单元测试](./test/unitTestReport.md)

- 应用情景

  卖出商品 a n 件, 限制单个用户购买 m 件, 限制购买时间段为: t1~t2

- 结构图

  - Docker 容器结构

  ![docker-structure](https://i.ibb.co/gW6NN4F/docker-structure.png)

- 部署方法

  - docker-compose 部署:
    ```bash
    # 测试过的docker-compose版本为: 1.29.2
    cd deploy && docker-compose up -d
    ```

- 性能测试

  - 测试场景

    系统: Ubuntu 20.04 LTS

    go version: go1.16.5

    配置: Intel 8250U 功率(网络数据: 10Watt~13Watt)

  ```bash
  ➜ ~/Source/goLearn/go-seckill/test git:(master) go run .
  2021/09/17 00:14:19 从 viper 读取到的 mysql 的配置是: root:12345678@tcp(db:3306)/seckill?charset=utf8mb4&parseTime=True&loc=Local
  2021/09/17 00:14:19 gorm connect to mysql: root:12345678@tcp(db:3306)/seckill?charset=utf8mb4&parseTime=True&loc=Local
  2021/09/17 00:14:19 Start test
  2021/09/17 00:14:20 客户端总共发送请求: 10000 个, 客户端角度的没有被服务器处理的请求数量:0
  2021/09/17 00:14:20 在 0~1 秒内服务器就有返回的请求数量是: 10000
  2021/09/17 00:14:20 在 1~2 秒内服务器就有返回的请求数量是: 0
  2021/09/17 00:14:20 在 2~3 秒内服务器就有返回的请求数量是: 0
  2021/09/17 00:14:20 在 3~4 秒内服务器就有返回的请求数量是: 0
  2021/09/17 00:14:20 在 4~5 秒内服务器就有返回的请求数量是: 0
  2021/09/17 00:14:20 大于 5 秒服务器返回的请求数量是: 0
  2021/09/17 00:14:20 最大响应时间: 962.9638ms, 最小响应时间: 1.5064ms, 平均响应时间: 570.8454ms, TPS: 10385
  2021/09/17 00:14:20 0~1s 内处理的请求数量: 10000, 占总体请求数量的 100.000%
  ➜ ~/Source/goLearn/go-seckill/test git:(master) go run .
  2021/09/17 00:14:27 从 viper 读取到的 mysql 的配置是: root:12345678@tcp(db:3306)/seckill?charset=utf8mb4&parseTime=True&loc=Local
  2021/09/17 00:14:27 gorm connect to mysql: root:12345678@tcp(db:3306)/seckill?charset=utf8mb4&parseTime=True&loc=Local
  2021/09/17 00:14:27 Start test
  2021/09/17 00:14:28 客户端总共发送请求: 10000 个, 客户端角度的没有被服务器处理的请求数量:0
  2021/09/17 00:14:28 在 0~1 秒内服务器就有返回的请求数量是: 10000
  2021/09/17 00:14:28 在 1~2 秒内服务器就有返回的请求数量是: 0
  2021/09/17 00:14:28 在 2~3 秒内服务器就有返回的请求数量是: 0
  2021/09/17 00:14:28 在 3~4 秒内服务器就有返回的请求数量是: 0
  2021/09/17 00:14:28 在 4~5 秒内服务器就有返回的请求数量是: 0
  2021/09/17 00:14:28 大于 5 秒服务器返回的请求数量是: 0
  2021/09/17 00:14:28 最大响应时间: 711.5618ms, 最小响应时间: 2.8884ms, 平均响应时间: 408.3676ms, TPS: 14054
  2021/09/17 00:14:28 0~1s 内处理的请求数量: 10000, 占总体请求数量的 100.000%
  ➜ ~/Source/goLearn/go-seckill/test git:(master) go run .
  2021/09/17 00:14:33 从 viper 读取到的 mysql 的配置是: root:12345678@tcp(db:3306)/seckill?charset=utf8mb4&parseTime=True&loc=Local
  2021/09/17 00:14:33 gorm connect to mysql: root:12345678@tcp(db:3306)/seckill?charset=utf8mb4&parseTime=True&loc=Local
  2021/09/17 00:14:33 Start test
  2021/09/17 00:14:34 客户端总共发送请求: 10000 个, 客户端角度的没有被服务器处理的请求数量:0
  2021/09/17 00:14:34 在 0~1 秒内服务器就有返回的请求数量是: 10000
  2021/09/17 00:14:34 在 1~2 秒内服务器就有返回的请求数量是: 0
  2021/09/17 00:14:34 在 2~3 秒内服务器就有返回的请求数量是: 0
  2021/09/17 00:14:34 在 3~4 秒内服务器就有返回的请求数量是: 0
  2021/09/17 00:14:34 在 4~5 秒内服务器就有返回的请求数量是: 0
  2021/09/17 00:14:34 大于 5 秒服务器返回的请求数量是: 0
  2021/09/17 00:14:34 最大响应时间: 709.0308ms, 最小响应时间: 1.2449ms, 平均响应时间: 449.0126ms, TPS: 14104
  2021/09/17 00:14:34 0~1s 内处理的请求数量: 10000, 占总体请求数量的 100.000%
  ```

  ```bash
  容器资源使用:
  CONTAINER ID   NAME                         CPU %     MEM USAGE / LIMIT     MEM %     NET I/O           BLOCK I/O       PIDS
  26cc3fb78688   deploy_db_1                  0.47%     349.3MiB / 15.37GiB   2.22%     9.15MB / 37.5MB   786kB / 162MB   46
  afbe29ec43f9   deploy_go-seckill-config_1   0.00%     1.855MiB / 15.37GiB   0.01%     27kB / 33.2kB     0B / 0B         5
  b3308de42fe5   deploy_go-seckill_1          0.05%     754.8MiB / 15.37GiB   4.80%     204MB / 152MB     0B / 0B         45
  fb2e8bbf791c   deploy_goodRedis_1           0.94%     56.21MiB / 15.37GiB   0.36%     27.7MB / 15.2MB   0B / 16.4kB     5
  585849b4d36e   deploy_orderRedis_1          0.76%     43.76MiB / 15.37GiB   0.28%     3MB / 1.78MB      0B / 737kB      5
  811faf1d5ffd   deploy_rabbitmq-receiver_1   0.03%     5.352MiB / 15.37GiB   0.03%     1.48MB / 1.11MB   0B / 0B         11
  f34be970bc81   deploy_rabbitmq-server_1     1.79%     117.7MiB / 15.37GiB   0.75%     839kB / 716kB     0B / 2.84MB     37
  46918dcaa780   deploy_tokenRedis_1          1.63%     84.7MiB / 15.37GiB    0.54%     36.2MB / 42.6MB   0B / 7.27MB     5
  ```
