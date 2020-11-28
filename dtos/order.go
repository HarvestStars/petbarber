package dtos

// ToRequirement 门店派单
type ToRequirement struct {
	ID         uint  `gorm:"primary_key"`
	CreatedAt  int64 // utc时间戳 精确到毫秒
	UpdatedAt  int64
	StartedAt  int64 // 订单开始时间
	FinishedAt int64 // 完成时间

	ServiceBits      int64  // 1~6 个服务类型
	ServiceItemsDesc string `gorm:"type:varchar(512)"`

	PayMode     int    // 付费模式
	PayModeDesc string `gorm:"type:varchar(512)"`
	Desc        string `gorm:"type:varchar(512)"`
	OrderType   int    // 订单类型 洗剪吹, 遛狗
	Status      int    // 订单状态
	UserID      uint   // 门店id
}

// ToMatch 美容师接单
type ToMatch struct {
	ID              uint  `gorm:"primary_key"`
	CreatedAt       int64 // utc时间戳 精确到毫秒
	UpdatedAt       int64
	Status          int   // 订单状态
	PethouseOrderID int64 // 门店订单号
	UserID          uint  // 美容师id
}
