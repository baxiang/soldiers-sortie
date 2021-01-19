package term

import (
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

type ActivationState struct {
	ActiveCode string         `gorm:"column:ACTIVE_CODE"`
	ActiveType string         `gorm:"column:ACTIVE_TYPE"`
	MchtCd     string         `gorm:"column:MCHT_CD"`
	TermId     string         `gorm:"column:TERM_ID"`
	NewKsn     string         `gorm:"column:NEW_KSN"`
	OldKsn     string         `gorm:"column:OLD_KSN"`
	IsActive   string         `gorm:"column:IS_ACTIVE"`
	RecOprId   string         `gorm:"column:REC_OPR_ID"`
	RecUpdOpr  string         `gorm:"column:REC_UPD_OPR"`
	ActiveDate mysql.NullTime `gorm:"column:ACTIVE_DATE"`
	CreateDate mysql.NullTime `gorm:"column:CREATE_DATE"`
	CreatedAt  time.Time      `gorm:"column:REC_CRT_TS"`
	UpdatedAt  time.Time      `gorm:"column:REC_UPD_TS"`
}

func (ActivationState) TableName() string {
	return "TBL_EDIT_TERM_ACTIVATION_STATE"
}

type ActivationStateMain struct {
	ActivationState
}

func (ActivationStateMain) TableName() string {
	return "TBL_TERM_ACTIVATION_STATE"
}

func SaveActivationState(db *gorm.DB, data *ActivationState) error {
	return db.Save(data).Error
}

func SaveActivationStateMain(db *gorm.DB, data *ActivationStateMain) error {
	return db.Save(data).Error
}

func QueryActivationState(db *gorm.DB, query *ActivationState, page int32, size int32) ([]*ActivationState, int32, error) {
	out := make([]*ActivationState, 0)
	var count int32
	db.Model(&ActivationState{}).Where(query).Count(&count)
	err := db.Where(query).Offset((page - 1) * size).Limit(size).Find(&out).Error
	return out, count, err
}

func QueryActivationStateMain(db *gorm.DB, query *ActivationStateMain, page int32, size int32) ([]*ActivationStateMain, int32, error) {
	out := make([]*ActivationStateMain, 0)
	var count int32
	db.Model(&ActivationStateMain{}).Where(query).Count(&count)
	err := db.Where(query).Offset((page - 1) * size).Limit(size).Find(&out).Error
	return out, count, err
}

func DeleteActivationState(db *gorm.DB, query *ActivationState) error {
	return db.Where(query).Delete(&ActivationState{}).Error
}
