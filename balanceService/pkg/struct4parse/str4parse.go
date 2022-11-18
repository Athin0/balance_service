package struct4parse

import "time"

type Balance struct {
	UserId int64   `json:"user_id"`
	Value  float64 `json:"value"`
}

type BalanceWithDesc struct {
	UserId      int64   `json:"user_id"`
	Value       float64 `json:"value"`
	Time        time.Time
	Description string `json:"description,omitempty"`
}

type Transaction struct {
	Id          int64   `json:"id"`
	UserId      int64   `json:"user_id"`
	ServiceId   int64   `json:"service_id"`
	OrderId     int64   `json:"order_id"`
	Value       float64 `json:"value"`
	Time        time.Time
	Description string `json:"description,omitempty"`
	Replenish   bool   `json:"replenish"`
}

type Reserve struct {
	Id        int64   `json:"id"`
	UserId    int64   `json:"user_id"`
	ServiceId int64   `json:"service_id"`
	OrderId   int64   `json:"order_id"`
	Value     float64 `json:"value"`
}
