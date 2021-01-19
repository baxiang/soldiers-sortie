package merchantservice

import (
	"context"
	"userService/pkg/kit"
	"userService/pkg/pb"

	grpctransport "github.com/go-kit/kit/transport/grpc"
)

type merchantServer struct {
	ListMerchantHandler              grpctransport.Handler
	ListGroupMerchantHandler         grpctransport.Handler
	SaveMerchantHandler              grpctransport.Handler
	SaveMerchantBankAccountHandler   grpctransport.Handler
	SaveMerchantBizDealHandler       grpctransport.Handler
	SaveMerchantBizFeeHandler        grpctransport.Handler
	SaveMerchantBusinessHandler      grpctransport.Handler
	SaveMerchantPictureHandler       grpctransport.Handler
	GetMerchantBankAccountHandler    grpctransport.Handler
	GetMerchantBizDealHandler        grpctransport.Handler
	GetMerchantBizFeeHandler         grpctransport.Handler
	GetMerchantBusinessHandler       grpctransport.Handler
	GetMerchantPictureHandler        grpctransport.Handler
	GetMerchantByIdHandler           grpctransport.Handler
	SaveMerchantBizDealAndFeeHandler grpctransport.Handler
	GenerateMchtCdHandler            grpctransport.Handler
	MerchantInfoQueryHandler         grpctransport.Handler
	MerchantForceChangeStatusHandler grpctransport.Handler
}

func (m *merchantServer) MerchantForceChangeStatus(ctx context.Context, in *pb.MerchantForceChangeStatusRequest) (*pb.MerchantForceChangeStatusReply, error) {
	_, res, err := m.MerchantForceChangeStatusHandler.ServeGRPC(ctx, in)
	if err != nil {
		return nil, err
	}
	return res.(*pb.MerchantForceChangeStatusReply), nil
}

func (m *merchantServer) MerchantInfoQuery(ctx context.Context, in *pb.MerchantInfoQueryRequest) (*pb.MerchantInfoQueryReply, error) {
	_, res, err := m.MerchantInfoQueryHandler.ServeGRPC(ctx, in)
	if err != nil {
		return nil, err
	}
	return res.(*pb.MerchantInfoQueryReply), nil
}

func (m *merchantServer) GenerateMchtCd(ctx context.Context, in *pb.GenerateMchtCdRequest) (*pb.GenerateMchtCdReply, error) {
	_, res, err := m.GenerateMchtCdHandler.ServeGRPC(ctx, in)
	if err != nil {
		return nil, err
	}
	return res.(*pb.GenerateMchtCdReply), nil
}

func (m *merchantServer) SaveMerchantBizDealAndFee(ctx context.Context, in *pb.SaveMerchantBizDealAndFeeRequest) (*pb.SaveMerchantBizDealAndFeeReply, error) {
	_, res, err := m.SaveMerchantBizDealAndFeeHandler.ServeGRPC(ctx, in)
	if err != nil {
		return nil, err
	}
	return res.(*pb.SaveMerchantBizDealAndFeeReply), nil
}

func (m *merchantServer) GetMerchantById(ctx context.Context, in *pb.GetMerchantByIdRequest) (*pb.GetMerchantByIdReply, error) {
	_, res, err := m.GetMerchantByIdHandler.ServeGRPC(ctx, in)
	if err != nil {
		return nil, err
	}
	return res.(*pb.GetMerchantByIdReply), nil
}

func (m *merchantServer) ListMerchant(ctx context.Context, in *pb.ListMerchantRequest) (*pb.ListMerchantReply, error) {
	_, res, err := m.ListMerchantHandler.ServeGRPC(ctx, in)
	if err != nil {
		return nil, err
	}
	reply, ok := res.(*pb.ListMerchantReply)
	if !ok {
		return nil, kit.ErrReplyTypeInvalid
	}
	return reply, nil
}
func (m *merchantServer) ListGroupMerchant(ctx context.Context, in *pb.ListGroupMerchantRequest) (*pb.ListGroupMerchantReply, error) {
	_, res, err := m.ListGroupMerchantHandler.ServeGRPC(ctx, in)
	if err != nil {
		return nil, err
	}
	reply, ok := res.(*pb.ListGroupMerchantReply)
	if !ok {
		return nil, kit.ErrReplyTypeInvalid
	}
	return reply, nil
}
func (m *merchantServer) SaveMerchant(ctx context.Context, in *pb.SaveMerchantRequest) (*pb.SaveMerchantReply, error) {
	_, res, err := m.SaveMerchantHandler.ServeGRPC(ctx, in)
	if err != nil {
		return nil, err
	}
	return res.(*pb.SaveMerchantReply), nil
}
func (m *merchantServer) SaveMerchantBankAccount(ctx context.Context, in *pb.SaveMerchantBankAccountRequest) (*pb.SaveMerchantBankAccountReply, error) {
	_, res, err := m.SaveMerchantBankAccountHandler.ServeGRPC(ctx, in)
	if err != nil {
		return nil, err
	}
	return res.(*pb.SaveMerchantBankAccountReply), nil
}

func (m *merchantServer) SaveMerchantBusiness(ctx context.Context, in *pb.SaveMerchantBusinessRequest) (*pb.SaveMerchantBusinessReply, error) {
	_, res, err := m.SaveMerchantBusinessHandler.ServeGRPC(ctx, in)
	if err != nil {
		return nil, err
	}
	return res.(*pb.SaveMerchantBusinessReply), nil
}

func (m *merchantServer) SaveMerchantPicture(ctx context.Context, in *pb.SaveMerchantPictureRequest) (*pb.SaveMerchantPictureReply, error) {
	_, res, err := m.SaveMerchantPictureHandler.ServeGRPC(ctx, in)
	if err != nil {
		return nil, err
	}
	return res.(*pb.SaveMerchantPictureReply), nil
}

func (m *merchantServer) GetMerchantBankAccount(ctx context.Context, in *pb.GetMerchantBankAccountRequest) (*pb.GetMerchantBankAccountReply, error) {
	_, res, err := m.GetMerchantBankAccountHandler.ServeGRPC(ctx, in)
	if err != nil {
		return nil, err
	}
	return res.(*pb.GetMerchantBankAccountReply), nil
}

func (m *merchantServer) GetMerchantBizDeal(ctx context.Context, in *pb.GetMerchantBizDealRequest) (*pb.GetMerchantBizDealReply, error) {
	_, res, err := m.GetMerchantBizDealHandler.ServeGRPC(ctx, in)
	if err != nil {
		return nil, err
	}
	return res.(*pb.GetMerchantBizDealReply), nil
}

func (m *merchantServer) GetMerchantBizFee(ctx context.Context, in *pb.GetMerchantBizFeeRequest) (*pb.GetMerchantBizFeeReply, error) {
	_, res, err := m.GetMerchantBizFeeHandler.ServeGRPC(ctx, in)
	if err != nil {
		return nil, err
	}
	return res.(*pb.GetMerchantBizFeeReply), nil
}

func (m *merchantServer) GetMerchantBusiness(ctx context.Context, in *pb.GetMerchantBusinessRequest) (*pb.GetMerchantBusinessReply, error) {
	_, res, err := m.GetMerchantBusinessHandler.ServeGRPC(ctx, in)
	if err != nil {
		return nil, err
	}
	return res.(*pb.GetMerchantBusinessReply), nil
}

func (m *merchantServer) GetMerchantPicture(ctx context.Context, in *pb.GetMerchantPictureRequest) (*pb.GetMerchantPictureReply, error) {
	_, res, err := m.GetMerchantPictureHandler.ServeGRPC(ctx, in)
	if err != nil {
		return nil, err
	}
	return res.(*pb.GetMerchantPictureReply), nil
}

func New(tracer grpctransport.ServerOption) pb.MerchantServer {
	svr := &merchantServer{}
	service := &merchantService{}
	options := make([]grpctransport.ServerOption, 0)
	if tracer != nil {
		options = append(options, tracer)
	}

	{
		endpoint := MakeSaveMerchantBizDealAndFeeEndpoint(service)
		endpoint = kit.LogginMiddleware(endpoint)
		svr.SaveMerchantBizDealAndFeeHandler = grpctransport.NewServer(
			endpoint,
			kit.DecodeRequest,
			kit.EncodeResponse,
			options...,
		)
	}

	{
		endpoint := MakeListMerchantEndpoint(service)
		endpoint = kit.LogginMiddleware(endpoint)
		svr.ListMerchantHandler = grpctransport.NewServer(
			endpoint,
			kit.DecodeRequest,
			kit.EncodeResponse,
			options...,
		)
	}

	{
		endpoint := MakeListGroupMerchantEndpoint(service)
		endpoint = kit.LogginMiddleware(endpoint)
		svr.ListGroupMerchantHandler = grpctransport.NewServer(
			endpoint,
			kit.DecodeRequest,
			kit.EncodeResponse,
			options...,
		)
	}

	{
		endpoint := MakeSaveMerchantEndpoint(service)
		endpoint = kit.LogginMiddleware(endpoint)
		svr.SaveMerchantHandler = grpctransport.NewServer(
			endpoint,
			kit.DecodeRequest,
			kit.EncodeResponse,
			options...,
		)
	}

	{
		endpoint := MakeSaveMerchantBankAccountEndpoint(service)
		endpoint = kit.LogginMiddleware(endpoint)
		svr.SaveMerchantBankAccountHandler = grpctransport.NewServer(
			endpoint,
			kit.DecodeRequest,
			kit.EncodeResponse,
			options...,
		)
	}

	{
		endpoint := MakeSaveMerchantBusinessEndpoint(service)
		endpoint = kit.LogginMiddleware(endpoint)
		svr.SaveMerchantBusinessHandler = grpctransport.NewServer(
			endpoint,
			kit.DecodeRequest,
			kit.EncodeResponse,
			options...,
		)
	}

	{
		endpoint := MakeSaveMerchantPictureEndpoint(service)
		endpoint = kit.LogginMiddleware(endpoint)
		svr.SaveMerchantPictureHandler = grpctransport.NewServer(
			endpoint,
			kit.DecodeRequest,
			kit.EncodeResponse,
			options...,
		)
	}

	{
		endpoint := MakeGetMerchantBankAccountEndpoint(service)
		endpoint = kit.LogginMiddleware(endpoint)
		svr.GetMerchantBankAccountHandler = grpctransport.NewServer(
			endpoint,
			kit.DecodeRequest,
			kit.EncodeResponse,
			options...,
		)
	}

	{
		endpoint := MakeGetMerchantBizDealEndpoint(service)
		endpoint = kit.LogginMiddleware(endpoint)
		svr.GetMerchantBizDealHandler = grpctransport.NewServer(
			endpoint,
			kit.DecodeRequest,
			kit.EncodeResponse,
			options...,
		)
	}

	{
		endpoint := MakeGetMerchantBizFeeEndpoint(service)
		endpoint = kit.LogginMiddleware(endpoint)
		svr.GetMerchantBizFeeHandler = grpctransport.NewServer(
			endpoint,
			kit.DecodeRequest,
			kit.EncodeResponse,
			options...,
		)
	}

	{
		endpoint := MakeGetMerchantBusinessEndpoint(service)
		endpoint = kit.LogginMiddleware(endpoint)
		svr.GetMerchantBusinessHandler = grpctransport.NewServer(
			endpoint,
			kit.DecodeRequest,
			kit.EncodeResponse,
			options...,
		)
	}

	{
		endpoint := MakeGetMerchantPictureEndpoint(service)
		endpoint = kit.LogginMiddleware(endpoint)
		svr.GetMerchantPictureHandler = grpctransport.NewServer(
			endpoint,
			kit.DecodeRequest,
			kit.EncodeResponse,
			options...,
		)
	}

	{
		endpoint := MakeGetMerchantByIdEndpoint(service)
		endpoint = kit.LogginMiddleware(endpoint)
		svr.GetMerchantByIdHandler = grpctransport.NewServer(
			endpoint,
			kit.DecodeRequest,
			kit.EncodeResponse,
			options...,
		)
	}

	{
		endpoint := MakeGenerateMchtCdEndpoint(service)
		endpoint = kit.LogginMiddleware(endpoint)
		svr.GenerateMchtCdHandler = grpctransport.NewServer(
			endpoint,
			kit.DecodeRequest,
			kit.EncodeResponse,
			options...,
		)
	}

	{
		endpoint := MakeMerchantInfoQueryEndpoint(service)
		endpoint = kit.LogginMiddleware(endpoint)
		svr.MerchantInfoQueryHandler = grpctransport.NewServer(
			endpoint,
			kit.DecodeRequest,
			kit.EncodeResponse,
			options...,
		)
	}

	{
		endpoint := MakeMerchantForceChangeStatusEndpoint(service)
		endpoint = kit.LogginMiddleware(endpoint)
		svr.MerchantForceChangeStatusHandler = grpctransport.NewServer(
			endpoint,
			kit.DecodeRequest,
			kit.EncodeResponse,
			options...,
		)
	}

	return svr
}
