package db

import "github.com/jinzhu/gorm"

// AccountInfo is the top sheet
type AccountInfo struct {
	gorm.Model
	Account     string // 手机号
	Avatar      string // 头像
	NickName    string // 昵称
	IsSuperUser bool   // 超级管理员
	UserType    int    // 美容师 or 门店 or both

	// 身份证信息
	Name         string // 实名
	IDCardNumber string // 身份证号
	IDCardFront  string // 身份证正面
	IDCardBack   string // 身份证背面
	IsVerified   bool   // 实名认证
}

// PetGroomer 宠物美容师
type PetGroomer struct {
	gorm.Model
	Rating int // 星级

	Qualification      int    // 资质
	IsCertifiedGroomer bool   // 人工审核认证美容师
	CertificateFrond   string // 资质证明正面
	CertificateBack    string // 资质证明背面

	Specialty string // 专业擅长: 猫咪清理，洁牙等
	AccountID uint
	Account   AccountInfo
}

// PetHouse 宠物门店
type PetHouse struct {
	gorm.Model
	Rating int

	Qualification     int    // 资质
	IsCertifiedHouse  bool   // 人工审核认证门店
	EnvironmentFrond  string // 门店门面照片
	EnvironmentInside string // 门店内部环境照片
	license           string // 营业执照照片

	Location  string // 门店地址
	WorkScope string // 业务范围: 猫咪洗护，猫狗寄养等
	AccountID uint
	Account   AccountInfo
}
