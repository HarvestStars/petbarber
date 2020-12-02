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

	PayMode      int     // 付费模式
	Basic        float32 // 底薪
	Commission   int     // 提成
	PayModeDesc  string  `gorm:"type:varchar(512)"`
	TotalPayment float32 // 总费用
	Desc         string  `gorm:"type:varchar(512)"`
	OrderType    int     // 订单类型 洗剪吹, 遛狗
	Status       int     // 订单状态
	MatchOrderID uint    // 美容师接单号
	UserID       uint    // account id
}

// 订单类型
const (
	WCB        = 1 // 洗剪吹
	WalkTheDog = 2 // 遛狗
	PickUp     = 3 // 接送
)

// 订单状态
const (
	CANCELORDER    = 1 // 商家取消订单
	CANCELPETHOUSE = 2 // 美容师取消订单
	CANCELGROOMER  = 3 // 商家取消美容师
	NEW            = 4 // 等待接单
	RUNNING        = 5 // 正在进行
	FINISHED       = 6 // 订单完成
)

// ToMatch 美容师接单
type ToMatch struct {
	ID              uint  `gorm:"primary_key" json:"id"`
	CreatedAt       int64 `json:"created_at"` // utc时间戳 精确到毫秒
	UpdatedAt       int64 `json:"updated_at"`
	Status          int   `json:"status"`            // 订单状态
	PethouseOrderID uint  `json:"pethouse_order_id"` // 门店订单号
	UserID          uint  `json:"user_id"`           // account id
}
