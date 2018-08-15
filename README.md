# 通过docker 的 swarm 功能，可以将应用分布式部曙到几个机器上（物理/虚拟机均可）
# 此处通过虚拟机部曙
# 一. 创建虚拟机并配置节点
## 1. 创建两台虚拟机
### $ docker-machine create --driver virtualbox myvm1
### $ docker-machine create --driver virtualbox myvm2

## 2. 查看 vm machine 信息
### $ docker-machine ls
### NAME    ACTIVE   DRIVER       STATE     URL                         SWARM   DOCKER        ERRORS
### myvm1   -        virtualbox   Running   tcp://192.168.99.100:2376           v17.06.2-ce
### myvm2   -        virtualbox   Running   tcp://192.168.99.101:2376           v17.06.2-ce

## 3. 配置myvm1为 manager 节点
### $ docker-machine ssh myvm1 "docker swarm init --advertise-addr <myvm1 ip>"
### Swarm initialized: current node <node ID> is now a manager.
###
### To add a worker to this swarm, run the following command:
###
###  docker swarm join \
###  --token <token> \
###  <myvm ip>:<port>
###
### To add a manager to this swarm, run 'docker swarm join-token manager' and follow the instructions.

## 4. 添加 worker 节点
### $ docker-machine ssh myvm2 "docker swarm join \
### --token <token> \
### <ip>:2377"
### 
### This node joined a swarm as a worker.

## 5. 查看节点信息
### $ docker-machine ssh myvm1 "docker node ls"
### ID                            HOSTNAME            STATUS              AVAILABILITY        MANAGER STATUS
### i8dvxx0ct3p7au19oaknd5nsq *   myvm1                   Ready               Active              Leader              18.06.0-ce
### vkbvfxkoijs0mv1w7fg4ru8rp     myvm2                   Ready               Active                                  18.06.0-ce

# 二. 制作并部曙服务
## 1. 制作好image
### 1.1 docker build -t order_mysql ../db
### 1.2 docker build -t orderservice .

## 2. 登录 docker
###$ docker login

## 3. 上传image到docker
###$ docker tag orderservice mingzhi198/orderservice

## 4. 部曙到 swarm
### $ docker stack deploy -c docker-compose.yml order_svr

## 5. 关闭swarm服务
### $ docker stack rm order_svr

## 6. 也可不创建vm部曙服务，区别为container都跑在本机上

## 7. 为减小image的尺寸，用scratch代替golang基础包, 区别请对比image_golang/Dockerfile

================================================================================
    因为国内网络及google配额问题，用腾讯地图api
    建议在linux/mac 下运行该服务
    1. 确保已经安装好docker
    2. 运行脚本 sh start.sh 部署服务
    3. 本地测试命令
    
        curl 'localhost:8080/order' -d '
        {
            "origin": ["39.983171", "116.308479"],
            "destination": ["39.99606", "116.353455"]
        }'
        
        curl -X PUT 'localhost:8080/order/10' -d '
        {
            "status":"taken"
        }'
    
        curl 'localhost:8080/orders?page=1&limit=5'
    
    
