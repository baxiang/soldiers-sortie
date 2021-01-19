package user

import (
	"time"

	"github.com/jinzhu/gorm"
)

const (
	TableLeaguer = "TBL_LEAGUER"
)

type Leaguer struct {
	LeaguerNo     string    `gorm:"column:LEAGUER_NO"`
	LeaguerName   string    `gorm:"column:LEAGUER_NAME"`
	LeaguerType   string    `gorm:"column:LEAGUER_TYPE"`
	LeaguerInfo   string    `gorm:"column:LEAGUER_INFO"`
	LeaguerStatus int32     `gorm:"column:LEAGUER_STATUS"`
	Created       time.Time `gorm:"column:REC_CRT_TS"`
	Updated       time.Time `gorm:"column:REC_UPD_TS"`
}

func (l Leaguer) TableName() string {
	return TableLeaguer
}

func ListLeaguers(db *gorm.DB, query *Leaguer, page int32, size int32) ([]*Leaguer, int32, error) {
	leaguers := make([]*Leaguer, 0)
	var count int32 = 0
	db.Model(query).Where(query).Count(&count)
	err := db.Where(query).Limit(size).Offset((page - 1) * size).Find(&leaguers).Error
	if err == gorm.ErrRecordNotFound {
		return leaguers, count, nil
	}
	return leaguers, count, err
}
