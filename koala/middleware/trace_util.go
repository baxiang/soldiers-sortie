package middleware

import (
	"context"
	"encoding/base64"
	"fmt"
	"github.com/baxiang/koala/logs"
	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go/config"
	"github.com/uber/jaeger-client-go/transport/zipkin"
	"google.golang.org/grpc/metadata"
	"github.com/uber/jaeger-client-go"
	"strings"
)

const (
	binHdrSuffix = "-bin"
)

type metadataTextMap metadata.MD


func (m metadataTextMap)Set(key,val string){
	encodedKey, encodedVal := encodeKeyValue(key, val)
	m[encodedKey] = []string{encodedVal}
}


func (m metadataTextMap) ForeachKey(callback func(key, val string) error) error {
	for k, vv := range m {
		for _, v := range vv {
			if decodedKey, decodedVal, err := metadata.DecodeKeyValue(k, v); err == nil {
				if err = callback(decodedKey, decodedVal); err != nil {
					return err
				}
			} else {
				return fmt.Errorf("failed decoding opentracing from gRPC metadata: %v", err)
			}
		}
	}
	return nil
}

func encodeKeyValue(key,val string)(string,string){
	key = strings.ToLower(key)
	if strings.HasSuffix(key,binHdrSuffix){
		val =base64.StdEncoding.EncodeToString([]byte(val))
	}
	return key,val
}

func InitTrace(serviceName,reportAddr,sampleType string,rate float64)(err error){
	transport, err := zipkin.NewHTTPTransport(reportAddr, zipkin.HTTPBatchSize(16), zipkin.HTTPLogger(jaeger.StdLogger))
	if err != nil {
		logs.Error(context.TODO(), "ERROR: cannot init zipkin: %v\n", err)
		return
	}
	cfg := &config.Configuration{
		ServiceName:serviceName,
		Sampler: &config.SamplerConfig{
			Type:  sampleType,
			Param: rate,
		},
		Reporter: &config.ReporterConfig{
			LogSpans: true,
		},
	}

	r := jaeger.NewRemoteReporter(transport)
	tracer, closer, err := cfg.NewTracer(
		config.Logger(jaeger.StdLogger),
		config.Reporter(r))
	if err != nil {
		logs.Error(context.TODO(), "ERROR: cannot init Jaeger: %v\n", err)
		return
	}
	_ = closer
	opentracing.SetGlobalTracer(tracer)
	return
}