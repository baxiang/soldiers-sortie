package institutionservice

import (
	"context"
	"userService/pkg/pb"

	grpctransport "github.com/go-kit/kit/transport/grpc"
)

type server struct {
	TnxHisDownloadHandler                grpctransport.Handler
	GetTfrTrnLogsHandler                 grpctransport.Handler
	GetTfrTrnLogHandler                  grpctransport.Handler
	DownloadTfrTrnLogsHandler            grpctransport.Handler
	listGroupsHandler                    grpctransport.Handler
	SaveGroupHandler                     grpctransport.Handler
	listInstitutionsHandler              grpctransport.Handler
	addInstitutionHandler                grpctransport.Handler
	GetInstitutionByIdHandler            grpctransport.Handler
	SaveInstitutionFeeControlCashHandler grpctransport.Handler
	GetInstitutionControlHandler         grpctransport.Handler
	GetInstitutionCashHandler            grpctransport.Handler
	GetInstitutionFeeHandler             grpctransport.Handler
	BindGroupHandler                     grpctransport.Handler
	ListBindGroupHandler                 grpctransport.Handler
	RemoveBindGroupHandler               grpctransport.Handler
}

func (s *server) RemoveBindGroup(ctx context.Context, in *pb.RemoveBindGroupRequest) (*pb.RemoveBindGroupReply, error) {
	_, res, err := s.RemoveBindGroupHandler.ServeGRPC(ctx, in)
	if err != nil {
		return nil, err
	}
	return res.(*pb.RemoveBindGroupReply), nil
}

func (s *server) ListBindGroup(ctx context.Context, in *pb.ListBindGroupRequest) (*pb.ListBindGroupReply, error) {
	_, res, err := s.ListBindGroupHandler.ServeGRPC(ctx, in)
	if err != nil {
		return nil, err
	}
	return res.(*pb.ListBindGroupReply), nil
}

func (s *server) BindGroup(ctx context.Context, in *pb.BindGroupRequest) (*pb.BindGroupReply, error) {
	_, res, err := s.BindGroupHandler.ServeGRPC(ctx, in)
	if err != nil {
		return nil, err
	}
	return res.(*pb.BindGroupReply), nil
}

func (s *server) SaveGroup(ctx context.Context, in *pb.SaveGroupRequest) (*pb.SaveGroupReply, error) {
	_, res, err := s.SaveGroupHandler.ServeGRPC(ctx, in)
	if err != nil {
		return nil, err
	}
	return res.(*pb.SaveGroupReply), nil
}

func (s *server) GetInstitutionControl(ctx context.Context, in *pb.GetInstitutionControlRequest) (*pb.GetInstitutionControlReply, error) {
	_, res, err := s.GetInstitutionControlHandler.ServeGRPC(ctx, in)
	if err != nil {
		return nil, err
	}
	return res.(*pb.GetInstitutionControlReply), nil
}

func (s *server) GetInstitutionCash(ctx context.Context, in *pb.GetInstitutionCashRequest) (*pb.GetInstitutionCashReply, error) {
	_, res, err := s.GetInstitutionCashHandler.ServeGRPC(ctx, in)
	if err != nil {
		return nil, err
	}
	return res.(*pb.GetInstitutionCashReply), nil
}

func (s *server) GetInstitutionFee(ctx context.Context, in *pb.GetInstitutionFeeRequest) (*pb.GetInstitutionFeeReply, error) {
	_, res, err := s.GetInstitutionFeeHandler.ServeGRPC(ctx, in)
	if err != nil {
		return nil, err
	}
	return res.(*pb.GetInstitutionFeeReply), nil
}

func (s *server) SaveInstitutionFeeControlCash(ctx context.Context, in *pb.SaveInstitutionFeeControlCashRequest) (*pb.SaveInstitutionFeeControlCashReply, error) {
	_, res, err := s.SaveInstitutionFeeControlCashHandler.ServeGRPC(ctx, in)
	if err != nil {
		return nil, err
	}
	return res.(*pb.SaveInstitutionFeeControlCashReply), nil
}

func (s *server) GetInstitutionById(ctx context.Context, in *pb.GetInstitutionByIdRequest) (*pb.GetInstitutionByIdReply, error) {
	_, res, err := s.GetInstitutionByIdHandler.ServeGRPC(ctx, in)
	if err != nil {
		return nil, err
	}
	return res.(*pb.GetInstitutionByIdReply), nil
}

func (s *server) TnxHisDownload(ctx context.Context, in *pb.InstitutionTnxHisDownloadReq) (*pb.InstitutionTnxHisDownloadResp, error) {
	_, res, err := s.TnxHisDownloadHandler.ServeGRPC(ctx, in)
	if err != nil {
		return nil, err
	}
	return res.(*pb.InstitutionTnxHisDownloadResp), nil
}

func (s *server) GetTfrTrnLogs(ctx context.Context, in *pb.GetTfrTrnLogsReq) (*pb.GetTfrTrnLogsResp, error) {
	_, res, err := s.GetTfrTrnLogsHandler.ServeGRPC(ctx, in)
	if err != nil {
		return nil, err
	}
	return res.(*pb.GetTfrTrnLogsResp), nil
}

func (s *server) GetTfrTrnLog(ctx context.Context, in *pb.GetTfrTrnLogReq) (*pb.GetTfrTrnLogResp, error) {
	_, res, err := s.GetTfrTrnLogHandler.ServeGRPC(ctx, in)
	if err != nil {
		return nil, err
	}
	return res.(*pb.GetTfrTrnLogResp), nil
}

func (s *server) DownloadTfrTrnLogs(ctx context.Context, in *pb.DownloadTfrTrnLogsReq) (*pb.DownloadTfrTrnLogsResp, error) {
	_, res, err := s.DownloadTfrTrnLogsHandler.ServeGRPC(ctx, in)
	if err != nil {
		return nil, err
	}
	return res.(*pb.DownloadTfrTrnLogsResp), nil
}

func (s *server) ListGroups(ctx context.Context, in *pb.ListGroupsRequest) (*pb.ListGroupsReply, error) {
	_, res, err := s.listGroupsHandler.ServeGRPC(ctx, in)
	if err != nil {
		return nil, err
	}
	reply, ok := res.(*pb.ListGroupsReply)
	if !ok {
		return nil, ErrReplyTypeInvalid
	}
	return reply, nil
}

func (s *server) ListInstitutions(ctx context.Context, in *pb.ListInstitutionsRequest) (*pb.ListInstitutionsReply, error) {
	_, res, err := s.listInstitutionsHandler.ServeGRPC(ctx, in)
	if err != nil {
		return nil, err
	}
	reply, ok := res.(*pb.ListInstitutionsReply)
	if !ok {
		return nil, ErrReplyTypeInvalid
	}
	return reply, nil
}

func (s *server) SaveInstitution(ctx context.Context, in *pb.SaveInstitutionRequest) (*pb.SaveInstitutionReply, error) {
	_, res, err := s.addInstitutionHandler.ServeGRPC(ctx, in)
	if err != nil {
		return nil, err
	}
	reply, ok := res.(*pb.SaveInstitutionReply)
	if !ok {
		return nil, ErrReplyTypeInvalid
	}
	return reply, nil

}

func New(tracer grpctransport.ServerOption) pb.InstitutionServer {
	svc := new(service)
	svr := new(server)
	options := make([]grpctransport.ServerOption, 0)
	if tracer != nil {
		options = append(options, tracer)
	}

	{
		e := MakeGetInstitutionControlEndpoint(svc)
		svr.GetInstitutionControlHandler = grpctransport.NewServer(
			e,
			grpcDecode,
			grpcEncode,
			options...,
		)
	}

	{
		e := MakeGetInstitutionCashEndpoint(svc)
		svr.GetInstitutionCashHandler = grpctransport.NewServer(
			e,
			grpcDecode,
			grpcEncode,
			options...,
		)
	}

	{
		e := MakeGetInstitutionFeeEndpoint(svc)
		svr.GetInstitutionFeeHandler = grpctransport.NewServer(
			e,
			grpcDecode,
			grpcEncode,
			options...,
		)
	}

	{
		e := MakeGetInstitutionByIdEndpoint(svc)
		svr.GetInstitutionByIdHandler = grpctransport.NewServer(
			e,
			grpcDecode,
			grpcEncode,
			options...,
		)
	}

	{
		e := MakeTnxHisDownloadEndpoint(svc)
		svr.TnxHisDownloadHandler = grpctransport.NewServer(
			e,
			grpcDecode,
			grpcEncode,
			options...,
		)
	}
	{
		e := MakeGetTfrTrnLogsEndpoint(svc)
		svr.GetTfrTrnLogsHandler = grpctransport.NewServer(
			e,
			grpcDecode,
			grpcEncode,
			options...,
		)
	}
	{
		e := MakeGetTfrTrnLogEndpoint(svc)
		svr.GetTfrTrnLogHandler = grpctransport.NewServer(
			e,
			grpcDecode,
			grpcEncode,
			options...,
		)
	}
	{
		e := MakeDownloadTfrTrnLogsEndpoint(svc)
		svr.DownloadTfrTrnLogsHandler = grpctransport.NewServer(
			e,
			grpcDecode,
			grpcEncode,
			options...,
		)
	}

	{
		e := MakeListGroupsEndpoint(svc)
		svr.listGroupsHandler = grpctransport.NewServer(
			e,
			grpcDecode,
			grpcEncode,
			options...,
		)
	}

	{
		e := MakeListInstitutionsEndpoint(svc)
		svr.listInstitutionsHandler = grpctransport.NewServer(
			e,
			grpcDecode,
			grpcEncode,
			options...,
		)
	}

	{
		e := MakeSaveInstitutionEndpoint(svc)
		svr.addInstitutionHandler = grpctransport.NewServer(
			e,
			grpcDecode,
			grpcEncode,
			options...,
		)
	}

	{
		e := MakeSaveInstitutionFeeControlCashEndpoint(svc)
		svr.SaveInstitutionFeeControlCashHandler = grpctransport.NewServer(
			e,
			grpcDecode,
			grpcEncode,
			options...,
		)
	}

	{
		e := MakeSaveGroupEndpoint(svc)
		svr.SaveGroupHandler = grpctransport.NewServer(
			e,
			grpcDecode,
			grpcEncode,
			options...,
		)
	}

	{
		e := MakeBindGroupEndpoint(svc)
		svr.BindGroupHandler = grpctransport.NewServer(
			e,
			grpcDecode,
			grpcEncode,
			options...,
		)
	}

	{
		e := MakeListBindGroupEndpoint(svc)
		svr.ListBindGroupHandler = grpctransport.NewServer(
			e,
			grpcDecode,
			grpcEncode,
			options...,
		)
	}

	{
		e := MakeRemoveBindGroupEndpoint(svc)
		svr.RemoveBindGroupHandler = grpctransport.NewServer(
			e,
			grpcDecode,
			grpcEncode,
			options...,
		)
	}

	return svr
}
