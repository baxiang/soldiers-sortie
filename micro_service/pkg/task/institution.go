package task

import (
	"fmt"
	camundamodel "userService/pkg/model/camunda"
	"userService/pkg/model/institution"
	"userService/pkg/model/static"

	"github.com/jinzhu/gorm"
)

func institutionRegister(db *gorm.DB, instance *camundamodel.ProcessInstance) error {
	// 查询机构信息
	info, err := institution.FindInstitutionInfoById(db, instance.DataId)
	if err != nil {
		return err
	}
	if info == nil {
		return fmt.Errorf("institution %s not found", instance.DataId)
	}
	// 查询fee，cash，control
	fees, _, err := institution.FindInstitutionFee(db, &institution.Fee{
		InsIdCd: info.InsIdCd,
	}, 1, 9999)
	if err != nil {
		return err
	}
	cashes, _, err := institution.FindInstitutionCash(db, &institution.Cash{
		InsIdCd: info.InsIdCd,
	}, 1, 9999)
	if err != nil {
		return err
	}
	controls, _, err := institution.FindInstitutionControl(db, &institution.Control{
		InsIdCd: info.InsIdCd,
	}, 1, 9999)
	if err != nil {
		return err
	}

	// 修改机构状态
	err = institution.UpdateInstitution(db, &institution.InstitutionInfo{
		InsIdCd: info.InsIdCd,
	}, &institution.InstitutionInfo{
		InsSta: "01",
	})
	if err != nil {
		return err
	}

	// 入正式表 institution fee cash control
	info.InsSta = "01"

	err = institution.SaveInstitutionMain(db, &institution.InstitutionInfoMain{
		InstitutionInfo: *info,
	})
	if err != nil {
		return err
	}
	for _, fee := range fees {
		err = institution.SaveInstitutionFeeMain(db, &institution.FeeMain{
			Fee: *fee,
		})
		if err != nil {
			return err
		}
	}
	for _, cash := range cashes {
		err = institution.SaveInstitutionCashMain(db, &institution.CashMain{
			Cash: *cash,
		})
		if err != nil {
			return err
		}
	}
	for _, control := range controls {
		err = institution.SaveInstitutionControlMain(db, &institution.ControlMain{
			Control: *control,
		})
		if err != nil {
			return err
		}
	}

	// 保存到字典表
	item := new(static.DictionaryItem)
	if info.InsType == "0" {
		item.DicType = "INS_COMPANY_CD"
		item.Memo = "所属机构编码"
	} else {
		item.DicType = "INS_ID_CD"
		item.Memo = "收单机构编码"
	}

	item.DicCode = info.InsIdCd
	item.DicName = info.InsIdCd + "_" + info.InsName
	item.DispOrder = info.InsIdCd

	err = static.SaveDictionaryItem(db, item)
	return err
}

func deleteInstitution(db *gorm.DB, instance *camundamodel.ProcessInstance) error {
	var err error
	// 删除 institution fee cash control
	err = institution.DeleteInstitution(db, &institution.InstitutionInfo{
		InsIdCd: instance.DataId,
	})
	if err != nil {
		return err
	}
	err = institution.DeleteInstitutionFee(db, &institution.Fee{
		InsIdCd: instance.DataId,
	})
	if err != nil {
		return err
	}
	err = institution.DeleteInstitutionCash(db, &institution.Cash{
		InsIdCd: instance.DataId,
	})
	if err != nil {
		return err
	}
	err = institution.DeleteInstitutionControl(db, &institution.Control{
		InsIdCd: instance.DataId,
	})
	if err != nil {
		return err
	}
	return nil
}

func institutionUpdate(db *gorm.DB, instance *camundamodel.ProcessInstance) error {
	// 查询机构信息
	info, err := institution.FindInstitutionInfoById(db, instance.DataId)
	if err != nil {
		return err
	}
	if info == nil {
		return fmt.Errorf("institution %s not found", instance.DataId)
	}
	// 查询fee，cash，control
	fees, _, err := institution.FindInstitutionFee(db, &institution.Fee{
		InsIdCd: info.InsIdCd,
	}, 1, 9999)
	if err != nil {
		return err
	}
	cashes, _, err := institution.FindInstitutionCash(db, &institution.Cash{
		InsIdCd: info.InsIdCd,
	}, 1, 9999)
	if err != nil {
		return err
	}
	controls, _, err := institution.FindInstitutionControl(db, &institution.Control{
		InsIdCd: info.InsIdCd,
	}, 1, 9999)
	if err != nil {
		return err
	}

	err = institution.UpdateInstitution(
		db,
		&institution.InstitutionInfo{InsIdCd: info.InsIdCd},
		&institution.InstitutionInfo{InsSta: "01"},
	)
	if err != nil {
		return err
	}
	info.InsSta = "01"
	// 入正式表 institution fee cash control
	err = institution.SaveInstitutionMain(db, &institution.InstitutionInfoMain{
		InstitutionInfo: *info,
	})
	if err != nil {
		return err
	}
	for _, fee := range fees {
		err = institution.SaveInstitutionFeeMain(db, &institution.FeeMain{
			Fee: *fee,
		})
		if err != nil {
			return err
		}
	}
	for _, cash := range cashes {
		err = institution.SaveInstitutionCashMain(db, &institution.CashMain{
			Cash: *cash,
		})
		if err != nil {
			return err
		}
	}
	for _, control := range controls {
		err = institution.SaveInstitutionControlMain(db, &institution.ControlMain{
			Control: *control,
		})
		if err != nil {
			return err
		}
	}

	// 保存到字典表
	item := new(static.DictionaryItem)
	if info.InsType == "0" {
		item.DicType = "INS_COMPANY_CD"
		item.Memo = "所属机构编码"
	} else {
		item.DicType = "INS_ID_CD"
		item.Memo = "收单机构编码"
	}

	item.DicCode = info.InsIdCd
	item.DicName = info.InsIdCd + "_" + info.InsName
	item.DispOrder = info.InsIdCd

	err = static.SaveDictionaryItem(db, item)
	return nil
}

func institutionUpdateCancel(db *gorm.DB, instance *camundamodel.ProcessInstance) error {
	// 查询正式表机构信息
	info, err := institution.FindInstitutionInfoMainById(db, instance.DataId)
	if err != nil {
		return err
	}
	if info == nil {
		return fmt.Errorf("institution %s not found", instance.DataId)
	}
	// 查询正式表fee，cash，control
	fees, _, err := institution.FindInstitutionFeeMain(db, &institution.FeeMain{
		Fee: institution.Fee{
			InsIdCd: info.InsIdCd,
		},
	}, 1, 9999)
	if err != nil {
		return err
	}
	cashes, _, err := institution.FindInstitutionCashMain(db, &institution.CashMain{
		Cash: institution.Cash{
			InsIdCd: info.InsIdCd,
		},
	}, 1, 9999)
	if err != nil {
		return err
	}
	controls, _, err := institution.FindInstitutionControlMain(db, &institution.ControlMain{
		Control: institution.Control{
			InsIdCd: info.InsIdCd,
		},
	}, 1, 9999)
	if err != nil {
		return err
	}

	// 入编辑表 institution fee cash control
	err = institution.SaveInstitution(db, &info.InstitutionInfo)
	if err != nil {
		return err
	}
	for _, fee := range fees {
		err = institution.SaveInstitutionFee(db, &fee.Fee)
		if err != nil {
			return err
		}
	}
	for _, cash := range cashes {
		err = institution.SaveInstitutionCash(db, &cash.Cash)
		if err != nil {
			return err
		}
	}
	for _, control := range controls {
		err = institution.SaveInstitutionControl(db, &control.Control)
		if err != nil {
			return err
		}
	}
	return nil
}

func institutionUnRegister(db *gorm.DB, instance *camundamodel.ProcessInstance) error {
	// 查询机构信息
	info, err := institution.FindInstitutionInfoById(db, instance.DataId)
	if err != nil {
		return err
	}
	if info == nil {
		return fmt.Errorf("institution %s not found", instance.DataId)
	}

	err = institution.UpdateInstitution(db, &institution.InstitutionInfo{InsIdCd: info.InsIdCd}, &institution.InstitutionInfo{InsSta: "00"})
	if err != nil {
		return err
	}

	info.InsSta = "00"
	err = institution.SaveInstitutionMain(db, &institution.InstitutionInfoMain{
		InstitutionInfo: *info,
	})
	return err
}

func institutionCancelUnRegister(db *gorm.DB, instance *camundamodel.ProcessInstance) error {
	// 查询机构信息
	info, err := institution.FindInstitutionInfoMainById(db, instance.DataId)
	if err != nil {
		return err
	}
	if info == nil {
		return fmt.Errorf("institution %s not found", instance.DataId)
	}

	err = institution.SaveInstitution(db, &info.InstitutionInfo)
	return err
}
