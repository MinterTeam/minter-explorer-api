package order

import "github.com/MinterTeam/minter-explorer-extender/v2/models"

var (
	statuses = map[models.OrderType]string{
		models.OrderTypeNew: "active",
		models.OrderTypePartiallyFilled: "partially_filled",
		models.OrderTypeFilled: "filled",
		models.OrderTypeUserCanceled: "canceled",
		models.OrderTypeExpired: "expired",
	}
)

func MakeOrderStatus(status models.OrderType) string {
	return statuses[status]
}
