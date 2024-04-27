package cahce

import "rcsp/internal/model"

type Cache interface {
	AddOrder(order model.Order) error
	GetOrder(key string) (model.Order, error)
	DeleteOrder(key string) error
}
