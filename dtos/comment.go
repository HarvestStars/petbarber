package dtos

// TComment 订单评论
type TComment struct {
	ID              uint    `gorm:"primary_key" json:"id"`
	CreatedAt       int64   `json:"created_at"` // utc时间戳 精确到毫秒
	Status          int     `json:"status"`     //1=新评论 2=已审核 3=已拒绝
	FromUserID      uint    `json:"from_user_id"`
	ToUserID        uint    `json:"to_user_id"`
	CommentType     int     `json:"comment_type"` // 评论类型 1=门店评价美容师 2=美容师评价门店
	Favor           float32 `json:"favor"`        // 评分
	Content         string  `gorm:"type:text" json:"content"`
	GroomerOrderID  uint    `json:"groomer_order_id"`  // 美容师接单号
	PethouseOrderID uint    `json:"pethouse_order_id"` // 门店订单号
}
