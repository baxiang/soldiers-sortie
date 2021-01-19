package institution

import (
	"time"

	"github.com/jinzhu/gorm"
)

type Fee struct {
	InsIdCd        string    `gorm:"column:INS_ID_CD;primary_key"`
	ProdCd         string    `gorm:"column:PROD_CD;primary_key"`
	BizCd          string    `gorm:"column:BIZ_CD;primary_key"`
	SubBizCd       string    `gorm:"column:SUB_BIZ_CD;primary_key"`
	InsFeeBizCd    string    `gorm:"column:INS_FEE_BIZ_CD;primary_key"`
	InsFeeCd       string    `gorm:"column:INS_FEE_CD"`
	InsFeeTp       string    `gorm:"column:INS_FEE_TP"`
	InsFeeParam    string    `gorm:"column:INS_FEE_PARAM"`
	InsFeePercent  float64   `gorm:"column:INS_FEE_PERCENT"`
	InsFeePct      float64   `gorm:"column:INS_FEE_PCT"`
	InsFeePctMin   float64   `gorm:"column:INS_FEE_PCT_MIN"`
	InsFeePctMax   float64   `gorm:"column:INS_FEE_PCT_MAX"`
	InsAFeeSame    string    `gorm:"column:INS_A_FEE_SAME"`
	InsAFeeParam   string    `gorm:"column:INS_A_FEE_PARAM"`
	InsAFeePercent float64   `gorm:"column:INS_A_FEE_PERCENT"`
	InsAFeePct     float64   `gorm:"column:INS_A_FEE_PCT"`
	InsAFeePctMin  float64   `gorm:"column:INS_A_FEE_PCT_MIN"`
	InsAFeePctMax  float64   `gorm:"column:INS_A_FEE_PCT_MAX"`
	MsgResvFld1    string    `gorm:"column:MSG_RESV_FLD1"`
	MsgResvFld2    string    `gorm:"column:MSG_RESV_FLD2"`
	MsgResvFld3    string    `gorm:"column:MSG_RESV_FLD3"`
	MsgResvFld4    string    `gorm:"column:MSG_RESV_FLD4"`
	MsgResvFld5    string    `gorm:"column:MSG_RESV_FLD5"`
	MsgResvFld6    string    `gorm:"column:MSG_RESV_FLD6"`
	MsgResvFld7    string    `gorm:"column:MSG_RESV_FLD7"`
	MsgResvFld8    string    `gorm:"column:MSG_RESV_FLD8"`
	MsgResvFld9    string    `gorm:"column:MSG_RESV_FLD9"`
	MsgResvFld10   string    `gorm:"column:MSG_RESV_FLD10"`
	RecOprId       string    `gorm:"column:REC_OPR_ID"`
	RecUpdOpr      string    `gorm:"column:REC_UPD_OPR"`
	CreatedAt      time.Time `gorm:"column:REC_CRT_TS"`
	UpdatedAt      time.Time `gorm:"column:REC_UPD_TS"`
}

type FeeMain struct {
	Fee
}

func (f FeeMain) TableName() string {
	return "TBL_INS_FEE_INF"
}

func (f Fee) TableName() string {
	return "TBL_EDIT_INS_FEE_INF"
}

func SaveInstitutionFee(db *gorm.DB, fee *Fee) error {
	return db.Save(fee).Error
}

func FindInstitutionFee(db *gorm.DB, query *Fee, page int32, size int32) ([]*Fee, int32, error) {
	out := make([]*Fee, 0)
	var count int32
	db.Model(&Fee{}).Where(query).Count(&count)
	err := db.Where(query).Offset((page - 1) * size).Limit(size).Find(&out).Error
	return out, count, err
}

func FindInstitutionFeeMain(db *gorm.DB, query *FeeMain, page int32, size int32) ([]*FeeMain, int32, error) {
	out := make([]*FeeMain, 0)
	var count int32
	db.Model(&FeeMain{}).Where(query).Count(&count)
	err := db.Where(query).Offset((page - 1) * size).Limit(size).Find(&out).Error
	return out, count, err
}

func SaveInstitutionFeeMain(db *gorm.DB, fee *FeeMain) error {
	return db.Save(fee).Error
}

func DeleteInstitutionFee(db *gorm.DB, query *Fee) error {
	return db.Where(query).Delete(&Fee{}).Error
}
