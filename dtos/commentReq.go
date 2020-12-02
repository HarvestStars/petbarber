package dtos

// TComment 订单评论
type CommentReq struct {
	OrderID uint    `json:"order_id"`
	Favor   float32 `json:"favor"`
	Content string  `json:"content"`
}
