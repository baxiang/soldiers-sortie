

```shell
cd pb&&protoc ./*.proto  --go_out=plugins=grpc:.
```

source
https://github.com/pascallin/go-kit-application

## docker 
```shell
docker run -d --name zipkin -p 9411:9411 openzipkin/zipkin
```