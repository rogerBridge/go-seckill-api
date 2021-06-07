# 单体redis商品抢购Demo

- 思路

    用户请求的三个参数分别为: 用户Id, 商品Id, 商品数量, ~~其中, 用户Id默认是经过网关校验的, 这里不做验证.~~(JWT校验)
    
- 特征:

    - [x] 限制商品单个用户可购买数量, 可购买时间段
    - [x] 支持订单取消
    - [x] 生成订单后内容传输给rabbitmq, 写入Mysql数据库, 去峰
    - [x] 用户ID(JWT)校验
    - [x] 一键docker-compose部署
    - [ ] vue单页应用 

- 应用情景

    1. 卖出商品a 100件, 限制每个人只能购买同款商品1件;
    2. 卖出商品a 100件, 商品b 100件, 商品a限制每个人只能购买2件, 商品b限制每个人只能购买1件;
    3. 多件商品, 设置特定时间段, 例如: a商品限制在xxxx.xx.xx xx:xx:xx ~ yyyy:yy.yy yy:yy:yy时间段内进行购买,
    每个人限制购买2件, 商品b购买特定时间段, 每个人限制购买5件, 商品c不做限制;
    4. 卖出商品x件, 靠延迟和运气抢;
    
- 结构图
    

- 部署方法
  - docker-compose部署:
    ```bash
    # 测试过的docker-compose版本为: 1.29.2
    cd deploy && docker-compose up
    ```

- 性能测试

    - 测试场景
    
        系统: Ubuntu 20.04 LTS
    
        实例 : Amazon Lightsail 2H4G $20/month
        
        go version: go1.16.5