package db

import "github.com/jinzhu/gorm"

// AccountInfo is the top sheet
type AccountInfo struct {
	gorm.Model
	Account         string `gorm:"not null;unique"` // 手机号
	Hashed_password string // 密码
	IsActive        bool   // 是否激活
	IsSuperUser     bool   // 超级管理员
	UserType        int    // 美容师 or 门店 or both
}

// PetGroomer 宠物美容师
type PetGroomer struct {
	gorm.Model
	Avatar   string // 头像
	NickName string // 昵称
	Rating   int    // 星级

	// 身份证信息
	Name         string // 实名
	IDCardNumber string // 身份证号
	IDCardFront  string // 身份证正面
	IDCardBack   string // 身份证背面
	IsVerified   bool   // 实名认证

	// 专业信息
	Qualification      int    // 资质
	IsCertifiedGroomer bool   // 人工审核认证美容师
	CertificateFrond   string // 资质证明正面
	CertificateBack    string // 资质证明背面

	Specialty string // 专业擅长: 猫咪清理，洁牙等
	AccountID uint   `gorm:"not null;unique"`
	Account   AccountInfo
}

// PetHouse 宠物门店
type PetHouse struct {
	gorm.Model
	Avatar   string // 头像
	NickName string // 昵称
	Rating   int    //星级

	// 身份证信息
	Name         string // 实名
	IDCardNumber string // 身份证号
	IDCardFront  string // 身份证正面
	IDCardBack   string // 身份证背面
	IsVerified   bool   // 实名认证

	// 专业信息
	Qualification     int    // 资质
	IsCertifiedHouse  bool   // 人工审核认证门店
	EnvironmentFrond  string // 门店门面照片
	EnvironmentInside string // 门店内部环境照片
	license           string // 营业执照照片
	Location          string // 门店地址

	WorkScope string // 业务范围: 猫咪洗护，猫狗寄养等
	AccountID uint   `gorm:"not null;unique"`
	Account   AccountInfo
}
