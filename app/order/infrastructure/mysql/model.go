package mysql

import "github.com/west2-online/DomTok/pkg/constants"

type Order struct {
	Id                    int64
	Status                int8
	Uid                   int64
	TotalAmountOfGoods    float64
	TotalAmountOfFreight  float64
	TotalAmountOfDiscount float64
	PaymentAmount         float64
	PaymentStatus         int8
	PaymentAt             int64
	PaymentStyle          string
	OrderedAt             int64
	DeletedAt             int64
	DeliveryAt            int64
	AddressID             int64
	AddressInfo           string
}

func (Order) TableName() string {
	return constants.OrderTableName
}

type OrderGoods struct {
	OrderID          int64
	MerchantID       int64
	GoodsID          int64
	GoodsVersion     int
	GoodsName        string
	GoodsHeadDrawing string
	StyleID          int8
	StyleName        string
	StyleHeadDrawing string
	OriginCast       float64
	SaleCast         float64
	PurchaseQuantity int64
	PaymentAmount    float64
	FreightAmount    float64
	SettlementAmount float64
	DiscountAmount   float64
	SingleCast       float64
	CouponID         int64
}

func (OrderGoods) TableName() string {
	return constants.OrderGoodsTableName
}
