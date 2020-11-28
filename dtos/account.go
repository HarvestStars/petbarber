package dtos

// Account
type TuAccount struct {
	ID        uint  `gorm:"primary_key"`
	CreatedAt int64 // utc时间戳 精确到毫秒
	UpdatedAt int64

	Account        string `gorm:"not null;unique"` // 电话号码
	HashedPassword string // 密码
	IsActive       bool   // 是否激活
	IsSuperuser    bool   // 是否超管
	UserType       int    // 用户类型: 0 未知, 1 宠物店, 2 美容师
}
