#! /usr/bin/sh

# add read and write privilege for all users(mysql)
mkdir db_data && chmod 777 db_data
# build mysql image
docker build -t order_mysql ./db/

# build app service, the image has been pushed to docker 
#docker build -t orderservice .

# get image from docker and deploy
docker stack deploy -c docker-compose.yml order_svr


