package dtos

import (
	"errors"
	"math"
)

type CreatePetHousePCOrderReq struct {
	RequestedAt  int64   `json:"requested_at"`
	StartedAt    int64   `json:"started_at"`
	FinishedAt   int64   `json:"finished_at"`
	ServiceItems []int   `json:"service_items"`
	Basic        float32 `json:"basic"`
	Commission   int     `json:"commission"`
}

func ToServiceBits(serviceItems []int) int64 {
	var bits int64
	for _, v := range serviceItems {
		bits |= 1 << v
	}
	return bits
}

func ToServiceItems(serviceBits int64) []int {
	var serviceItems []int
	for i := 0; i < 6; i++ {
		serviceBits &= int64(math.Pow(2, float64(i)))
		if serviceBits != int64(0) {
			serviceItems = append(serviceItems, i)
		}
	}
	return serviceItems
}

func ToServiceDesc(serviceItems []int) string {
	var serviceDesc string
	for _, v := range serviceItems {
		switch v {
		case 1:
			serviceDesc += "Shearing" + ","
		case 2:
			serviceDesc += "TeethClean" + ","
		case 3:
			serviceDesc += "SPA" + ","
		case 4:
			serviceDesc += "Cat" + ","
		case 5:
			serviceDesc += "BigSizeDog" + ","
		case 6:
			serviceDesc += "ExoticPet" + ","
		default:
		}
	}
	return serviceDesc
}

func ToPayMode(basic float32, commission int) (int, error) {
	if basic > float32(0.0) && commission > 0 {
		return 3, nil
	}

	if basic > float32(0.0) {
		return 1, nil
	}

	if commission > 0 {
		return 2, nil
	}
	return 0, errors.New("ORDER_PAYMENT_DATA_MISSION")
}

func ToPayModeDesc(basic float32, commission int) (string, error) {
	if basic > float32(0.0) && commission > 0 {
		return "MIXED", nil
	}

	if basic > float32(0.0) {
		return "BASIC", nil
	}

	if commission > 0 {
		return "Commission", nil
	}
	return "", errors.New("ORDER_PAYMENT_DATA_MISSION")
}
