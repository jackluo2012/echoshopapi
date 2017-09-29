#!/usr/bin/env bash
cd /root/go/src/ShopApi/server && go build && ./server
ps -ef|grep nginx| awk '{print $2}' | xargs kill -9
/usr/sbin/nginx
