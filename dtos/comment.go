package dtos

// TComment 订单评论
type TComment struct {
	ID          uint  `gorm:"primary_key"`
	CreatedAt   int64 // utc时间戳 精确到毫秒
	Status      int
	FromUserID  uint
	ToUserID    uint
	CommentType int
	Favor       float32
	Content     string `gorm:"type:text"`
}
