package user

import (
	"encoding/json"
	"time"

	"github.com/jinzhu/gorm"
)

const (
	TableUser           = "TBL_USER"
	TableAuthAssignment = "TBL_AUTH_ASSIGNMENT"
	TableAuthItemChild  = "TBL_AUTH_ITEM_CHILD"
	TableMenu           = "TBL_MENU"
	TableAuthItem       = "TBL_AUTH_ITEM"
	TableRole           = "TBL_ROLE"
	TableRoute          = "TBL_ROUTE"
	TablePermission     = "TBL_PERMISSION"
	TableUserInfo       = "TBL_USER_INFO"
)

type User struct {
	UserID             int64     `gorm:"column:USER_ID;primary_key"`
	LeaguerNO          string    `gorm:"column:LEAGUER_NO"`
	UserName           string    `gorm:"column:USER_NAME;unique"`
	AuthKey            string    `gorm:"column:AUTH_KEY"`
	PasswordHash       string    `gorm:"column:PASSWORD_HASH"`
	PasswordHashNew    string    `gorm:"column:PASSWORD_HASH_NEW"`
	PasswordResetToken string    `gorm:"column:PASSWORD_RESET_TOKEN"`
	Email              string    `gorm:"column:EMAIL"`
	UserType           string    `gorm:"column:USER_TYPE"`
	UserGroupNo        string    `gorm:"column:USER_GROUP_NO"`
	UserInfo           string    `gorm:"column:USER_INFO"`
	UserStatus         int64     `gorm:"column:USER_STATUS"`
	UserNotice         string    `gorm:"column:USER_NOTICE"`
	CreatedAt          time.Time `gorm:"column:REC_CRT_TS"`
	UpdatedAt          time.Time `gorm:"column:REC_UPD_TS"`
	ParentUserName     *string   `gorm:"column:PARENT_USER_NAME"`
}

func (u User) TableName() string {
	return TableUser
}

type AuthAssignment struct {
	ItemName string `gorm:"column:ITEM_NAME"`
	UserID   int64  `gorm:"column:USER_ID"`
	UserName string `gorm:"column:USER_NAME"`
}

func (a AuthAssignment) TableName() string {
	return TableAuthAssignment
}

type AuthItemChild struct {
	Parent string `gorm:"column:PARENT"`
	Child  string `gorm:"column:CHILD"`
}

func (a AuthItemChild) TableName() string {
	return TableAuthItemChild
}

type Menu struct {
	ID        int32  `gorm:"column:ID;primary_key"`
	Name      string `gorm:"column:NAME"`
	Parent    int32  `gorm:"column:PARENT"`
	MenuRoute string `gorm:"column:MENU_ROUTE"`
	MenuOrder int32  `gorm:"column:MENU_ORDER"`
	MenuData  string `gorm:"column:MENU_DATA"`
}

func (m Menu) TableName() string {
	return TableMenu
}

type AuthItem struct {
	Name        string `gorm:"column:NAME"`
	Type        int64  `gorm:"column:TYPE"`
	Description string `gorm:"column:DESCRIPTION"`
	RuleName    string `gorm:"column:RULE_NAME"`
	ItemData    string `gorm:"column:ITEM_DATA"`
	ItemCode    string `gorm:"column:ITEM_CODE"`
	ParentItem  string `gorm:"column:PARENT_ITEM"`
	CreateUser  int64  `gorm:"column:CREATE_USER"`
}

func (a AuthItem) TableName() string {
	return TableAuthItem
}

type Role struct {
	ID        int64     `gorm:"column:ROLE_ID;primary_key"`
	Role      string    `gorm:"column:ROLE_NAME"`
	CreatedAt time.Time `gorm:"column:CREATED_AT"`
	UpdatedAt time.Time `gorm:"column:UPDATED_AT"`
}

func (a Role) TableName() string {
	return TableRole
}

type Permissions struct {
	Menus []*Menu
	Items []*AuthItem
}

func (p Permissions) MarshalBinary() ([]byte, error) {
	return json.Marshal(p)
}

type Route struct {
	ID        int32     `gorm:"column:ROUTE_ID;primary_key"`
	Name      string    `gorm:"column:ROUTE_NAME"`
	CreatedAt time.Time `gorm:"column:CREATED_AT"`
	UpdatedAt time.Time `gorm:"column:UPDATED_AT"`
}

func (r Route) TableName() string {
	return TableRoute
}

type Permission struct {
	ID        int64     `gorm:"column:PERMISSION_ID;primary_key"`
	Name      string    `gorm:"column:PERMISSION_NAME"`
	CreatedAt time.Time `gorm:"column:CREATED_AT"`
	UpdatedAt time.Time `gorm:"column:UPDATED_AT"`
}

func (r Permission) TableName() string {
	return TablePermission
}

// 根据用户名查询用户信息
// 入参
// userName: TBL_USER表中的USER_NAME字段
// 返回
// *User: TBL_USER表中的用户信息
func FindUserByUserName(db *gorm.DB, userName string) (*User, error) {
	user := &User{
		UserName: userName,
	}
	err := db.Where(&User{UserName: userName}).First(user).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return user, err
}

func getAuthRoutes(db *gorm.DB, itemNames []string) ([]*AuthItemChild, error) {
	results := make([]*AuthItemChild, 0)
	err := db.Where("PARENT in (?)", itemNames).Find(&results).Error
	if err == gorm.ErrRecordNotFound {
		err = nil
		return results, err
	}
	names := make([]string, 0)
	for _, result := range results {
		names = append(names, result.Child)
	}
	if len(names) != 0 {
		childResults, err := getAuthRoutes(db, names)
		if err != nil {
			return nil, err
		}
		results = append(results, childResults...)
	}
	return results, err
}

func GetAuthMenu(db *gorm.DB, items []string, like bool) ([]*Menu, error) {
	menus := make([]*Menu, 0)
	var err error
	if like {
		for _, item := range items {
			ms := make([]*Menu, 0)
			err = db.Where("MENU_ROUTE like ?", item).Find(&ms).Error
			if err != nil {
				return nil, err
			}
			menus = append(menus, ms...)
		}
	} else {
		err = db.Where("MENU_ROUTE in (?)", items).Find(&menus).Error
	}
	if err != nil {
		return nil, err
	}

	parents := make(map[int32]bool)
	if len(menus) != 0 {
		for _, menu := range menus {
			if menu.Parent != 0 {
				parents[menu.Parent] = true
			}
		}
	}

	parentIDs := make([]int32, 0, len(parents))
	for id := range parents {
		parentIDs = append(parentIDs, id)
	}

	rootMenus := make([]*Menu, 0)
	err = db.Where("ID in (?)", parentIDs).Find(&rootMenus).Error
	if err == gorm.ErrRecordNotFound {
		err = nil
	}
	if err != nil {
		return nil, err
	}

	menus = append(menus, rootMenus...)
	return menus, nil
}

func SaveUser(db *gorm.DB, user *User) (*User, error) {
	err := db.Create(user).Error
	if err != nil {
		return nil, err
	}

	var findUser User
	err = db.Where(&User{UserName: user.UserName}).First(&findUser).Error
	return &findUser, err
}

func FindRole(db *gorm.DB, role string) (*Role, error) {
	var r Role
	err := db.Where(&Role{Role: role}).First(&r).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &r, err
}

func SaveRole(db *gorm.DB, role *Role) error {
	return db.Create(role).Error
}

func FindRoutesByNames(db *gorm.DB, names []string) ([]*Route, error) {
	rs := make([]*Route, 0)
	err := db.Where("ROUTE_NAME in (?)", names).Find(&rs).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return rs, nil
}

func SaveRoutes(db *gorm.DB, names []string) error {
	db = db.Begin()
	var err error
	defer db.Rollback()
	for _, name := range names {
		err = db.Create(&Route{
			Name: name,
		}).Error
		if err != nil {
			return err
		}
	}
	return db.Commit().Error
}

func ListRoutes(db *gorm.DB) ([]*Route, error) {
	rs := make([]*Route, 0)
	err := db.Find(&rs).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	return rs, nil
}

func SavePermission(db *gorm.DB, name string) error {
	return db.Create(&Permission{
		Name: name,
	}).Error
}

func UpdatePermission(db *gorm.DB, id int64, permission *Permission) error {
	return db.Model(&Permission{
		ID: id,
	}).Update(permission).Error
}

func FindPermissionByName(db *gorm.DB, name string) (*Permission, error) {
	p := new(Permission)
	err := db.Where(&Permission{Name: name}).First(p).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}

	return p, err
}

func FindPermissionByID(db *gorm.DB, id int64) (*Permission, error) {
	p := new(Permission)
	err := db.Where(&Permission{ID: id}).First(p).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return p, err
}

func FindRouteByName(db *gorm.DB, route string) (*Route, error) {
	r := new(Route)

	err := db.Where(&Route{Name: route}).First(r).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return r, err
}

func FindRoutesByName(db *gorm.DB, routes []string) ([]*Route, error) {
	rs := make([]*Route, 0)

	err := db.Where("ROUTE_NAME in (?)", routes).Find(&rs).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	return rs, nil
}

func DeletePermission(db *gorm.DB, permission *Permission) error {
	return db.Delete(permission).Error
}

func ListPermissions(db *gorm.DB, query *Permission, page int32, size int32) ([]*Permission, int32, error) {
	ps := make([]*Permission, 0)
	var count int32 = 0
	db.Model(query).Where(query).Count(&count)

	err := db.Where(query).Limit(size).Offset((page - 1) * size).Find(&ps).Error
	if err == gorm.ErrRecordNotFound {
		return ps, count, nil
	}
	return ps, count, err
}

func ListRole(db *gorm.DB, query *Role, page int32, size int32) ([]*Role, int32, error) {
	roles := make([]*Role, 0)
	var count int32 = 0
	db.Model(query).Where(query).Count(&count)

	err := db.Where(query).Limit(size).Offset((page - 1) * size).Find(&roles).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, 0, err
	}

	return roles, count, nil
}

func FindRoleByID(db *gorm.DB, id int64) (*Role, error) {
	role := new(Role)
	err := db.Where(&Role{ID: id}).First(role).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return role, err
}

func UpdateRole(db *gorm.DB, id int64, role *Role) error {
	return db.Model(&Role{
		ID: id,
	}).Update(role).Error
}

func DeleteRole(db *gorm.DB, role *Role) error {
	return db.Delete(role).Error
}

func ListUsers(db *gorm.DB, query *User, page, size int32) ([]*User, int32, error) {
	us := make([]*User, 0)
	var count int32 = 0
	db.Model(&User{}).Where(query).Count(&count)
	err := db.Where(query).Limit(size).Offset((page - 1) * size).Find(&us).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, 0, err
	}
	return us, count, nil
}

func UpdateUser(db *gorm.DB, id int64, user *User) error {
	return db.Model(&User{UserID: id}).Update(user).Error
}

func FindUserByID(db *gorm.DB, id int64) (*User, error) {
	user := &User{}
	err := db.Where("USER_ID = ?", id).First(user).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return user, err
}

func ListMenus(db *gorm.DB, query *Menu, page int32, size int32) ([]*Menu, int32, error) {

	menus := make([]*Menu, 0)
	var count int32 = 0
	db.Model(query).Where(query).Count(&count)

	err := db.Where(query).Limit(size).Offset((page - 1) * size).Find(&menus).Error
	if err == gorm.ErrRecordNotFound {
		return menus, count, nil
	}
	return menus, count, err
}

func ListMenusByIDs(db *gorm.DB, ids []int32) ([]*Menu, error) {
	menus := make([]*Menu, 0)

	err := db.Where("ID in (?)", ids).Find(&menus).Error
	if err == gorm.ErrRecordNotFound {
		return menus, nil
	}
	return menus, err
}

func FindMenuByName(db *gorm.DB, name string) (*Menu, error) {
	menu := new(Menu)
	err := db.Where(&Menu{Name: name}).Find(menu).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return menu, err
}

func FindMenuByID(db *gorm.DB, id int32) (*Menu, error) {
	menu := new(Menu)
	err := db.Where(&Menu{ID: id}).Find(menu).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return menu, err
}

func SaveMenu(db *gorm.DB, menu *Menu) error {
	return db.Create(menu).Error
}

func DeleteMenu(db *gorm.DB, id int32) error {
	return db.Delete(&Menu{ID: id}).Error
}

type UserInfo struct {
	UserType string `gorm:"column:USER_TYPE"`
	UserDesc string `gorm:"column:USER_DESC"`
	UserInfo string `gorm:"column:USER_INFO"`
}

func (t UserInfo) TableName() string {
	return TableUserInfo
}

func ListUserInfo(db *gorm.DB) ([]*UserInfo, error) {
	infos := make([]*UserInfo, 0)
	err := db.Find(&infos).Error
	if err == gorm.ErrRecordNotFound {
		err = nil
	}
	return infos, err
}

func FindRolesByIds(db *gorm.DB, ids []int64) ([]*Role, error) {
	roles := make([]*Role, 0)
	err := db.Find(&roles, "ROLE_ID in (?)", ids).Error
	if err == gorm.ErrRecordNotFound {
		return roles, nil
	}
	return roles, err
}

func FindPermissionsByIds(db *gorm.DB, ids []int64) ([]*Permission, error) {
	permissions := make([]*Permission, 0)
	err := db.Find(&permissions, "PERMISSION_ID in (?)", ids).Error
	if err == gorm.ErrRecordNotFound {
		return permissions, nil
	}
	return permissions, err
}

func DeleteRoute(db *gorm.DB, data *Route) error {
	return db.Where(data).Delete(data).Error
}
