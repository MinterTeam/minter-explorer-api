package order

type Type string
type Status string

const (
	OrderTypeBuy  Type = "buy"
	OrderTypeSell Type = "sell"

	OrderStatusActive          Status = "active"
	OrderStatusPartiallyFilled Status = "partially_filled"
	OrderStatusFilled          Status = "filled"
	OrderStatusCanceled        Status = "canceled"
	OrderStatusExpired         Status = "expired"
)
