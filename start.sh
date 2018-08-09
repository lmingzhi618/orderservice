#! /usr/bin/sh

# 1. get mysql
docker pull mysql

# 2. build mysql image
docker build -t order_mysql db

# 3. run mysql container
#docker run --rm -p 3306:3306 -e MYSQL_ROOT_PASSWORD=abc123456 --name order_mysql -d order_mysql
docker run --rm -p 3306:3306 --net=host -e MYSQL_ROOT_PASSWORD=abc123456 --name order_mysql -d order_mysql

# 4. build app service
docker build -t orderservice .

# 5. run app service and link mysql db
#docker run --rm -p 8080:8080 --link order_mysql:mysql --name orderservice -d orderservice 
docker run --rm -p 8080:8080 --net=host --name orderservice -d orderservice 


