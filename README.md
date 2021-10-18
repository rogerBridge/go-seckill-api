<!-- # 秒杀小游戏 -->

之前在掘金看过一篇关于商品秒杀系统的架构, 感觉挺有意思, 自己模仿着搞了一个, 感觉还可以, 就当一个 playground 吧~~

- [x] 限制商品单个用户可购买数量, 可购买时间段
- [x] 支持订单取消
- [x] 生成订单后内容传输给 rabbitmq, 写入 Mysql 数据库
- [x] 用户 ID(JWT)校验, 权限校验(URL 匹配), 单点登录
- [x] 一键 docker-compose 部署
- [ ] React SPA
- [ ] [单元测试](./test/unitTestReport.md)
- [x] [接口文档](https://documenter.getpostman.com/view/2753268/UUxtGBMW)

---

- [应用情景](#应用情景)
- [目录结构](#目录结构)
- [测试地址](#测试地址)
- [应用结构图](#应用结构图)
- [部署方法](#部署方法)
- [性能测试](#性能测试)

---

## 应用情景

卖出商品 a n 件, 限制单个用户购买 m 件, 限制购买时间段为: t1~t2

## 目录结构

```
go-seckill/
|-- api(包含postman项目导出文件, swagger.json(2.0))
|-- assets
    |-- disgrams(流程图文件)
|-- build
    |-- config-server(config-server容器制作)
    |-- go-seckill(go-seckill容器制作)
    |-- rabbitmqReceiver(rabbitmq-receiver容器制作)
    |-- rebuild.sh(docker-compose下线, 上述三个容器重新编译后, 删除旧image, 重新上线)
|-- cmd
    |-- go-seckill 主文件, App的入口
|-- deploy
    |-- mysql
    |-- redis(redis相关的配置文件, docker-compose容器启动需要使用)
    |-- docker-compose.yml
|-- interval
    |-- auth(JWT验证)
    |-- config(读取config-server的方法)
    |-- controllers2(fasthttp.RequestHandler, http请求处理)
    |-- db(gorm操作)
    |-- logconf(logrus相关设置)
    |-- rabbitmq(rabbitmq 通用channel和send函数)
    |-- redisconf(redis函数)
    |-- router(初始化路由设置, 部分路由加入auth middleware)
    |-- utils(通用设置)
|-- test
    |-- pressuremaker(压力测试, 包括测试用户生成, 测试用户token生成, 存储到sqlite, 开始测试, 测试数据分析等)
```

## 测试地址

示例:

获取所有商品信息: https://demo.fenr.men/api/v0/goodList

> 具体接口文档请查看: [接口文档](https://documenter.getpostman.com/view/2753268/UUxtGBMW)

测试用户:

普通用户: user:12345678

管理员: manager:12345678

## 应用结构图

- Docker 容器结构

![docker-structure](https://i.ibb.co/g9CYFKf/docker-structure.png)

| 容器名称           | 作用                                                                                              |
| ------------------ | ------------------------------------------------------------------------------------------------- |
| go-seckill         | 主应用, 监听特定端口, 处理请求                                                                    |
| rabbitmq-server    | amqp 服务, 用于接收应用发送的消息, 并由 rabbitmq-receiver 处理                                    |
| rabbitmq-receiver  | amqp 接收端, 处理 rabbitmq-server 发送过来的消息, 这里处理 app 发送过来的订单请求, 并存储在 mysql |
| tokenRedis         | 用户登录之后的 token 存储, 如果没有单点登录的需求,可以不用                                        |
| goodRedis          | 存储商品信息, 名称, 价格, 库存数量等　                                                            |
| orderRedis         | 存储用户产生的订单, 商品购买数量                                                                  |
| mysql              | 数据存储                                                                                          |
| config-file server | 存储 mysql, redis, rabbitmq 的配置文件                                                            |

## 部署方法

- docker-compose 部署:
  ```bash
  # 测试过的docker-compose版本为: 1.29.2
  # 如果服务器安装了golang sdk, 那么执行下面:
  cd go-seckill/build && bash rebuild.sh
  # 如果服务器没有安装golang sdk
  cd go-seckill/deploy && docker-compose -f docker-compose.yml up -d
  ```

## 性能测试

- 测试场景

  系统: Ubuntu 20.04 LTS

  go version: go1.16.5

  配置: Intel 8250U 功率(网络数据: 10Watt~13Watt)

```bash
➜  ~/Source/goLearn/go-seckill/test git:(master) go run .
2021/09/18 21:48:11 从viper读取到的mysql的配置是: root:12345678@tcp(db:3306)/seckill?charset=utf8mb4&parseTime=True&loc=Local
2021/09/18 21:48:11 gorm connect to mysql:  root:12345678@tcp(db:3306)/seckill?charset=utf8mb4&parseTime=True&loc=Local
2021/09/18 21:48:12 客户端总共发送请求: 10000个, 客户端角度的没有被服务器处理的请求数量:0
2021/09/18 21:48:12 在0~1秒内服务器就有返回的请求数量是: 9999
2021/09/18 21:48:12 在1~2秒内服务器就有返回的请求数量是: 1
2021/09/18 21:48:12 在2~3秒内服务器就有返回的请求数量是: 0
2021/09/18 21:48:12 在3~4秒内服务器就有返回的请求数量是: 0
2021/09/18 21:48:12 在4~5秒内服务器就有返回的请求数量是: 0
2021/09/18 21:48:12 大于5秒服务器返回的请求数量是: 0
2021/09/18 21:48:12 最大响应时间: 1012.7917ms, 最小响应时间: 1.4261ms, 平均响应时间: 501.6913ms, TPS: 9874
2021/09/18 21:48:12 0~1s 内处理的请求数量: 9999, 占总体请求数量的99.990%
➜  ~/Source/goLearn/go-seckill/test git:(master) go run .
2021/09/18 21:48:17 从viper读取到的mysql的配置是: root:12345678@tcp(db:3306)/seckill?charset=utf8mb4&parseTime=True&loc=Local
2021/09/18 21:48:17 gorm connect to mysql:  root:12345678@tcp(db:3306)/seckill?charset=utf8mb4&parseTime=True&loc=Local
2021/09/18 21:48:18 客户端总共发送请求: 10000个, 客户端角度的没有被服务器处理的请求数量:0
2021/09/18 21:48:18 在0~1秒内服务器就有返回的请求数量是: 10000
2021/09/18 21:48:18 在1~2秒内服务器就有返回的请求数量是: 0
2021/09/18 21:48:18 在2~3秒内服务器就有返回的请求数量是: 0
2021/09/18 21:48:18 在3~4秒内服务器就有返回的请求数量是: 0
2021/09/18 21:48:18 在4~5秒内服务器就有返回的请求数量是: 0
2021/09/18 21:48:18 大于5秒服务器返回的请求数量是: 0
2021/09/18 21:48:18 最大响应时间: 845.8818ms, 最小响应时间: 62.1917ms, 平均响应时间: 526.4314ms, TPS: 11822
2021/09/18 21:48:18 0~1s 内处理的请求数量: 10000, 占总体请求数量的100.000%
➜  ~/Source/goLearn/go-seckill/test git:(master) go run .
2021/09/18 21:49:13 从viper读取到的mysql的配置是: root:12345678@tcp(db:3306)/seckill?charset=utf8mb4&parseTime=True&loc=Local
2021/09/18 21:49:13 gorm connect to mysql:  root:12345678@tcp(db:3306)/seckill?charset=utf8mb4&parseTime=True&loc=Local
2021/09/18 21:49:14 客户端总共发送请求: 10000个, 客户端角度的没有被服务器处理的请求数量:0
2021/09/18 21:49:14 在0~1秒内服务器就有返回的请求数量是: 10000
2021/09/18 21:49:14 在1~2秒内服务器就有返回的请求数量是: 0
2021/09/18 21:49:14 在2~3秒内服务器就有返回的请求数量是: 0
2021/09/18 21:49:14 在3~4秒内服务器就有返回的请求数量是: 0
2021/09/18 21:49:14 在4~5秒内服务器就有返回的请求数量是: 0
2021/09/18 21:49:14 大于5秒服务器返回的请求数量是: 0
2021/09/18 21:49:14 最大响应时间: 877.8193ms, 最小响应时间: 0.4899ms, 平均响应时间: 499.2294ms, TPS: 11392
2021/09/18 21:49:14 0~1s 内处理的请求数量: 10000, 占总体请求数量的100.000%
➜  ~/Source/goLearn/go-seckill/test git:(master) go run .
2021/09/18 21:49:27 从viper读取到的mysql的配置是: root:12345678@tcp(db:3306)/seckill?charset=utf8mb4&parseTime=True&loc=Local
2021/09/18 21:49:27 gorm connect to mysql:  root:12345678@tcp(db:3306)/seckill?charset=utf8mb4&parseTime=True&loc=Local
2021/09/18 21:49:28 客户端总共发送请求: 10000个, 客户端角度的没有被服务器处理的请求数量:0
2021/09/18 21:49:28 在0~1秒内服务器就有返回的请求数量是: 10000
2021/09/18 21:49:28 在1~2秒内服务器就有返回的请求数量是: 0
2021/09/18 21:49:28 在2~3秒内服务器就有返回的请求数量是: 0
2021/09/18 21:49:28 在3~4秒内服务器就有返回的请求数量是: 0
2021/09/18 21:49:28 在4~5秒内服务器就有返回的请求数量是: 0
2021/09/18 21:49:28 大于5秒服务器返回的请求数量是: 0
2021/09/18 21:49:28 最大响应时间: 841.1800ms, 最小响应时间: 67.9813ms, 平均响应时间: 465.2974ms, TPS: 11888
2021/09/18 21:49:28 0~1s 内处理的请求数量: 10000, 占总体请求数量的100.000%
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
