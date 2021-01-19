package camunda

import (
	"time"

	"github.com/jinzhu/gorm"
)

type Action struct {
	ActionId   int64     `gorm:"column:action_id;primary_key"`
	Action     string    `gorm:"column:action"`
	Comment    string    `gorm:"column:comment"`
	TaskId     int64     `gorm:"column:task_id"`
	InstanceId int64     `gorm:"column:instance_id"`
	UserId     int64     `gorm:"column:user_id"`
	RoleId     int64     `gorm:"column:role_id"`
	CreatedAt  time.Time `gorm:"column:created_at"`
	UpdatedAt  time.Time `gorm:"column:updated_at"`
}

type Remark struct {
	Action
	UserName string `gorm:"column:USER_NAME"`
}

func (Action) TableName() string {
	return "TBL_CAMUNDA_ACTION"
}

func SaveRemark(db *gorm.DB, action *Action) error {
	return db.Create(action).Error
}

func QueryRemark(db *gorm.DB, query *Action, page int32, size int32) ([]*Remark, int32, error) {
	out := make([]*Remark, 0)
	var count int32
	db.Model(query).Where(query).Count(&count)
	err := db.Table("TBL_CAMUNDA_ACTION").Select("TBL_CAMUNDA_ACTION.*, u.USER_NAME").Joins("left join TBL_USER u on u.USER_ID = TBL_CAMUNDA_ACTION.user_id").Where(query).Offset((page - 1) * size).Limit(size).Scan(&out).Error
	return out, count, err
}
