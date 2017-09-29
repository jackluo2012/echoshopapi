# ShopApi
Some api for a shop app, use golang and echo framework.

## 运行
```shell
$ cd ShopApi/sever
$ go build
$ ./server
```
```shell
curl -X POST -H 'Content-Type: application/json' \
-d '{"mobile":"13528191831","password":"jackluo123"}' \
https://api.menu.zchd.ltd/api/v1/user/login


curl https://api.menu.zchd.ltd/api/v1/user/getinfo \
-H "x-access-token:eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1MDY2MDQwOTAsInVpZCI6Mn0.uaR_w4hFb3tVCGSDfRRUnvVPHnT76w1i4rzXANOoeJc"

```


- [x] 基本框架
- [x] 数据库
- [x] 用户系统
- [x] 商品系统
- [x] 交易系统
- [x] 地址管理系统
