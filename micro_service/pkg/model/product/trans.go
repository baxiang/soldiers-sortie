package product

import "github.com/jinzhu/gorm"

type Trans struct {
	ProdCd      string `gorm:"column:PROD_CD"`
	BizCd       string `gorm:"column:BIZ_CD"`
	TransCd     string `gorm:"column:TRANS_CD"`
	UpdateDate  string `gorm:"column:UPDATE_DATE"`
	Description string `gorm:"column:DESCRIPTION"`
	ResvFld1    string `gorm:"column:RESV_FLD1"`
	ResvFld2    string `gorm:"column:RESV_FLD2"`
	ResvFld3    string `gorm:"column:RESV_FLD3"`
}

func (Trans) TableName() string {
	return "TBL_PROD_BIZ_TRANS_MAP"
}

func ListTrans(db *gorm.DB, query *Trans) ([]*Trans, error) {
	out := make([]*Trans, 0)
	err := db.Where(query).Find(&out).Error
	return out, err
}
