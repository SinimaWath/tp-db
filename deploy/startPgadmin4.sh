#!/usr/bin/env bash

docker run -p 8080:80 --link db-forum:db --name pgadmin -e "PGADMIN_DEFAULT_EMAIL=sinimawath@gmail.com" -e "PGADMIN_DEFAULT_PASSWORD=admin" -d dpage/pgadmin4