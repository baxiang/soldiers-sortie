package institutionservice

import (
	"context"
	"fmt"
	"strconv"
	"userService/pkg/common"
	cleartxnM "userService/pkg/model/cleartxn"
	insmodel "userService/pkg/model/institution"
	"userService/pkg/pb"
	"userService/pkg/util"

	"net/http"

	"github.com/jinzhu/copier"
)

type service struct{}

func (s *service) RemoveBindGroup(ctx context.Context, in *pb.RemoveBindGroupRequest) (*pb.RemoveBindGroupReply, error) {
	reply := new(pb.RemoveBindGroupReply)
	if in.Item == nil || in.Item.GroupId == "" || in.Item.InsIdCd == "" {
		reply.Err = &pb.Error{
			Code:        http.StatusBadRequest,
			Message:     InvalidParam,
			Description: "删除信息不能为空",
		}
		return reply, nil
	}
	data := new(insmodel.InsGroupBind)
	data.InsIdCd = in.Item.InsIdCd
	data.GroupId, _ = strconv.ParseInt(in.Item.GroupId, 10, 64)
	return reply, insmodel.RemoveInsGroupBind(common.DB, data)
}

func (s *service) ListBindGroup(ctx context.Context, in *pb.ListBindGroupRequest) (*pb.ListBindGroupReply, error) {
	reply := new(pb.ListBindGroupReply)
	if in.Item == nil {
		reply.Err = &pb.Error{
			Code:        http.StatusBadRequest,
			Message:     InvalidParam,
			Description: "绑定信息不能为空",
		}
		return reply, nil
	}

	groupId, _ := strconv.ParseInt(in.Item.GroupId, 10, 64)
	items, err := insmodel.ListInsGroupBind(common.DB, groupId)
	if err != nil {
		return nil, err
	}

	pbItems := make([]*pb.BindGroupField, 0, len(items))
	for _, item := range items {
		data := &pb.BindGroupField{
			GroupId: fmt.Sprintf("%d", item.GroupId),
			InsIdCd: item.InsIdCd,
		}
		if !item.CreatedAt.IsZero() {
			data.CreatedAt = item.CreatedAt.Format(util.TimePattern)
		}

		if !item.UpdatedAt.IsZero() {
			data.UpdatedAt = item.UpdatedAt.Format(util.TimePattern)
		}
		pbItems = append(pbItems, data)
	}

	reply.Items = pbItems

	return reply, err
}

func (s *service) BindGroup(ctx context.Context, in *pb.BindGroupRequest) (*pb.BindGroupReply, error) {
	reply := new(pb.BindGroupReply)
	if in.Item == nil || in.Item.GroupId == "" || in.Item.InsIdCd == "" {
		reply.Err = &pb.Error{
			Code:        http.StatusBadRequest,
			Message:     InvalidParam,
			Description: "绑定信息不能为空",
		}
		return reply, nil
	}

	group := new(insmodel.InsGroupBind)
	group.GroupId, _ = strconv.ParseInt(in.Item.GroupId, 10, 64)
	group.InsIdCd = in.Item.InsIdCd
	err := insmodel.SaveInsGroupBind(common.DB, group)
	return reply, err
}

func (s *service) SaveGroup(ctx context.Context, in *pb.SaveGroupRequest) (*pb.SaveGroupReply, error) {
	reply := new(pb.SaveGroupReply)
	if in.Item == nil || in.Item.Name == "" {
		reply.Err = &pb.Error{
			Code:        http.StatusBadRequest,
			Message:     InvalidParam,
			Description: "机构组数据不能为空",
		}
		return reply, nil
	}

	group := new(insmodel.Group)
	group.GroupId, _ = strconv.ParseInt(in.Item.GroupId, 10, 64)
	group.Name = in.Item.Name
	err := insmodel.SaveGroup(common.DB, group)
	return reply, err
}

func (s *service) GetInstitutionCash(ctx context.Context, in *pb.GetInstitutionCashRequest) (*pb.GetInstitutionCashReply, error) {
	reply := new(pb.GetInstitutionCashReply)

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
		query := new(insmodel.Cash)
		if in.Item != nil {
			query.InsIdCd = in.Item.InsIdCd
			query.ProdCd = in.Item.ProdCd
			query.InsDefaultFlag = in.Item.InsDefaultFlag
			query.InsDefaultCash, _ = strconv.ParseFloat(in.Item.InsDefaultCash, 64)
			query.InsCurrentCash, _ = strconv.ParseFloat(in.Item.InsCurrentCash, 64)
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
			query.RecOprId = in.Item.RecOprId
			query.RecUpdOpr = in.Item.RecUpdOpr
		}
		out, count, err := insmodel.FindInstitutionCash(db, query, in.Page, in.Size)
		if err != nil {
			return nil, err
		}

		items := make([]*pb.InstitutionCashField, 0, len(out))
		for _, data := range out {
			item := &pb.InstitutionCashField{
				InsIdCd:        data.InsIdCd,
				ProdCd:         data.ProdCd,
				InsDefaultFlag: data.InsDefaultFlag,
				InsDefaultCash: fmt.Sprintf("%f", data.InsDefaultCash),
				InsCurrentCash: fmt.Sprintf("%f", data.InsCurrentCash),
				MsgResvFld1:    data.MsgResvFld1,
				MsgResvFld2:    data.MsgResvFld2,
				MsgResvFld3:    data.MsgResvFld3,
				MsgResvFld4:    data.MsgResvFld4,
				MsgResvFld5:    data.MsgResvFld5,
				MsgResvFld6:    data.MsgResvFld6,
				MsgResvFld7:    data.MsgResvFld7,
				MsgResvFld8:    data.MsgResvFld8,
				MsgResvFld9:    data.MsgResvFld9,
				MsgResvFld10:   data.MsgResvFld10,
				RecOprId:       data.RecOprId,
				RecUpdOpr:      data.RecUpdOpr,
			}
			if !data.CreatedAt.IsZero() {
				item.CreatedAt = data.CreatedAt.Format(util.TimePattern)
			}
			if !data.UpdatedAt.IsZero() {
				item.UpdatedAt = data.UpdatedAt.Format(util.TimePattern)
			}
			items = append(items, item)
		}
		reply.Items = items
		reply.Count = count
		reply.Size = in.Size
		reply.Page = in.Page
	} else {
		query := new(insmodel.CashMain)
		if in.Item != nil {
			query.InsIdCd = in.Item.InsIdCd
			query.ProdCd = in.Item.ProdCd
			query.InsDefaultFlag = in.Item.InsDefaultFlag
			query.InsDefaultCash, _ = strconv.ParseFloat(in.Item.InsDefaultCash, 64)
			query.InsCurrentCash, _ = strconv.ParseFloat(in.Item.InsCurrentCash, 64)
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
			query.RecOprId = in.Item.RecOprId
			query.RecUpdOpr = in.Item.RecUpdOpr
		}
		out, count, err := insmodel.FindInstitutionCashMain(db, query, in.Page, in.Size)
		if err != nil {
			return nil, err
		}

		items := make([]*pb.InstitutionCashField, 0, len(out))
		for _, data := range out {
			item := &pb.InstitutionCashField{
				InsIdCd:        data.InsIdCd,
				ProdCd:         data.ProdCd,
				InsDefaultFlag: data.InsDefaultFlag,
				InsDefaultCash: fmt.Sprintf("%f", data.InsDefaultCash),
				InsCurrentCash: fmt.Sprintf("%f", data.InsCurrentCash),
				MsgResvFld1:    data.MsgResvFld1,
				MsgResvFld2:    data.MsgResvFld2,
				MsgResvFld3:    data.MsgResvFld3,
				MsgResvFld4:    data.MsgResvFld4,
				MsgResvFld5:    data.MsgResvFld5,
				MsgResvFld6:    data.MsgResvFld6,
				MsgResvFld7:    data.MsgResvFld7,
				MsgResvFld8:    data.MsgResvFld8,
				MsgResvFld9:    data.MsgResvFld9,
				MsgResvFld10:   data.MsgResvFld10,
				RecOprId:       data.RecOprId,
				RecUpdOpr:      data.RecUpdOpr,
			}
			if !data.CreatedAt.IsZero() {
				item.CreatedAt = data.CreatedAt.Format(util.TimePattern)
			}
			if !data.UpdatedAt.IsZero() {
				item.UpdatedAt = data.UpdatedAt.Format(util.TimePattern)
			}
			items = append(items, item)
		}
		reply.Items = items
		reply.Count = count
		reply.Size = in.Size
		reply.Page = in.Page
	}

	return reply, nil
}

func (s *service) GetInstitutionFee(ctx context.Context, in *pb.GetInstitutionFeeRequest) (*pb.GetInstitutionFeeReply, error) {
	reply := new(pb.GetInstitutionFeeReply)

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
		query := new(insmodel.Fee)
		if in.Item != nil {
			query.InsIdCd = in.Item.InsIdCd
			query.ProdCd = in.Item.ProdCd
			query.BizCd = in.Item.BizCd
			query.SubBizCd = in.Item.SubBizCd
			query.InsFeeBizCd = in.Item.InsFeeBizCd
			query.InsFeeCd = in.Item.InsFeeCd
			query.InsFeeTp = in.Item.InsFeeTp
			query.InsFeeParam = in.Item.InsFeeParam
			query.InsFeePercent, _ = strconv.ParseFloat(in.Item.InsFeePercent, 64)
			query.InsFeePct, _ = strconv.ParseFloat(in.Item.InsFeePct, 64)
			query.InsFeePctMin, _ = strconv.ParseFloat(in.Item.InsFeePctMin, 64)
			query.InsFeePctMax, _ = strconv.ParseFloat(in.Item.InsFeePctMax, 64)
			query.InsAFeeSame = in.Item.InsAFeeSame
			query.InsAFeeParam = in.Item.InsAFeeParam
			query.InsAFeePercent, _ = strconv.ParseFloat(in.Item.InsAFeePercent, 64)
			query.InsAFeePct, _ = strconv.ParseFloat(in.Item.InsAFeePct, 64)
			query.InsAFeePctMin, _ = strconv.ParseFloat(in.Item.InsAFeePctMin, 64)
			query.InsAFeePctMax, _ = strconv.ParseFloat(in.Item.InsAFeePctMax, 64)
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
			query.RecOprId = in.Item.RecOprId
			query.RecUpdOpr = in.Item.RecUpdOpr
		}
		out, count, err := insmodel.FindInstitutionFee(db, query, in.Page, in.Size)
		if err != nil {
			return nil, err
		}

		items := make([]*pb.InstitutionFeeField, 0, len(out))
		for _, data := range out {
			item := &pb.InstitutionFeeField{
				InsIdCd:        data.InsIdCd,
				ProdCd:         data.ProdCd,
				BizCd:          data.BizCd,
				SubBizCd:       data.SubBizCd,
				InsFeeBizCd:    data.InsFeeBizCd,
				InsFeeCd:       data.InsFeeCd,
				InsFeeTp:       data.InsFeeTp,
				InsFeeParam:    data.InsFeeParam,
				InsFeePercent:  fmt.Sprintf("%f", data.InsFeePercent),
				InsFeePct:      fmt.Sprintf("%f", data.InsFeePct),
				InsFeePctMin:   fmt.Sprintf("%f", data.InsFeePctMin),
				InsFeePctMax:   fmt.Sprintf("%f", data.InsFeePctMax),
				InsAFeePercent: fmt.Sprintf("%f", data.InsAFeePercent),
				InsAFeePct:     fmt.Sprintf("%f", data.InsAFeePct),
				InsAFeePctMin:  fmt.Sprintf("%f", data.InsAFeePctMin),
				InsAFeePctMax:  fmt.Sprintf("%f", data.InsAFeePctMax),
				InsAFeeSame:    data.InsAFeeSame,
				InsAFeeParam:   data.InsAFeeParam,
				MsgResvFld1:    data.MsgResvFld1,
				MsgResvFld2:    data.MsgResvFld2,
				MsgResvFld3:    data.MsgResvFld3,
				MsgResvFld4:    data.MsgResvFld4,
				MsgResvFld5:    data.MsgResvFld5,
				MsgResvFld6:    data.MsgResvFld6,
				MsgResvFld7:    data.MsgResvFld7,
				MsgResvFld8:    data.MsgResvFld8,
				MsgResvFld9:    data.MsgResvFld9,
				MsgResvFld10:   data.MsgResvFld10,
				RecOprId:       data.RecOprId,
				RecUpdOpr:      data.RecUpdOpr,
			}
			if !data.CreatedAt.IsZero() {
				item.CreatedAt = data.CreatedAt.Format(util.TimePattern)
			}
			if !data.UpdatedAt.IsZero() {
				item.UpdatedAt = data.UpdatedAt.Format(util.TimePattern)
			}
			items = append(items, item)
		}
		reply.Items = items
		reply.Count = count
		reply.Size = in.Size
		reply.Page = in.Page
	} else {
		query := new(insmodel.FeeMain)
		if in.Item != nil {
			query.InsIdCd = in.Item.InsIdCd
			query.ProdCd = in.Item.ProdCd
			query.BizCd = in.Item.BizCd
			query.SubBizCd = in.Item.SubBizCd
			query.InsFeeBizCd = in.Item.InsFeeBizCd
			query.InsFeeCd = in.Item.InsFeeCd
			query.InsFeeTp = in.Item.InsFeeTp
			query.InsFeeParam = in.Item.InsFeeParam
			query.InsFeePercent, _ = strconv.ParseFloat(in.Item.InsFeePercent, 64)
			query.InsFeePct, _ = strconv.ParseFloat(in.Item.InsFeePct, 64)
			query.InsFeePctMin, _ = strconv.ParseFloat(in.Item.InsFeePctMin, 64)
			query.InsFeePctMax, _ = strconv.ParseFloat(in.Item.InsFeePctMax, 64)
			query.InsAFeeSame = in.Item.InsAFeeSame
			query.InsAFeeParam = in.Item.InsAFeeParam
			query.InsAFeePercent, _ = strconv.ParseFloat(in.Item.InsAFeePercent, 64)
			query.InsAFeePct, _ = strconv.ParseFloat(in.Item.InsAFeePct, 64)
			query.InsAFeePctMin, _ = strconv.ParseFloat(in.Item.InsAFeePctMin, 64)
			query.InsAFeePctMax, _ = strconv.ParseFloat(in.Item.InsAFeePctMax, 64)
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
			query.RecOprId = in.Item.RecOprId
			query.RecUpdOpr = in.Item.RecUpdOpr
		}
		out, count, err := insmodel.FindInstitutionFeeMain(db, query, in.Page, in.Size)
		if err != nil {
			return nil, err
		}

		items := make([]*pb.InstitutionFeeField, 0, len(out))
		for _, data := range out {
			item := &pb.InstitutionFeeField{
				InsIdCd:        data.InsIdCd,
				ProdCd:         data.ProdCd,
				BizCd:          data.BizCd,
				SubBizCd:       data.SubBizCd,
				InsFeeBizCd:    data.InsFeeBizCd,
				InsFeeCd:       data.InsFeeCd,
				InsFeeTp:       data.InsFeeTp,
				InsFeeParam:    data.InsFeeParam,
				InsFeePercent:  fmt.Sprintf("%f", data.InsFeePercent),
				InsFeePct:      fmt.Sprintf("%f", data.InsFeePct),
				InsFeePctMin:   fmt.Sprintf("%f", data.InsFeePctMin),
				InsFeePctMax:   fmt.Sprintf("%f", data.InsFeePctMax),
				InsAFeePercent: fmt.Sprintf("%f", data.InsAFeePercent),
				InsAFeePct:     fmt.Sprintf("%f", data.InsAFeePct),
				InsAFeePctMin:  fmt.Sprintf("%f", data.InsAFeePctMin),
				InsAFeePctMax:  fmt.Sprintf("%f", data.InsAFeePctMax),
				InsAFeeSame:    data.InsAFeeSame,
				InsAFeeParam:   data.InsAFeeParam,
				MsgResvFld1:    data.MsgResvFld1,
				MsgResvFld2:    data.MsgResvFld2,
				MsgResvFld3:    data.MsgResvFld3,
				MsgResvFld4:    data.MsgResvFld4,
				MsgResvFld5:    data.MsgResvFld5,
				MsgResvFld6:    data.MsgResvFld6,
				MsgResvFld7:    data.MsgResvFld7,
				MsgResvFld8:    data.MsgResvFld8,
				MsgResvFld9:    data.MsgResvFld9,
				MsgResvFld10:   data.MsgResvFld10,
				RecOprId:       data.RecOprId,
				RecUpdOpr:      data.RecUpdOpr,
			}
			if !data.CreatedAt.IsZero() {
				item.CreatedAt = data.CreatedAt.Format(util.TimePattern)
			}
			if !data.UpdatedAt.IsZero() {
				item.UpdatedAt = data.UpdatedAt.Format(util.TimePattern)
			}
			items = append(items, item)
		}
		reply.Items = items
		reply.Count = count
		reply.Size = in.Size
		reply.Page = in.Page
	}

	return reply, nil
}

func (s *service) GetInstitutionControl(ctx context.Context, in *pb.GetInstitutionControlRequest) (*pb.GetInstitutionControlReply, error) {
	reply := new(pb.GetInstitutionControlReply)

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
		query := new(insmodel.Control)
		if in.Item != nil {
			query.InsIdCd = in.Item.InsIdCd
			query.InsCompanyCd = in.Item.InsCompanyCd
			query.ProdCd = in.Item.ProdCd
			query.BizCd = in.Item.BizCd
			query.CtrlSta = in.Item.CtrlSta
			query.InsBegTm = in.Item.InsBegTm
			query.InsEndTm = in.Item.InsEndTm
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
			query.RecOprId = in.Item.RecOprId
			query.RecUpdOpr = in.Item.RecUpdOpr
		}
		out, count, err := insmodel.FindInstitutionControl(db, query, in.Page, in.Size)
		if err != nil {
			return nil, err
		}

		items := make([]*pb.InstitutionControlField, 0, len(out))
		for _, data := range out {
			item := &pb.InstitutionControlField{

				InsIdCd:      data.InsIdCd,
				InsCompanyCd: data.InsCompanyCd,
				ProdCd:       data.ProdCd,
				BizCd:        data.BizCd,
				CtrlSta:      data.CtrlSta,
				InsBegTm:     data.InsBegTm,
				InsEndTm:     data.InsEndTm,
				MsgResvFld1:  data.MsgResvFld1,
				MsgResvFld2:  data.MsgResvFld2,
				MsgResvFld3:  data.MsgResvFld3,
				MsgResvFld4:  data.MsgResvFld4,
				MsgResvFld5:  data.MsgResvFld5,
				MsgResvFld6:  data.MsgResvFld6,
				MsgResvFld7:  data.MsgResvFld7,
				MsgResvFld8:  data.MsgResvFld8,
				MsgResvFld9:  data.MsgResvFld9,
				MsgResvFld10: data.MsgResvFld10,
				RecOprId:     data.RecOprId,
				RecUpdOpr:    data.RecUpdOpr,
			}
			if !data.CreatedAt.IsZero() {
				item.CreatedAt = data.CreatedAt.Format(util.TimePattern)
			}
			if !data.UpdatedAt.IsZero() {
				item.UpdatedAt = data.UpdatedAt.Format(util.TimePattern)
			}
			items = append(items, item)
		}
		reply.Items = items
		reply.Count = count
		reply.Size = in.Size
		reply.Page = in.Page
	} else {
		query := new(insmodel.ControlMain)
		if in.Item != nil {
			query.InsIdCd = in.Item.InsIdCd
			query.InsCompanyCd = in.Item.InsCompanyCd
			query.ProdCd = in.Item.ProdCd
			query.BizCd = in.Item.BizCd
			query.CtrlSta = in.Item.CtrlSta
			query.InsBegTm = in.Item.InsBegTm
			query.InsEndTm = in.Item.InsEndTm
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
			query.RecOprId = in.Item.RecOprId
			query.RecUpdOpr = in.Item.RecUpdOpr
		}
		out, count, err := insmodel.FindInstitutionControlMain(db, query, in.Page, in.Size)
		if err != nil {
			return nil, err
		}

		items := make([]*pb.InstitutionControlField, 0, len(out))
		for _, data := range out {
			item := &pb.InstitutionControlField{

				InsIdCd:      data.InsIdCd,
				InsCompanyCd: data.InsCompanyCd,
				ProdCd:       data.ProdCd,
				BizCd:        data.BizCd,
				CtrlSta:      data.CtrlSta,
				InsBegTm:     data.InsBegTm,
				InsEndTm:     data.InsEndTm,
				MsgResvFld1:  data.MsgResvFld1,
				MsgResvFld2:  data.MsgResvFld2,
				MsgResvFld3:  data.MsgResvFld3,
				MsgResvFld4:  data.MsgResvFld4,
				MsgResvFld5:  data.MsgResvFld5,
				MsgResvFld6:  data.MsgResvFld6,
				MsgResvFld7:  data.MsgResvFld7,
				MsgResvFld8:  data.MsgResvFld8,
				MsgResvFld9:  data.MsgResvFld9,
				MsgResvFld10: data.MsgResvFld10,
				RecOprId:     data.RecOprId,
				RecUpdOpr:    data.RecUpdOpr,
			}
			if !data.CreatedAt.IsZero() {
				item.CreatedAt = data.CreatedAt.Format(util.TimePattern)
			}
			if !data.UpdatedAt.IsZero() {
				item.UpdatedAt = data.UpdatedAt.Format(util.TimePattern)
			}
			items = append(items, item)
		}
		reply.Items = items
		reply.Count = count
		reply.Size = in.Size
		reply.Page = in.Page
	}

	return reply, nil
}

func (s *service) SaveInstitutionFeeControlCash(ctx context.Context, in *pb.SaveInstitutionFeeControlCashRequest) (*pb.SaveInstitutionFeeControlCashReply, error) {
	reply := new(pb.SaveInstitutionFeeControlCashReply)
	db := common.DB.Begin()
	defer db.Rollback()

	if in.Fees != nil {
		for _, fee := range in.Fees {
			ins := new(insmodel.Fee)
			{
				ins.InsIdCd = fee.InsIdCd
				ins.ProdCd = fee.ProdCd
				ins.BizCd = fee.BizCd
				ins.SubBizCd = fee.SubBizCd
				ins.InsFeeBizCd = fee.InsFeeBizCd
				ins.InsFeeCd = fee.InsFeeCd
				ins.InsFeeTp = fee.InsFeeTp
				ins.InsFeeParam = fee.InsFeeParam
				ins.InsFeePercent, _ = strconv.ParseFloat(fee.InsFeePercent, 64)
				ins.InsFeePct, _ = strconv.ParseFloat(fee.InsFeePct, 64)
				ins.InsFeePctMin, _ = strconv.ParseFloat(fee.InsFeePctMin, 64)
				ins.InsFeePctMax, _ = strconv.ParseFloat(fee.InsFeePctMax, 64)
				ins.InsAFeeSame = fee.InsAFeeSame
				ins.InsAFeeParam = fee.InsAFeeParam
				ins.InsAFeePercent, _ = strconv.ParseFloat(fee.InsAFeePercent, 64)
				ins.InsAFeePct, _ = strconv.ParseFloat(fee.InsAFeePct, 64)
				ins.InsAFeePctMin, _ = strconv.ParseFloat(fee.InsAFeePctMin, 64)
				ins.InsAFeePctMax, _ = strconv.ParseFloat(fee.InsAFeePctMax, 64)
				ins.MsgResvFld1 = fee.MsgResvFld1
				ins.MsgResvFld2 = fee.MsgResvFld2
				ins.MsgResvFld3 = fee.MsgResvFld3
				ins.MsgResvFld4 = fee.MsgResvFld4
				ins.MsgResvFld5 = fee.MsgResvFld5
				ins.MsgResvFld6 = fee.MsgResvFld6
				ins.MsgResvFld7 = fee.MsgResvFld7
				ins.MsgResvFld8 = fee.MsgResvFld8
				ins.MsgResvFld9 = fee.MsgResvFld9
				ins.MsgResvFld10 = fee.MsgResvFld10
				ins.RecOprId = fee.RecOprId
				ins.RecUpdOpr = fee.RecUpdOpr
			}
			err := insmodel.SaveInstitutionFee(db, ins)
			if err != nil {
				return nil, err
			}
		}
	}
	if in.Cashes != nil {
		for _, cash := range in.Cashes {
			ins := new(insmodel.Cash)
			{
				ins.InsIdCd = cash.InsIdCd
				ins.ProdCd = cash.ProdCd
				ins.InsDefaultFlag = cash.InsDefaultFlag
				ins.InsDefaultCash, _ = strconv.ParseFloat(cash.InsDefaultCash, 64)
				ins.InsCurrentCash, _ = strconv.ParseFloat(cash.InsCurrentCash, 64)
				ins.MsgResvFld1 = cash.MsgResvFld1
				ins.MsgResvFld2 = cash.MsgResvFld2
				ins.MsgResvFld3 = cash.MsgResvFld3
				ins.MsgResvFld4 = cash.MsgResvFld4
				ins.MsgResvFld5 = cash.MsgResvFld5
				ins.MsgResvFld6 = cash.MsgResvFld6
				ins.MsgResvFld7 = cash.MsgResvFld7
				ins.MsgResvFld8 = cash.MsgResvFld8
				ins.MsgResvFld9 = cash.MsgResvFld9
				ins.MsgResvFld10 = cash.MsgResvFld10
				ins.RecOprId = cash.RecOprId
				ins.RecUpdOpr = cash.RecUpdOpr
			}
			err := insmodel.SaveInstitutionCash(db, ins)
			if err != nil {
				return nil, err
			}
		}
	}
	if in.Controls != nil {
		for _, control := range in.Controls {
			ins := new(insmodel.Control)
			{
				ins.InsIdCd = control.InsIdCd
				ins.InsCompanyCd = control.InsCompanyCd
				ins.ProdCd = control.ProdCd
				ins.BizCd = control.BizCd
				ins.CtrlSta = control.CtrlSta
				ins.InsBegTm = control.InsBegTm
				ins.InsEndTm = control.InsEndTm
				ins.MsgResvFld1 = control.MsgResvFld1
				ins.MsgResvFld2 = control.MsgResvFld2
				ins.MsgResvFld3 = control.MsgResvFld3
				ins.MsgResvFld4 = control.MsgResvFld4
				ins.MsgResvFld5 = control.MsgResvFld5
				ins.MsgResvFld6 = control.MsgResvFld6
				ins.MsgResvFld7 = control.MsgResvFld7
				ins.MsgResvFld8 = control.MsgResvFld8
				ins.MsgResvFld9 = control.MsgResvFld9
				ins.MsgResvFld10 = control.MsgResvFld10
				ins.RecOprId = control.RecOprId
				ins.RecUpdOpr = control.RecUpdOpr
			}
			err := insmodel.SaveInstitutionControl(db, ins)
			if err != nil {
				return nil, err
			}
		}
	}

	db.Commit()
	return reply, nil
}

func (s *service) GetInstitutionById(ctx context.Context, in *pb.GetInstitutionByIdRequest) (*pb.GetInstitutionByIdReply, error) {
	reply := new(pb.GetInstitutionByIdReply)
	if in.Id == "" {
		reply.Err = &pb.Error{
			Code:        http.StatusBadRequest,
			Message:     InvalidParam,
			Description: "id不能为空",
		}
		return reply, nil
	}
	db := common.DB
	edit := true
	if in.Type == "main" {
		edit = false
	}

	if edit {
		info, err := insmodel.FindInstitutionInfoById(db, in.Id)
		if err != nil {
			return nil, err
		}
		if info != nil {
			item := new(pb.InstitutionField)
			item.InsIdCd = info.InsIdCd
			item.InsCompanyCd = info.InsCompanyCd
			item.InsType = info.InsType
			item.InsName = info.InsName
			item.InsProvCd = info.InsProvCd
			item.InsCityCd = info.InsCityCd
			item.InsRegionCd = info.InsRegionCd
			item.InsSta = info.InsSta
			item.InsStlmTp = info.InsStlmTp
			item.InsAloStlmCycle = info.InsAloStlmCycle
			item.InsAloStlmMd = info.InsAloStlmMd
			item.InsStlmCNm = info.InsStlmCNm
			item.InsStlmCAcct = info.InsStlmCAcct
			item.InsStlmCBkNo = info.InsStlmCBkNo
			item.InsStlmCBkNm = info.InsStlmCBkNm
			item.InsStlmDNm = info.InsStlmDNm
			item.InsStlmDAcct = info.InsStlmDAcct
			item.InsStlmDBkNo = info.InsStlmDBkNo
			item.InsStlmDBkNm = info.InsStlmDBkNm
			item.MsgResvFld1 = info.MsgResvFld1
			item.MsgResvFld2 = info.MsgResvFld2
			item.MsgResvFld3 = info.MsgResvFld3
			item.MsgResvFld4 = info.MsgResvFld4
			item.MsgResvFld5 = info.MsgResvFld5
			item.MsgResvFld6 = info.MsgResvFld6
			item.MsgResvFld7 = info.MsgResvFld7
			item.MsgResvFld8 = info.MsgResvFld8
			item.MsgResvFld9 = info.MsgResvFld9
			item.MsgResvFld10 = info.MsgResvFld10
			item.RecOprId = info.RecOprId
			item.RecUpdOpr = info.RecUpdOpr
			if !info.CreatedAt.IsZero() {
				item.CreatedAt = info.CreatedAt.Format(util.TimePattern)
			}
			if !info.UpdatedAt.IsZero() {
				item.UpdatedAt = info.UpdatedAt.Format(util.TimePattern)
			}
			reply.Item = item
		}
	} else {
		info, err := insmodel.FindInstitutionInfoMainById(db, in.Id)
		if err != nil {
			return nil, err
		}
		if info != nil {
			item := new(pb.InstitutionField)
			item.InsIdCd = info.InsIdCd
			item.InsCompanyCd = info.InsCompanyCd
			item.InsType = info.InsType
			item.InsName = info.InsName
			item.InsProvCd = info.InsProvCd
			item.InsCityCd = info.InsCityCd
			item.InsRegionCd = info.InsRegionCd
			item.InsSta = info.InsSta
			item.InsStlmTp = info.InsStlmTp
			item.InsAloStlmCycle = info.InsAloStlmCycle
			item.InsAloStlmMd = info.InsAloStlmMd
			item.InsStlmCNm = info.InsStlmCNm
			item.InsStlmCAcct = info.InsStlmCAcct
			item.InsStlmCBkNo = info.InsStlmCBkNo
			item.InsStlmCBkNm = info.InsStlmCBkNm
			item.InsStlmDNm = info.InsStlmDNm
			item.InsStlmDAcct = info.InsStlmDAcct
			item.InsStlmDBkNo = info.InsStlmDBkNo
			item.InsStlmDBkNm = info.InsStlmDBkNm
			item.MsgResvFld1 = info.MsgResvFld1
			item.MsgResvFld2 = info.MsgResvFld2
			item.MsgResvFld3 = info.MsgResvFld3
			item.MsgResvFld4 = info.MsgResvFld4
			item.MsgResvFld5 = info.MsgResvFld5
			item.MsgResvFld6 = info.MsgResvFld6
			item.MsgResvFld7 = info.MsgResvFld7
			item.MsgResvFld8 = info.MsgResvFld8
			item.MsgResvFld9 = info.MsgResvFld9
			item.MsgResvFld10 = info.MsgResvFld10
			item.RecOprId = info.RecOprId
			item.RecUpdOpr = info.RecUpdOpr
			if !info.CreatedAt.IsZero() {
				item.CreatedAt = info.CreatedAt.Format(util.TimePattern)
			}
			if !info.UpdatedAt.IsZero() {
				item.UpdatedAt = info.UpdatedAt.Format(util.TimePattern)
			}
			reply.Item = item
		}
	}

	return reply, nil
}

//Download .
func (s *service) TnxHisDownload(ctx context.Context, in *pb.InstitutionTnxHisDownloadReq) (*pb.InstitutionTnxHisDownloadResp, error) {
	if in.Name == "" {
		return nil, ErrDownloadFileNameEmpty
	}

	clearTxnEnty := cleartxnM.ClearTxn{}
	Txns, err := clearTxnEnty.GetWithTime(common.DB, in.StartTime, in.EndTime)
	if err != nil {
		return nil, err
	}

	fileDir, err := DownloadFileWithDay(Txns)
	if err != nil {
		return nil, err
	}

	err = Compress(fileDir, in.Name)
	if err != nil {
		return nil, err
	}

	return &pb.InstitutionTnxHisDownloadResp{Result: true}, nil
}

//GetTfrTrnLogs .
func (s *service) GetTfrTrnLogs(ctx context.Context, in *pb.GetTfrTrnLogsReq) (*pb.GetTfrTrnLogsResp, error) {
	var cond cleartxnM.TfrTrnLog

	cond.MchntCd = in.MchntCd
	cond.PriAcctNo = in.PriAcctNO
	cond.KeyRsp = in.KeyRsp
	cond.PriAcctNo = in.PriAcctNO
	cond.CardClass = in.CardClass
	cond.RoutInsIdCd = in.RoutIndustryInsIdCd
	cond.FwdInsIdCd = in.FwdInsIdCd
	cond.IssInsIdCd = in.IssInsIdCd
	cond.RespCd = in.RespCd
	cond.TermId = in.TermId
	cond.ProdCd = in.ProdCd
	cond.BizCd = in.BizCd
	cond.MaTransCd = in.MaTransCd
	cond.MsgResvFld2 = in.MsgResvFld2

	var accountRegion = ""
	if in.BeginAt != "" && in.EndAt != "" {
		accountRegion = fmt.Sprintf("'%s' < TRANS_DT AND TRANS_DT < '%s'", in.BeginAt, in.EndAt)
	} else if in.BeginAt != "" {
		accountRegion = fmt.Sprintf("'%s' < TRANS_DT", in.BeginAt)
	} else if in.EndAt != "" {
		accountRegion = fmt.Sprintf("TRANS_DT < '%s'", in.EndAt)
	}

	trfTrnLogsEnty := cleartxnM.TfrTrnLog{}
	results, count, total, err := trfTrnLogsEnty.GetWithLimit(common.DB, &cond, accountRegion, in.Limit, in.Page)
	if err != nil {
		return nil, err
	}

	var items []*pb.GetTfrTrnLogsItem
	for _, insTxn := range results {
		item := pb.GetTfrTrnLogsItem{
			KeyRsp:       insTxn.KeyRsp,
			MchntCd:      insTxn.MchntCd,
			CardAccptrNm: insTxn.CardAccptrNm,
			TransDt:      insTxn.TransDt,
			MaSettleDt:   insTxn.MaSettleDt,
			TransMt:      insTxn.TransMt,
			MaTransCd:    insTxn.MaTransCd,
			FwdInsIdCd:   insTxn.FwdInsIdCd,
			TransAt:      insTxn.TransAt,
			PriAcctNo:    insTxn.PriAcctNo,
			IssInsIdCd:   insTxn.IssInsIdCd,
			TermId:       insTxn.TermId,
			ProdCd:       insTxn.ProdCd,
			CardClass:    insTxn.CardClass,
			TransSt:      insTxn.TransSt,
			RespCd:       insTxn.RespCd,
		}
		items = append(items, &item)
	}

	return &pb.GetTfrTrnLogsResp{Items: items, Count: count, Total: total}, nil
}

//GetTfrTrnLog .
func (s *service) GetTfrTrnLog(ctx context.Context, in *pb.GetTfrTrnLogReq) (*pb.GetTfrTrnLogResp, error) {
	trfTrnLogsEnty := cleartxnM.TfrTrnLog{}
	resp := new(pb.GetTfrTrnLogResp)
	trfTrnLog, err := trfTrnLogsEnty.GetByKeyRsp(common.DB, in.KeyRsp)
	if err != nil {
		return nil, err
	}
	if trfTrnLog == nil {
		resp.Err = &pb.Error{
			Code:        http.StatusBadRequest,
			Message:     InvalidParam,
			Description: "输入的参数错误",
		}
		return resp, nil
	}
	copier.Copy(resp, trfTrnLog)
	return resp, nil
}

func (s *service) DownloadTfrTrnLogs(ctx context.Context, in *pb.DownloadTfrTrnLogsReq) (*pb.DownloadTfrTrnLogsResp, error) {
	if in.Name == "" {
		return nil, ErrDownloadFileNameEmpty
	}

	var cond cleartxnM.TfrTrnLog

	cond.MchntCd = in.MchntCd
	cond.PriAcctNo = in.PriAcctNO
	cond.KeyRsp = in.KeyRsp
	cond.PriAcctNo = in.PriAcctNO
	cond.CardClass = in.CardClass
	cond.RoutIndustryInsIdCd = in.RoutIndustryInsIdCd
	cond.FwdInsIdCd = in.FwdInsIdCd
	cond.IssInsIdCd = in.IssInsIdCd
	cond.RespCd = in.RespCd
	cond.TermId = in.TermId
	cond.ProdCd = in.ProdCd
	cond.BizCd = in.BizCd
	cond.MaTransCd = in.MaTransCd
	cond.MsgResvFld2 = in.MsgResvFld2

	var accountRegion = ""
	if in.BeginAt != "" && in.EndAt != "" {
		accountRegion = fmt.Sprintf("'%s' < TRANS_DT AND TRANS_DT < '%s'", in.BeginAt, in.EndAt)
	} else if in.BeginAt != "" {
		accountRegion = fmt.Sprintf("'%s' < TRANS_DT", in.BeginAt)
	} else if in.EndAt != "" {
		accountRegion = fmt.Sprintf("TRANS_DT < '%s'", in.EndAt)
	}

	trfTrnLogsEnty := cleartxnM.TfrTrnLog{}
	results, err := trfTrnLogsEnty.Get(common.DB, &cond, accountRegion)
	if err != nil {
		return nil, err
	}
	uid, err := DownloadTfrTrnLogs(results)
	if err != nil {
		return nil, err
	}

	err = Compress(uid, in.Name)
	if err != nil {
		return nil, err
	}
	return &pb.DownloadTfrTrnLogsResp{Code: true}, nil
}

func (s *service) ListGroups(ctx context.Context, in *pb.ListGroupsRequest) (*pb.ListGroupsReply, error) {
	reply := new(pb.ListGroupsReply)
	db := common.DB
	if in.Page == 0 {
		in.Page = 1
	}
	if in.Size == 0 {
		in.Size = 10
	}

	groups, count, err := insmodel.ListGroups(db, in.Page, in.Size)
	if err != nil {
		return nil, err
	}

	pbIns := make([]*pb.InstitutionGroupField, len(groups))
	for i := range groups {
		pbIns[i] = &pb.InstitutionGroupField{
			GroupId: fmt.Sprintf("%d", groups[i].GroupId),
			Name:    groups[i].Name,
		}
		if !groups[i].CreatedAt.IsZero() {
			pbIns[i].CreatedAt = groups[i].CreatedAt.Format(util.TimePattern)
		}

		if !groups[i].UpdatedAt.IsZero() {
			pbIns[i].UpdatedAt = groups[i].UpdatedAt.Format(util.TimePattern)
		}
	}

	reply.Items = pbIns
	reply.Count = count
	reply.Size = in.Size
	reply.Page = in.Page
	return reply, nil
}

func (s *service) ListInstitutions(ctx context.Context, in *pb.ListInstitutionsRequest) (*pb.ListInstitutionsReply, error) {
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
		query := new(insmodel.InstitutionInfo)
		if in.Item != nil {
			query.InsIdCd = in.Item.InsIdCd
			query.InsCompanyCd = in.Item.InsCompanyCd
			query.InsType = in.Item.InsType
			query.InsName = in.Item.InsName
			query.InsProvCd = in.Item.InsProvCd
			query.InsCityCd = in.Item.InsCityCd
			query.InsRegionCd = in.Item.InsRegionCd
			query.InsSta = in.Item.InsSta
			query.InsStlmTp = in.Item.InsStlmTp
			query.InsAloStlmCycle = in.Item.InsAloStlmCycle
			query.InsAloStlmMd = in.Item.InsAloStlmMd
			query.InsStlmCNm = in.Item.InsStlmCNm
			query.InsStlmCAcct = in.Item.InsStlmCAcct
			query.InsStlmCBkNo = in.Item.InsStlmCBkNo
			query.InsStlmCBkNm = in.Item.InsStlmCBkNm
			query.InsStlmDNm = in.Item.InsStlmDNm
			query.InsStlmDAcct = in.Item.InsStlmDAcct
			query.InsStlmDBkNo = in.Item.InsStlmDBkNo
			query.InsStlmDBkNm = in.Item.InsStlmDBkNm
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
			query.RecOprId = in.Item.RecOprId
			query.RecUpdOpr = in.Item.RecUpdOpr
		}
		ins, count, err := insmodel.QueryInstitutionInfo(db, query, in.Page, in.Size)
		if err != nil {
			return nil, err
		}

		pbIns := make([]*pb.InstitutionField, len(ins))

		for i := range ins {
			pbIns[i] = &pb.InstitutionField{
				InsIdCd:         ins[i].InsIdCd,
				InsCompanyCd:    ins[i].InsCompanyCd,
				InsType:         ins[i].InsType,
				InsName:         ins[i].InsName,
				InsProvCd:       ins[i].InsProvCd,
				InsCityCd:       ins[i].InsCityCd,
				InsRegionCd:     ins[i].InsRegionCd,
				InsSta:          ins[i].InsSta,
				InsStlmTp:       ins[i].InsStlmTp,
				InsAloStlmCycle: ins[i].InsAloStlmCycle,
				InsAloStlmMd:    ins[i].InsAloStlmMd,
				InsStlmCNm:      ins[i].InsStlmCNm,
				InsStlmCAcct:    ins[i].InsStlmCAcct,
				InsStlmCBkNo:    ins[i].InsStlmCBkNo,
				InsStlmCBkNm:    ins[i].InsStlmCBkNm,
				InsStlmDNm:      ins[i].InsStlmDNm,
				InsStlmDAcct:    ins[i].InsStlmDAcct,
				InsStlmDBkNo:    ins[i].InsStlmDBkNo,
				InsStlmDBkNm:    ins[i].InsStlmDBkNm,
				MsgResvFld1:     ins[i].MsgResvFld1,
				MsgResvFld2:     ins[i].MsgResvFld2,
				MsgResvFld3:     ins[i].MsgResvFld3,
				MsgResvFld4:     ins[i].MsgResvFld4,
				MsgResvFld5:     ins[i].MsgResvFld5,
				MsgResvFld6:     ins[i].MsgResvFld6,
				MsgResvFld7:     ins[i].MsgResvFld7,
				MsgResvFld8:     ins[i].MsgResvFld8,
				MsgResvFld9:     ins[i].MsgResvFld9,
				MsgResvFld10:    ins[i].MsgResvFld10,
				RecOprId:        ins[i].RecOprId,
				RecUpdOpr:       ins[i].RecUpdOpr,
			}
			if !ins[i].CreatedAt.IsZero() {
				pbIns[i].CreatedAt = ins[i].CreatedAt.Format(util.TimePattern)
			}

			if !ins[i].UpdatedAt.IsZero() {
				pbIns[i].UpdatedAt = ins[i].UpdatedAt.Format(util.TimePattern)
			}

		}

		return &pb.ListInstitutionsReply{
			Items: pbIns,
			Count: count,
			Page:  in.Page,
			Size:  in.Size,
		}, nil
	} else {
		query := new(insmodel.InstitutionInfoMain)
		if in.Item != nil {
			query.InsIdCd = in.Item.InsIdCd
			query.InsCompanyCd = in.Item.InsCompanyCd
			query.InsType = in.Item.InsType
			query.InsName = in.Item.InsName
			query.InsProvCd = in.Item.InsProvCd
			query.InsCityCd = in.Item.InsCityCd
			query.InsRegionCd = in.Item.InsRegionCd
			query.InsSta = in.Item.InsSta
			query.InsStlmTp = in.Item.InsStlmTp
			query.InsAloStlmCycle = in.Item.InsAloStlmCycle
			query.InsAloStlmMd = in.Item.InsAloStlmMd
			query.InsStlmCNm = in.Item.InsStlmCNm
			query.InsStlmCAcct = in.Item.InsStlmCAcct
			query.InsStlmCBkNo = in.Item.InsStlmCBkNo
			query.InsStlmCBkNm = in.Item.InsStlmCBkNm
			query.InsStlmDNm = in.Item.InsStlmDNm
			query.InsStlmDAcct = in.Item.InsStlmDAcct
			query.InsStlmDBkNo = in.Item.InsStlmDBkNo
			query.InsStlmDBkNm = in.Item.InsStlmDBkNm
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
			query.RecOprId = in.Item.RecOprId
			query.RecUpdOpr = in.Item.RecUpdOpr
		}
		ins, count, err := insmodel.QueryInstitutionInfoMain(db, query, in.Page, in.Size)
		if err != nil {
			return nil, err
		}

		pbIns := make([]*pb.InstitutionField, len(ins))

		for i := range ins {
			pbIns[i] = &pb.InstitutionField{
				InsIdCd:         ins[i].InsIdCd,
				InsCompanyCd:    ins[i].InsCompanyCd,
				InsType:         ins[i].InsType,
				InsName:         ins[i].InsName,
				InsProvCd:       ins[i].InsProvCd,
				InsCityCd:       ins[i].InsCityCd,
				InsRegionCd:     ins[i].InsRegionCd,
				InsSta:          ins[i].InsSta,
				InsStlmTp:       ins[i].InsStlmTp,
				InsAloStlmCycle: ins[i].InsAloStlmCycle,
				InsAloStlmMd:    ins[i].InsAloStlmMd,
				InsStlmCNm:      ins[i].InsStlmCNm,
				InsStlmCAcct:    ins[i].InsStlmCAcct,
				InsStlmCBkNo:    ins[i].InsStlmCBkNo,
				InsStlmCBkNm:    ins[i].InsStlmCBkNm,
				InsStlmDNm:      ins[i].InsStlmDNm,
				InsStlmDAcct:    ins[i].InsStlmDAcct,
				InsStlmDBkNo:    ins[i].InsStlmDBkNo,
				InsStlmDBkNm:    ins[i].InsStlmDBkNm,
				MsgResvFld1:     ins[i].MsgResvFld1,
				MsgResvFld2:     ins[i].MsgResvFld2,
				MsgResvFld3:     ins[i].MsgResvFld3,
				MsgResvFld4:     ins[i].MsgResvFld4,
				MsgResvFld5:     ins[i].MsgResvFld5,
				MsgResvFld6:     ins[i].MsgResvFld6,
				MsgResvFld7:     ins[i].MsgResvFld7,
				MsgResvFld8:     ins[i].MsgResvFld8,
				MsgResvFld9:     ins[i].MsgResvFld9,
				MsgResvFld10:    ins[i].MsgResvFld10,
				RecOprId:        ins[i].RecOprId,
				RecUpdOpr:       ins[i].RecUpdOpr,
			}
			if !ins[i].CreatedAt.IsZero() {
				pbIns[i].CreatedAt = ins[i].CreatedAt.Format(util.TimePattern)
			}

			if !ins[i].UpdatedAt.IsZero() {
				pbIns[i].UpdatedAt = ins[i].UpdatedAt.Format(util.TimePattern)
			}

		}

		return &pb.ListInstitutionsReply{
			Items: pbIns,
			Count: count,
			Page:  in.Page,
			Size:  in.Size,
		}, nil
	}
}

func (s *service) SaveInstitution(ctx context.Context, in *pb.SaveInstitutionRequest) (*pb.SaveInstitutionReply, error) {
	var reply pb.SaveInstitutionReply
	if in.Item == nil {
		reply.Err = &pb.Error{
			Code:        http.StatusBadRequest,
			Message:     "InvalidParamsError",
			Description: "机构信息为空",
		}
		return &reply, nil
	}

	if in.Item.InsIdCd == "" {
		reply.Err = &pb.Error{
			Code:        http.StatusBadRequest,
			Message:     "InvalidParamsError",
			Description: "id不能为空",
		}
		return &reply, nil
	}
	db := common.DB

	ins := new(insmodel.InstitutionInfo)
	{
		ins.InsIdCd = in.Item.InsIdCd
		ins.InsCompanyCd = in.Item.InsCompanyCd
		ins.InsType = in.Item.InsType
		ins.InsName = in.Item.InsName
		ins.InsProvCd = in.Item.InsProvCd
		ins.InsCityCd = in.Item.InsCityCd
		ins.InsRegionCd = in.Item.InsRegionCd
		ins.InsSta = in.Item.InsSta
		ins.InsStlmTp = in.Item.InsStlmTp
		ins.InsAloStlmCycle = in.Item.InsAloStlmCycle
		ins.InsAloStlmMd = in.Item.InsAloStlmMd
		ins.InsStlmCNm = in.Item.InsStlmCNm
		ins.InsStlmCAcct = in.Item.InsStlmCAcct
		ins.InsStlmCBkNo = in.Item.InsStlmCBkNo
		ins.InsStlmCBkNm = in.Item.InsStlmCBkNm
		ins.InsStlmDNm = in.Item.InsStlmDNm
		ins.InsStlmDAcct = in.Item.InsStlmDAcct
		ins.InsStlmDBkNo = in.Item.InsStlmDBkNo
		ins.InsStlmDBkNm = in.Item.InsStlmDBkNm
		ins.MsgResvFld1 = in.Item.MsgResvFld1
		ins.MsgResvFld2 = in.Item.MsgResvFld2
		ins.MsgResvFld3 = in.Item.MsgResvFld3
		ins.MsgResvFld4 = in.Item.MsgResvFld4
		ins.MsgResvFld5 = in.Item.MsgResvFld5
		ins.MsgResvFld6 = in.Item.MsgResvFld6
		ins.MsgResvFld7 = in.Item.MsgResvFld7
		ins.MsgResvFld8 = in.Item.MsgResvFld8
		ins.MsgResvFld9 = in.Item.MsgResvFld9
		ins.MsgResvFld10 = in.Item.MsgResvFld10
		ins.RecOprId = in.Item.RecOprId
		ins.RecUpdOpr = in.Item.RecUpdOpr
	}
	err := insmodel.SaveInstitution(db, ins)
	if err != nil {
		return nil, err
	}

	return &reply, nil
}
