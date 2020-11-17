package db

import "github.com/jinzhu/gorm"

// PetGroomer 宠物美容师
type TuGroomer struct {
	gorm.Model
	Avatar   string  // 头像图片路径
	NickName string  // 昵称
	Rating   float32 // 星级

	// 身份证信息
	Name         string // 实名
	IDCardNumber string // 身份证号
	IDCardFront  string // 身份证正面图片路径
	IDCardBack   string // 身份证背面图片路径
	IsVerified   bool   // 是否实名认证

	// 专业信息
	Qualification      int    // 资质
	IsCertifiedGroomer bool   // 美容师是否通过人工审核认证
	CertificateFront   string // 资质证明正面图片路径
	CertificateBack    string // 资质证明背面图片路径

	Specialty string // 专业擅长: 猫咪清理，洁牙等
	AccountID uint   `gorm:"not null;unique"`
}

// PetHouse 宠物门店
type TuPethouse struct {
	gorm.Model
	Avatar   string  // 头像图片路径
	NickName string  // 昵称
	Rating   float32 // 星级

	// 身份证信息
	Name         string // 实名
	IDCardNumber string // 身份证号
	IDCardFront  string // 身份证正面图片路径
	IDCardBack   string // 身份证背面图片路径
	IsVerified   bool   // 实名认证

	// 专业信息
	IsCertifiedHouse  bool   // 人工审核认证门店
	EnvironmentFront  string // 门店门面照片路径
	EnvironmentInside string // 门店内部环境照片路径
	License           string // 营业执照照片路径
	Location          string // 门店地址

	WorkScope string // 业务范围: 猫咪洗护，猫狗寄养等
	AccountID uint   `gorm:"not null;unique"`
}
