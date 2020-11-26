package db

import "github.com/jinzhu/gorm"

// PetGroomer 宠物美容师
type TuGroomer struct {
	gorm.Model
	Avatar   string  `gorm:"text"` // 头像图片路径
	NickName string  // 昵称
	Rating   float32 // 星级

	// 身份证信息
	Name         string // 实名
	IDCardNumber string // 身份证号
	IDCardFront  string `gorm:"text"` // 身份证正面图片路径
	IDCardBack   string `gorm:"text"` // 身份证背面图片路径
	IsVerified   bool   `gorm:"text"` // 是否实名认证

	// 专业信息
	Qualification      int    // 资质
	IsCertifiedGroomer bool   // 美容师是否通过人工审核认证
	CertificateFront   string `gorm:"text"` // 资质证明正面图片路径
	CertificateBack    string `gorm:"text"` // 资质证明背面图片路径

	Specialty string `gorm:"text"` // 专业擅长: 猫咪清理，洁牙等
	AccountID uint   `gorm:"not null;unique"`
}

// PetHouse 宠物门店
type TuPethouse struct {
	gorm.Model
	Avatar   string  `gorm:"text"` // 头像图片路径
	NickName string  // 昵称
	Rating   float32 // 星级

	// 身份证信息
	Name         string // 实名
	IDCardNumber string // 身份证号
	IDCardFront  string `gorm:"text"` // 身份证正面图片路径
	IDCardBack   string `gorm:"text"` // 身份证背面图片路径
	IsVerified   bool   `gorm:"text"` // 实名认证

	// 专业信息
	IsCertifiedHouse  bool   // 人工审核认证门店
	EnvironmentFront  string `gorm:"text"` // 门店门面照片路径
	EnvironmentInside string `gorm:"text"` // 门店内部环境照片路径
	License           string `gorm:"text"` // 营业执照照片路径
	Location          string `gorm:"text"` // 门店地址

	WorkScope string `gorm:"text"` // 业务范围: 猫咪洗护，猫狗寄养等
	AccountID uint   `gorm:"not null;unique"`
}
