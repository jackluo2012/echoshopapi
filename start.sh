#!/bin/bash
docker stop nginx
docker rm nginx
docker run -d --name nginx \
            -p 80:80 \
            -p 443:443 \
            --link mysql:mysql.io \
            --add-host api.menu.zchd.ltd:127.0.0.1 \
            -v /Users/jackluo/Works/golang:/root/go \
            -v $(pwd)/nginx:/etc/nginx/conf.d \
            wemall

#docker exec nginx /root/go/src/wemall/npm.sh
docker exec nginx /usr/sbin/nginx
docker exec -it nginx /bin/bash
