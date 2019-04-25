#! /usr/bin/sh

# build mysql image
docker build -t mysql_order./db

# run mysql container
docker run --name db --rm -v /db_data:/var/lib/mysql  -p 3306:3306 -e MYSQL_ROOT_PASSWORD=abc123456  -d mysql_order

# build app service
docker build -t orderservice .

# run app service and link mysql db
docker run --rm -p 8080:8080 --link db:mysql --name ordersvr -d orderservice 


