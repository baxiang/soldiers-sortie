package client
import (
	grpctransport "github.com/go-kit/kit/transport/grpc"
	kitgrpc "github.com/go-kit/kit/transport/grpc"
	"github.com/baxiang/soldiers-sortie/zipkin-go/pb"
	endpts "github.com/baxiang/soldiers-sortie/zipkin-go/endpoint"
	"github.com/baxiang/soldiers-sortie/zipkin-go/service"
	"google.golang.org/grpc"
)

func StringDiff(conn *grpc.ClientConn, clientTracer kitgrpc.ClientOption) service.Service {

	var ep = grpctransport.NewClient(conn,
		"pb.StringService",
		"Diff",
		EncodeGRPCStringRequest,
		DecodeGRPCStringResponse,
		pb.StringResponse{},
		clientTracer,
	).Endpoint()

	StringEp := endpts.StringEndpoints{
		StringEndpoint: ep,
	}
	return StringEp
}