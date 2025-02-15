package constants

const (
	PaymentOrderNotExist               = -1
	PaymentStatusPending               = 0 // 待支付
	PaymentStatusProcessing            = 1 // 处理中
	PaymentStatusSuccess               = 2 // 成功支付
	PaymentStatusFailed                = 3 // 支付失败
	PaymentOrderNotExistToken          = ""
	PaymentOrderNotExistExpirationTime = 0
	UserNotExist                       = -1
	UserNotExistToken                  = ""
	UserNotExistExpirationTime         = 0
	HavePaidToken                      = ""
	HavePaidExpirationTime             = 0
	ErrorToken                         = ""
	ErrorExpirationTime                = 0
)
