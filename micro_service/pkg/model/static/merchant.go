package static

import "github.com/jinzhu/gorm"

type MerchantFirstThree struct {
	DicCode string `gorm:"column:DIC_CODE"`
	DicName string `gorm:"column:DIC_NAME"`
}

func FindMerchantFirstThree(db *gorm.DB, orgCode string) ([]*MerchantFirstThree, error) {
	out := make([]*MerchantFirstThree, 0)
	err := db.Table(tableDictItem+" a").
		Joins("left join TBL_ORG_DICTIONARYITEM b on  a.DIC_CODE = b.ITEM_CODE").
		Select("DIC_CODE, DIC_NAME").
		Where("ORG_CODE = ? and DIC_TYPE = ?",
			orgCode,
			"MCHT_FIRST_3",
		).Scan(&out).Error
	return out, err
}
