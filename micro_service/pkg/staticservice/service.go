package staticservice

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"userService/pkg/common"
	productmodel "userService/pkg/model/product"
	"userService/pkg/model/static"
	"userService/pkg/pb"
	"userService/pkg/util"

	"github.com/hashicorp/consul/api"
)

var (
	DictionaryConsulKey          = "static/dictionaryItem"
	InsProdBizFeeMapInfConsulKey = "static/insProdBizFeeMapInf"
	ProdBizTransMapConsulKey     = "static/prodBizTransMap"
	InsInfConsulKey              = "static/insInf"
)

type StaticMapData struct {
	dicItem             []*static.DictionaryItem
	insProdBizFeeMapInf []*static.InsProdBizFeeMapInf
	prodBizTransMap     []*static.ProdBizTransMap
	insInf              []*static.InsInf
}

var MyMap = StaticMapData{
	dicItem:             make([]*static.DictionaryItem, 0),
	insProdBizFeeMapInf: make([]*static.InsProdBizFeeMapInf, 0),
	prodBizTransMap:     make([]*static.ProdBizTransMap, 0),
	insInf:              make([]*static.InsInf, 0),
}

//service .
type service struct{}

func (s *service) SaveInsProdBizFeeMapInfo(ctx context.Context, in *pb.SaveInsProdBizFeeMapInfoRequest) (*pb.SaveInsProdBizFeeMapInfoReply, error) {
	reply := new(pb.SaveInsProdBizFeeMapInfoReply)
	if in.Item == nil {
		reply.Err = &pb.Error{
			Code:        http.StatusBadRequest,
			Message:     InvalidParam,
			Description: "item 不能为空",
		}
		return reply, nil
	}
	db := common.DB
	var err error
	if in.InsFeeBizCd != "" {
		err = static.DeleteInsProdBizFeeMapInf(db, in.InsFeeBizCd)
		if err != nil {
			return nil, err
		}
	}

	data := new(static.InsProdBizFeeMapInf)
	data.ProdCd = in.Item.ProdCd
	data.BizCd = in.Item.BizCd
	data.MccMTp = in.Item.MccMTp
	data.MccSTp = in.Item.MccSTp
	data.InsFeeBizCd = in.Item.InsFeeBizCd
	data.InsFeeBizNm = in.Item.InsFeeBizNm
	data.MsgResvFld1 = in.Item.MsgResvFld1
	data.MsgResvFld2 = in.Item.MsgResvFld2
	data.MsgResvFld3 = in.Item.MsgResvFld3
	data.MsgResvFld4 = in.Item.MsgResvFld4
	data.MsgResvFld5 = in.Item.MsgResvFld5
	data.MsgResvFld6 = in.Item.MsgResvFld6
	data.MsgResvFld7 = in.Item.MsgResvFld7
	data.MsgResvFld8 = in.Item.MsgResvFld8
	data.MsgResvFld9 = in.Item.MsgResvFld9
	data.MsgResvFld10 = in.Item.MsgResvFld10
	data.RecOprID = in.Item.RecOprID
	data.RecUpdOpr = in.Item.RecUpdOpr
	err = static.SaveInsProdBizFeeMapInf(db, data)

	return reply, err
}

func (s *service) ListOrgDictionaryItem(ctx context.Context, in *pb.ListOrgDictionaryItemRequest) (*pb.ListOrgDictionaryItemReply, error) {
	reply := new(pb.ListOrgDictionaryItemReply)
	if in.Page == 0 {
		in.Page = 1
	}
	if in.Size == 0 {
		in.Size = 10
	}

	db := common.DB

	query := new(static.OrgDictionaryItem)
	if in.Item != nil {
		query.Id = in.Item.Id
		query.TypeCode = in.Item.TypeCode
		query.OrgCode = in.Item.OrgCode
		query.ItemCode = in.Item.ItemCode
		query.TypeParm1 = in.Item.TypeParm1
		query.TypeParm2 = in.Item.TypeParm2
		query.TypeParm3 = in.Item.TypeParm3
		query.Remarks = in.Item.Remarks
		query.MsgResvFld1 = in.Item.MsgResvFld1
		query.MsgResvFld2 = in.Item.MsgResvFld2
		query.MsgResvFld3 = in.Item.MsgResvFld3
		query.MsgResvFld4 = in.Item.MsgResvFld4
		query.MsgResvFld5 = in.Item.MsgResvFld5
	}

	out, count, err := static.ListOrgDictionaryItem(db, query, in.Page, in.Size)
	if err != nil {
		return nil, err
	}
	items := make([]*pb.OrgDictionaryItemField, 0, len(out))
	for _, o := range out {
		items = append(items, &pb.OrgDictionaryItemField{
			Id:          o.Id,
			TypeCode:    o.TypeCode,
			OrgCode:     o.OrgCode,
			ItemCode:    o.ItemCode,
			TypeParm1:   o.TypeParm1,
			TypeParm2:   o.TypeParm2,
			TypeParm3:   o.TypeParm3,
			Remarks:     o.Remarks,
			MsgResvFld1: o.MsgResvFld1,
			MsgResvFld2: o.MsgResvFld2,
			MsgResvFld3: o.MsgResvFld3,
			MsgResvFld4: o.MsgResvFld4,
			MsgResvFld5: o.MsgResvFld5,
		})
	}
	reply.Size = in.Size
	reply.Page = in.Page
	reply.Count = count
	reply.Items = items

	return reply, nil
}

func (s *service) SaveOrgDictionaryItem(ctx context.Context, in *pb.SaveOrgDictionaryItemRequest) (*pb.SaveOrgDictionaryItemReply, error) {
	reply := new(pb.SaveOrgDictionaryItemReply)
	var err error
	db := common.DB.Begin()
	defer db.Rollback()

	if in.DictionaryItem != nil {
		item := new(static.DictionaryItem)
		item.DicType = in.DictionaryItem.DicType
		item.DicCode = in.DictionaryItem.DicCode
		item.DicName = in.DictionaryItem.DicName
		item.DispOrder = in.DictionaryItem.DispOrder
		item.Memo = in.DictionaryItem.Memo
		err = static.SaveDictionaryItem(db, item)
		if err != nil {
			return nil, err
		}
	}

	if in.OrgDictionary != nil {
		if in.OrgDictionary.OrgCode == "" || in.OrgDictionary.TypeCode == "" {
			reply.Err = &pb.Error{
				Code:        http.StatusBadRequest,
				Message:     InvalidParam,
				Description: "orgCode, typeCode 不能为空",
			}
			return reply, nil
		}

		err = static.DeleteOrgDictionaryItem(
			db,
			&static.OrgDictionaryItem{
				OrgCode:  in.OrgDictionary.OrgCode,
				TypeCode: in.OrgDictionary.TypeCode,
			},
		)
		if err != nil {
			return nil, err
		}

		for _, item := range in.OrgDictionary.Items {

			data := new(static.OrgDictionaryItem)
			data.Id = item.Id
			data.TypeCode = item.TypeCode
			data.OrgCode = item.OrgCode
			data.ItemCode = item.ItemCode
			data.TypeParm1 = item.TypeParm1
			data.TypeParm2 = item.TypeParm2
			data.TypeParm3 = item.TypeParm3
			data.Remarks = item.Remarks
			data.MsgResvFld1 = item.MsgResvFld1
			data.MsgResvFld2 = item.MsgResvFld2
			data.MsgResvFld3 = item.MsgResvFld3
			data.MsgResvFld4 = item.MsgResvFld4
			data.MsgResvFld5 = item.MsgResvFld5
			err = static.SaveOrgDictionaryItem(db, data)
			if err != nil {
				return nil, err
			}
		}
	}

	db.Commit()
	return reply, nil
}

func (s *service) FindMerchantFirstThreeCode(ctx context.Context, in *pb.FindMerchantFirstThreeCodeRequest) (*pb.FindMerchantFirstThreeCodeReply, error) {
	reply := new(pb.FindMerchantFirstThreeCodeReply)
	if in.Code == "" {
		reply.Err = &pb.Error{
			Code:        http.StatusBadRequest,
			Message:     InvalidParam,
			Description: "code不能为空",
		}
		return reply, nil
	}
	db := common.DB
	out, err := static.FindMerchantFirstThree(db, in.Code)
	if err != nil {
		return nil, err
	}
	items := make([]*pb.MerchantFirstThreeField, 0, len(out))
	for _, o := range out {
		items = append(items, &pb.MerchantFirstThreeField{
			DicCode: o.DicCode,
			DicName: o.DicName,
		})
	}
	reply.Items = items
	return reply, nil

}

func (s *service) FindArea(ctx context.Context, in *pb.FindAreaRequest) (*pb.FindAreaReply, error) {
	reply := new(pb.FindAreaReply)
	if in.Code == "" {
		reply.Err = &pb.Error{
			Code:        http.StatusBadRequest,
			Message:     InvalidParam,
			Description: "code不能为空",
		}
		return reply, nil
	}
	level, _ := strconv.Atoi(in.Level)
	db := common.DB
	if level > 1 {
		// 查询城市
		areas, err := static.FindCity(db, in.Code, in.Level)
		if err != nil {
			return nil, err
		}
		items := make([]*pb.Area, 0, len(areas))
		for _, area := range areas {
			items = append(items, &pb.Area{
				Name:    area.Name,
				DicCode: area.DicCode,
			})
		}
		reply.Items = items
	} else {
		// 查询省
		areas, err := static.FindProvince(db, in.Code)
		if err != nil {
			return nil, err
		}
		items := make([]*pb.Area, 0, len(areas))
		for _, area := range areas {
			items = append(items, &pb.Area{
				Name:    area.Name,
				DicCode: area.DicCode,
			})
		}
		reply.Items = items
	}
	return reply, nil
}

func (s *service) ListFeeMap(ctx context.Context, in *pb.ListFeeMapRequest) (*pb.ListFeeMapReply, error) {
	reply := new(pb.ListFeeMapReply)
	query := new(productmodel.Fee)
	if in.Item != nil {
		query.ProdCd = in.Item.ProdCd
		query.BizCd = in.Item.BizCd
		query.SubBizCd = in.Item.SubBizCd
		query.UpdateDate = in.Item.UpdateDate
		query.Description = in.Item.Description
		query.ResvFld1 = in.Item.ResvFld1
		query.ResvFld2 = in.Item.ResvFld2
		query.ResvFld3 = in.Item.ResvFld3
	}
	db := common.DB

	items, err := productmodel.ListFee(db, query)
	if err != nil {
		return nil, err
	}

	pbItems := make([]*pb.ProductBizFeeMapField, len(items))

	for i := range items {
		pbItems[i] = &pb.ProductBizFeeMapField{
			ProdCd:      items[i].ProdCd,
			BizCd:       items[i].BizCd,
			SubBizCd:    items[i].SubBizCd,
			UpdateDate:  items[i].UpdateDate,
			Description: items[i].Description,
			ResvFld1:    items[i].ResvFld1,
			ResvFld2:    items[i].ResvFld2,
			ResvFld3:    items[i].ResvFld3,
		}
	}

	reply.Items = pbItems
	return reply, nil
}

func (s *service) ListTransMap(ctx context.Context, in *pb.ListTransMapRequest) (*pb.ListTransMapReply, error) {
	reply := new(pb.ListTransMapReply)
	query := new(productmodel.Trans)
	if in.Item != nil {
		query.ProdCd = in.Item.ProdCd
		query.BizCd = in.Item.BizCd
		query.TransCd = in.Item.TransCd
		query.UpdateDate = in.Item.UpdateDate
		query.Description = in.Item.Description
		query.ResvFld1 = in.Item.ResvFld1
		query.ResvFld2 = in.Item.ResvFld2
		query.ResvFld3 = in.Item.ResvFld3
	}
	db := common.DB

	items, err := productmodel.ListTrans(db, query)
	if err != nil {
		return nil, err
	}

	pbItems := make([]*pb.ProductBizTransMapField, len(items))

	for i := range items {
		pbItems[i] = &pb.ProductBizTransMapField{
			ProdCd:      items[i].ProdCd,
			BizCd:       items[i].BizCd,
			TransCd:     items[i].TransCd,
			UpdateDate:  items[i].UpdateDate,
			Description: items[i].Description,
			ResvFld1:    items[i].ResvFld1,
			ResvFld2:    items[i].ResvFld2,
			ResvFld3:    items[i].ResvFld3,
		}
	}

	reply.Items = pbItems
	return reply, nil
}

func (s *service) GetInsProdBizFeeMapInfo(ctx context.Context, in *pb.GetInsProdBizFeeMapInfoRequest) (*pb.GetInsProdBizFeeMapInfoReply, error) {
	reply := new(pb.GetInsProdBizFeeMapInfoReply)
	query := new(static.InsProdBizFeeMapInf)
	if in.Item != nil {
		query.ProdCd = in.Item.ProdCd
		query.BizCd = in.Item.BizCd
		query.MccMTp = in.Item.MccMTp
		query.MccSTp = in.Item.MccSTp
		query.InsFeeBizCd = in.Item.InsFeeBizCd
		query.InsFeeBizNm = in.Item.InsFeeBizNm
		query.MsgResvFld1 = in.Item.MsgResvFld1
		query.MsgResvFld2 = in.Item.MsgResvFld2
		query.MsgResvFld3 = in.Item.MsgResvFld3
		query.MsgResvFld4 = in.Item.MsgResvFld4
		query.MsgResvFld5 = in.Item.MsgResvFld5
		query.MsgResvFld6 = in.Item.MsgResvFld6
		query.MsgResvFld7 = in.Item.MsgResvFld7
		query.MsgResvFld8 = in.Item.MsgResvFld8
		query.MsgResvFld9 = in.Item.MsgResvFld9
		query.MsgResvFld10 = in.Item.MsgResvFld10
		query.RecOprID = in.Item.RecOprID
		query.RecUpdOpr = in.Item.RecUpdOpr
	}
	infos, err := static.GetInsProdBizFeeMapInf(common.DB, query)
	if err != nil {
		return nil, err
	}
	items := make([]*pb.InsProdBizFeeMapInfoField, len(infos))
	for i := range infos {
		items[i] = &pb.InsProdBizFeeMapInfoField{
			ProdCd:       infos[i].ProdCd,
			BizCd:        infos[i].BizCd,
			MccMTp:       infos[i].MccMTp,
			MccSTp:       infos[i].MccSTp,
			InsFeeBizCd:  infos[i].InsFeeBizCd,
			InsFeeBizNm:  infos[i].InsFeeBizNm,
			MsgResvFld1:  infos[i].MsgResvFld1,
			MsgResvFld2:  infos[i].MsgResvFld2,
			MsgResvFld3:  infos[i].MsgResvFld3,
			MsgResvFld4:  infos[i].MsgResvFld4,
			MsgResvFld5:  infos[i].MsgResvFld5,
			MsgResvFld6:  infos[i].MsgResvFld6,
			MsgResvFld7:  infos[i].MsgResvFld7,
			MsgResvFld8:  infos[i].MsgResvFld8,
			MsgResvFld9:  infos[i].MsgResvFld9,
			MsgResvFld10: infos[i].MsgResvFld10,
			RecOprID:     infos[i].RecOprID,
			RecUpdOpr:    infos[i].RecUpdOpr,
		}
		if !infos[i].CreatedAt.IsZero() {
			items[i].CreatedAt = infos[i].CreatedAt.Format(util.TimePattern)
		}
		if !infos[i].UpdatedAt.IsZero() {
			items[i].UpdatedAt = infos[i].UpdatedAt.Format(util.TimePattern)
		}
	}
	reply.Items = items
	return reply, nil
}

func (s *service) FindUnionPayMccList(ctx context.Context, in *pb.FindUnionPayMccListRequest) (*pb.FindUnionPayMccListReply, error) {
	reply := new(pb.FindUnionPayMccListReply)
	if in.Page == 0 {
		in.Page = 1
	}
	if in.Size == 0 {
		in.Size = 10
	}
	db := common.DB
	query := new(static.Mcc)
	if in.Item != nil {
		query.Id = in.Item.Id
		query.Code = in.Item.Code
		query.Name = in.Item.Name
		query.Category = in.Item.Category
		query.CategoryType = in.Item.CategoryType
		query.Industry = in.Item.Industry
		query.Status = in.Item.Status
	}

	items, count, err := static.QueryMcc(db, query, in.Page, in.Size)
	if err != nil {
		return nil, err
	}

	pbItems := make([]*pb.UnionPayMccField, len(items))
	for i := range items {
		pbItems[i] = &pb.UnionPayMccField{
			Id:           items[i].Id,
			Code:         items[i].Code,
			Name:         items[i].Name,
			Category:     items[i].Category,
			CategoryType: items[i].CategoryType,
			Industry:     items[i].Industry,
			Status:       items[i].Status,
		}
		if !items[i].UpdatedAt.IsZero() {
			pbItems[i].UpdatedAt = items[i].UpdatedAt.Format(util.TimePattern)
		}
	}
	reply.Size = in.Size
	reply.Page = in.Page
	reply.Items = pbItems
	reply.Count = count
	return reply, nil
}

func (s *service) GetUnionPayBankList(ctx context.Context, in *pb.GetUnionPayBankListRequest) (*pb.GetUnionPayBankListReply, error) {
	reply := new(pb.GetUnionPayBankListReply)
	if in.Item == nil {
		reply.Err = &pb.Error{
			Code:        http.StatusBadRequest,
			Message:     InvalidParam,
			Description: "item不能为空",
		}
		return reply, nil
	}
	if in.Page == 0 {
		in.Page = 1
	}
	if in.Size == 0 {
		in.Size = 10
	}
	db := common.DB

	items, count, err := static.FindBankList(db, &static.BankList{
		Code: in.Item.Code,
		Name: in.Item.Name,
	}, in.Page, in.Size)
	if err != nil {
		return nil, err
	}

	pbItems := make([]*pb.UnionPayBankListField, len(items))
	for i := range items {
		pbItems[i] = &pb.UnionPayBankListField{
			Id:   items[i].Id,
			Code: items[i].Code,
			Name: items[i].Name,
		}
		if !items[i].UpdatedAt.IsZero() {
			pbItems[i].UpdatedAt = items[i].UpdatedAt.Format(util.TimePattern)
		}
	}
	reply.Size = in.Size
	reply.Page = in.Page
	reply.Items = pbItems
	reply.Count = count
	return reply, nil

}

//Download .
func (s *service) SyncData(ctx context.Context, in *pb.StaticSyncDataReq) (*pb.StaticSyncDataResp, error) {
	insProdBizFeeMapInfs, _ := static.GetInsProdBizFeeMapInf(common.DB, nil)
	prodBizTransMaps := static.GetProdBizTransMap(common.DB)
	insInfs := static.GetInsInf(common.DB)

	val2, err := json.Marshal(insProdBizFeeMapInfs)
	if err != nil {
		return nil, err
	}
	kv2 := &api.KVPair{
		Key:   InsProdBizFeeMapInfConsulKey,
		Flags: 0,
		Value: val2,
	}

	val3, err := json.Marshal(prodBizTransMaps)
	if err != nil {
		return nil, err
	}
	kv3 := &api.KVPair{
		Key:   ProdBizTransMapConsulKey,
		Flags: 0,
		Value: val3,
	}

	val4, err := json.Marshal(insInfs)
	if err != nil {
		return nil, err
	}
	kv4 := &api.KVPair{
		Key:   InsInfConsulKey,
		Flags: 0,
		Value: val4,
	}

	common.ConsulClient.KV().Put(kv2, nil)
	common.ConsulClient.KV().Put(kv3, nil)
	common.ConsulClient.KV().Put(kv4, nil)

	return &pb.StaticSyncDataResp{Result: true}, nil
}

func (s *service) GetDictionaryItem(ctx context.Context, in *pb.StaticGetDictionaryItemReq) (*pb.StaticGetDictionaryItemResp, error) {

	db := common.DB
	query := new(static.DictionaryItem)
	if in.Item != nil {
		query.DicType = in.Item.DicType
		query.DicCode = in.Item.DicCode
		query.DicName = in.Item.DicName
		query.DispOrder = in.Item.DispOrder
		query.Memo = in.Item.Memo
	}
	items := make([]*pb.DictionaryItemField, 0)
	results, err := static.QueryDictionaryItem(db, query)
	if err != nil {
		return nil, err
	}
	for _, result := range results {
		items = append(items, &pb.DictionaryItemField{
			DicType:   result.DicType,
			DicCode:   result.DicCode,
			DicName:   result.DicName,
			DispOrder: result.DispOrder,
			Memo:      result.Memo,
		})
	}
	return &pb.StaticGetDictionaryItemResp{Items: items}, nil
}

func (s *service) GetDicByProdAndBiz(ctx context.Context, in *pb.StaticGetDicByProdAndBizReq) (*pb.StaticGetDicByProdAndBizResp, error) {
	if in.ProdCd == "" {
		return &pb.StaticGetDicByProdAndBizResp{Err: &pb.Error{
			Code:        http.StatusBadRequest,
			Message:     InvalidParam,
			Description: "必须传入ProdCd",
		}}, nil
	}
	bizCdArr := make([]string, 0)
	prodCdArr := make([]string, 0)
	dicTypeCondition := make([]string, 0)
	dicNameCondition := make([]string, 0)

	if in.BizCd == "" {
		prodCdArr = append(prodCdArr, in.ProdCd)
		bizCdArr = getBizCdByProdCd(prodCdArr)
		if len(bizCdArr) == 0 {
			return &pb.StaticGetDicByProdAndBizResp{}, nil
		}

		dicTypeCondition = append(dicTypeCondition, "BIZ_CD")
		results := getDicItemByCondition(dicTypeCondition, dicNameCondition, bizCdArr)
		if len(results) == 0 {
			return &pb.StaticGetDicByProdAndBizResp{}, nil
		}

		items := make([]*pb.StaticGetDictionaryItem, 0)
		for i := 0; i < len(results); i++ {
			items = append(items, &pb.StaticGetDictionaryItem{
				DicCode: results[i].DicCode,
				DicName: results[i].DicName,
				DicType: results[i].DicType,
			})
		}
		return &pb.StaticGetDicByProdAndBizResp{Items: items}, nil
	}

	transCds := getTransCdByProdAndBiz(in.ProdCd, in.BizCd)
	if len(transCds) == 0 {
		return &pb.StaticGetDicByProdAndBizResp{}, nil
	}

	dicTypeCondition = append(dicTypeCondition, "TRANS_CD")
	results := getDicItemByCondition(dicTypeCondition, dicNameCondition, transCds)

	items := make([]*pb.StaticGetDictionaryItem, 0)
	for i := 0; i < len(results); i++ {
		items = append(items, &pb.StaticGetDictionaryItem{
			DicCode: results[i].DicCode,
			DicName: results[i].DicName,
			DicType: results[i].DicType,
		})
	}
	return &pb.StaticGetDicByProdAndBizResp{Items: items}, nil

}

func (s *service) CheckValues(ctx context.Context, in *pb.StaticCheckValuesReq) (*pb.StaticCheckValuesResp, error) {
	if in.ProdCd == "" && in.BizCd == "" && in.TransCd == "" && in.InsCompanyCd == "" && in.FwdInsIdCd == "" {
		return &pb.StaticCheckValuesResp{Err: &pb.Error{
			Code:        http.StatusBadRequest,
			Message:     InvalidParam,
			Description: "必须传入参数",
		}}, nil
	}
	result := checkProdBizAndTrans(in.ProdCd, in.BizCd, in.TransCd) && checkCpyCdAndFwd(in.InsCompanyCd, in.FwdInsIdCd)

	return &pb.StaticCheckValuesResp{Result: result}, nil
}

func checkCpyCdAndFwd(insCompanyCd, fwdInsIdCd string) bool {
	if insCompanyCd == "" && fwdInsIdCd == "" {
		return true
	}

	dicCode := make([]string, 0)
	if insCompanyCd != "" {
		dicCode = append(dicCode, insCompanyCd)
	}
	if fwdInsIdCd != "" {
		dicCode = append(dicCode, fwdInsIdCd)
	}

	results := getDicItemByCondition([]string{}, []string{}, dicCode)
	if len(results) != len(dicCode) {
		return false
	}
	if insCompanyCd != "" && fwdInsIdCd != "" {
		cpyCdArr := getDicByInsCmpCd(insCompanyCd)
		return inArr(cpyCdArr, fwdInsIdCd)
	}
	return true
}

func checkProdBizAndTrans(prodCd, bizCd, transCd string) bool {
	if prodCd == "" && bizCd == "" && transCd == "" {
		return true
	}

	dicCode := make([]string, 0)
	if prodCd != "" {
		dicCode = append(dicCode, prodCd)
	}
	if bizCd != "" {
		dicCode = append(dicCode, bizCd)
	}
	if transCd != "" {
		dicCode = append(dicCode, transCd)
	}

	results := getDicItemByCondition([]string{}, []string{}, dicCode)
	if len(results) != len(dicCode) {
		return false
	}

	if bizCd == "" || prodCd == "" {
		return true
	}

	if transCd == "" {
		bizArr := getBizCdByProdCd([]string{prodCd})
		return inArr(bizArr, bizCd)
	}

	tranArr := getTransCdByProdAndBiz(prodCd, bizCd)
	return inArr(tranArr, transCd)
}

func inArr(arr []string, val string) bool {
	for _, v := range arr {
		if v == val {
			return true
		}
	}
	return false
}

func getDicItemByCondition(dicTypeArr, dicNameArr, dicCodeArr []string) []*static.DictionaryItem {
	returnDic := make([]*static.DictionaryItem, 0)
	for i := 0; i < len(MyMap.dicItem); i++ {
		if includes(dicTypeArr, MyMap.dicItem[i].DicType) &&
			includes(dicNameArr, MyMap.dicItem[i].DicName) &&
			includes(dicCodeArr, MyMap.dicItem[i].DicCode) {
			returnDic = append(returnDic, MyMap.dicItem[i])
		}
	}
	return returnDic
}

func getBizCdByProdCd(prodCd []string) []string {
	bizCdArr := make([]string, 0)
	for i := 0; i < len(MyMap.insProdBizFeeMapInf); i++ {
		if includes(prodCd, (*MyMap.insProdBizFeeMapInf[i]).ProdCd) {
			bizCdArr = append(bizCdArr, (*MyMap.insProdBizFeeMapInf[i]).BizCd)
		}
	}
	return bizCdArr
}

func getTransCdByProdAndBiz(prodCd, bizCd string) []string {
	transCdArr := make([]string, 0)
	for i := 0; i < len(MyMap.prodBizTransMap); i++ {
		if (*MyMap.prodBizTransMap[i]).ProdCd == prodCd &&
			(*MyMap.prodBizTransMap[i]).BizCd == bizCd {
			transCdArr = append(transCdArr, (*MyMap.prodBizTransMap[i]).TransCd)
		}
	}
	return transCdArr
}

func getDicByInsCmpCd(insCmpCd string) []string {
	insIds := make([]string, 0)
	for i := 0; i < len(MyMap.insInf); i++ {
		if (*MyMap.insInf[i]).InsCompanyCd == insCmpCd {
			insIds = append(insIds, (*MyMap.insInf[i]).InsIDCd)
		}
	}
	return insIds
}

func includes(arr []string, val string) bool {
	if len(arr) == 0 {
		return true
	}
	for _, v := range arr {
		if v == val {
			return true
		}
	}
	return false
}

func (s *service) GetDictionaryLayerItem(ctx context.Context, in *pb.GetDictionaryLayerItemReq) (*pb.GetDictionaryLayerItemResp, error) {
	out := make([]*pb.DictionaryLayerItem, 0)
	db := common.DB
	items := static.GetDictionaryLayerItem(db, &static.DictionaryLayerItem{
		DicType:  in.DicType,
		DicPCode: in.DicPCode,
		DicCode:  in.DicCode,
		DicLevel: in.DicLevel,
	})
	for i := range items {
		out = append(out, &pb.DictionaryLayerItem{
			DicType:   items[i].DicType,
			DicCode:   items[i].DicCode,
			DicPCode:  items[i].DicPCode,
			DicLevel:  items[i].DicLevel,
			DisPOrder: items[i].DisPOrder,
			Name:      items[i].Name,
			Memo:      items[i].Memo,
		})
	}
	return &pb.GetDictionaryLayerItemResp{
		Items: out,
	}, nil
}

func (s *service) GetDictionaryItemByPk(ctx context.Context, in *pb.GetDictionaryItemByPkReq) (*pb.GetDictionaryItemByPkResp, error) {
	reply := new(pb.GetDictionaryItemByPkResp)
	if in.DicCode == "" || in.DicType == "" {
		reply.Err = &pb.Error{
			Code:        http.StatusBadRequest,
			Message:     InvalidParam,
			Description: "dicCode 和 dicType不能为空",
		}
		return reply, nil
	}
	db := common.DB

	dic := static.GetDictionaryItemByPk(db, &static.DictionaryItem{
		DicType: in.DicType,
		DicCode: in.DicCode,
	})
	if dic != nil {
		reply.Item = &pb.StaticGetDictionaryItem{
			DicType: dic.DicType,
			DicCode: dic.DicCode,
			DicName: dic.DicName,
		}
	}

	return reply, nil
}
