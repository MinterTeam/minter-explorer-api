package order

import "github.com/MinterTeam/minter-explorer-extender/v2/models"

var (
	statuses = map[models.OrderType]Status{
		models.OrderTypeActive:          OrderStatusActive,
		models.OrderTypeNew:             OrderStatusActive,
		models.OrderTypePartiallyFilled: OrderStatusPartiallyFilled,
		models.OrderTypeFilled:          OrderStatusFilled,
		models.OrderTypeCanceled:        OrderStatusCanceled,
		models.OrderTypeExpired:         OrderStatusExpired,
	}
)

func MakeOrderStatus(status models.OrderType) Status {
	return statuses[status]
}
