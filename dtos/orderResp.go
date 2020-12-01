package dtos

// 门店派单响应
type PCOrderResp struct {
	ID           uint        `json:"id"` // torequirement 订单号
	StartedAt    int64       `json:"started_at"`
	FinishedAt   int64       `json:"finished_at"`
	CreateAt     int64       `json:"created_at"`
	OrderType    int         `json:"order_type"`
	Status       int         `json:"status"`
	ServiceItems []int       `json:"service_items"`
	Payment      PaymentInfo `json:"payment"`
	Children     Children    `json:"children"`
	UserID       uint        `json:"user_id"`
}

type PaymentInfo struct {
	Mode   int    `json:"mode"`
	Detail Detail `json:"detail"`
}

type Detail struct {
	Basic      float32 `json:"basic"`
	Commission int     `json:"commission"`
}

type Children struct {
	MatchOrder ToMatch   `json:"match_order"`
	Groomer    TuGroomer `json:"groomer"`
}

func (order *PCOrderResp) RespTransfer(requirementOrder ToRequirement, matchOrder ToMatch, groomer TuGroomer) error {
	order.ID = requirementOrder.ID
	order.StartedAt = requirementOrder.StartedAt
	order.FinishedAt = requirementOrder.FinishedAt
	order.CreateAt = requirementOrder.CreatedAt
	order.OrderType = requirementOrder.OrderType
	order.Status = requirementOrder.Status
	order.ServiceItems = ToServiceItems(requirementOrder.ServiceBits)
	payModeInt, err := ToPayMode(requirementOrder.Basic, requirementOrder.Commission)
	if err != nil {
		return err
	}
	detail := Detail{Basic: requirementOrder.Basic, Commission: requirementOrder.Commission}
	order.Payment = PaymentInfo{Mode: payModeInt, Detail: detail}
	order.Children = Children{MatchOrder: matchOrder, Groomer: groomer}
	order.UserID = requirementOrder.UserID
	return nil
}

// 美容师接单响应
type PCMatchResp struct {
	ID              uint   `json:"id"`
	Status          int    `json:"status"`
	CreateAt        int64  `json:"created_at"`
	UpdatedAt       int64  `json:"updated_at"`
	Parent          Parent `json:"parent"`
	PethouseOrderID uint   `json:"pethouse_order_id"`
	UserID          uint   `json:"user_id"`
}

type Parent struct {
	ID           uint        `json:"id"`
	StartedAt    int64       `json:"started_at"`
	FinishedAt   int64       `json:"finished_at"`
	OrderType    int         `json:"order_type"`
	CreateAt     int64       `json:"created_at"`
	UpdatedAt    int64       `json:"updated_at"`
	ServiceItems []int       `json:"service_items"`
	Payment      PaymentInfo `json:"payment"`
}

func (order *PCMatchResp) RespTransfer(matchOrder ToMatch, requirementOrder ToRequirement) error {
	order.ID = matchOrder.ID
	order.Status = matchOrder.Status
	order.CreateAt = matchOrder.CreatedAt
	order.UpdatedAt = matchOrder.UpdatedAt
	payModeInt, err := ToPayMode(requirementOrder.Basic, requirementOrder.Commission)
	if err != nil {
		return err
	}
	detail := Detail{Basic: requirementOrder.Basic, Commission: requirementOrder.Commission}
	order.Parent = Parent{
		ID:           requirementOrder.ID,
		StartedAt:    requirementOrder.StartedAt,
		FinishedAt:   requirementOrder.FinishedAt,
		OrderType:    requirementOrder.OrderType,
		CreateAt:     requirementOrder.CreatedAt,
		UpdatedAt:    requirementOrder.UpdatedAt,
		ServiceItems: ToServiceItems(requirementOrder.ServiceBits),
		Payment:      PaymentInfo{Mode: payModeInt, Detail: detail},
	}
	order.PethouseOrderID = requirementOrder.ID
	order.UserID = matchOrder.UserID
	return nil
}
