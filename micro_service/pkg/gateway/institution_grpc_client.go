package gateway

import (
	"context"
	"userService/pkg/pb"

	"github.com/go-kit/kit/endpoint"
	grpctransport "github.com/go-kit/kit/transport/grpc"
	kitgrpc "github.com/go-kit/kit/transport/grpc"
	"google.golang.org/grpc"
)

type InstitutionEndpoints struct {
	TnxHisDownloadEndpoint                endpoint.Endpoint
	GetTfrTrnLogsEndpoint                 endpoint.Endpoint
	GetTfrTrnLogEndpoint                  endpoint.Endpoint
	DownloadTfrTrnLogsEndpoint            endpoint.Endpoint
	ListGroupsEndpoint                    endpoint.Endpoint
	SaveGroupEndpoint                     endpoint.Endpoint
	ListInstitutionsEndpoint              endpoint.Endpoint
	SaveInstitutionEndpoint               endpoint.Endpoint
	GetInstitutionByIdEndpoint            endpoint.Endpoint
	SaveInstitutionFeeControlCashEndpoint endpoint.Endpoint
	GetInstitutionControlEndpoint         endpoint.Endpoint
	GetInstitutionCashEndpoint            endpoint.Endpoint
	GetInstitutionFeeEndpoint             endpoint.Endpoint
	BindGroupEndpoint                     endpoint.Endpoint
	ListBindGroupEndpoint                 endpoint.Endpoint
	RemoveBindGroupEndpoint               endpoint.Endpoint
}

func (s *InstitutionEndpoints) RemoveBindGroup(ctx context.Context, in *pb.RemoveBindGroupRequest) (*pb.RemoveBindGroupReply, error) {
	res, err := s.RemoveBindGroupEndpoint(ctx, in)
	if err != nil {
		return nil, err
	}
	return res.(*pb.RemoveBindGroupReply), nil
}

func (s *InstitutionEndpoints) ListBindGroup(ctx context.Context, in *pb.ListBindGroupRequest) (*pb.ListBindGroupReply, error) {
	res, err := s.ListBindGroupEndpoint(ctx, in)
	if err != nil {
		return nil, err
	}
	return res.(*pb.ListBindGroupReply), nil
}

func (s *InstitutionEndpoints) BindGroup(ctx context.Context, in *pb.BindGroupRequest) (*pb.BindGroupReply, error) {
	res, err := s.BindGroupEndpoint(ctx, in)
	if err != nil {
		return nil, err
	}
	return res.(*pb.BindGroupReply), nil
}

func (s *InstitutionEndpoints) SaveGroup(ctx context.Context, in *pb.SaveGroupRequest) (*pb.SaveGroupReply, error) {
	res, err := s.SaveGroupEndpoint(ctx, in)
	if err != nil {
		return nil, err
	}
	return res.(*pb.SaveGroupReply), nil
}

func (s *InstitutionEndpoints) GetInstitutionControl(ctx context.Context, in *pb.GetInstitutionControlRequest) (*pb.GetInstitutionControlReply, error) {
	res, err := s.GetInstitutionControlEndpoint(ctx, in)
	if err != nil {
		return nil, err
	}
	return res.(*pb.GetInstitutionControlReply), nil
}

func (s *InstitutionEndpoints) GetInstitutionCash(ctx context.Context, in *pb.GetInstitutionCashRequest) (*pb.GetInstitutionCashReply, error) {
	res, err := s.GetInstitutionCashEndpoint(ctx, in)
	if err != nil {
		return nil, err
	}
	return res.(*pb.GetInstitutionCashReply), nil
}

func (s *InstitutionEndpoints) GetInstitutionFee(ctx context.Context, in *pb.GetInstitutionFeeRequest) (*pb.GetInstitutionFeeReply, error) {
	res, err := s.GetInstitutionFeeEndpoint(ctx, in)
	if err != nil {
		return nil, err
	}
	return res.(*pb.GetInstitutionFeeReply), nil
}

func (s *InstitutionEndpoints) SaveInstitutionFeeControlCash(ctx context.Context, in *pb.SaveInstitutionFeeControlCashRequest) (*pb.SaveInstitutionFeeControlCashReply, error) {
	res, err := s.SaveInstitutionFeeControlCashEndpoint(ctx, in)
	if err != nil {
		return nil, err
	}
	return res.(*pb.SaveInstitutionFeeControlCashReply), nil
}

func (s *InstitutionEndpoints) GetInstitutionById(ctx context.Context, in *pb.GetInstitutionByIdRequest) (*pb.GetInstitutionByIdReply, error) {
	res, err := s.GetInstitutionByIdEndpoint(ctx, in)
	if err != nil {
		return nil, err
	}
	return res.(*pb.GetInstitutionByIdReply), nil
}

func (s *InstitutionEndpoints) TnxHisDownload(ctx context.Context, in *pb.InstitutionTnxHisDownloadReq) (*pb.InstitutionTnxHisDownloadResp, error) {
	res, err := s.TnxHisDownloadEndpoint(ctx, in)
	if err != nil {
		return nil, err
	}
	return res.(*pb.InstitutionTnxHisDownloadResp), nil
}

func (s *InstitutionEndpoints) GetTfrTrnLogs(ctx context.Context, in *pb.GetTfrTrnLogsReq) (*pb.GetTfrTrnLogsResp, error) {
	res, err := s.GetTfrTrnLogsEndpoint(ctx, in)
	if err != nil {
		return nil, err
	}
	return res.(*pb.GetTfrTrnLogsResp), nil
}

func (s *InstitutionEndpoints) GetTfrTrnLog(ctx context.Context, in *pb.GetTfrTrnLogReq) (*pb.GetTfrTrnLogResp, error) {
	res, err := s.GetTfrTrnLogEndpoint(ctx, in)
	if err != nil {
		return nil, err
	}
	return res.(*pb.GetTfrTrnLogResp), nil
}

func (s *InstitutionEndpoints) DownloadTfrTrnLogs(ctx context.Context, in *pb.DownloadTfrTrnLogsReq) (*pb.DownloadTfrTrnLogsResp, error) {
	res, err := s.DownloadTfrTrnLogsEndpoint(ctx, in)
	if err != nil {
		return nil, err
	}
	return res.(*pb.DownloadTfrTrnLogsResp), nil
}

func (s *InstitutionEndpoints) ListGroups(ctx context.Context, in *pb.ListGroupsRequest) (*pb.ListGroupsReply, error) {
	res, err := s.ListGroupsEndpoint(ctx, in)
	if err != nil {
		return nil, err
	}
	return res.(*pb.ListGroupsReply), nil
}

func (s *InstitutionEndpoints) ListInstitutions(ctx context.Context, in *pb.ListInstitutionsRequest) (*pb.ListInstitutionsReply, error) {
	res, err := s.ListInstitutionsEndpoint(ctx, in)
	if err != nil {
		return nil, err
	}
	return res.(*pb.ListInstitutionsReply), nil
}

func (s *InstitutionEndpoints) SaveInstitution(ctx context.Context, in *pb.SaveInstitutionRequest) (*pb.SaveInstitutionReply, error) {
	res, err := s.SaveInstitutionEndpoint(ctx, in)
	if err != nil {
		return nil, err
	}
	return res.(*pb.SaveInstitutionReply), nil
}

func NewInstitutionServiceGRPCClient(conn *grpc.ClientConn, tracer kitgrpc.ClientOption) *InstitutionEndpoints {
	endpoints := new(InstitutionEndpoints)
	options := make([]kitgrpc.ClientOption, 0)
	if tracer != nil {
		options = append(options, tracer)
	}
	{
		e := grpctransport.NewClient(
			conn,
			"pb.Institution",
			"TnxHisDownload",
			encodeRequest,
			decodeResponse,
			pb.InstitutionTnxHisDownloadResp{},
			options...,
		).Endpoint()
		endpoints.TnxHisDownloadEndpoint = e
	}

	{
		e := grpctransport.NewClient(
			conn,
			"pb.Institution",
			"GetTfrTrnLogs",
			encodeRequest,
			decodeResponse,
			pb.GetTfrTrnLogsResp{},
			options...,
		).Endpoint()
		endpoints.GetTfrTrnLogsEndpoint = e
	}

	{
		e := grpctransport.NewClient(
			conn,
			"pb.Institution",
			"GetTfrTrnLog",
			encodeRequest,
			decodeResponse,
			pb.GetTfrTrnLogResp{},
			options...,
		).Endpoint()
		endpoints.GetTfrTrnLogEndpoint = e
	}

	{
		e := grpctransport.NewClient(
			conn,
			"pb.Institution",
			"DownloadTfrTrnLogs",
			encodeRequest,
			decodeResponse,
			pb.DownloadTfrTrnLogsResp{},
			options...,
		).Endpoint()
		endpoints.DownloadTfrTrnLogsEndpoint = e
	}

	{
		e := grpctransport.NewClient(
			conn,
			"pb.Institution",
			"ListGroups",
			encodeRequest,
			decodeResponse,
			pb.ListGroupsReply{},
			options...,
		).Endpoint()
		endpoints.ListGroupsEndpoint = e
	}

	{
		e := grpctransport.NewClient(
			conn,
			"pb.Institution",
			"ListInstitutions",
			encodeRequest,
			decodeResponse,
			pb.ListInstitutionsReply{},
			options...,
		).Endpoint()
		endpoints.ListInstitutionsEndpoint = e
	}

	{
		e := grpctransport.NewClient(
			conn,
			"pb.Institution",
			"SaveInstitution",
			encodeRequest,
			decodeResponse,
			pb.SaveInstitutionReply{},
			options...,
		).Endpoint()
		endpoints.SaveInstitutionEndpoint = e
	}

	{
		e := grpctransport.NewClient(
			conn,
			"pb.Institution",
			"GetInstitutionById",
			encodeRequest,
			decodeResponse,
			pb.GetInstitutionByIdReply{},
			options...,
		).Endpoint()
		endpoints.GetInstitutionByIdEndpoint = e
	}

	{
		e := grpctransport.NewClient(
			conn,
			"pb.Institution",
			"SaveInstitutionFeeControlCash",
			encodeRequest,
			decodeResponse,
			pb.SaveInstitutionFeeControlCashReply{},
			options...,
		).Endpoint()
		endpoints.SaveInstitutionFeeControlCashEndpoint = e
	}

	{
		e := grpctransport.NewClient(
			conn,
			"pb.Institution",
			"GetInstitutionControl",
			encodeRequest,
			decodeResponse,
			pb.GetInstitutionControlReply{},
			options...,
		).Endpoint()
		endpoints.GetInstitutionControlEndpoint = e
	}

	{
		e := grpctransport.NewClient(
			conn,
			"pb.Institution",
			"GetInstitutionCash",
			encodeRequest,
			decodeResponse,
			pb.GetInstitutionCashReply{},
			options...,
		).Endpoint()
		endpoints.GetInstitutionCashEndpoint = e
	}

	{
		e := grpctransport.NewClient(
			conn,
			"pb.Institution",
			"GetInstitutionFee",
			encodeRequest,
			decodeResponse,
			pb.GetInstitutionFeeReply{},
			options...,
		).Endpoint()
		endpoints.GetInstitutionFeeEndpoint = e
	}

	{
		e := grpctransport.NewClient(
			conn,
			"pb.Institution",
			"SaveGroup",
			encodeRequest,
			decodeResponse,
			pb.SaveGroupReply{},
			options...,
		).Endpoint()
		endpoints.SaveGroupEndpoint = e
	}

	{
		e := grpctransport.NewClient(
			conn,
			"pb.Institution",
			"BindGroup",
			encodeRequest,
			decodeResponse,
			pb.BindGroupReply{},
			options...,
		).Endpoint()
		endpoints.BindGroupEndpoint = e
	}

	{
		e := grpctransport.NewClient(
			conn,
			"pb.Institution",
			"ListBindGroup",
			encodeRequest,
			decodeResponse,
			pb.ListBindGroupReply{},
			options...,
		).Endpoint()
		endpoints.ListBindGroupEndpoint = e
	}

	{
		e := grpctransport.NewClient(
			conn,
			"pb.Institution",
			"RemoveBindGroup",
			encodeRequest,
			decodeResponse,
			pb.RemoveBindGroupReply{},
			options...,
		).Endpoint()
		endpoints.RemoveBindGroupEndpoint = e
	}

	return endpoints
}
