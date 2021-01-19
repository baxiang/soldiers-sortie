package main

import (
	"fmt"
	"userService/pkg/model/user"

	"github.com/casbin/casbin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"

	"userService/pkg/gormadapter"
)

var adapter = gormadapter.NewAdapter("mysql", "root:root@tcp(127.0.0.1:3306)/")
var enforcer = casbin.NewEnforcer("configs/rbac.conf", adapter)

type AuthModel struct {
	Role string
	Data string
}

type UserModel struct {
	User string
	Role string
}

func (a AuthModel) String() string {
	return fmt.Sprintf(`{"role": "%s", "data": "%s"}`, a.Role, a.Data)
}

func main() {
	initCasbin()
	//checkValid()
}

func slice2Map(ss []string) map[string]bool {
	m := make(map[string]bool)
	for _, s := range ss {
		m[s] = true
	}
	return m
}

func initCasbin() {
	db, err := gorm.Open("mysql", "root:root@tcp(127.0.0.1:3306)/apsmgm?charset=utf8&parseTime=true")
	if err != nil {
		logrus.Fatal(err)
	}

	names := make([]string, 0)
	db.Table("TBL_AUTH_ITEM").Pluck("NAME", &names)
	logrus.Infoln(names)
	logrus.Infoln(len(names))

	roles := make([]string, 0)
	db.Table("TBL_AUTH_ITEM_CHILD").Pluck("distinct PARENT", &roles)
	logrus.Infoln(roles)
	logrus.Infoln(len(roles))
	roleMap := slice2Map(roles)

	assigns := make([]user.AuthAssignment, 0)
	db.Table("TBL_AUTH_ASSIGNMENT").
		Where("TBL_USER.USER_ID is not null").
		Select("TBL_AUTH_ASSIGNMENT.ITEM_NAME, TBL_USER.USER_ID, TBL_USER.USER_NAME").
		Joins("left join TBL_USER on TBL_AUTH_ASSIGNMENT.USER_ID = TBL_USER.USER_ID").Scan(&assigns)

	ms := make([]AuthModel, 0)
	us := make([]UserModel, 0)

	for _, role := range roles {
		children := make([]string, 0)
		db.Table("TBL_AUTH_ITEM_CHILD").Where("PARENT = ?", role).Pluck("CHILD", &children)
		logrus.Infoln(role, children)
		for _, c := range children {
			if roleMap[c] {
				us = append(us, UserModel{User: role, Role: c})
			} else {
				ms = append(ms, AuthModel{Role: role, Data: c})
			}
		}
	}

	for _, assign := range assigns {
		us = append(us, UserModel{Role: assign.ItemName, User: assign.UserName})
	}

	logrus.Infoln(ms)
	logrus.Infoln(len(ms))

	for _, m := range ms {
		logrus.Infoln(m.Role)
		logrus.Infoln(m.Data)
		enforcer.AddPolicy(m.Role, m.Data)

	}

	db = db.Debug()
	db = db.Begin()
	defer db.Rollback()
	for _, role := range roles {
		err = user.SaveRole(db, &user.Role{
			Role: role,
		})
		if err != nil {
			return
		}
	}
	db.Commit()

	//for _ := range us {
	//logrus.Infoln(u.User)
	//logrus.Infoln(u.Role)
	//enforcer.AddRoleForUser(u.User, u.Role)
	//}
}

func checkValid() {
	//enforcer.AddRoleForUser("jingx", "dsafdsa")
	//enforcer.AddRoleForUser("jingx", "管理员权限")
	//ok := enforcer.Enforce("zymcxy", "/clear/downinsdz")
	//logrus.Infoln(ok)
	roles := enforcer.GetAllRoles()
	logrus.Infoln(len(roles))
}
