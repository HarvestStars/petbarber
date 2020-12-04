package dtos

type PageInfo struct {
	TotalItems int `json:"total_items"`
	TotalPages int `json:"total_pages"`
	PageSize   int `json:"page_size"`
	PageIndex  int `json:"page_index"`
}

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
	Basic        float32 `json:"basic"`
	Commission   int     `json:"commission"`
	TotalPayment float32 `json:"total_pay"`
}

type Children struct {
	MatchOrder ToMatch   `json:"match_order"`
	Groomer    TuGroomer `json:"groomer"`
}

type PCOrderListResp struct {
	List     []PCOrderResp `json:"lists"`
	PageInfo PageInfo      `json:"pagination"`
}

// 门店产生派单的响应
type PCOrderRespOnlyCreate struct {
	ID           uint        `json:"id"` // torequirement 订单号
	StartedAt    int64       `json:"started_at"`
	FinishedAt   int64       `json:"finished_at"`
	CreateAt     int64       `json:"created_at"`
	OrderType    int         `json:"order_type"`
	Status       int         `json:"status"`
	ServiceItems []int       `json:"service_items"`
	Payment      PaymentInfo `json:"payment"`
	UserID       uint        `json:"user_id"`
}

func (order *PCOrderRespOnlyCreate) RespCreateOrder(requirementOrder ToRequirement) error {
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
	detail := Detail{Basic: requirementOrder.Basic, Commission: requirementOrder.Commission, TotalPayment: requirementOrder.TotalPayment}
	order.Payment = PaymentInfo{Mode: payModeInt, Detail: detail}
	order.UserID = requirementOrder.UserID
	return nil
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
	detail := Detail{Basic: requirementOrder.Basic, Commission: requirementOrder.Commission, TotalPayment: requirementOrder.TotalPayment}
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
	RequirementOrder PCOrderRespOnlyCreate `json:"requirement_order"`
	PetHouse         TuPethouse            `json:"tu_pethouse"`
}

func (order *PCMatchResp) RespTransfer(matchOrder ToMatch, requirementOrder ToRequirement, petHouse TuPethouse) {
	order.ID = matchOrder.ID
	order.Status = matchOrder.Status
	order.CreateAt = matchOrder.CreatedAt
	order.UpdatedAt = matchOrder.UpdatedAt
	var requirementResp PCOrderRespOnlyCreate
	requirementResp.RespCreateOrder(requirementOrder)
	order.Parent = Parent{
		RequirementOrder: requirementResp,
		PetHouse:         petHouse,
	}
	order.PethouseOrderID = requirementOrder.ID
	order.UserID = matchOrder.UserID
}

type PCMatchListResp struct {
	List     []PCMatchResp `json:"lists"`
	PageInfo PageInfo      `json:"pagination"`
}

// 等待接单表响应
type PCActiveOrderResp struct {
	ID           uint        `json:"id"` // torequirement 订单号
	StartedAt    int64       `json:"started_at"`
	FinishedAt   int64       `json:"finished_at"`
	CreateAt     int64       `json:"created_at"`
	OrderType    int         `json:"order_type"`
	Status       int         `json:"status"`
	ServiceItems []int       `json:"service_items"`
	Payment      PaymentInfo `json:"payment"`
	PetHouse     TuPethouse  `json:"pethouse"`
	UserID       uint        `json:"user_id"`
}

func (order *PCActiveOrderResp) RespTransfer(requirementOrder ToRequirement, petHouse TuPethouse) error {
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
	detail := Detail{Basic: requirementOrder.Basic, Commission: requirementOrder.Commission, TotalPayment: requirementOrder.TotalPayment}
	order.Payment = PaymentInfo{Mode: payModeInt, Detail: detail}
	order.PetHouse = petHouse
	order.UserID = requirementOrder.UserID
	return nil
}

type PCActiveListResp struct {
	List     []PCActiveOrderResp `json:"lists"`
	PageInfo PageInfo            `json:"pagination"`
}
