package order

import "github.com/phrara/mallive/common/genproto/orderpb"

type Order struct {
	OrderID     string
	CustomerID  string
	Status      string
	PaymentLink string
	Items       []*orderpb.Item
}
