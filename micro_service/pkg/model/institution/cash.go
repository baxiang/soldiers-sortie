package institution

import (
	"time"

	"github.com/jinzhu/gorm"
)

type Cash struct {
	InsIdCd        string    `gorm:"column:INS_ID_CD;primary_key"`
	ProdCd         string    `gorm:"column:PROD_CD;primary_key"`
	InsDefaultFlag string    `gorm:"column:INS_DEFAULT_FLAG"`
	InsDefaultCash float64   `gorm:"column:INS_DEFAULT_CASH"`
	InsCurrentCash float64   `gorm:"column:INS_CURRENT_CASH"`
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

type CashMain struct {
	Cash
}

func (c CashMain) TableName() string {
	return "TBL_INS_CASH_INF"
}

func (c Cash) TableName() string {
	return "TBL_EDIT_INS_CASH_INF"
}

func SaveInstitutionCash(db *gorm.DB, cash *Cash) error {
	return db.Save(cash).Error
}

func FindInstitutionCashMain(db *gorm.DB, query *CashMain, page int32, size int32) ([]*CashMain, int32, error) {
	out := make([]*CashMain, 0)
	var count int32
	db.Model(&CashMain{}).Where(query).Count(&count)
	err := db.Where(query).Offset((page - 1) * size).Limit(size).Find(&out).Error
	return out, count, err
}

func FindInstitutionCash(db *gorm.DB, query *Cash, page int32, size int32) ([]*Cash, int32, error) {
	out := make([]*Cash, 0)
	var count int32
	db.Model(&Cash{}).Where(query).Count(&count)
	err := db.Where(query).Offset((page - 1) * size).Limit(size).Find(&out).Error
	return out, count, err
}

func SaveInstitutionCashMain(db *gorm.DB, cash *CashMain) error {
	return db.Save(cash).Error
}

func DeleteInstitutionCash(db *gorm.DB, query *Cash) error {
	return db.Where(query).Delete(&Cash{}).Error
}
