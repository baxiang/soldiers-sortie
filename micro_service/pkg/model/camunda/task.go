package camunda

import (
	"time"

	"github.com/jinzhu/gorm"
)

type Task struct {
	TaskId        int64     `gorm:"column:task_id;primary_key"`
	InstanceId    int64     `gorm:"column:instance_id"`
	Title         string    `gorm:"column:title"`
	RoleId        int64     `gorm:"column:role_id"`
	CurrentNode   string    `gorm:"column:current_node"`
	CamundaTaskId string    `gorm:"column:camunda_task_id"`
	EndFlag       *bool     `gorm:"column:end_flag"`
	UserName      string    `gorm:"column:username"`
	WorkflowName  string    `gorm:"column:workflow_name"`
	DataId        string    `gorm:"column:data_id"`
	CreatedAt     time.Time `gorm:"column:created_at"`
	UpdatedAt     time.Time `gorm:"column:updated_at"`
}

func (t Task) TableName() string {
	return "TBL_CAMUNDA_TASK"
}

func QueryTask(db *gorm.DB, query *Task, page int32, size int32) ([]*Task, int32, error) {
	out := make([]*Task, 0)
	var count int32
	db.Model(&Task{}).Where(query).Count(&count)
	err := db.Where(query).Offset((page - 1) * size).Limit(size).Find(&out).Error
	if err == gorm.ErrRecordNotFound {
		return out, count, nil
	}
	return out, count, err
}

func SaveTask(db *gorm.DB, task *Task) error {
	return db.Save(task).Error
}

func UpdateTask(db *gorm.DB, query *Task, task *Task) error {
	return db.Model(task).Where(query).Updates(task).Error
}

func FindTaskById(db *gorm.DB, id int64) (*Task, error) {
	out := new(Task)
	err := db.Where("task_id = ?", id).Take(out).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return out, err
}

func FindTaskByCamundaId(db *gorm.DB, id string) (*Task, error) {
	out := new(Task)
	err := db.Where("camunda_task_id = ?", id).Take(out).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return out, err
}

func FindTaskByRoles(db *gorm.DB, roles []int64, page int32, size int32) ([]*Task, int32, error) {
	out := make([]*Task, 0)
	var count int32
	db.Model(&Task{}).Where("role_id in (?) and end_flag = false", roles).Count(&count)
	err := db.Order("updated_at desc").Where("role_id in (?) and end_flag = false", roles).Offset((page - 1) * size).Limit(size).Find(&out).Error
	return out, count, err
}
