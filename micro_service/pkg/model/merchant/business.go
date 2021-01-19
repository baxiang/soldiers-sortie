package merchant

import (
	"time"

	"github.com/jinzhu/gorm"
)

type Business struct {
	MchtCd                 string    `gorm:"column:MCHT_CD;primary_key"`
	ProdCd                 string    `gorm:"column:PROD_CD"`
	ProdCdText             string    `gorm:"column:PROD_CD_TEXT"`
	FeeMoneyCd             string    `gorm:"column:FEE_MONEY_CD"`
	FeeModeType            string    `gorm:"column:FEE_MODE_TYPE"`
	FeeSettlementType      string    `gorm:"column:FEE_SETTLEMENT_TYPE"`
	FeeHoliday             string    `gorm:"column:FEE_HOLIDAY"`
	ServiceFeeType         string    `gorm:"column:SERVICE_FEE_TYPE"`
	ServiceFeeStaticAmount float64   `gorm:"column:SERVICE_FEE_STATIC_AMMOUNT"`
	ServiceFeeLevelCount   int64     `gorm:"column:SERVICE_FEE_LEVEL_COUNT"`
	ServiceFeeMode         string    `gorm:"column:SERVICE_FEE_MODE"`
	ServiceFeeUnit         string    `gorm:"column:SERVICE_FEE_UNIT"`
	ServiceFeeTerm         string    `gorm:"column:SERVICE_FEE_TERM"`
	ServiceFeeSumto        string    `gorm:"column:SERVICE_FEE_SUMTO"`
	ServiceFeeCircle       string    `gorm:"column:SERVICE_FEE_CIRCLE"`
	ServiceFeeOthers       string    `gorm:"column:SERVICE_FEE_OTHERS"`
	ServiceFeeStart        string    `gorm:"column:SERVICE_FEE_START"`
	ServiceFeeClct         string    `gorm:"column:SERVICE_FEE_CLCT"`
	ServiceFeeClctOthers   string    `gorm:"column:SERVICE_FEE_CLCT_OTHERS"`
	SystemFlag             string    `gorm:"column:SYSTEMFLAG"`
	Ext1                   string    `gorm:"column:EXT1"`
	Ext2                   string    `gorm:"column:EXT2"`
	Ext3                   string    `gorm:"column:EXT3"`
	Ext4                   string    `gorm:"column:EXT4"`
	ServiceFeeYesNo        string    `gorm:"column:SERVICE_FEE_YESNO"`
	RecOprId               string    `gorm:"column:REC_OPR_ID"`
	RecUpdOpr              string    `gorm:"column:REC_UPD_OPR"`
	OperIn                 string    `gorm:"column:OPER_IN"`
	CreatedAt              time.Time `gorm:"column:REC_CRT_TS"`
	UpdatedAt              time.Time `gorm:"column:REC_UPD_TS"`
}

type BusinessMain struct {
	Business
}

func (Business) TableName() string {
	return "TBL_EDIT_MCHT_BUSINESS"
}

func (BusinessMain) TableName() string {
	return "TBL_MCHT_BUSINESS"
}

func SaveBusiness(db *gorm.DB, data *Business) error {
	return db.Save(data).Error
}

func SaveBusinessMain(db *gorm.DB, data *BusinessMain) error {
	return db.Save(data).Error
}

func QueryBusiness(db *gorm.DB, query *Business, page int32, size int32) ([]*Business, int32, error) {
	out := make([]*Business, 0)
	var count int32
	db.Model(&Business{}).Where(query).Count(&count)
	err := db.Where(query).Offset((page - 1) * size).Limit(size).Find(&out).Error
	return out, count, err
}

func QueryBusinessMain(db *gorm.DB, query *BusinessMain, page int32, size int32) ([]*BusinessMain, int32, error) {
	out := make([]*BusinessMain, 0)
	var count int32
	db.Model(&BusinessMain{}).Where(query).Count(&count)
	err := db.Where(query).Offset((page - 1) * size).Limit(size).Find(&out).Error
	return out, count, err
}

func DeleteBusiness(db *gorm.DB, query *Business) error {
	return db.Where(query).Delete(&Business{}).Error
}
