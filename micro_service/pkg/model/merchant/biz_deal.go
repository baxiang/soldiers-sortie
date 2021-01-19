package merchant

import (
	"time"

	"github.com/jinzhu/gorm"
)

type BizDeal struct {
	MchtCd    string    `gorm:"column:MCHT_CD;primary_key"`
	ProdCd    string    `gorm:"column:PROD_CD;primary_key"`
	BizCd     string    `gorm:"column:BIZ_CD;primary_key"`
	TransCd   string    `gorm:"column:TRANS_CD;primary_key"`
	OperIn    string    `gorm:"column:OPER_IN"`
	RecOprId  string    `gorm:"column:REC_OPR_ID"`
	RecUpdOpr string    `gorm:"column:REC_UPD_OPR"`
	CreatedAt time.Time `gorm:"column:REC_CRT_TS"`
	UpdatedAt time.Time `gorm:"column:REC_UPD_TS"`
}

type BizDealMain struct {
	BizDeal
}

func (BizDeal) TableName() string {
	return "TBL_EDIT_MCHT_BIZ_DEAL"
}

func (BizDealMain) TableName() string {
	return "TBL_MCHT_BIZ_DEAL"
}

func SaveBizDeal(db *gorm.DB, data *BizDeal) error {
	return db.Save(data).Error
}

func SaveBizDealMain(db *gorm.DB, data *BizDealMain) error {
	return db.Save(data).Error
}

func QueryBizDeal(db *gorm.DB, query *BizDeal, page int32, size int32) ([]*BizDeal, int32, error) {
	out := make([]*BizDeal, 0)
	var count int32
	db.Model(&BizDeal{}).Where(query).Count(&count)
	err := db.Where(query).Offset((page - 1) * size).Limit(size).Find(&out).Error
	return out, count, err
}

func QueryBizDealMain(db *gorm.DB, query *BizDealMain, page int32, size int32) ([]*BizDealMain, int32, error) {
	out := make([]*BizDealMain, 0)
	var count int32
	db.Model(&BizDealMain{}).Where(query).Count(&count)
	err := db.Where(query).Offset((page - 1) * size).Limit(size).Find(&out).Error
	return out, count, err
}

func DeleteBizDeal(db *gorm.DB, query *BizDeal) error {
	return db.Where(query).Delete(&BizDeal{}).Error
}
