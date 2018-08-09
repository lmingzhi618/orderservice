因为国内网络原因，google maps api 访问失败, 且 google maps api 的调用配额有限, 因此用 腾讯地图 api 代替
建议在ubuntu 14.04 下运行该服务
1. 确保已经安装好docker
2. 运行脚本 sh start.sh 部署服务
3. 如需单元测试，在 orderserver 目录下运行命令： go test

