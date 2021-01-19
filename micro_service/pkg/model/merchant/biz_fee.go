package merchant

import (
	"time"

	"github.com/jinzhu/gorm"
)

type BizFee struct {
	MchtCd          string    `gorm:"column:MCHT_CD;primary_key"`
	ProdCd          string    `gorm:"column:PROD_CD;primary_key"`
	BizCd           string    `gorm:"column:BIZ_CD;primary_key"`
	SubBizCd        string    `gorm:"column:SUB_BIZ_CD;primary_key"`
	MchtFeeMd       string    `gorm:"column:MCHT_FEE_MD"`
	MchtFeePercent  float64   `gorm:"column:MCHT_FEE_PERCENT"`
	MchtFeePctMin   float64   `gorm:"column:MCHT_FEE_PCT_MIN"`
	MchtFeePctMax   float64   `gorm:"column:MCHT_FEE_PCT_MAX"`
	MchtFeeSingle   float64   `gorm:"column:MCHT_FEE_SINGLE"`
	MchtAFeeSame    string    `gorm:"column:MCHT_A_FEE_SAME"`
	MchtAFeeMd      string    `gorm:"column:MCHT_A_FEE_MD"`
	MchtAFeePercent float64   `gorm:"column:MCHT_A_FEE_PERCENT"`
	MchtAFeePctMin  float64   `gorm:"column:MCHT_A_FEE_PCT_MIN"`
	MchtAFeePctMax  float64   `gorm:"column:MCHT_A_FEE_PCT_MAX"`
	MchtAFeeSingle  float64   `gorm:"column:MCHT_A_FEE_SINGLE"`
	OperIn          string    `gorm:"column:OPER_IN"`
	RecOprId        string    `gorm:"column:REC_OPR_ID"`
	RecUpdOpr       string    `gorm:"column:REC_UPD_OPR"`
	CreatedAt       time.Time `gorm:"column:REC_CRT_TS"`
	UpdatedAt       time.Time `gorm:"column:REC_UPD_TS"`
}

type BizFeeMain struct {
	BizFee
}

func (BizFee) TableName() string {
	return "TBL_EDIT_MCHT_BIZ_FEE"
}

func (BizFeeMain) TableName() string {
	return "TBL_MCHT_BIZ_FEE"
}

func SaveBizFee(db *gorm.DB, data *BizFee) error {
	return db.Save(data).Error
}

func SaveBizFeeMain(db *gorm.DB, data *BizFeeMain) error {
	return db.Save(data).Error
}

func QueryBizFee(db *gorm.DB, query *BizFee, page int32, size int32) ([]*BizFee, int32, error) {
	out := make([]*BizFee, 0)
	var count int32
	db.Model(&BizFee{}).Where(query).Count(&count)
	err := db.Where(query).Offset((page - 1) * size).Limit(size).Find(&out).Error
	return out, count, err
}

func QueryBizFeeMain(db *gorm.DB, query *BizFeeMain, page int32, size int32) ([]*BizFeeMain, int32, error) {
	out := make([]*BizFeeMain, 0)
	var count int32
	db.Model(&BizFeeMain{}).Where(query).Count(&count)
	err := db.Where(query).Offset((page - 1) * size).Limit(size).Find(&out).Error
	return out, count, err
}

func DeleteBizFee(db *gorm.DB, query *BizFee) error {
	return db.Where(query).Delete(&BizFee{}).Error
}
