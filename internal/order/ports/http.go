package ports

import (
	"fmt"

	client "github.com/phrara/mallive/common/client/order"

	"github.com/gin-gonic/gin"
	"github.com/phrara/mallive/common/consts"
	"github.com/phrara/mallive/common/handler/errors"
	"github.com/phrara/mallive/common/server"
	"github.com/phrara/mallive/order/app"
	"github.com/phrara/mallive/order/app/command"
	"github.com/phrara/mallive/order/app/dto"
	"github.com/phrara/mallive/order/app/query"
	"github.com/phrara/mallive/order/convertor"
)

// 强制检查接口实现
var _ ServerInterface = (*OrderHTTPServer)(nil)

type OrderHTTPServer struct {
	app app.Application

	server.BaseResponse
	// 匿名字段注入，实现接口
	ServerInterface
}


func (o *OrderHTTPServer) PostCustomerCustomerIDOrder(c *gin.Context, customerID string) {
	var (
		req  client.CreateOrderRequest
		resp dto.CreateOrderResponse
		err  error
	)
	defer func() {
		o.Response(c, err, &resp)
	}()

	if err = c.ShouldBindJSON(&req); err != nil {
		err = errors.NewWithError(consts.ErrnoBindRequestError, err)
		return
	}
	if err = o.validate(req); err != nil {
		err = errors.NewWithError(consts.ErrnoRequestValidateError, err)
		return
	}
	r, err := o.app.Commands.CreateOrder.Handle(c.Request.Context(), command.CreateOrder{
		CustomerID: req.CustomerId,
		Items:      convertor.NewItemWithQuantityConvertor().ClientsToEntities(req.Items),
	})
	if err != nil {
		//err = errors.NewWithError()
		return
	}
	resp = dto.CreateOrderResponse{
		OrderID:     r.OrderID,
		CustomerID:  req.CustomerId,
		RedirectURL: fmt.Sprintf("http://localhost:8282/success?customerID=%s&orderID=%s", req.CustomerId, r.OrderID),
	}
}

func (o *OrderHTTPServer) GetCustomerCustomerIDOrderOrderID(c *gin.Context, customerID string, orderID string) {
	var (
		err  error
		resp interface{}
	)
	defer func() {
		o.Response(c, err, resp)
	}()

	order, err := o.app.Queries.GetCustomerOrder.Handle(c.Request.Context(), query.GetCustomerOrder{
		OrderID:    orderID,
		CustomerID: customerID,
	})
	if err != nil {
		return
	}

	resp = convertor.NewOrderConvertor().EntityToClient(order)
}


func (H OrderHTTPServer) validate(req client.CreateOrderRequest) error {
	for _, v := range req.Items {
		if v.Quantity <= 0 {
			return fmt.Errorf("quantity must be positive, got %d from %s", v.Quantity, v.Id)
		}
	}
	return nil
}

func NewOrderHTTPServer(app app.Application) *OrderHTTPServer {
	return &OrderHTTPServer{
		app: app,
	}
}