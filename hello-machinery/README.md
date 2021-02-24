
发送端
```shell
 go run cmd/main.go send
```

接收端
```shell
go run cmd/main.go worker
```

docker
```shell
docker run -d --name redis -p 6379:6379  redis:4.0
```

test
```shell
curl localhost:8080/add
```

demo
https://github.com/RichardKnop/machinery/tree/master/v2/example