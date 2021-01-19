package merchantservice

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"
	"userService/pkg/common"
	"userService/pkg/pb"
	"userService/pkg/util"

	insmodel "userService/pkg/model/institution"
	merchantmodel "userService/pkg/model/merchant"
	usermodel "userService/pkg/model/user"

	"google.golang.org/grpc/metadata"
)

type merchantService struct{}

func (m *merchantService) MerchantForceChangeStatus(ctx context.Context, in *pb.MerchantForceChangeStatusRequest) (*pb.MerchantForceChangeStatusReply, error) {
	reply := new(pb.MerchantForceChangeStatusReply)
	if in.MchtCd == "" {
		reply.Err = &pb.Error{
			Code:        http.StatusBadRequest,
			Message:     InvalidParam,
			Description: "商户号不能为空",
		}
		return reply, nil
	}
	status := "01"
	systemFlag := "01"
	switch in.Operate {
	case "freeze":
		status = "13"
		systemFlag = "13"
	case "unfreeze":
		status = "01"
		systemFlag = "01"
	case "cancel":
		status = "00"
		systemFlag = "00"
	default:
		reply.Err = &pb.Error{
			Code:        http.StatusBadRequest,
			Message:     InvalidParam,
			Description: "商户号不能为空",
		}
		return reply, nil
	}
	db := common.DB.Begin()
	defer db.Rollback()

	err := merchantmodel.UpdateMerchant(
		db,
		&merchantmodel.MerchantInfo{MchtCd: in.MchtCd},
		&merchantmodel.MerchantInfo{Status: status, SystemFlag: systemFlag},
	)
	if err != nil {
		return nil, err
	}
	err = merchantmodel.UpdateMerchantMain(
		db,
		&merchantmodel.MerchantInfoMain{
			MerchantInfo: merchantmodel.MerchantInfo{
				MchtCd: in.MchtCd,
			}},
		&merchantmodel.MerchantInfoMain{
			MerchantInfo: merchantmodel.MerchantInfo{
				Status:     status,
				SystemFlag: systemFlag,
			}},
	)
	if err != nil {
		return nil, err
	}

	db.Commit()

	return reply, nil
}

func (m *merchantService) MerchantInfoQuery(ctx context.Context, in *pb.MerchantInfoQueryRequest) (*pb.MerchantInfoQueryReply, error) {
	reply := new(pb.MerchantInfoQueryReply)
	if in.Size == 0 {
		in.Size = 10
	}
	if in.Page == 0 {
		in.Page = 1
	}

	db := common.DB

	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		reply.Err = &pb.Error{
			Code:        http.StatusBadRequest,
			Message:     InvalidParam,
			Description: "找不到用户信息",
		}
		return reply, nil
	}
	ids := md.Get("userid")
	if !ok || len(ids) == 0 {
		reply.Err = &pb.Error{
			Code:        http.StatusBadRequest,
			Message:     InvalidParam,
			Description: "找不到用户信息",
		}
		return reply, nil
	}
	id, _ := strconv.ParseInt(ids[0], 10, 64)
	user, err := usermodel.FindUserByID(db, id)
	if err != nil {
		return nil, err
	}
	if user == nil {
		reply.Err = &pb.Error{
			Code:        http.StatusBadRequest,
			Message:     InvalidParam,
			Description: "找不到用户信息",
		}
		return reply, nil
	}

	insIds := make([]string, 0)
	if user.UserType != "admin" && user.UserType != "institution" && user.UserType != "institution_group" {
		reply.Items = make([]*pb.MerchantInfoField, 0)
		reply.Count = 0
		reply.Page = in.Page
		reply.Size = in.Size
		return reply, nil
	}
	if user.UserType == "institution_group" {
		groupId, _ := strconv.ParseInt(user.UserGroupNo, 10, 64)
		groups, err := insmodel.ListInsGroupBind(db, groupId)
		if err != nil {
			return nil, err
		}
		for _, group := range groups {
			insIds = append(insIds, group.InsIdCd)
		}
		if len(insIds) == 0 {
			reply.Items = make([]*pb.MerchantInfoField, 0)
			reply.Count = 0
			reply.Page = in.Page
			reply.Size = in.Size
			return reply, nil
		}
	}

	edit := true
	if in.Type == "main" {
		edit = false
	}
	if edit {
		query := new(merchantmodel.MerchantAccount)
		if in.Item != nil {
			query.MchtCd = in.Item.MchtCd
			if user.UserType == "institution" {
				query.AipBranCd = user.UserGroupNo
			}
			if in.Item.AipBranCd != "" {
				query.AipBranCd = in.Item.AipBranCd
			}
			query.GroupCd = in.Item.GroupCd
			query.Name = in.Item.Name
			query.NameBusi = in.Item.NameBusi
			query.Status = in.Item.Status
			query.SystemFlag = in.Item.SystemFlag
			query.BankBelongCd = in.Item.BankBelongCd
			query.GroupCd = in.Item.GroupCd
			query.AccountName = in.Item.AccountName
			query.Account = in.Item.Account

		}
		merchants, count, err := merchantmodel.QueryMerchantAccountInfos(db, query, insIds, in.Page, in.Size)
		if err != nil {
			return nil, err
		}

		pbMerchants := make([]*pb.MerchantInfoField, len(merchants))
		for i := range merchants {
			pbMerchants[i] = &pb.MerchantInfoField{
				MchtCd:       merchants[i].MchtCd,
				Name:         merchants[i].Name,
				AipBranCd:    merchants[i].AipBranCd,
				BankBelongCd: merchants[i].BankBelongCd,
				NameBusi:     merchants[i].NameBusi,
				GroupCd:      merchants[i].GroupCd,
				AccountName:  merchants[i].AccountName,
				Account:      merchants[i].Account,
				Status:       merchants[i].Status,
				SystemFlag:   merchants[i].SystemFlag,
			}
			if !merchants[i].UpdatedAt.IsZero() {
				pbMerchants[i].UpdatedAt = merchants[i].UpdatedAt.Format(util.TimePattern)
			}
		}
		reply.Items = pbMerchants
		reply.Count = count
		reply.Page = in.Page
		reply.Size = in.Size
		return reply, nil

	} else {
		query := new(merchantmodel.MerchantAccountMain)
		if in.Item != nil {
			query.MchtCd = in.Item.MchtCd
			if user.UserType == "institution" {
				query.AipBranCd = user.UserGroupNo
			}
			if in.Item.AipBranCd != "" {
				query.AipBranCd = in.Item.AipBranCd
			}
			query.GroupCd = in.Item.GroupCd
			query.Name = in.Item.Name
			query.NameBusi = in.Item.NameBusi
			query.Status = in.Item.Status
			query.SystemFlag = in.Item.SystemFlag
			query.BankBelongCd = in.Item.BankBelongCd
			query.GroupCd = in.Item.GroupCd
			query.AccountName = in.Item.AccountName
			query.Account = in.Item.Account
		}
		merchants, count, err := merchantmodel.QueryMerchantAccountInfosMain(db, query, insIds, in.Page, in.Size)
		if err != nil {
			return nil, err
		}

		pbMerchants := make([]*pb.MerchantInfoField, len(merchants))
		for i := range merchants {
			pbMerchants[i] = &pb.MerchantInfoField{
				MchtCd:       merchants[i].MchtCd,
				Name:         merchants[i].Name,
				AipBranCd:    merchants[i].AipBranCd,
				BankBelongCd: merchants[i].BankBelongCd,
				NameBusi:     merchants[i].NameBusi,
				GroupCd:      merchants[i].GroupCd,
				AccountName:  merchants[i].AccountName,
				Account:      merchants[i].Account,
				Status:       merchants[i].Status,
				SystemFlag:   merchants[i].SystemFlag,
			}
			if !merchants[i].UpdatedAt.IsZero() {
				pbMerchants[i].UpdatedAt = merchants[i].UpdatedAt.Format(util.TimePattern)
			}
		}
		reply.Items = pbMerchants
		reply.Count = count
		reply.Page = in.Page
		reply.Size = in.Size
		return reply, nil
	}
}

func (m *merchantService) GenerateMchtCd(ctx context.Context, in *pb.GenerateMchtCdRequest) (*pb.GenerateMchtCdReply, error) {
	reply := new(pb.GenerateMchtCdReply)
	if len(in.MccCd) != 4 || len(in.MchtCd3) != 3 || len(in.MchtCdPc) != 4 {
		reply.Err = &pb.Error{
			Code:        http.StatusBadRequest,
			Message:     "InvalidParamsError",
			Description: "参数错误",
		}
		return reply, nil
	}

	db := common.DB
	prefix := in.MchtCd3 + in.MchtCdPc + in.MccCd
	used, err := merchantmodel.FindMerchantCdByPrefix(db, prefix)
	if err != nil {
		return nil, err
	}

	id := prefix
	cds := make([]bool, 10000)
	for _, d := range used {
		if len(d.Id) == 15 {
			i, _ := strconv.ParseInt(d.Id[11:], 10, 64)
			if i >= 0 && i < 10000 {
				cds[i] = true
			}
		}
	}
	cd := 0
	sn := ""
	hasSkip := false
	for i := 0; i < 10000; i++ {
		if !cds[i] {
			hasSkip = true
			cd = i
			break
		}
	}
	if hasSkip {
		// 如果有跳跃，使用跳跃序号
		sn = fmt.Sprintf("%0.4d", cd)
		id += sn
	} else {
		// 如果没有跳跃
		reply.Err = &pb.Error{
			Code:        http.StatusBadRequest,
			Message:     "InvalidParamsError",
			Description: "没有可用的商户号",
		}
		return reply, nil
	}
	err = merchantmodel.SaveMerchantCd(db, &merchantmodel.UsedMerchantCd{Id: id})
	if err != nil {
		return nil, err
	}

	reply.Sn = sn
	return reply, nil
}

// 商户产品和费率保存
func (m *merchantService) SaveMerchantBizDealAndFee(ctx context.Context, in *pb.SaveMerchantBizDealAndFeeRequest) (*pb.SaveMerchantBizDealAndFeeReply, error) {
	reply := new(pb.SaveMerchantBizDealAndFeeReply)
	if in.Fees == nil || in.Deals == nil {
		reply.Err = &pb.Error{
			Code:        http.StatusBadRequest,
			Message:     "InvalidParamsError",
			Description: "保存信息为空",
		}
		return reply, nil
	}
	db := common.DB.Begin()
	defer db.Rollback()

	{
		for _, deal := range in.Deals {
			data := new(merchantmodel.BizDeal)
			data.MchtCd = deal.MchtCd
			data.ProdCd = deal.ProdCd
			data.BizCd = deal.BizCd
			data.TransCd = deal.TransCd
			data.OperIn = deal.OperIn
			data.RecOprId = deal.RecOprId
			data.RecUpdOpr = deal.RecUpdOpr
			err := merchantmodel.SaveBizDeal(db, data)
			if err != nil {
				return nil, err
			}
		}
	}
	{
		for _, fee := range in.Fees {
			data := new(merchantmodel.BizFee)
			data.MchtCd = fee.MchtCd
			data.ProdCd = fee.ProdCd
			data.BizCd = fee.BizCd
			data.SubBizCd = fee.SubBizCd
			data.MchtFeeMd = fee.MchtFeeMd
			data.MchtFeePercent, _ = strconv.ParseFloat(fee.MchtFeePercent, 64)
			data.MchtFeePctMin, _ = strconv.ParseFloat(fee.MchtFeePctMin, 64)
			data.MchtFeePctMax, _ = strconv.ParseFloat(fee.MchtFeePctMax, 64)
			data.MchtFeeSingle, _ = strconv.ParseFloat(fee.MchtFeeSingle, 64)
			data.MchtAFeeSame = fee.MchtAFeeSame
			data.MchtAFeeMd = fee.MchtAFeeMd
			data.MchtAFeePercent, _ = strconv.ParseFloat(fee.MchtAFeePercent, 64)
			data.MchtAFeePctMin, _ = strconv.ParseFloat(fee.MchtAFeePctMin, 64)
			data.MchtAFeePctMax, _ = strconv.ParseFloat(fee.MchtAFeePctMax, 64)
			data.MchtAFeeSingle, _ = strconv.ParseFloat(fee.MchtAFeeSingle, 64)
			data.OperIn = fee.OperIn
			data.RecOprId = fee.RecOprId
			data.RecUpdOpr = fee.RecUpdOpr
			err := merchantmodel.SaveBizFee(db, data)
			if err != nil {
				return nil, err
			}
		}
	}
	db.Commit()
	return reply, nil
}

func (m *merchantService) GetMerchantById(ctx context.Context, in *pb.GetMerchantByIdRequest) (*pb.GetMerchantByIdReply, error) {
	reply := new(pb.GetMerchantByIdReply)
	if in.Id == "" {
		reply.Err = &pb.Error{
			Code:        http.StatusBadRequest,
			Message:     InvalidParam,
			Description: "id不能为空",
		}
		return reply, nil
	}
	edit := true
	if in.Type == "main" {
		edit = false
	}
	db := common.DB

	if edit {
		info, err := merchantmodel.FindMerchantInfoById(db, in.Id)
		if err != nil {
			return nil, err
		}
		if info != nil {
			item := new(pb.MerchantField)
			item.MchtCd = info.MchtCd
			item.Sn = info.Sn
			item.AipBranCd = info.AipBranCd
			item.GroupCd = info.GroupCd
			item.OriChnl = info.OriChnl
			item.OriChnlDesc = info.OriChnlDesc
			item.BankBelongCd = info.BankBelongCd
			item.DvpBy = info.DvpBy
			item.MccCd18 = info.MccCd18
			item.ApplDate = info.ApplDate
			item.UpBcCd = info.UpBcCd
			item.UpAcCd = info.UpAcCd
			item.UpMccCd = info.UpMccCd
			item.Name = info.Name
			item.NameBusi = info.NameBusi
			item.BusiLiceNo = info.BusiLiceNo
			item.BusiRang = info.BusiRang
			item.BusiMain = info.BusiMain
			item.Certif = info.Certif
			item.CertifType = info.CertifType
			item.CertifNo = info.CertifNo
			item.ProvCd = info.ProvCd
			item.CityCd = info.CityCd
			item.AreaCd = info.AreaCd
			item.RegAddr = info.RegAddr
			item.ContactName = info.ContactName
			item.ContactPhoneNo = info.ContactPhoneNo
			item.IsGroup = info.IsGroup
			item.MoneyToGroup = info.MoneyToGroup
			item.StlmWay = info.StlmWay
			item.StlmWayDesc = info.StlmWayDesc
			item.StlmInsCircle = info.StlmInsCircle
			item.Status = info.Status
			item.UcBcCd32 = info.UcBcCd32
			item.K2WorkflowId = info.K2WorkflowId
			item.SystemFlag = info.SystemFlag
			item.ApprovalUsername = info.ApprovalUsername
			item.FinalApprovalUsername = info.FinalApprovalUsername
			item.IsUpStandard = info.IsUpStandard
			item.BillingType = info.BillingType
			item.BillingLevel = info.BillingLevel
			item.Slogan = info.Slogan
			item.Ext1 = info.Ext1
			item.Ext2 = info.Ext2
			item.Ext3 = info.Ext3
			item.Ext4 = info.Ext4
			item.AreaStandard = info.AreaStandard
			item.MchtCdAreaCd = info.MchtCdAreaCd
			item.UcBcCdArea = info.UcBcCdArea
			item.RecOprId = info.RecOprId
			item.RecUpdOpr = info.RecUpdOpr
			item.OperIn = info.OperIn
			item.OemOrgCode = info.OemOrgCode
			item.IsEleInvoice = info.IsEleInvoice
			item.DutyParagraph = info.DutyParagraph
			item.TaxMachineBrand = info.TaxMachineBrand
			item.Ext5 = info.Ext5
			item.Ext6 = info.Ext6
			item.Ext7 = info.Ext7
			item.Ext8 = info.Ext8
			item.Ext9 = info.Ext9
			item.BusiLiceSt = info.BusiLiceSt
			item.BusiLiceDt = info.BusiLiceDt
			item.CertifSt = info.CertifSt
			item.CertifDt = info.CertifDt

			if !info.CreatedAt.IsZero() {
				item.CreatedAt = info.CreatedAt.Format(util.TimePattern)
			}
			if !info.UpdatedAt.IsZero() {
				item.UpdatedAt = info.UpdatedAt.Format(util.TimePattern)
			}
			if info.RecApllyTs.Valid {
				item.RecApllyTs = info.RecApllyTs.Time.Format(util.TimePattern)
			}
			if info.ApprDate.Valid {
				item.ApprDate = info.ApprDate.Time.Format(util.TimePattern)
			}
			if info.DeleteDate.Valid {
				item.DeleteDate = info.DeleteDate.Time.Format(util.TimePattern)
			}
			reply.Item = item
		}
	} else {
		info, err := merchantmodel.FindMerchantInfoMainById(db, in.Id)
		if err != nil {
			return nil, err
		}
		if info != nil {
			item := new(pb.MerchantField)
			item.MchtCd = info.MchtCd
			item.Sn = info.Sn
			item.AipBranCd = info.AipBranCd
			item.GroupCd = info.GroupCd
			item.OriChnl = info.OriChnl
			item.OriChnlDesc = info.OriChnlDesc
			item.BankBelongCd = info.BankBelongCd
			item.DvpBy = info.DvpBy
			item.MccCd18 = info.MccCd18
			item.ApplDate = info.ApplDate
			item.UpBcCd = info.UpBcCd
			item.UpAcCd = info.UpAcCd
			item.UpMccCd = info.UpMccCd
			item.Name = info.Name
			item.NameBusi = info.NameBusi
			item.BusiLiceNo = info.BusiLiceNo
			item.BusiRang = info.BusiRang
			item.BusiMain = info.BusiMain
			item.Certif = info.Certif
			item.CertifType = info.CertifType
			item.CertifNo = info.CertifNo
			item.ProvCd = info.ProvCd
			item.CityCd = info.CityCd
			item.AreaCd = info.AreaCd
			item.RegAddr = info.RegAddr
			item.ContactName = info.ContactName
			item.ContactPhoneNo = info.ContactPhoneNo
			item.IsGroup = info.IsGroup
			item.MoneyToGroup = info.MoneyToGroup
			item.StlmWay = info.StlmWay
			item.StlmWayDesc = info.StlmWayDesc
			item.StlmInsCircle = info.StlmInsCircle
			item.Status = info.Status
			item.UcBcCd32 = info.UcBcCd32
			item.K2WorkflowId = info.K2WorkflowId
			item.SystemFlag = info.SystemFlag
			item.ApprovalUsername = info.ApprovalUsername
			item.FinalApprovalUsername = info.FinalApprovalUsername
			item.IsUpStandard = info.IsUpStandard
			item.BillingType = info.BillingType
			item.BillingLevel = info.BillingLevel
			item.Slogan = info.Slogan
			item.Ext1 = info.Ext1
			item.Ext2 = info.Ext2
			item.Ext3 = info.Ext3
			item.Ext4 = info.Ext4
			item.AreaStandard = info.AreaStandard
			item.MchtCdAreaCd = info.MchtCdAreaCd
			item.UcBcCdArea = info.UcBcCdArea
			item.RecOprId = info.RecOprId
			item.RecUpdOpr = info.RecUpdOpr
			item.OperIn = info.OperIn
			item.OemOrgCode = info.OemOrgCode
			item.IsEleInvoice = info.IsEleInvoice
			item.DutyParagraph = info.DutyParagraph
			item.TaxMachineBrand = info.TaxMachineBrand
			item.Ext5 = info.Ext5
			item.Ext6 = info.Ext6
			item.Ext7 = info.Ext7
			item.Ext8 = info.Ext8
			item.Ext9 = info.Ext9
			item.BusiLiceSt = info.BusiLiceSt
			item.BusiLiceDt = info.BusiLiceDt
			item.CertifSt = info.CertifSt
			item.CertifDt = info.CertifDt

			if !info.CreatedAt.IsZero() {
				item.CreatedAt = info.CreatedAt.Format(util.TimePattern)
			}
			if !info.UpdatedAt.IsZero() {
				item.UpdatedAt = info.UpdatedAt.Format(util.TimePattern)
			}
			if info.RecApllyTs.Valid {
				item.RecApllyTs = info.RecApllyTs.Time.Format(util.TimePattern)
			}
			if info.ApprDate.Valid {
				item.ApprDate = info.ApprDate.Time.Format(util.TimePattern)
			}
			if info.DeleteDate.Valid {
				item.DeleteDate = info.DeleteDate.Time.Format(util.TimePattern)
			}
			reply.Item = item
		}
	}

	return reply, nil
}

func (m *merchantService) GetMerchantPicture(ctx context.Context, in *pb.GetMerchantPictureRequest) (*pb.GetMerchantPictureReply, error) {
	reply := new(pb.GetMerchantPictureReply)
	if in.Item == nil {
		reply.Err = &pb.Error{
			Code:        http.StatusBadRequest,
			Message:     "InvalidParamsError",
			Description: "查询条件不能为空",
		}
		return reply, nil
	}

	if in.Page == 0 {
		in.Page = 1
	}
	if in.Size == 0 {
		in.Size = 10
	}

	edit := true
	if in.Type == "main" {
		edit = false
	}
	db := common.DB

	if edit {
		query := new(merchantmodel.Picture)
		{
			query.FileId = in.Item.FileId
			query.MchtCd = in.Item.MchtCd
			query.DocType = in.Item.DocType
			query.FileType = in.Item.FileType
			query.FileName = in.Item.FileName
			query.PIndex = in.Item.PIndex
			query.PCode = in.Item.PCode
			query.Url = in.Item.Url
			query.SystemFlag = in.Item.SystemFlag
			query.Status = in.Item.Status
			query.RecOprId = in.Item.RecOprId
			query.RecUpdOpr = in.Item.RecUpdOpr

		}
		items, count, err := merchantmodel.QueryPicture(db, query, in.Page, in.Size)
		if err != nil {
			return nil, err
		}
		pbItems := make([]*pb.MerchantPictureField, len(items))

		for i := range items {
			pbItems[i] = &pb.MerchantPictureField{
				FileId:     items[i].FileId,
				MchtCd:     items[i].MchtCd,
				DocType:    items[i].DocType,
				FileType:   items[i].FileType,
				FileName:   items[i].FileName,
				PIndex:     items[i].PIndex,
				PCode:      items[i].PCode,
				Url:        items[i].Url,
				SystemFlag: items[i].SystemFlag,
				Status:     items[i].Status,
				RecOprId:   items[i].RecOprId,
				RecUpdOpr:  items[i].RecUpdOpr,
			}
			if !items[i].CreatedAt.IsZero() {
				pbItems[i].CreatedAt = items[i].CreatedAt.Format(util.TimePattern)
			}
			if !items[i].UpdatedAt.IsZero() {
				pbItems[i].UpdatedAt = items[i].UpdatedAt.Format(util.TimePattern)
			}
		}

		reply.Count = count
		reply.Items = pbItems

	} else {
		query := new(merchantmodel.PictureMain)
		{
			query.FileId = in.Item.FileId
			query.MchtCd = in.Item.MchtCd
			query.DocType = in.Item.DocType
			query.FileType = in.Item.FileType
			query.FileName = in.Item.FileName
			query.PIndex = in.Item.PIndex
			query.PCode = in.Item.PCode
			query.Url = in.Item.Url
			query.SystemFlag = in.Item.SystemFlag
			query.Status = in.Item.Status
			query.RecOprId = in.Item.RecOprId
			query.RecUpdOpr = in.Item.RecUpdOpr

		}
		items, count, err := merchantmodel.QueryPictureMain(db, query, in.Page, in.Size)
		if err != nil {
			return nil, err
		}
		pbItems := make([]*pb.MerchantPictureField, len(items))

		for i := range items {
			pbItems[i] = &pb.MerchantPictureField{
				FileId:     items[i].FileId,
				MchtCd:     items[i].MchtCd,
				DocType:    items[i].DocType,
				FileType:   items[i].FileType,
				FileName:   items[i].FileName,
				PIndex:     items[i].PIndex,
				PCode:      items[i].PCode,
				Url:        items[i].Url,
				SystemFlag: items[i].SystemFlag,
				Status:     items[i].Status,
				RecOprId:   items[i].RecOprId,
				RecUpdOpr:  items[i].RecUpdOpr,
			}
			if !items[i].CreatedAt.IsZero() {
				pbItems[i].CreatedAt = items[i].CreatedAt.Format(util.TimePattern)
			}
			if !items[i].UpdatedAt.IsZero() {
				pbItems[i].UpdatedAt = items[i].UpdatedAt.Format(util.TimePattern)
			}
		}

		reply.Count = count
		reply.Items = pbItems

	}

	reply.Page = in.Page
	reply.Size = in.Size
	return reply, nil
}

func (m *merchantService) GetMerchantBusiness(ctx context.Context, in *pb.GetMerchantBusinessRequest) (*pb.GetMerchantBusinessReply, error) {
	reply := new(pb.GetMerchantBusinessReply)
	if in.Item == nil {
		reply.Err = &pb.Error{
			Code:        http.StatusBadRequest,
			Message:     "InvalidParamsError",
			Description: "查询条件不能为空",
		}
		return reply, nil
	}

	if in.Page == 0 {
		in.Page = 1
	}
	if in.Size == 0 {
		in.Size = 10
	}

	edit := true
	if in.Type == "main" {
		edit = false
	}
	db := common.DB

	if edit {
		query := new(merchantmodel.Business)
		{
			query.MchtCd = in.Item.MchtCd
			query.ProdCd = in.Item.ProdCd
			query.ProdCdText = in.Item.ProdCdText
			query.FeeMoneyCd = in.Item.FeeMoneyCd
			query.FeeModeType = in.Item.FeeModeType
			query.FeeSettlementType = in.Item.FeeSettlementType
			query.FeeHoliday = in.Item.FeeHoliday
			query.ServiceFeeType = in.Item.ServiceFeeType
			query.ServiceFeeStaticAmount, _ = strconv.ParseFloat(in.Item.ServiceFeeStaticAmount, 64)
			query.ServiceFeeLevelCount, _ = strconv.ParseInt(in.Item.ServiceFeeLevelCount, 10, 64)
			query.ServiceFeeMode = in.Item.ServiceFeeMode
			query.ServiceFeeUnit = in.Item.ServiceFeeUnit
			query.ServiceFeeTerm = in.Item.ServiceFeeTerm
			query.ServiceFeeSumto = in.Item.ServiceFeeSumto
			query.ServiceFeeCircle = in.Item.ServiceFeeCircle
			query.ServiceFeeOthers = in.Item.ServiceFeeOthers
			query.ServiceFeeStart = in.Item.ServiceFeeStart
			query.ServiceFeeClct = in.Item.ServiceFeeClct
			query.ServiceFeeClctOthers = in.Item.ServiceFeeClctOthers
			query.SystemFlag = in.Item.SystemFlag
			query.Ext1 = in.Item.Ext1
			query.Ext2 = in.Item.Ext2
			query.Ext3 = in.Item.Ext3
			query.Ext4 = in.Item.Ext4
			query.ServiceFeeYesNo = in.Item.ServiceFeeYesNo
			query.RecOprId = in.Item.RecOprId
			query.RecUpdOpr = in.Item.RecUpdOpr
			query.OperIn = in.Item.OperIn

		}
		items, count, err := merchantmodel.QueryBusiness(db, query, in.Page, in.Size)
		if err != nil {
			return nil, err
		}
		pbItems := make([]*pb.MerchantBusinessField, len(items))

		for i := range items {
			pbItems[i] = &pb.MerchantBusinessField{
				MchtCd:                 items[i].MchtCd,
				ProdCd:                 items[i].ProdCd,
				ProdCdText:             items[i].ProdCdText,
				FeeMoneyCd:             items[i].FeeMoneyCd,
				FeeModeType:            items[i].FeeModeType,
				FeeSettlementType:      items[i].FeeSettlementType,
				FeeHoliday:             items[i].FeeHoliday,
				ServiceFeeType:         items[i].ServiceFeeType,
				ServiceFeeStaticAmount: fmt.Sprintf("%f", items[i].ServiceFeeStaticAmount),
				ServiceFeeLevelCount:   fmt.Sprintf("%d", items[i].ServiceFeeLevelCount),
				ServiceFeeMode:         items[i].ServiceFeeMode,
				ServiceFeeUnit:         items[i].ServiceFeeUnit,
				ServiceFeeTerm:         items[i].ServiceFeeTerm,
				ServiceFeeSumto:        items[i].ServiceFeeSumto,
				ServiceFeeCircle:       items[i].ServiceFeeCircle,
				ServiceFeeOthers:       items[i].ServiceFeeOthers,
				ServiceFeeStart:        items[i].ServiceFeeStart,
				ServiceFeeClct:         items[i].ServiceFeeClct,
				ServiceFeeClctOthers:   items[i].ServiceFeeClctOthers,
				SystemFlag:             items[i].SystemFlag,
				Ext1:                   items[i].Ext1,
				Ext2:                   items[i].Ext2,
				Ext3:                   items[i].Ext3,
				Ext4:                   items[i].Ext4,
				ServiceFeeYesNo:        items[i].ServiceFeeYesNo,
				RecOprId:               items[i].RecOprId,
				RecUpdOpr:              items[i].RecUpdOpr,
				OperIn:                 items[i].OperIn,
			}
			if !items[i].CreatedAt.IsZero() {
				pbItems[i].CreatedAt = items[i].CreatedAt.Format(util.TimePattern)
			}
			if !items[i].UpdatedAt.IsZero() {
				pbItems[i].UpdatedAt = items[i].UpdatedAt.Format(util.TimePattern)
			}
		}

		reply.Count = count
		reply.Items = pbItems
	} else {
		query := new(merchantmodel.BusinessMain)
		{
			query.MchtCd = in.Item.MchtCd
			query.ProdCd = in.Item.ProdCd
			query.ProdCdText = in.Item.ProdCdText
			query.FeeMoneyCd = in.Item.FeeMoneyCd
			query.FeeModeType = in.Item.FeeModeType
			query.FeeSettlementType = in.Item.FeeSettlementType
			query.FeeHoliday = in.Item.FeeHoliday
			query.ServiceFeeType = in.Item.ServiceFeeType
			query.ServiceFeeStaticAmount, _ = strconv.ParseFloat(in.Item.ServiceFeeStaticAmount, 64)
			query.ServiceFeeLevelCount, _ = strconv.ParseInt(in.Item.ServiceFeeLevelCount, 10, 64)
			query.ServiceFeeMode = in.Item.ServiceFeeMode
			query.ServiceFeeUnit = in.Item.ServiceFeeUnit
			query.ServiceFeeTerm = in.Item.ServiceFeeTerm
			query.ServiceFeeSumto = in.Item.ServiceFeeSumto
			query.ServiceFeeCircle = in.Item.ServiceFeeCircle
			query.ServiceFeeOthers = in.Item.ServiceFeeOthers
			query.ServiceFeeStart = in.Item.ServiceFeeStart
			query.ServiceFeeClct = in.Item.ServiceFeeClct
			query.ServiceFeeClctOthers = in.Item.ServiceFeeClctOthers
			query.SystemFlag = in.Item.SystemFlag
			query.Ext1 = in.Item.Ext1
			query.Ext2 = in.Item.Ext2
			query.Ext3 = in.Item.Ext3
			query.Ext4 = in.Item.Ext4
			query.ServiceFeeYesNo = in.Item.ServiceFeeYesNo
			query.RecOprId = in.Item.RecOprId
			query.RecUpdOpr = in.Item.RecUpdOpr
			query.OperIn = in.Item.OperIn

		}
		items, count, err := merchantmodel.QueryBusinessMain(db, query, in.Page, in.Size)
		if err != nil {
			return nil, err
		}
		pbItems := make([]*pb.MerchantBusinessField, len(items))

		for i := range items {
			pbItems[i] = &pb.MerchantBusinessField{
				MchtCd:                 items[i].MchtCd,
				ProdCd:                 items[i].ProdCd,
				ProdCdText:             items[i].ProdCdText,
				FeeMoneyCd:             items[i].FeeMoneyCd,
				FeeModeType:            items[i].FeeModeType,
				FeeSettlementType:      items[i].FeeSettlementType,
				FeeHoliday:             items[i].FeeHoliday,
				ServiceFeeType:         items[i].ServiceFeeType,
				ServiceFeeStaticAmount: fmt.Sprintf("%f", items[i].ServiceFeeStaticAmount),
				ServiceFeeLevelCount:   fmt.Sprintf("%d", items[i].ServiceFeeLevelCount),
				ServiceFeeMode:         items[i].ServiceFeeMode,
				ServiceFeeUnit:         items[i].ServiceFeeUnit,
				ServiceFeeTerm:         items[i].ServiceFeeTerm,
				ServiceFeeSumto:        items[i].ServiceFeeSumto,
				ServiceFeeCircle:       items[i].ServiceFeeCircle,
				ServiceFeeOthers:       items[i].ServiceFeeOthers,
				ServiceFeeStart:        items[i].ServiceFeeStart,
				ServiceFeeClct:         items[i].ServiceFeeClct,
				ServiceFeeClctOthers:   items[i].ServiceFeeClctOthers,
				SystemFlag:             items[i].SystemFlag,
				Ext1:                   items[i].Ext1,
				Ext2:                   items[i].Ext2,
				Ext3:                   items[i].Ext3,
				Ext4:                   items[i].Ext4,
				ServiceFeeYesNo:        items[i].ServiceFeeYesNo,
				RecOprId:               items[i].RecOprId,
				RecUpdOpr:              items[i].RecUpdOpr,
				OperIn:                 items[i].OperIn,
			}
			if !items[i].CreatedAt.IsZero() {
				pbItems[i].CreatedAt = items[i].CreatedAt.Format(util.TimePattern)
			}
			if !items[i].UpdatedAt.IsZero() {
				pbItems[i].UpdatedAt = items[i].UpdatedAt.Format(util.TimePattern)
			}
		}

		reply.Count = count
		reply.Items = pbItems
	}

	reply.Page = in.Page
	reply.Size = in.Size
	return reply, nil
}

func (m *merchantService) GetMerchantBizFee(ctx context.Context, in *pb.GetMerchantBizFeeRequest) (*pb.GetMerchantBizFeeReply, error) {
	reply := new(pb.GetMerchantBizFeeReply)
	if in.Item == nil {
		reply.Err = &pb.Error{
			Code:        http.StatusBadRequest,
			Message:     "InvalidParamsError",
			Description: "查询条件不能为空",
		}
		return reply, nil
	}

	if in.Page == 0 {
		in.Page = 1
	}
	if in.Size == 0 {
		in.Size = 10
	}

	edit := true
	if in.Type == "main" {
		edit = false
	}
	db := common.DB
	if edit {
		query := new(merchantmodel.BizFee)
		{
			query.MchtCd = in.Item.MchtCd
			query.ProdCd = in.Item.ProdCd
			query.BizCd = in.Item.BizCd
			query.SubBizCd = in.Item.SubBizCd
			query.MchtFeeMd = in.Item.MchtFeeMd
			query.MchtFeePercent, _ = strconv.ParseFloat(in.Item.MchtFeePercent, 64)
			query.MchtFeePctMin, _ = strconv.ParseFloat(in.Item.MchtFeePctMin, 64)
			query.MchtFeePctMax, _ = strconv.ParseFloat(in.Item.MchtFeePctMax, 64)
			query.MchtFeeSingle, _ = strconv.ParseFloat(in.Item.MchtFeeSingle, 64)
			query.MchtAFeePercent, _ = strconv.ParseFloat(in.Item.MchtAFeePercent, 64)
			query.MchtAFeePctMin, _ = strconv.ParseFloat(in.Item.MchtAFeePctMin, 64)
			query.MchtAFeePctMax, _ = strconv.ParseFloat(in.Item.MchtAFeePctMax, 64)
			query.MchtAFeeSingle, _ = strconv.ParseFloat(in.Item.MchtAFeeSingle, 64)
			query.MchtAFeeSame = in.Item.MchtAFeeSame
			query.MchtAFeeMd = in.Item.MchtAFeeMd
			query.OperIn = in.Item.OperIn
			query.RecOprId = in.Item.RecOprId
			query.RecUpdOpr = in.Item.RecUpdOpr

		}
		items, count, err := merchantmodel.QueryBizFee(db, query, in.Page, in.Size)
		if err != nil {
			return nil, err
		}
		pbItems := make([]*pb.MerchantBizFeeField, len(items))

		for i := range items {
			pbItems[i] = &pb.MerchantBizFeeField{
				MchtCd:          items[i].MchtCd,
				ProdCd:          items[i].ProdCd,
				BizCd:           items[i].BizCd,
				SubBizCd:        items[i].SubBizCd,
				MchtFeeMd:       items[i].MchtFeeMd,
				MchtFeePercent:  fmt.Sprintf("%f", items[i].MchtFeePercent),
				MchtFeePctMin:   fmt.Sprintf("%f", items[i].MchtFeePctMin),
				MchtFeePctMax:   fmt.Sprintf("%f", items[i].MchtFeePctMax),
				MchtFeeSingle:   fmt.Sprintf("%f", items[i].MchtFeeSingle),
				MchtAFeeSame:    items[i].MchtAFeeSame,
				MchtAFeeMd:      items[i].MchtAFeeMd,
				MchtAFeePercent: fmt.Sprintf("%f", items[i].MchtAFeePercent),
				MchtAFeePctMin:  fmt.Sprintf("%f", items[i].MchtAFeePctMin),
				MchtAFeePctMax:  fmt.Sprintf("%f", items[i].MchtAFeePctMax),
				MchtAFeeSingle:  fmt.Sprintf("%f", items[i].MchtAFeeSingle),
				OperIn:          items[i].OperIn,
				RecOprId:        items[i].RecOprId,
				RecUpdOpr:       items[i].RecUpdOpr,
			}
			if !items[i].CreatedAt.IsZero() {
				pbItems[i].CreatedAt = items[i].CreatedAt.Format(util.TimePattern)
			}
			if !items[i].UpdatedAt.IsZero() {
				pbItems[i].UpdatedAt = items[i].UpdatedAt.Format(util.TimePattern)
			}
		}

		reply.Count = count
		reply.Items = pbItems
	} else {
		query := new(merchantmodel.BizFeeMain)
		{
			query.MchtCd = in.Item.MchtCd
			query.ProdCd = in.Item.ProdCd
			query.BizCd = in.Item.BizCd
			query.SubBizCd = in.Item.SubBizCd
			query.MchtFeeMd = in.Item.MchtFeeMd
			query.MchtFeePercent, _ = strconv.ParseFloat(in.Item.MchtFeePercent, 64)
			query.MchtFeePctMin, _ = strconv.ParseFloat(in.Item.MchtFeePctMin, 64)
			query.MchtFeePctMax, _ = strconv.ParseFloat(in.Item.MchtFeePctMax, 64)
			query.MchtFeeSingle, _ = strconv.ParseFloat(in.Item.MchtFeeSingle, 64)
			query.MchtAFeePercent, _ = strconv.ParseFloat(in.Item.MchtAFeePercent, 64)
			query.MchtAFeePctMin, _ = strconv.ParseFloat(in.Item.MchtAFeePctMin, 64)
			query.MchtAFeePctMax, _ = strconv.ParseFloat(in.Item.MchtAFeePctMax, 64)
			query.MchtAFeeSingle, _ = strconv.ParseFloat(in.Item.MchtAFeeSingle, 64)
			query.MchtAFeeSame = in.Item.MchtAFeeSame
			query.MchtAFeeMd = in.Item.MchtAFeeMd
			query.OperIn = in.Item.OperIn
			query.RecOprId = in.Item.RecOprId
			query.RecUpdOpr = in.Item.RecUpdOpr

		}
		items, count, err := merchantmodel.QueryBizFeeMain(db, query, in.Page, in.Size)
		if err != nil {
			return nil, err
		}
		pbItems := make([]*pb.MerchantBizFeeField, len(items))

		for i := range items {
			pbItems[i] = &pb.MerchantBizFeeField{
				MchtCd:          items[i].MchtCd,
				ProdCd:          items[i].ProdCd,
				BizCd:           items[i].BizCd,
				SubBizCd:        items[i].SubBizCd,
				MchtFeeMd:       items[i].MchtFeeMd,
				MchtFeePercent:  fmt.Sprintf("%f", items[i].MchtFeePercent),
				MchtFeePctMin:   fmt.Sprintf("%f", items[i].MchtFeePctMin),
				MchtFeePctMax:   fmt.Sprintf("%f", items[i].MchtFeePctMax),
				MchtFeeSingle:   fmt.Sprintf("%f", items[i].MchtFeeSingle),
				MchtAFeeSame:    items[i].MchtAFeeSame,
				MchtAFeeMd:      items[i].MchtAFeeMd,
				MchtAFeePercent: fmt.Sprintf("%f", items[i].MchtAFeePercent),
				MchtAFeePctMin:  fmt.Sprintf("%f", items[i].MchtAFeePctMin),
				MchtAFeePctMax:  fmt.Sprintf("%f", items[i].MchtAFeePctMax),
				MchtAFeeSingle:  fmt.Sprintf("%f", items[i].MchtAFeeSingle),
				OperIn:          items[i].OperIn,
				RecOprId:        items[i].RecOprId,
				RecUpdOpr:       items[i].RecUpdOpr,
			}
			if !items[i].CreatedAt.IsZero() {
				pbItems[i].CreatedAt = items[i].CreatedAt.Format(util.TimePattern)
			}
			if !items[i].UpdatedAt.IsZero() {
				pbItems[i].UpdatedAt = items[i].UpdatedAt.Format(util.TimePattern)
			}
		}

		reply.Count = count
		reply.Items = pbItems
	}

	reply.Page = in.Page
	reply.Size = in.Size
	return reply, nil
}

func (m *merchantService) GetMerchantBizDeal(ctx context.Context, in *pb.GetMerchantBizDealRequest) (*pb.GetMerchantBizDealReply, error) {
	reply := new(pb.GetMerchantBizDealReply)
	if in.Item == nil {
		reply.Err = &pb.Error{
			Code:        http.StatusBadRequest,
			Message:     "InvalidParamsError",
			Description: "查询条件不能为空",
		}
		return reply, nil
	}

	if in.Page == 0 {
		in.Page = 1
	}
	if in.Size == 0 {
		in.Size = 10
	}

	edit := true
	if in.Type == "main" {
		edit = false
	}
	db := common.DB

	if edit {
		query := new(merchantmodel.BizDeal)
		{
			query.MchtCd = in.Item.MchtCd
			query.ProdCd = in.Item.ProdCd
			query.BizCd = in.Item.BizCd
			query.TransCd = in.Item.TransCd
			query.OperIn = in.Item.OperIn
			query.RecOprId = in.Item.RecOprId
			query.RecUpdOpr = in.Item.RecUpdOpr
		}
		items, count, err := merchantmodel.QueryBizDeal(db, query, in.Page, in.Size)
		if err != nil {
			return nil, err
		}
		pbItems := make([]*pb.MerchantBizDealField, len(items))

		for i := range items {
			pbItems[i] = &pb.MerchantBizDealField{
				MchtCd:    items[i].MchtCd,
				ProdCd:    items[i].ProdCd,
				BizCd:     items[i].BizCd,
				TransCd:   items[i].TransCd,
				OperIn:    items[i].OperIn,
				RecOprId:  items[i].RecOprId,
				RecUpdOpr: items[i].RecUpdOpr,
			}
			if !items[i].CreatedAt.IsZero() {
				pbItems[i].CreatedAt = items[i].CreatedAt.Format(util.TimePattern)
			}
			if !items[i].UpdatedAt.IsZero() {
				pbItems[i].UpdatedAt = items[i].UpdatedAt.Format(util.TimePattern)
			}
		}

		reply.Count = count
		reply.Items = pbItems

	} else {
		query := new(merchantmodel.BizDealMain)
		{
			query.MchtCd = in.Item.MchtCd
			query.ProdCd = in.Item.ProdCd
			query.BizCd = in.Item.BizCd
			query.TransCd = in.Item.TransCd
			query.OperIn = in.Item.OperIn
			query.RecOprId = in.Item.RecOprId
			query.RecUpdOpr = in.Item.RecUpdOpr
		}
		items, count, err := merchantmodel.QueryBizDealMain(db, query, in.Page, in.Size)
		if err != nil {
			return nil, err
		}
		pbItems := make([]*pb.MerchantBizDealField, len(items))

		for i := range items {
			pbItems[i] = &pb.MerchantBizDealField{
				MchtCd:    items[i].MchtCd,
				ProdCd:    items[i].ProdCd,
				BizCd:     items[i].BizCd,
				TransCd:   items[i].TransCd,
				OperIn:    items[i].OperIn,
				RecOprId:  items[i].RecOprId,
				RecUpdOpr: items[i].RecUpdOpr,
			}
			if !items[i].CreatedAt.IsZero() {
				pbItems[i].CreatedAt = items[i].CreatedAt.Format(util.TimePattern)
			}
			if !items[i].UpdatedAt.IsZero() {
				pbItems[i].UpdatedAt = items[i].UpdatedAt.Format(util.TimePattern)
			}
		}

		reply.Count = count
		reply.Items = pbItems
	}

	reply.Page = in.Page
	reply.Size = in.Size
	return reply, nil
}

func (m *merchantService) GetMerchantBankAccount(ctx context.Context, in *pb.GetMerchantBankAccountRequest) (*pb.GetMerchantBankAccountReply, error) {
	reply := new(pb.GetMerchantBankAccountReply)
	if in.Item == nil {
		reply.Err = &pb.Error{
			Code:        http.StatusBadRequest,
			Message:     "InvalidParamsError",
			Description: "查询条件不能为空",
		}
		return reply, nil
	}

	if in.Page == 0 {
		in.Page = 1
	}
	if in.Size == 0 {
		in.Size = 10
	}

	edit := true
	if in.Type == "main" {
		edit = false
	}
	db := common.DB

	if edit {
		query := new(merchantmodel.BankAccount)
		{
			query.OwnerCd = in.Item.OwnerCd
			query.OwnerCd = in.Item.OwnerCd
			query.AccountType = in.Item.AccountType
			query.Name = in.Item.Name
			query.Account = in.Item.Account
			query.UcBcCd = in.Item.UcBcCd
			query.Province = in.Item.Province
			query.City = in.Item.City
			query.BankCode = in.Item.BankCode
			query.BankName = in.Item.BankName
			query.OperIn = in.Item.OperIn
			query.RecOprId = in.Item.RecOprId
			query.RecUpdOpr = in.Item.RecUpdOpr
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
		}
		items, count, err := merchantmodel.QueryBankAccount(db, query, in.Page, in.Size)
		if err != nil {
			return nil, err
		}
		pbItems := make([]*pb.MerchantBankAccountField, len(items))

		for i := range items {
			pbItems[i] = &pb.MerchantBankAccountField{
				OwnerCd:      items[i].OwnerCd,
				AccountType:  items[i].AccountType,
				Name:         items[i].Name,
				Account:      items[i].Account,
				UcBcCd:       items[i].UcBcCd,
				Province:     items[i].Province,
				City:         items[i].City,
				BankCode:     items[i].BankCode,
				BankName:     items[i].BankName,
				OperIn:       items[i].OperIn,
				RecOprId:     items[i].RecOprId,
				RecUpdOpr:    items[i].RecUpdOpr,
				MsgResvFld1:  items[i].MsgResvFld1,
				MsgResvFld2:  items[i].MsgResvFld2,
				MsgResvFld3:  items[i].MsgResvFld3,
				MsgResvFld4:  items[i].MsgResvFld4,
				MsgResvFld5:  items[i].MsgResvFld5,
				MsgResvFld6:  items[i].MsgResvFld6,
				MsgResvFld7:  items[i].MsgResvFld7,
				MsgResvFld8:  items[i].MsgResvFld8,
				MsgResvFld9:  items[i].MsgResvFld9,
				MsgResvFld10: items[i].MsgResvFld10,
			}
			if !items[i].CreatedAt.IsZero() {
				pbItems[i].CreatedAt = items[i].CreatedAt.Format(util.TimePattern)
			}
			if !items[i].UpdatedAt.IsZero() {
				pbItems[i].UpdatedAt = items[i].UpdatedAt.Format(util.TimePattern)
			}
		}

		reply.Count = count
		reply.Items = pbItems

	} else {
		query := new(merchantmodel.BankAccountMain)
		{
			query.OwnerCd = in.Item.OwnerCd
			query.OwnerCd = in.Item.OwnerCd
			query.AccountType = in.Item.AccountType
			query.Name = in.Item.Name
			query.Account = in.Item.Account
			query.UcBcCd = in.Item.UcBcCd
			query.Province = in.Item.Province
			query.City = in.Item.City
			query.BankCode = in.Item.BankCode
			query.BankName = in.Item.BankName
			query.OperIn = in.Item.OperIn
			query.RecOprId = in.Item.RecOprId
			query.RecUpdOpr = in.Item.RecUpdOpr
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
		}
		items, count, err := merchantmodel.QueryBankAccountMain(db, query, in.Page, in.Size)
		if err != nil {
			return nil, err
		}
		pbItems := make([]*pb.MerchantBankAccountField, len(items))

		for i := range items {
			pbItems[i] = &pb.MerchantBankAccountField{
				OwnerCd:      items[i].OwnerCd,
				AccountType:  items[i].AccountType,
				Name:         items[i].Name,
				Account:      items[i].Account,
				UcBcCd:       items[i].UcBcCd,
				Province:     items[i].Province,
				City:         items[i].City,
				BankCode:     items[i].BankCode,
				BankName:     items[i].BankName,
				OperIn:       items[i].OperIn,
				RecOprId:     items[i].RecOprId,
				RecUpdOpr:    items[i].RecUpdOpr,
				MsgResvFld1:  items[i].MsgResvFld1,
				MsgResvFld2:  items[i].MsgResvFld2,
				MsgResvFld3:  items[i].MsgResvFld3,
				MsgResvFld4:  items[i].MsgResvFld4,
				MsgResvFld5:  items[i].MsgResvFld5,
				MsgResvFld6:  items[i].MsgResvFld6,
				MsgResvFld7:  items[i].MsgResvFld7,
				MsgResvFld8:  items[i].MsgResvFld8,
				MsgResvFld9:  items[i].MsgResvFld9,
				MsgResvFld10: items[i].MsgResvFld10,
			}
			if !items[i].CreatedAt.IsZero() {
				pbItems[i].CreatedAt = items[i].CreatedAt.Format(util.TimePattern)
			}
			if !items[i].UpdatedAt.IsZero() {
				pbItems[i].UpdatedAt = items[i].UpdatedAt.Format(util.TimePattern)
			}
		}

		reply.Count = count
		reply.Items = pbItems
	}

	reply.Page = in.Page
	reply.Size = in.Size
	return reply, nil
}

func (m *merchantService) SaveMerchant(ctx context.Context, in *pb.SaveMerchantRequest) (*pb.SaveMerchantReply, error) {
	var reply pb.SaveMerchantReply
	if in.Merchant == nil {
		reply.Err = &pb.Error{
			Code:        http.StatusBadRequest,
			Message:     "InvalidParamsError",
			Description: "保存信息为空",
		}
		return &reply, nil
	}

	if in.Merchant.MchtCd == "" {
		reply.Err = &pb.Error{
			Code:        http.StatusBadRequest,
			Message:     "InvalidParamsError",
			Description: "id不能为空",
		}
		return &reply, nil
	}
	db := common.DB.Begin()
	defer db.Rollback()
	var err error

	mtch := new(merchantmodel.MerchantInfo)
	{
		mtch.MchtCd = in.Merchant.MchtCd
		mtch.Sn = in.Merchant.Sn
		mtch.AipBranCd = in.Merchant.AipBranCd
		mtch.GroupCd = in.Merchant.GroupCd
		mtch.OriChnl = in.Merchant.OriChnl
		mtch.OriChnlDesc = in.Merchant.OriChnlDesc
		mtch.BankBelongCd = in.Merchant.BankBelongCd
		mtch.DvpBy = in.Merchant.DvpBy
		mtch.MccCd18 = in.Merchant.MccCd18
		mtch.ApplDate = in.Merchant.ApplDate
		mtch.UpBcCd = in.Merchant.UpBcCd
		mtch.UpAcCd = in.Merchant.UpAcCd
		mtch.UpMccCd = in.Merchant.UpMccCd
		mtch.Name = in.Merchant.Name
		mtch.NameBusi = in.Merchant.NameBusi
		mtch.BusiLiceNo = in.Merchant.BusiLiceNo
		mtch.BusiRang = in.Merchant.BusiRang
		mtch.BusiMain = in.Merchant.BusiMain
		mtch.Certif = in.Merchant.Certif
		mtch.CertifType = in.Merchant.CertifType
		mtch.CertifNo = in.Merchant.CertifNo
		mtch.ProvCd = in.Merchant.ProvCd
		mtch.CityCd = in.Merchant.CityCd
		mtch.AreaCd = in.Merchant.AreaCd
		mtch.RegAddr = in.Merchant.RegAddr
		mtch.ContactName = in.Merchant.ContactName
		mtch.ContactPhoneNo = in.Merchant.ContactPhoneNo
		mtch.IsGroup = in.Merchant.IsGroup
		mtch.MoneyToGroup = in.Merchant.MoneyToGroup
		mtch.StlmWay = in.Merchant.StlmWay
		mtch.StlmWayDesc = in.Merchant.StlmWayDesc
		mtch.StlmInsCircle = in.Merchant.StlmInsCircle
		mtch.Status = in.Merchant.Status
		mtch.UcBcCd32 = in.Merchant.UcBcCd32
		mtch.K2WorkflowId = in.Merchant.K2WorkflowId
		mtch.SystemFlag = in.Merchant.SystemFlag
		mtch.ApprovalUsername = in.Merchant.ApprovalUsername
		mtch.FinalApprovalUsername = in.Merchant.FinalApprovalUsername
		mtch.IsUpStandard = in.Merchant.IsUpStandard
		mtch.BillingType = in.Merchant.BillingType
		mtch.BillingLevel = in.Merchant.BillingLevel
		mtch.Slogan = in.Merchant.Slogan
		mtch.Ext1 = in.Merchant.Ext1
		mtch.Ext2 = in.Merchant.Ext2
		mtch.Ext3 = in.Merchant.Ext3
		mtch.Ext4 = in.Merchant.Ext4
		mtch.AreaStandard = in.Merchant.AreaStandard
		mtch.MchtCdAreaCd = in.Merchant.MchtCdAreaCd
		mtch.UcBcCdArea = in.Merchant.UcBcCdArea
		mtch.RecOprId = in.Merchant.RecOprId
		mtch.RecUpdOpr = in.Merchant.RecUpdOpr
		mtch.OperIn = in.Merchant.OperIn
		mtch.IsEleInvoice = in.Merchant.IsEleInvoice
		mtch.DutyParagraph = in.Merchant.DutyParagraph
		mtch.TaxMachineBrand = in.Merchant.TaxMachineBrand
		mtch.Ext5 = in.Merchant.Ext5
		mtch.Ext6 = in.Merchant.Ext6
		mtch.Ext7 = in.Merchant.Ext7
		mtch.Ext8 = in.Merchant.Ext8
		mtch.Ext9 = in.Merchant.Ext9
		mtch.BusiLiceSt = in.Merchant.BusiLiceSt
		mtch.BusiLiceDt = in.Merchant.BusiLiceDt
		mtch.CertifSt = in.Merchant.CertifSt
		mtch.CertifDt = in.Merchant.CertifDt
		mtch.OemOrgCode = in.Merchant.OemOrgCode

		if in.Merchant.ApprDate != "" {
			mtch.ApprDate.Time, err = time.Parse(util.TimePattern, in.Merchant.ApprDate)
			if err == nil {
				mtch.ApprDate.Valid = true
			}
		}
		if in.Merchant.DeleteDate != "" {
			mtch.DeleteDate.Time, err = time.Parse(util.TimePattern, in.Merchant.ApprDate)
			if err == nil {
				mtch.DeleteDate.Valid = true
			}
		}
		if in.Merchant.RecApllyTs != "" {
			mtch.DeleteDate.Time, err = time.Parse(util.TimePattern, in.Merchant.ApprDate)
			if err == nil {
				mtch.DeleteDate.Valid = true
			}
		}
	}
	err = merchantmodel.SaveMerchant(db, mtch)
	if err != nil {
		return nil, err
	}

	if in.Account != nil {
		data := new(merchantmodel.BankAccount)
		{
			data.OwnerCd = in.Account.OwnerCd
			data.AccountType = in.Account.AccountType
			data.Name = in.Account.Name
			data.Account = in.Account.Account
			data.UcBcCd = in.Account.UcBcCd
			data.Province = in.Account.Province
			data.City = in.Account.City
			data.BankCode = in.Account.BankCode
			data.BankName = in.Account.BankName
			data.OperIn = in.Account.OperIn
			data.RecOprId = in.Account.RecOprId
			data.RecUpdOpr = in.Account.RecUpdOpr
			data.MsgResvFld1 = in.Account.MsgResvFld1
			data.MsgResvFld2 = in.Account.MsgResvFld2
			data.MsgResvFld3 = in.Account.MsgResvFld3
			data.MsgResvFld4 = in.Account.MsgResvFld4
			data.MsgResvFld5 = in.Account.MsgResvFld5
			data.MsgResvFld6 = in.Account.MsgResvFld6
			data.MsgResvFld7 = in.Account.MsgResvFld7
			data.MsgResvFld8 = in.Account.MsgResvFld8
			data.MsgResvFld9 = in.Account.MsgResvFld9
			data.MsgResvFld10 = in.Account.MsgResvFld10
		}
		err = merchantmodel.SaveBankAccount(db, data)
		if err != nil {
			return nil, err
		}
	}

	db.Commit()
	return &reply, nil
}

func (m *merchantService) ListMerchant(ctx context.Context, in *pb.ListMerchantRequest) (*pb.ListMerchantReply, error) {
	reply := new(pb.ListMerchantReply)
	if in.Size == 0 {
		in.Size = 10
	}
	if in.Page == 0 {
		in.Page = 1
	}

	db := common.DB

	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		reply.Err = &pb.Error{
			Code:        http.StatusBadRequest,
			Message:     InvalidParam,
			Description: "找不到用户信息",
		}
		return reply, nil
	}
	ids := md.Get("userid")
	if !ok || len(ids) == 0 {
		reply.Err = &pb.Error{
			Code:        http.StatusBadRequest,
			Message:     InvalidParam,
			Description: "找不到用户信息",
		}
		return reply, nil
	}
	id, _ := strconv.ParseInt(ids[0], 10, 64)
	user, err := usermodel.FindUserByID(db, id)
	if err != nil {
		return nil, err
	}
	if user == nil {
		reply.Err = &pb.Error{
			Code:        http.StatusBadRequest,
			Message:     InvalidParam,
			Description: "找不到用户信息",
		}
		return reply, nil
	}

	insIds := make([]string, 0)
	if user.UserType != "admin" && user.UserType != "institution" && user.UserType != "institution_group" {
		reply.Items = make([]*pb.MerchantField, 0)
		reply.Count = 0
		reply.Page = in.Page
		reply.Size = in.Size
		return reply, nil
	}
	if user.UserType == "institution_group" {
		groupId, _ := strconv.ParseInt(user.UserGroupNo, 10, 64)
		groups, err := insmodel.ListInsGroupBind(db, groupId)
		if err != nil {
			return nil, err
		}
		for _, group := range groups {
			insIds = append(insIds, group.InsIdCd)
		}
		if len(insIds) == 0 {
			reply.Items = make([]*pb.MerchantField, 0)
			reply.Count = 0
			reply.Page = in.Page
			reply.Size = in.Size
			return reply, nil
		}
	}

	edit := true
	if in.Type == "main" {
		edit = false
	}

	if edit {
		query := new(merchantmodel.MerchantInfo)
		if in.Item != nil {
			query.MchtCd = in.Item.MchtCd
			query.Sn = in.Item.Sn
			if user.UserType == "institution" {
				query.AipBranCd = user.UserGroupNo
			}
			if in.Item.AipBranCd != "" {
				query.AipBranCd = in.Item.AipBranCd
			}
			query.GroupCd = in.Item.GroupCd
			query.OriChnl = in.Item.OriChnl
			query.OriChnlDesc = in.Item.OriChnlDesc
			query.BankBelongCd = in.Item.BankBelongCd
			query.DvpBy = in.Item.DvpBy
			query.MccCd18 = in.Item.MccCd18
			query.ApplDate = in.Item.ApplDate
			query.UpBcCd = in.Item.UpBcCd
			query.UpAcCd = in.Item.UpAcCd
			query.UpMccCd = in.Item.UpMccCd
			query.Name = in.Item.Name
			query.NameBusi = in.Item.NameBusi
			query.BusiLiceNo = in.Item.BusiLiceNo
			query.BusiRang = in.Item.BusiRang
			query.BusiMain = in.Item.BusiMain
			query.Certif = in.Item.Certif
			query.CertifType = in.Item.CertifType
			query.CertifNo = in.Item.CertifNo
			query.ProvCd = in.Item.ProvCd
			query.CityCd = in.Item.CityCd
			query.AreaCd = in.Item.AreaCd
			query.RegAddr = in.Item.RegAddr
			query.ContactName = in.Item.ContactName
			query.ContactPhoneNo = in.Item.ContactPhoneNo
			query.IsGroup = in.Item.IsGroup
			query.MoneyToGroup = in.Item.MoneyToGroup
			query.StlmWay = in.Item.StlmWay
			query.StlmWayDesc = in.Item.StlmWayDesc
			query.StlmInsCircle = in.Item.StlmInsCircle
			query.Status = in.Item.Status
			query.UcBcCd32 = in.Item.UcBcCd32
			query.K2WorkflowId = in.Item.K2WorkflowId
			query.SystemFlag = in.Item.SystemFlag
			query.ApprovalUsername = in.Item.ApprovalUsername
			query.FinalApprovalUsername = in.Item.FinalApprovalUsername
			query.IsUpStandard = in.Item.IsUpStandard
			query.BillingType = in.Item.BillingType
			query.BillingLevel = in.Item.BillingLevel
			query.Slogan = in.Item.Slogan
			query.Ext1 = in.Item.Ext1
			query.Ext2 = in.Item.Ext2
			query.Ext3 = in.Item.Ext3
			query.Ext4 = in.Item.Ext4
			query.AreaStandard = in.Item.AreaStandard
			query.MchtCdAreaCd = in.Item.MchtCdAreaCd
			query.UcBcCdArea = in.Item.UcBcCdArea
			query.RecOprId = in.Item.RecOprId
			query.RecUpdOpr = in.Item.RecUpdOpr
			query.OperIn = in.Item.OperIn
			query.OemOrgCode = in.Item.OemOrgCode
			query.IsEleInvoice = in.Item.IsEleInvoice
			query.DutyParagraph = in.Item.DutyParagraph
			query.TaxMachineBrand = in.Item.TaxMachineBrand
			query.Ext5 = in.Item.Ext5
			query.Ext6 = in.Item.Ext6
			query.Ext7 = in.Item.Ext7
			query.Ext8 = in.Item.Ext8
			query.Ext9 = in.Item.Ext9
			query.BusiLiceSt = in.Item.BusiLiceSt
			query.BusiLiceDt = in.Item.BusiLiceDt
			query.CertifSt = in.Item.CertifSt
			query.CertifDt = in.Item.CertifDt
		}

		merchants, count, err := merchantmodel.QueryMerchantInfos(db, query, insIds, in.Page, in.Size)
		if err != nil {
			return nil, err
		}

		pbMerchants := make([]*pb.MerchantField, len(merchants))
		for i := range merchants {
			pbMerchants[i] = &pb.MerchantField{
				MchtCd:                merchants[i].MchtCd,
				Sn:                    merchants[i].Sn,
				AipBranCd:             merchants[i].AipBranCd,
				GroupCd:               merchants[i].GroupCd,
				OriChnl:               merchants[i].OriChnl,
				OriChnlDesc:           merchants[i].OriChnlDesc,
				BankBelongCd:          merchants[i].BankBelongCd,
				DvpBy:                 merchants[i].DvpBy,
				MccCd18:               merchants[i].MccCd18,
				ApplDate:              merchants[i].ApplDate,
				UpBcCd:                merchants[i].UpBcCd,
				UpAcCd:                merchants[i].UpAcCd,
				UpMccCd:               merchants[i].UpMccCd,
				Name:                  merchants[i].Name,
				NameBusi:              merchants[i].NameBusi,
				BusiLiceNo:            merchants[i].BusiLiceNo,
				BusiRang:              merchants[i].BusiRang,
				BusiMain:              merchants[i].BusiMain,
				Certif:                merchants[i].Certif,
				CertifType:            merchants[i].CertifType,
				CertifNo:              merchants[i].CertifNo,
				ProvCd:                merchants[i].ProvCd,
				CityCd:                merchants[i].CityCd,
				AreaCd:                merchants[i].AreaCd,
				RegAddr:               merchants[i].RegAddr,
				ContactName:           merchants[i].ContactName,
				ContactPhoneNo:        merchants[i].ContactPhoneNo,
				IsGroup:               merchants[i].IsGroup,
				MoneyToGroup:          merchants[i].MoneyToGroup,
				StlmWay:               merchants[i].StlmWay,
				StlmWayDesc:           merchants[i].StlmWayDesc,
				StlmInsCircle:         merchants[i].StlmInsCircle,
				Status:                merchants[i].Status,
				UcBcCd32:              merchants[i].UcBcCd32,
				K2WorkflowId:          merchants[i].K2WorkflowId,
				SystemFlag:            merchants[i].SystemFlag,
				ApprovalUsername:      merchants[i].ApprovalUsername,
				FinalApprovalUsername: merchants[i].FinalApprovalUsername,
				IsUpStandard:          merchants[i].IsUpStandard,
				BillingType:           merchants[i].BillingType,
				BillingLevel:          merchants[i].BillingLevel,
				Slogan:                merchants[i].Slogan,
				Ext1:                  merchants[i].Ext1,
				Ext2:                  merchants[i].Ext2,
				Ext3:                  merchants[i].Ext3,
				Ext4:                  merchants[i].Ext4,
				AreaStandard:          merchants[i].AreaStandard,
				MchtCdAreaCd:          merchants[i].MchtCdAreaCd,
				UcBcCdArea:            merchants[i].UcBcCdArea,
				RecOprId:              merchants[i].RecOprId,
				RecUpdOpr:             merchants[i].RecUpdOpr,
				OperIn:                merchants[i].OperIn,
				OemOrgCode:            merchants[i].OemOrgCode,
				IsEleInvoice:          merchants[i].IsEleInvoice,
				DutyParagraph:         merchants[i].DutyParagraph,
				TaxMachineBrand:       merchants[i].TaxMachineBrand,
				Ext5:                  merchants[i].Ext5,
				Ext6:                  merchants[i].Ext6,
				Ext7:                  merchants[i].Ext7,
				Ext8:                  merchants[i].Ext8,
				Ext9:                  merchants[i].Ext9,
				BusiLiceSt:            merchants[i].BusiLiceSt,
				BusiLiceDt:            merchants[i].BusiLiceDt,
				CertifSt:              merchants[i].CertifSt,
				CertifDt:              merchants[i].CertifDt,
			}
			if merchants[i].ApprDate.Valid {
				pbMerchants[i].ApprDate = merchants[i].ApprDate.Time.Format(util.TimePattern)
			}
			if merchants[i].DeleteDate.Valid {
				pbMerchants[i].DeleteDate = merchants[i].DeleteDate.Time.Format(util.TimePattern)
			}
			if !merchants[i].CreatedAt.IsZero() {
				pbMerchants[i].CreatedAt = merchants[i].CreatedAt.Format(util.TimePattern)
			}
			if !merchants[i].UpdatedAt.IsZero() {
				pbMerchants[i].UpdatedAt = merchants[i].UpdatedAt.Format(util.TimePattern)
			}
			if merchants[i].RecApllyTs.Valid {
				pbMerchants[i].RecApllyTs = merchants[i].RecApllyTs.Time.Format(util.TimePattern)
			}
		}

		return &pb.ListMerchantReply{
			Items: pbMerchants,
			Count: count,
			Page:  in.Page,
			Size:  in.Size,
		}, nil
	} else {
		query := new(merchantmodel.MerchantInfoMain)
		if in.Item != nil {
			query.MchtCd = in.Item.MchtCd
			query.Sn = in.Item.Sn
			if user.UserType == "institution" {
				query.AipBranCd = user.UserGroupNo
			}
			if in.Item.AipBranCd != "" {
				query.AipBranCd = in.Item.AipBranCd
			}
			query.GroupCd = in.Item.GroupCd
			query.OriChnl = in.Item.OriChnl
			query.OriChnlDesc = in.Item.OriChnlDesc
			query.BankBelongCd = in.Item.BankBelongCd
			query.DvpBy = in.Item.DvpBy
			query.MccCd18 = in.Item.MccCd18
			query.ApplDate = in.Item.ApplDate
			query.UpBcCd = in.Item.UpBcCd
			query.UpAcCd = in.Item.UpAcCd
			query.UpMccCd = in.Item.UpMccCd
			query.Name = in.Item.Name
			query.NameBusi = in.Item.NameBusi
			query.BusiLiceNo = in.Item.BusiLiceNo
			query.BusiRang = in.Item.BusiRang
			query.BusiMain = in.Item.BusiMain
			query.Certif = in.Item.Certif
			query.CertifType = in.Item.CertifType
			query.CertifNo = in.Item.CertifNo
			query.ProvCd = in.Item.ProvCd
			query.CityCd = in.Item.CityCd
			query.AreaCd = in.Item.AreaCd
			query.RegAddr = in.Item.RegAddr
			query.ContactName = in.Item.ContactName
			query.ContactPhoneNo = in.Item.ContactPhoneNo
			query.IsGroup = in.Item.IsGroup
			query.MoneyToGroup = in.Item.MoneyToGroup
			query.StlmWay = in.Item.StlmWay
			query.StlmWayDesc = in.Item.StlmWayDesc
			query.StlmInsCircle = in.Item.StlmInsCircle
			query.Status = in.Item.Status
			query.UcBcCd32 = in.Item.UcBcCd32
			query.K2WorkflowId = in.Item.K2WorkflowId
			query.SystemFlag = in.Item.SystemFlag
			query.ApprovalUsername = in.Item.ApprovalUsername
			query.FinalApprovalUsername = in.Item.FinalApprovalUsername
			query.IsUpStandard = in.Item.IsUpStandard
			query.BillingType = in.Item.BillingType
			query.BillingLevel = in.Item.BillingLevel
			query.Slogan = in.Item.Slogan
			query.Ext1 = in.Item.Ext1
			query.Ext2 = in.Item.Ext2
			query.Ext3 = in.Item.Ext3
			query.Ext4 = in.Item.Ext4
			query.AreaStandard = in.Item.AreaStandard
			query.MchtCdAreaCd = in.Item.MchtCdAreaCd
			query.UcBcCdArea = in.Item.UcBcCdArea
			query.RecOprId = in.Item.RecOprId
			query.RecUpdOpr = in.Item.RecUpdOpr
			query.OperIn = in.Item.OperIn
			query.OemOrgCode = in.Item.OemOrgCode
			query.IsEleInvoice = in.Item.IsEleInvoice
			query.DutyParagraph = in.Item.DutyParagraph
			query.TaxMachineBrand = in.Item.TaxMachineBrand
			query.Ext5 = in.Item.Ext5
			query.Ext6 = in.Item.Ext6
			query.Ext7 = in.Item.Ext7
			query.Ext8 = in.Item.Ext8
			query.Ext9 = in.Item.Ext9
			query.BusiLiceSt = in.Item.BusiLiceSt
			query.BusiLiceDt = in.Item.BusiLiceDt
			query.CertifSt = in.Item.CertifSt
			query.CertifDt = in.Item.CertifDt
		}

		merchants, count, err := merchantmodel.QueryMerchantInfosMain(db, query, insIds, in.Page, in.Size)
		if err != nil {
			return nil, err
		}

		pbMerchants := make([]*pb.MerchantField, len(merchants))
		for i := range merchants {
			pbMerchants[i] = &pb.MerchantField{
				MchtCd:                merchants[i].MchtCd,
				Sn:                    merchants[i].Sn,
				AipBranCd:             merchants[i].AipBranCd,
				GroupCd:               merchants[i].GroupCd,
				OriChnl:               merchants[i].OriChnl,
				OriChnlDesc:           merchants[i].OriChnlDesc,
				BankBelongCd:          merchants[i].BankBelongCd,
				DvpBy:                 merchants[i].DvpBy,
				MccCd18:               merchants[i].MccCd18,
				ApplDate:              merchants[i].ApplDate,
				UpBcCd:                merchants[i].UpBcCd,
				UpAcCd:                merchants[i].UpAcCd,
				UpMccCd:               merchants[i].UpMccCd,
				Name:                  merchants[i].Name,
				NameBusi:              merchants[i].NameBusi,
				BusiLiceNo:            merchants[i].BusiLiceNo,
				BusiRang:              merchants[i].BusiRang,
				BusiMain:              merchants[i].BusiMain,
				Certif:                merchants[i].Certif,
				CertifType:            merchants[i].CertifType,
				CertifNo:              merchants[i].CertifNo,
				ProvCd:                merchants[i].ProvCd,
				CityCd:                merchants[i].CityCd,
				AreaCd:                merchants[i].AreaCd,
				RegAddr:               merchants[i].RegAddr,
				ContactName:           merchants[i].ContactName,
				ContactPhoneNo:        merchants[i].ContactPhoneNo,
				IsGroup:               merchants[i].IsGroup,
				MoneyToGroup:          merchants[i].MoneyToGroup,
				StlmWay:               merchants[i].StlmWay,
				StlmWayDesc:           merchants[i].StlmWayDesc,
				StlmInsCircle:         merchants[i].StlmInsCircle,
				Status:                merchants[i].Status,
				UcBcCd32:              merchants[i].UcBcCd32,
				K2WorkflowId:          merchants[i].K2WorkflowId,
				SystemFlag:            merchants[i].SystemFlag,
				ApprovalUsername:      merchants[i].ApprovalUsername,
				FinalApprovalUsername: merchants[i].FinalApprovalUsername,
				IsUpStandard:          merchants[i].IsUpStandard,
				BillingType:           merchants[i].BillingType,
				BillingLevel:          merchants[i].BillingLevel,
				Slogan:                merchants[i].Slogan,
				Ext1:                  merchants[i].Ext1,
				Ext2:                  merchants[i].Ext2,
				Ext3:                  merchants[i].Ext3,
				Ext4:                  merchants[i].Ext4,
				AreaStandard:          merchants[i].AreaStandard,
				MchtCdAreaCd:          merchants[i].MchtCdAreaCd,
				UcBcCdArea:            merchants[i].UcBcCdArea,
				RecOprId:              merchants[i].RecOprId,
				RecUpdOpr:             merchants[i].RecUpdOpr,
				OperIn:                merchants[i].OperIn,
				OemOrgCode:            merchants[i].OemOrgCode,
				IsEleInvoice:          merchants[i].IsEleInvoice,
				DutyParagraph:         merchants[i].DutyParagraph,
				TaxMachineBrand:       merchants[i].TaxMachineBrand,
				Ext5:                  merchants[i].Ext5,
				Ext6:                  merchants[i].Ext6,
				Ext7:                  merchants[i].Ext7,
				Ext8:                  merchants[i].Ext8,
				Ext9:                  merchants[i].Ext9,
				BusiLiceSt:            merchants[i].BusiLiceSt,
				BusiLiceDt:            merchants[i].BusiLiceDt,
				CertifSt:              merchants[i].CertifSt,
				CertifDt:              merchants[i].CertifDt,
			}
			if merchants[i].ApprDate.Valid {
				pbMerchants[i].ApprDate = merchants[i].ApprDate.Time.Format(util.TimePattern)
			}
			if merchants[i].DeleteDate.Valid {
				pbMerchants[i].DeleteDate = merchants[i].DeleteDate.Time.Format(util.TimePattern)
			}
			if !merchants[i].CreatedAt.IsZero() {
				pbMerchants[i].CreatedAt = merchants[i].CreatedAt.Format(util.TimePattern)
			}
			if !merchants[i].UpdatedAt.IsZero() {
				pbMerchants[i].UpdatedAt = merchants[i].UpdatedAt.Format(util.TimePattern)
			}
			if merchants[i].RecApllyTs.Valid {
				pbMerchants[i].RecApllyTs = merchants[i].RecApllyTs.Time.Format(util.TimePattern)
			}
		}

		return &pb.ListMerchantReply{
			Items: pbMerchants,
			Count: count,
			Page:  in.Page,
			Size:  in.Size,
		}, nil

	}
}

func (m *merchantService) ListGroupMerchant(ctx context.Context, in *pb.ListGroupMerchantRequest) (*pb.ListGroupMerchantReply, error) {
	if in.Size == 0 {
		in.Size = 10
	}
	if in.Page == 0 {
		in.Page = 1
	}
	db := common.DB

	edit := true
	if in.Type == "main" {
		edit = false
	}

	if edit {
		query := new(merchantmodel.Group)
		if in.Item != nil {
			query.JtMchtCd = in.Item.JtMchtCd
			query.JtMchtNm = in.Item.JtMchtNm
			query.JtArea = in.Item.JtArea
			query.MchtStlmCNm = in.Item.MchtStlmCNm
			query.MchtStlmCAcct = in.Item.MchtStlmCAcct
			query.ChtStlmInsIdCd = in.Item.ChtStlmInsIdCd
			query.MchtStlmInsNm = in.Item.MchtStlmInsNm
			query.MchtPaySysAcct = in.Item.MchtPaySysAcct
			query.ProvCd = in.Item.ProvCd
			query.CityCd = in.Item.CityCd
			query.AipBranCd = in.Item.AipBranCd
			query.SystemFlag = in.Item.SystemFlag
			query.Status = in.Item.Status
			query.RecOprId = in.Item.RecOprId
			query.RecUpdOpr = in.Item.RecUpdOpr
			query.JtAddr = in.Item.JtAddr
		}

		groups, count, err := merchantmodel.QueryMerchantGroups(db, query, in.Page, in.Size)
		if err != nil {
			return nil, err
		}

		pbGroups := make([]*pb.MerchantGroupField, len(groups))

		for i := range groups {
			pbGroups[i] = &pb.MerchantGroupField{
				JtMchtCd:       groups[i].JtMchtCd,
				JtMchtNm:       groups[i].JtMchtNm,
				JtArea:         groups[i].JtArea,
				MchtStlmCNm:    groups[i].MchtStlmCNm,
				MchtStlmCAcct:  groups[i].MchtStlmCAcct,
				ChtStlmInsIdCd: groups[i].ChtStlmInsIdCd,
				MchtStlmInsNm:  groups[i].MchtStlmInsNm,
				MchtPaySysAcct: groups[i].MchtPaySysAcct,
				ProvCd:         groups[i].ProvCd,
				CityCd:         groups[i].CityCd,
				AipBranCd:      groups[i].AipBranCd,
				SystemFlag:     groups[i].SystemFlag,
				Status:         groups[i].Status,
				RecOprId:       groups[i].RecOprId,
				RecUpdOpr:      groups[i].RecUpdOpr,
				JtAddr:         groups[i].JtAddr,
			}
			if !groups[i].CreatedAt.IsZero() {
				pbGroups[i].CreatedAt = groups[i].CreatedAt.Format(util.TimePattern)
			}
			if !groups[i].UpdatedAt.IsZero() {
				pbGroups[i].UpdatedAt = groups[i].UpdatedAt.Format(util.TimePattern)
			}
		}
		return &pb.ListGroupMerchantReply{
			Items: pbGroups,
			Count: count,
			Page:  in.Page,
			Size:  in.Size,
		}, nil
	} else {
		query := new(merchantmodel.GroupMain)
		if in.Item != nil {
			query.JtMchtCd = in.Item.JtMchtCd
			query.JtMchtNm = in.Item.JtMchtNm
			query.JtArea = in.Item.JtArea
			query.MchtStlmCNm = in.Item.MchtStlmCNm
			query.MchtStlmCAcct = in.Item.MchtStlmCAcct
			query.ChtStlmInsIdCd = in.Item.ChtStlmInsIdCd
			query.MchtStlmInsNm = in.Item.MchtStlmInsNm
			query.MchtPaySysAcct = in.Item.MchtPaySysAcct
			query.ProvCd = in.Item.ProvCd
			query.CityCd = in.Item.CityCd
			query.AipBranCd = in.Item.AipBranCd
			query.SystemFlag = in.Item.SystemFlag
			query.Status = in.Item.Status
			query.RecOprId = in.Item.RecOprId
			query.RecUpdOpr = in.Item.RecUpdOpr
			query.JtAddr = in.Item.JtAddr
		}

		groups, count, err := merchantmodel.QueryMerchantGroupsMain(db, query, in.Page, in.Size)
		if err != nil {
			return nil, err
		}

		pbGroups := make([]*pb.MerchantGroupField, len(groups))

		for i := range groups {
			pbGroups[i] = &pb.MerchantGroupField{
				JtMchtCd:       groups[i].JtMchtCd,
				JtMchtNm:       groups[i].JtMchtNm,
				JtArea:         groups[i].JtArea,
				MchtStlmCNm:    groups[i].MchtStlmCNm,
				MchtStlmCAcct:  groups[i].MchtStlmCAcct,
				ChtStlmInsIdCd: groups[i].ChtStlmInsIdCd,
				MchtStlmInsNm:  groups[i].MchtStlmInsNm,
				MchtPaySysAcct: groups[i].MchtPaySysAcct,
				ProvCd:         groups[i].ProvCd,
				CityCd:         groups[i].CityCd,
				AipBranCd:      groups[i].AipBranCd,
				SystemFlag:     groups[i].SystemFlag,
				Status:         groups[i].Status,
				RecOprId:       groups[i].RecOprId,
				RecUpdOpr:      groups[i].RecUpdOpr,
				JtAddr:         groups[i].JtAddr,
			}
			if !groups[i].CreatedAt.IsZero() {
				pbGroups[i].CreatedAt = groups[i].CreatedAt.Format(util.TimePattern)
			}
			if !groups[i].UpdatedAt.IsZero() {
				pbGroups[i].UpdatedAt = groups[i].UpdatedAt.Format(util.TimePattern)
			}
		}
		return &pb.ListGroupMerchantReply{
			Items: pbGroups,
			Count: count,
			Page:  in.Page,
			Size:  in.Size,
		}, nil
	}
}

func (m *merchantService) SaveMerchantBankAccount(ctx context.Context, in *pb.SaveMerchantBankAccountRequest) (*pb.SaveMerchantBankAccountReply, error) {
	var reply pb.SaveMerchantBankAccountReply
	if in.Item == nil {
		reply.Err = &pb.Error{
			Code:        http.StatusBadRequest,
			Message:     "InvalidParamsError",
			Description: "保存信息为空",
		}
		return &reply, nil
	}
	db := common.DB

	data := new(merchantmodel.BankAccount)
	{
		data.OwnerCd = in.Item.OwnerCd
		data.AccountType = in.Item.AccountType
		data.Name = in.Item.Name
		data.Account = in.Item.Account
		data.UcBcCd = in.Item.UcBcCd
		data.Province = in.Item.Province
		data.City = in.Item.City
		data.BankCode = in.Item.BankCode
		data.BankName = in.Item.BankName
		data.OperIn = in.Item.OperIn
		data.RecOprId = in.Item.RecOprId
		data.RecUpdOpr = in.Item.RecUpdOpr
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
	}
	err := merchantmodel.SaveBankAccount(db, data)
	return &reply, err
}

func (m *merchantService) SaveMerchantBusiness(ctx context.Context, in *pb.SaveMerchantBusinessRequest) (*pb.SaveMerchantBusinessReply, error) {
	var reply pb.SaveMerchantBusinessReply
	if in.Item == nil {
		reply.Err = &pb.Error{
			Code:        http.StatusBadRequest,
			Message:     "InvalidParamsError",
			Description: "保存信息为空",
		}
		return &reply, nil
	}
	db := common.DB

	data := new(merchantmodel.Business)
	{
		data.MchtCd = in.Item.MchtCd
		data.ProdCd = in.Item.ProdCd
		data.ProdCdText = in.Item.ProdCdText
		data.FeeMoneyCd = in.Item.FeeMoneyCd
		data.FeeModeType = in.Item.FeeModeType
		data.FeeSettlementType = in.Item.FeeSettlementType
		data.FeeHoliday = in.Item.FeeHoliday
		data.ServiceFeeType = in.Item.ServiceFeeType
		data.ServiceFeeStaticAmount, _ = strconv.ParseFloat(in.Item.ServiceFeeStaticAmount, 64)
		data.ServiceFeeLevelCount, _ = strconv.ParseInt(in.Item.ServiceFeeLevelCount, 10, 64)
		data.ServiceFeeMode = in.Item.ServiceFeeMode
		data.ServiceFeeUnit = in.Item.ServiceFeeUnit
		data.ServiceFeeTerm = in.Item.ServiceFeeTerm
		data.ServiceFeeSumto = in.Item.ServiceFeeSumto
		data.ServiceFeeCircle = in.Item.ServiceFeeCircle
		data.ServiceFeeOthers = in.Item.ServiceFeeOthers
		data.ServiceFeeStart = in.Item.ServiceFeeStart
		data.ServiceFeeClct = in.Item.ServiceFeeClct
		data.ServiceFeeClctOthers = in.Item.ServiceFeeClctOthers
		data.SystemFlag = in.Item.SystemFlag
		data.Ext1 = in.Item.Ext1
		data.Ext2 = in.Item.Ext2
		data.Ext3 = in.Item.Ext3
		data.Ext4 = in.Item.Ext4
		data.ServiceFeeYesNo = in.Item.ServiceFeeYesNo
		data.RecOprId = in.Item.RecOprId
		data.RecUpdOpr = in.Item.RecUpdOpr
		data.OperIn = in.Item.OperIn
	}
	err := merchantmodel.SaveBusiness(db, data)
	return &reply, err
}

func (m *merchantService) SaveMerchantPicture(ctx context.Context, in *pb.SaveMerchantPictureRequest) (*pb.SaveMerchantPictureReply, error) {
	var reply pb.SaveMerchantPictureReply
	if in.Item == nil {
		reply.Err = &pb.Error{
			Code:        http.StatusBadRequest,
			Message:     "InvalidParamsError",
			Description: "保存信息为空",
		}
		return &reply, nil
	}
	db := common.DB

	data := new(merchantmodel.Picture)
	{
		data.FileId = in.Item.FileId
		data.MchtCd = in.Item.MchtCd
		data.DocType = in.Item.DocType
		data.FileType = in.Item.FileType
		data.FileName = in.Item.FileName
		data.PIndex = in.Item.PIndex
		data.PCode = in.Item.PCode
		data.Url = in.Item.Url
		data.SystemFlag = in.Item.SystemFlag
		data.Status = in.Item.Status
		data.RecOprId = in.Item.RecOprId
		data.RecUpdOpr = in.Item.RecUpdOpr
	}
	err := merchantmodel.SavePicture(db, data)
	return &reply, err
}
