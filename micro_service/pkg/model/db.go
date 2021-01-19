package model

import (
	"fmt"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/sirupsen/logrus"
)

type Options struct {
	User     string
	Password string
	DB       string
	Addr     string
}

func NewDB(opt *Options) (*gorm.DB, error) {
	url := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=true&loc=Local",
		opt.User,
		opt.Password,
		opt.Addr,
		opt.DB,
	)
	logrus.Debugln("mysql url:", url)
	db, err := gorm.Open(
		"mysql",
		url,
	)
	if err != nil {
		return nil, err
	}

	db.DB().SetMaxIdleConns(10)
	db.DB().SetMaxOpenConns(100)
	return db, nil
}
