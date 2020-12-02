package dtos

// TComment 订单评论
type TComment struct {
	ID              uint  `gorm:"primary_key"`
	CreatedAt       int64 // utc时间戳 精确到毫秒
	Status          int   //1=新评论 2=已审核 3=已拒绝
	FromUserID      uint
	ToUserID        uint
	CommentType     int     // 评论类型 1=评价门店 2=评价美容师
	Favor           float32 // 评分
	Content         string  `gorm:"type:text"`
	GroomerOrderID  uint    // 美容师接单号
	PethouseOrderID uint    // 门店订单号
}
