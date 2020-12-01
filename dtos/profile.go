package dtos

// PetHouse 宠物门店
type TuPethouse struct {
	//gorm.Model
	ID        uint  `gorm:"primary_key" json:"id"`
	CreatedAt int64 `json:"created_at"` // utc时间戳 精确到毫秒
	UpdatedAt int64 `json:"updated_at"`

	Avatar   string  `gorm:"type:text" json:"avatar"` // 头像图片路径
	NickName string  `json:"nick_name"`               // 昵称
	Favor    float32 `json:"favor"`                   // 星级评分
	Status   int     `json:"status"`                  // 状态

	// 身份证信息
	Name         string `json:"name"`                           // 实名
	IDCardNumber string `json:"id_card_number"`                 // 身份证号
	IDCardFront  string `gorm:"type:text" json:"id_card_front"` // 身份证正面图片路径
	IDCardBack   string `gorm:"type:text" json:"id_card_back"`  // 身份证背面图片路径
	IsVerified   bool   `json:"is_verified"`                    // 实名认证

	// 专业信息
	IsCertifiedHouse  bool   `json:"is_certified_house"`                  // 人工审核认证门店
	EnvironmentFront  string `gorm:"type:text" json:"environment_front"`  // 门店门面照片路径
	EnvironmentInside string `gorm:"type:text" json:"environment_inside"` // 门店内部环境照片路径
	License           string `gorm:"type:text" json:"license"`            // 营业执照照片路径
	Location          string `gorm:"type:text" json:"location"`           // 门店地址

	WorkScope string `gorm:"type:text" json:"work_scope"`       // 业务范围: 猫咪洗护，猫狗寄养等
	AccountID uint   `gorm:"not null;unique" json:"account_id"` // account_user表主键
}

// PetGroomer 宠物美容师
type TuGroomer struct {
	//gorm.Model
	ID        uint  `gorm:"primary_key" json:"id"`
	CreatedAt int64 `json:"created_at"` // utc时间戳 精确到毫秒
	UpdatedAt int64 `json:"updated_at"`

	Avatar   string  `gorm:"type:text" json:"avatar"` // 头像图片路径
	NickName string  `json:"nick_name"`               // 昵称
	Favor    float32 `json:"favor"`                   // 星级评分
	Status   int     `json:"status"`                  // 状态

	// 身份证信息
	Name         string `json:"name"`                           // 实名
	IDCardNumber string `json:"id_card_number"`                 // 身份证号
	IDCardFront  string `gorm:"type:text" json:"id_card_front"` // 身份证正面图片路径
	IDCardBack   string `gorm:"type:text" json:"id_card_back"`  // 身份证背面图片路径
	IsVerified   bool   `json:"is_verified"`                    // 是否实名认证

	// 专业信息
	Qualification      int    `json:"qualification"`                      // 资质
	IsCertifiedGroomer bool   `json:"is_certified_groomer"`               // 美容师是否通过人工审核认证
	CertificateFront   string `gorm:"type:text" json:"certificate_front"` // 资质证明正面图片路径
	CertificateBack    string `gorm:"type:text" json:"certificate_back"`  // 资质证明背面图片路径

	Specialty string `gorm:"type:text" json:"specialty"`        // 专业擅长: 猫咪清理，洁牙等
	AccountID uint   `gorm:"not null;unique" json:"account_id"` // account_user表主键
}
