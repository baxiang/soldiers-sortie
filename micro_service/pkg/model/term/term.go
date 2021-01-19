package term

import (
	"time"
	"userService/pkg/model/merchant"

	"github.com/jinzhu/gorm"
)

type Info struct {
	MchtCd          string    `gorm:"column:MCHT_CD;primary_key"`
	TermId          string    `gorm:"column:TERM_ID;primary_key"`
	TermTp          string    `gorm:"column:TERM_TP"`
	Belong          string    `gorm:"column:BELONG"`
	BelongSub       string    `gorm:"column:BELONG_SUB"`
	TmnlMoneyIntype string    `gorm:"column:TMNL_MONEY_INTYPE"`
	TmnlMoney       float64   `gorm:"column:TMNL_MONEY"`
	TmnlBrand       string    `gorm:"column:TMNL_BRAND"`
	TmnlModelNo     string    `gorm:"column:TMNL_MODEL_NO"`
	TmnlBarcode     string    `gorm:"column:TMNL_BARCODE"`
	DeviceCd        string    `gorm:"column:DEVICE_CD"`
	InstallLocation string    `gorm:"column:INSTALLLOCATION"`
	TmnlIntype      string    `gorm:"column:TMNL_INTYPE"`
	DialOut         string    `gorm:"column:DIAL_OUT"`
	DealTypes       string    `gorm:"column:DEAL_TYPES"`
	RecOprId        string    `gorm:"column:REC_OPR_ID"`
	RecUpdOpr       string    `gorm:"column:REC_UPD_OPR"`
	AppCd           string    `gorm:"column:APP_CD"`
	SystemFlag      string    `gorm:"column:SYSTEMFLAG"`
	Status          string    `gorm:"column:STATUS"`
	ActiveCode      string    `gorm:"column:ACTIVE_CODE"`
	NoFlag          string    `gorm:"column:NO_FLAG"`
	MsgResvFld1     string    `gorm:"column:MSG_RESV_FLD1"`
	MsgResvFld2     string    `gorm:"column:MSG_RESV_FLD2"`
	MsgResvFld3     string    `gorm:"column:MSG_RESV_FLD3"`
	MsgResvFld4     string    `gorm:"column:MSG_RESV_FLD4"`
	MsgResvFld5     string    `gorm:"column:MSG_RESV_FLD5"`
	MsgResvFld6     string    `gorm:"column:MSG_RESV_FLD6"`
	MsgResvFld7     string    `gorm:"column:MSG_RESV_FLD7"`
	MsgResvFld8     string    `gorm:"column:MSG_RESV_FLD8"`
	MsgResvFld9     string    `gorm:"column:MSG_RESV_FLD9"`
	MsgResvFld10    string    `gorm:"column:MSG_RESV_FLD10"`
	CreatedAt       time.Time `gorm:"column:REC_CRT_TS"`
	UpdatedAt       time.Time `gorm:"column:REC_UPD_TS"`
}

func (Info) TableName() string {
	return "TBL_EDIT_TERM_INF"
}

type InfoMain struct {
	Info
}

func (InfoMain) TableName() string {
	return "TBL_TERM_INF"
}

func QueryTermInfo(db *gorm.DB, query *Info, insIds []string, merchantInfo *merchant.MerchantInfo, page int32, size int32) ([]*Info, int32, error) {
	out := make([]*Info, 0)
	var count int32
	var err error
	if len(insIds) == 0 {
		q := db.Table("TBL_EDIT_TERM_INF").Select("TBL_EDIT_TERM_INF.*")
		if merchantInfo != nil {
			if merchantInfo.AipBranCd != "" || merchantInfo.BankBelongCd != "" {
				q = q.Joins("left join TBL_EDIT_MCHT_INF b on TBL_EDIT_TERM_INF.MCHT_CD = b.MCHT_CD")
			}
			if merchantInfo.AipBranCd != "" {
				q = q.Where("b.AIP_BRAN_CD = ?", merchantInfo.AipBranCd)
			}
			if merchantInfo.BankBelongCd != "" {
				q = q.Where("b.BANK_BELONG_CD = ?", merchantInfo.BankBelongCd)
			}
		}
		q = q.Where(query)
		q.Count(&count)
		err = q.Offset((page - 1) * size).Limit(size).Scan(&out).Error
	} else {
		q := db.Table("TBL_EDIT_TERM_INF").
			Select("TBL_EDIT_TERM_INF.*").
			Joins("left join TBL_EDIT_MCHT_INF b on TBL_EDIT_TERM_INF.MCHT_CD = b.MCHT_CD").
			Where("b.AIP_BRAN_CD in (?)", insIds)
		if merchantInfo != nil && merchantInfo.BankBelongCd != "" && merchantInfo.AipBranCd != "" {
			if merchantInfo.AipBranCd != "" {
				q = q.Where("b.AIP_BRAN_CD = ?", merchantInfo.AipBranCd)
			}
			if merchantInfo.BankBelongCd != "" {
				q = q.Where("b.BANK_BELONG_CD = ?", merchantInfo.BankBelongCd)
			}
		}
		q = q.Where(query)
		q.Count(&count)
		err = q.Offset((page - 1) * size).Limit(size).Scan(&out).Error
	}
	return out, count, err
}

func SaveTermInfo(db *gorm.DB, data *Info) error {
	return db.Save(data).Error
}

func SaveTermInfoMain(db *gorm.DB, data *InfoMain) error {
	return db.Save(data).Error
}

func QueryTermInfoMain(db *gorm.DB, query *InfoMain, insIds []string, merchantInfo *merchant.MerchantInfo, page int32, size int32) ([]*InfoMain, int32, error) {
	out := make([]*InfoMain, 0)
	var count int32
	var err error
	if len(insIds) == 0 {
		q := db.Table("TBL_TERM_INF").Select("TBL_TERM_INF.*")
		if merchantInfo != nil {
			if merchantInfo.AipBranCd != "" || merchantInfo.BankBelongCd != "" {
				q = q.Joins("left join TBL_MCHT_INF b on TBL_TERM_INF.MCHT_CD = b.MCHT_CD")
			}
			if merchantInfo.AipBranCd != "" {
				q = q.Where("b.AIP_BRAN_CD = ?", merchantInfo.AipBranCd)
			}
			if merchantInfo.BankBelongCd != "" {
				q = q.Where("b.BANK_BELONG_CD = ?", merchantInfo.BankBelongCd)
			}
		}
		q = q.Where(query)
		q.Count(&count)
		err = q.Offset((page - 1) * size).Limit(size).Scan(&out).Error
	} else {
		q := db.Table("TBL_TERM_INF").
			Select("TBL_TERM_INF.*").
			Joins("left join TBL_MCHT_INF b on TBL_TERM_INF.MCHT_CD = b.MCHT_CD").
			Where("b.AIP_BRAN_CD in (?)", insIds)
		if merchantInfo != nil && merchantInfo.BankBelongCd != "" && merchantInfo.AipBranCd != "" {
			if merchantInfo.AipBranCd != "" {
				q = q.Where("b.AIP_BRAN_CD = ?", merchantInfo.AipBranCd)
			}
			if merchantInfo.BankBelongCd != "" {
				q = q.Where("b.BANK_BELONG_CD = ?", merchantInfo.BankBelongCd)
			}
		}
		q = q.Where(query)
		q.Count(&count)
		err = q.Offset((page - 1) * size).Limit(size).Scan(&out).Error
	}
	return out, count, err
}

func DeleteTerm(db *gorm.DB, query *Info) error {
	return db.Where(query).Delete(&Info{}).Error
}

func UpdateTerm(db *gorm.DB, query *Info, data *Info) error {
	return db.Model(&Info{}).Where(query).Updates(data).Error
}

func UpdateTermMain(db *gorm.DB, query *InfoMain, data *InfoMain) error {
	return db.Model(&InfoMain{}).Where(query).Updates(data).Error
}

func FindTermByPk(db *gorm.DB, mchtCd string, termId string) (*Info, error) {
	out := new(Info)
	err := db.Where(&Info{MchtCd: mchtCd, TermId: termId}).Take(out).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return out, err
}
