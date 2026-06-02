package models

type OrderFilter struct {
	State *OrderState `json:"state"`
	Sort  string      `json:"sort"` // asc or desc
}
