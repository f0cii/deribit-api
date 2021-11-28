package models

type SubaccountsDetails struct {
	OpenOrders []Order    `json:"open_orders"`
	Positions  []Position `json:"positions"`
	UID        int        `json:"uid"`
}
