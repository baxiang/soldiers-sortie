## 启动 consul
```shell
docker run --rm --name consul  -p 8500:8500 consul agent -dev -ui -client=0.0.0.0 -bind=0.0.0.0
```