package db

import "github.com/jinzhu/gorm"

// AccountInfo is the top sheet
type AccountInfo struct {
	gorm.Model
	Account     string
	HashedPWD   string
	IsActive    bool
	IsSuperUser bool
	UserType    int
}

// PetGroomer 宠物美容师
type PetGroomer struct {
	gorm.Model
	Name          string
	Avatar        string // 头像
	Rating        int
	Qualification int
	Status        int
	AccountID     uint
	Account       AccountInfo
}

// PetHouse 宠物门店
type PetHouse struct {
	gorm.Model
	Name          string
	Avatar        string // 头像
	Rating        int
	Qualification int
	Status        int
	Location      string // 门店地址
	AccountID     uint
	Account       AccountInfo
}
