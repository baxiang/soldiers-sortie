package task

import (
	camundamodel "userService/pkg/model/camunda"
	"userService/pkg/model/merchant"
	"userService/pkg/model/term"

	"github.com/jinzhu/gorm"
)

func termAdd(db *gorm.DB, instance *camundamodel.ProcessInstance) error {
	// 查询终端信息
	infos, _, err := term.QueryTermInfo(db, &term.Info{MchtCd: instance.DataId}, nil, nil, 1, 9999)
	if err != nil {
		return err
	}

	for _, info := range infos {
		if info.SystemFlag != "00" && info.SystemFlag != "01" {
			err = term.UpdateTerm(db, &term.Info{MchtCd: info.MchtCd, TermId: info.TermId}, &term.Info{Status: "01", SystemFlag: "01"})
			if err != nil {
				return err
			}

			info.Status = "01"
			info.SystemFlag = "01"
			err = term.SaveTermInfoMain(db, &term.InfoMain{Info: *info})
			if err != nil {
				return err
			}
		}

	}

	err = merchant.UpdateMerchant(db, &merchant.MerchantInfo{MchtCd: instance.DataId}, &merchant.MerchantInfo{SystemFlag: "01"})

	return err
}

func termDelete(db *gorm.DB, instance *camundamodel.ProcessInstance) error {
	infos, _, err := term.QueryTermInfo(db, &term.Info{MchtCd: instance.DataId}, nil, nil, 1, 9999)
	if err != nil {
		return err
	}

	for _, info := range infos {
		if info.SystemFlag != "00" && info.SystemFlag != "01" {
			err = term.DeleteTerm(db, &term.Info{MchtCd: info.MchtCd, TermId: info.TermId})
			if err != nil {
				return err
			}
		}
	}

	err = merchant.UpdateMerchant(db, &merchant.MerchantInfo{MchtCd: instance.DataId}, &merchant.MerchantInfo{SystemFlag: "01"})

	return err
}

func termInfoUnregister(db *gorm.DB, instance *camundamodel.ProcessInstance) error {
	infos, _, err := term.QueryTermInfo(db, &term.Info{MchtCd: instance.DataId}, nil, nil, 1, 9999)
	if err != nil {
		return err
	}

	for _, info := range infos {
		if info.SystemFlag != "00" && info.SystemFlag != "01" {
			err = term.UpdateTerm(
				db,
				&term.Info{MchtCd: info.MchtCd, TermId: info.TermId},
				&term.Info{Status: "00", SystemFlag: "00"},
			)
			if err != nil {
				return err
			}

			err = term.UpdateTermMain(
				db,
				&term.InfoMain{Info: term.Info{MchtCd: info.MchtCd, TermId: info.TermId}},
				&term.InfoMain{Info: term.Info{Status: "00", SystemFlag: "00"}},
			)
			if err != nil {
				return err
			}
		}

	}

	err = merchant.UpdateMerchant(db, &merchant.MerchantInfo{MchtCd: instance.DataId}, &merchant.MerchantInfo{SystemFlag: "01"})

	return err
}

func cancelTermInfoUnregister(db *gorm.DB, instance *camundamodel.ProcessInstance) error {
	infos, _, err := term.QueryTermInfo(db, &term.Info{MchtCd: instance.DataId}, nil, nil, 1, 9999)
	if err != nil {
		return err
	}

	for _, info := range infos {
		if info.SystemFlag != "00" && info.SystemFlag != "01" {
			err = term.UpdateTerm(db, &term.Info{MchtCd: info.MchtCd, TermId: info.TermId}, &term.Info{Status: "01", SystemFlag: "01"})
			if err != nil {
				return err
			}
		}

	}

	err = merchant.UpdateMerchant(db, &merchant.MerchantInfo{MchtCd: instance.DataId}, &merchant.MerchantInfo{SystemFlag: "01"})

	return err
}
