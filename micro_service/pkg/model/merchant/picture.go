package merchant

import (
	"time"

	"github.com/jinzhu/gorm"
)

type Picture struct {
	FileId     string    `gorm:"column:FILE_ID;primary_key"`
	MchtCd     string    `gorm:"column:MCHT_CD"`
	DocType    string    `gorm:"column:DOC_TYPE"`
	FileType   string    `gorm:"column:FILE_TYPE"`
	FileName   string    `gorm:"column:FILE_NAME"`
	PIndex     int64     `gorm:"column:PINDEX"`
	PCode      string    `gorm:"column:PCODE"`
	Url        string    `gorm:"column:URL"`
	SystemFlag string    `gorm:"column:SYSTEMFLAG"`
	Status     string    `gorm:"column:STATUS"`
	RecOprId   string    `gorm:"column:REC_OPR_ID"`
	RecUpdOpr  string    `gorm:"column:REC_UPD_OPR"`
	CreatedAt  time.Time `gorm:"column:REC_CRT_TS"`
	UpdatedAt  time.Time `gorm:"column:REC_UPD_TS"`
}

type PictureMain struct {
	Picture
}

func (Picture) TableName() string {
	return "TBL_EDIT_MCHT_PICTURE"
}

func (PictureMain) TableName() string {
	return "TBL_MCHT_PICTURE"
}

func SavePicture(db *gorm.DB, data *Picture) error {
	return db.Save(data).Error
}

func SavePictureMain(db *gorm.DB, data *PictureMain) error {
	return db.Save(data).Error
}

func QueryPicture(db *gorm.DB, query *Picture, page int32, size int32) ([]*Picture, int32, error) {
	out := make([]*Picture, 0)
	var count int32
	db.Model(&Picture{}).Where(query).Count(&count)
	err := db.Where(query).Offset((page - 1) * size).Limit(size).Find(&out).Error
	return out, count, err
}

func QueryPictureMain(db *gorm.DB, query *PictureMain, page int32, size int32) ([]*PictureMain, int32, error) {
	out := make([]*PictureMain, 0)
	var count int32
	db.Model(&PictureMain{}).Where(query).Count(&count)
	err := db.Where(query).Offset((page - 1) * size).Limit(size).Find(&out).Error
	return out, count, err
}

func DeletePicture(db *gorm.DB, query *Picture) error {
	return db.Where(query).Delete(&Picture{}).Error
}
