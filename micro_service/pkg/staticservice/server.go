package staticservice

import (
	"context"
	"userService/pkg/kit"
	"userService/pkg/pb"

	grpctransport "github.com/go-kit/kit/transport/grpc"
)

type server struct {
	SyncDataHandler                   grpctransport.Handler
	GetDictionaryItemHandler          grpctransport.Handler
	getDicByProdAndBiz                grpctransport.Handler
	checkValues                       grpctransport.Handler
	getDictionaryLayerItem            grpctransport.Handler
	GetDictionaryItemByPkHandler      grpctransport.Handler
	GetUnionPayBankListHandler        grpctransport.Handler
	FindUnionPayMccListHandler        grpctransport.Handler
	GetInsProdBizFeeMapInfoHandler    grpctransport.Handler
	ListTransMapHandler               grpctransport.Handler
	ListFeeMapHandler                 grpctransport.Handler
	FindAreaHandler                   grpctransport.Handler
	FindMerchantFirstThreeCodeHandler grpctransport.Handler
	SaveOrgDictionaryItemHandler      grpctransport.Handler
	ListOrgDictionaryItemHandler      grpctransport.Handler
	SaveInsProdBizFeeMapInfoHandler   grpctransport.Handler
}

func (s *server) SaveInsProdBizFeeMapInfo(ctx context.Context, in *pb.SaveInsProdBizFeeMapInfoRequest) (*pb.SaveInsProdBizFeeMapInfoReply, error) {
	_, res, err := s.SaveInsProdBizFeeMapInfoHandler.ServeGRPC(ctx, in)
	if err != nil {
		return nil, err
	}
	return res.(*pb.SaveInsProdBizFeeMapInfoReply), nil
}

func (s *server) ListOrgDictionaryItem(ctx context.Context, in *pb.ListOrgDictionaryItemRequest) (*pb.ListOrgDictionaryItemReply, error) {
	_, res, err := s.ListOrgDictionaryItemHandler.ServeGRPC(ctx, in)
	if err != nil {
		return nil, err
	}
	return res.(*pb.ListOrgDictionaryItemReply), nil
}

func (s *server) SaveOrgDictionaryItem(ctx context.Context, in *pb.SaveOrgDictionaryItemRequest) (*pb.SaveOrgDictionaryItemReply, error) {
	_, res, err := s.SaveOrgDictionaryItemHandler.ServeGRPC(ctx, in)
	if err != nil {
		return nil, err
	}
	return res.(*pb.SaveOrgDictionaryItemReply), nil
}

func (s *server) FindMerchantFirstThreeCode(ctx context.Context, in *pb.FindMerchantFirstThreeCodeRequest) (*pb.FindMerchantFirstThreeCodeReply, error) {
	_, res, err := s.FindMerchantFirstThreeCodeHandler.ServeGRPC(ctx, in)
	if err != nil {
		return nil, err
	}
	return res.(*pb.FindMerchantFirstThreeCodeReply), nil
}

func (s *server) FindArea(ctx context.Context, in *pb.FindAreaRequest) (*pb.FindAreaReply, error) {
	_, res, err := s.FindAreaHandler.ServeGRPC(ctx, in)
	if err != nil {
		return nil, err
	}
	return res.(*pb.FindAreaReply), nil
}

func (s *server) ListFeeMap(ctx context.Context, in *pb.ListFeeMapRequest) (*pb.ListFeeMapReply, error) {
	_, res, err := s.ListFeeMapHandler.ServeGRPC(ctx, in)
	if err != nil {
		return nil, err
	}
	return res.(*pb.ListFeeMapReply), nil
}

func (s *server) ListTransMap(ctx context.Context, in *pb.ListTransMapRequest) (*pb.ListTransMapReply, error) {
	_, res, err := s.ListTransMapHandler.ServeGRPC(ctx, in)
	if err != nil {
		return nil, err
	}

	return res.(*pb.ListTransMapReply), nil
}

func (g *server) GetInsProdBizFeeMapInfo(ctx context.Context, in *pb.GetInsProdBizFeeMapInfoRequest) (*pb.GetInsProdBizFeeMapInfoReply, error) {
	_, res, err := g.GetInsProdBizFeeMapInfoHandler.ServeGRPC(ctx, in)
	if err != nil {
		return nil, err
	}
	return res.(*pb.GetInsProdBizFeeMapInfoReply), nil
}

func (g *server) FindUnionPayMccList(ctx context.Context, in *pb.FindUnionPayMccListRequest) (*pb.FindUnionPayMccListReply, error) {
	_, res, err := g.FindUnionPayMccListHandler.ServeGRPC(ctx, in)
	if err != nil {
		return nil, err
	}
	return res.(*pb.FindUnionPayMccListReply), nil
}

func (g *server) GetUnionPayBankList(ctx context.Context, in *pb.GetUnionPayBankListRequest) (*pb.GetUnionPayBankListReply, error) {
	_, res, err := g.GetUnionPayBankListHandler.ServeGRPC(ctx, in)
	if err != nil {
		return nil, err
	}
	return res.(*pb.GetUnionPayBankListReply), nil
}

func (g *server) GetDictionaryItemByPk(ctx context.Context, in *pb.GetDictionaryItemByPkReq) (*pb.GetDictionaryItemByPkResp, error) {
	_, res, err := g.GetDictionaryItemByPkHandler.ServeGRPC(ctx, in)
	if err != nil {
		return nil, err
	}
	return res.(*pb.GetDictionaryItemByPkResp), nil
}

func (g *server) GetDictionaryLayerItem(ctx context.Context, in *pb.GetDictionaryLayerItemReq) (*pb.GetDictionaryLayerItemResp, error) {
	_, res, err := g.getDictionaryLayerItem.ServeGRPC(ctx, in)
	if err != nil {
		return nil, err
	}
	return res.(*pb.GetDictionaryLayerItemResp), nil
}

func (g *server) SyncData(ctx context.Context, in *pb.StaticSyncDataReq) (*pb.StaticSyncDataResp, error) {
	_, res, err := g.SyncDataHandler.ServeGRPC(ctx, in)
	if err != nil {
		return nil, err
	}
	return res.(*pb.StaticSyncDataResp), nil
}

func (g *server) GetDictionaryItem(ctx context.Context, in *pb.StaticGetDictionaryItemReq) (*pb.StaticGetDictionaryItemResp, error) {
	_, res, err := g.GetDictionaryItemHandler.ServeGRPC(ctx, in)
	if err != nil {
		return nil, err
	}
	return res.(*pb.StaticGetDictionaryItemResp), nil
}

func (g *server) GetDicByProdAndBiz(ctx context.Context, in *pb.StaticGetDicByProdAndBizReq) (*pb.StaticGetDicByProdAndBizResp, error) {
	_, res, err := g.getDicByProdAndBiz.ServeGRPC(ctx, in)
	if err != nil {
		return nil, err
	}
	return res.(*pb.StaticGetDicByProdAndBizResp), nil
}

func (g *server) CheckValues(ctx context.Context, in *pb.StaticCheckValuesReq) (*pb.StaticCheckValuesResp, error) {
	_, res, err := g.checkValues.ServeGRPC(ctx, in)
	if err != nil {
		return nil, err
	}
	return res.(*pb.StaticCheckValuesResp), nil
}

//Newserver .
func New(tracer grpctransport.ServerOption) pb.StaticServer {
	svr := new(server)
	svc := new(service)
	options := make([]grpctransport.ServerOption, 0)
	if tracer != nil {
		options = append(options, tracer)
	}

	{
		e := MakeSyncDataEndpoint(svc)
		svr.SyncDataHandler = grpctransport.NewServer(
			e,
			grpcDecode,
			grpcEncode,
			options...,
		)
	}

	{
		e := MakeGetDictionaryItemEndpoint(svc)
		svr.GetDictionaryItemHandler = grpctransport.NewServer(
			e,
			grpcDecode,
			grpcEncode,
			options...,
		)
	}

	{
		e := MakeGetDicByProdAndBizEndpoint(svc)
		svr.getDicByProdAndBiz = grpctransport.NewServer(
			e,
			grpcDecode,
			grpcEncode,
			options...,
		)
	}

	{
		e := MakeCheckValuesEndpoint(svc)
		svr.checkValues = grpctransport.NewServer(
			e,
			grpcDecode,
			grpcEncode,
			options...,
		)
	}

	{
		e := MakeGetDictionaryLayerItemEndpoint(svc)
		svr.getDictionaryLayerItem = grpctransport.NewServer(
			e,
			grpcDecode,
			grpcEncode,
			options...,
		)
	}

	{
		e := MakeGetDictionaryItemByPkEndpoint(svc)
		svr.GetDictionaryItemByPkHandler = grpctransport.NewServer(
			e,
			grpcDecode,
			grpcEncode,
			options...,
		)
	}

	{
		e := MakeGetUnionPayBankListEndpoint(svc)
		svr.GetUnionPayBankListHandler = grpctransport.NewServer(
			e,
			grpcDecode,
			grpcEncode,
			options...,
		)
	}

	{
		e := MakeFindUnionPayMccListEndpoint(svc)
		svr.FindUnionPayMccListHandler = grpctransport.NewServer(
			e,
			grpcDecode,
			grpcEncode,
			options...,
		)
	}

	{
		e := MakeGetInsProdBizFeeMapInfoEndpoint(svc)
		svr.GetInsProdBizFeeMapInfoHandler = grpctransport.NewServer(
			e,
			grpcDecode,
			grpcEncode,
			options...,
		)
	}

	{
		endpoint := MakeListTransMapEndpoint(svc)
		endpoint = kit.LogginMiddleware(endpoint)
		svr.ListTransMapHandler = grpctransport.NewServer(
			endpoint,
			kit.DecodeRequest,
			kit.EncodeResponse,
			options...,
		)
	}

	{
		endpoint := MakeListFeeMapEndpoint(svc)
		endpoint = kit.LogginMiddleware(endpoint)
		svr.ListFeeMapHandler = grpctransport.NewServer(
			endpoint,
			kit.DecodeRequest,
			kit.EncodeResponse,
			options...,
		)
	}

	{
		endpoint := MakeFindAreaEndpoint(svc)
		endpoint = kit.LogginMiddleware(endpoint)
		svr.FindAreaHandler = grpctransport.NewServer(
			endpoint,
			kit.DecodeRequest,
			kit.EncodeResponse,
			options...,
		)
	}

	{
		endpoint := MakeFindMerchantFirstThreeCodeEndpoint(svc)
		endpoint = kit.LogginMiddleware(endpoint)
		svr.FindMerchantFirstThreeCodeHandler = grpctransport.NewServer(
			endpoint,
			kit.DecodeRequest,
			kit.EncodeResponse,
			options...,
		)
	}

	{
		endpoint := MakeSaveOrgDictionaryItemEndpoint(svc)
		endpoint = kit.LogginMiddleware(endpoint)
		svr.SaveOrgDictionaryItemHandler = grpctransport.NewServer(
			endpoint,
			kit.DecodeRequest,
			kit.EncodeResponse,
			options...,
		)
	}

	{
		endpoint := MakeListOrgDictionaryItemEndpoint(svc)
		endpoint = kit.LogginMiddleware(endpoint)
		svr.ListOrgDictionaryItemHandler = grpctransport.NewServer(
			endpoint,
			kit.DecodeRequest,
			kit.EncodeResponse,
			options...,
		)
	}

	{
		endpoint := MakeSaveInsProdBizFeeMapInfoEndpoint(svc)
		endpoint = kit.LogginMiddleware(endpoint)
		svr.SaveInsProdBizFeeMapInfoHandler = grpctransport.NewServer(
			endpoint,
			kit.DecodeRequest,
			kit.EncodeResponse,
			options...,
		)
	}

	return svr
}
