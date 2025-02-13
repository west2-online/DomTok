package repository

import (
	"context"
	"github.com/west2-online/DomTok/app/payment/domain/model"
)

type PaymentDB interface {
	ProcessPayment(ctx context.Context, paymentID int64) (model.Payment, error)
}
