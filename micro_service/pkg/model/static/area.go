package static

import (
	"time"

	"github.com/jinzhu/gorm"
)

const (
	orgDictionaryItem   = "TBL_ORG_DICTIONARYITEM"
	dictionaryLayerItem = "TBL_DICTIONARYLAYERITEM"
)

type Area struct {
	Name    string `gorm:"column:NAME"`
	DicCode string `gorm:"column:DIC_CODE"`
}

func FindProvince(db *gorm.DB, orgCode string) ([]*Area, error) {
	out := make([]*Area, 0)
	err := db.Table(dictionaryLayerItem).
		Select("NAME, DIC_CODE").
		Joins("left join TBL_ORG_DICTIONARYITEM b on b.ITEM_CODE = DIC_CODE").
		Where(
			"ORG_CODE = ? and TYPE_CODE = ? and DIC_TYPE = ? and DIC_LEVEL = ?",
			orgCode,
			"UnionPay_Bank_32Area",
			"UnionPay_Area",
			"1",
		).Order("DIC_CODE", true).Scan(&out).Error
	return out, err
}

func FindCity(db *gorm.DB, dicCode string, level string) ([]*Area, error) {
	out := make([]*Area, 0)
	err := db.Table(dictionaryLayerItem).
		Select("NAME, DIC_CODE").
		Where("DIC_PCODE = ? and DIC_LEVEL = ? and dic_type = ?",
			dicCode,
			level,
			"UnionPay_Area",
		).
		Scan(&out).Error
	return out, err
}

type OrgDictionaryItem struct {
	Id          int64     `gorm:"column:ID;primary_key"`
	TypeCode    string    `gorm:"column:TYPE_CODE"`
	OrgCode     string    `gorm:"column:ORG_CODE"`
	ItemCode    string    `gorm:"column:ITEM_CODE"`
	TypeParm1   string    `gorm:"column:TYPE_PARM1"`
	TypeParm2   string    `gorm:"column:TYPE_PARM2"`
	TypeParm3   string    `gorm:"column:TYPE_PARM3"`
	Remarks     string    `gorm:"column:REMARKS"`
	MsgResvFld1 string    `gorm:"column:MSG_RESV_FLD1"`
	MsgResvFld2 string    `gorm:"column:MSG_RESV_FLD2"`
	MsgResvFld3 string    `gorm:"column:MSG_RESV_FLD3"`
	MsgResvFld4 string    `gorm:"column:MSG_RESV_FLD4"`
	MsgResvFld5 string    `gorm:"column:MSG_RESV_FLD5"`
	CreatedAt   time.Time `gorm:"column:REC_CRT_TS"`
}

func (OrgDictionaryItem) TableName() string {
	return orgDictionaryItem
}

func SaveOrgDictionaryItem(db *gorm.DB, data *OrgDictionaryItem) error {
	return db.Save(data).Error
}

func DeleteOrgDictionaryItem(db *gorm.DB, query *OrgDictionaryItem) error {
	return db.Delete(&OrgDictionaryItem{}, query).Error
}

func ListOrgDictionaryItem(db *gorm.DB, query *OrgDictionaryItem, page int32, size int32) ([]*OrgDictionaryItem, int32, error) {
	var count int32
	var err error
	out := make([]*OrgDictionaryItem, 0)
	db.Model(&OrgDictionaryItem{}).Where(query).Count(&count)
	err = db.Where(query).Offset((page - 1) * size).Limit(size).Find(&out).Error
	return out, count, err
}
