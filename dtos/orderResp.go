package dtos

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

func (order *PCOrderResp) RespTransfer(requirementOrder ToRequirement, matchOrder ToMatch, groomer TuGroomer, pcOrderReq CreatePetHousePCOrderReq) error {
	order.ID = requirementOrder.ID
	order.StartedAt = requirementOrder.StartedAt
	order.FinishedAt = requirementOrder.FinishedAt
	order.CreateAt = requirementOrder.CreatedAt
	order.OrderType = requirementOrder.OrderType
	order.Status = requirementOrder.Status
	order.ServiceItems = pcOrderReq.ServiceItems
	payModeInt, err := ToPayMode(pcOrderReq.Basic, pcOrderReq.Commission)
	if err != nil {
		return err
	}
	detail := Detail{Basic: pcOrderReq.Basic, Commission: pcOrderReq.Commission}
	order.Payment = PaymentInfo{Mode: payModeInt, Detail: detail}
	order.Children = Children{MatchOrder: matchOrder, Groomer: groomer}
	order.UserID = requirementOrder.UserID
	return nil
}
