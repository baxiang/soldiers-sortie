
## 参考
https://github.com/eddycjy/go-gin-example

## swag 
https://github.com/swaggo/swag
```shell
go get -u github.com/swaggo/swag/cmd/swag
swag -v
swag init -dir cmd
```
## docker
```shell
docker run -d --name mysql -p3306:3306 -e MYSQL_DATABASE=blog -e MYSQL_ROOT_PASSWORD=123456 mysql:5.7
```

redis
```shell
docker run -d --name redis -p6379:6379  redis:4.0
```