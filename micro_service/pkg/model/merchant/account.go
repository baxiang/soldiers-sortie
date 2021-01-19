package merchant

import (
	"time"

	"github.com/jinzhu/gorm"
)

type BankAccount struct {
	OwnerCd      string    `gorm:"column:OWNER_CD;primary_key"`
	AccountType  string    `gorm:"column:ACCOUNTTYPE;primary_key"`
	Name         string    `gorm:"column:NAME"`
	Account      string    `gorm:"column:ACCOUNT"`
	UcBcCd       string    `gorm:"column:UC_BC_CD"`
	Province     string    `gorm:"column:PROVINCE"`
	City         string    `gorm:"column:CITY"`
	BankCode     string    `gorm:"column:BANK_CODE"`
	BankName     string    `gorm:"column:BANK_NAME"`
	OperIn       string    `gorm:"column:OPER_IN"`
	RecOprId     string    `gorm:"column:REC_OPR_ID"`
	RecUpdOpr    string    `gorm:"column:REC_UPD_OPR"`
	MsgResvFld1  string    `gorm:"column:MSG_RESV_FLD1"`
	MsgResvFld2  string    `gorm:"column:MSG_RESV_FLD2"`
	MsgResvFld3  string    `gorm:"column:MSG_RESV_FLD3"`
	MsgResvFld4  string    `gorm:"column:MSG_RESV_FLD4"`
	MsgResvFld5  string    `gorm:"column:MSG_RESV_FLD5"`
	MsgResvFld6  string    `gorm:"column:MSG_RESV_FLD6"`
	MsgResvFld7  string    `gorm:"column:MSG_RESV_FLD7"`
	MsgResvFld8  string    `gorm:"column:MSG_RESV_FLD8"`
	MsgResvFld9  string    `gorm:"column:MSG_RESV_FLD9"`
	MsgResvFld10 string    `gorm:"column:MSG_RESV_FLD10"`
	CreatedAt    time.Time `gorm:"column:REC_CRT_TS"`
	UpdatedAt    time.Time `gorm:"column:REC_UPD_TS"`
}

func (b BankAccount) TableName() string {
	return "TBL_EDIT_MCHT_BANKACCOUNT"
}

type BankAccountMain struct {
	BankAccount
}

func (b BankAccountMain) TableName() string {
	return "TBL_MCHT_BANKACCOUNT"
}

func SaveBankAccount(db *gorm.DB, data *BankAccount) error {
	return db.Save(data).Error
}

func SaveBankAccountMain(db *gorm.DB, data *BankAccountMain) error {
	return db.Save(data).Error
}

func QueryBankAccount(db *gorm.DB, query *BankAccount, page int32, size int32) ([]*BankAccount, int32, error) {
	out := make([]*BankAccount, 0)
	var count int32
	db.Model(&BankAccount{}).Where(query).Count(&count)
	err := db.Where(query).Offset((page - 1) * size).Limit(size).Find(&out).Error
	return out, count, err
}

func QueryBankAccountMain(db *gorm.DB, query *BankAccountMain, page int32, size int32) ([]*BankAccountMain, int32, error) {
	out := make([]*BankAccountMain, 0)
	var count int32
	db.Model(&BankAccountMain{}).Where(query).Count(&count)
	err := db.Where(query).Offset((page - 1) * size).Limit(size).Find(&out).Error
	return out, count, err
}

func DeleteBankAccount(db *gorm.DB, query *BankAccount) error {
	return db.Where(query).Delete(&BankAccount{}).Error
}
