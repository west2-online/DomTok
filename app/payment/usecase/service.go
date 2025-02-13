package usecase

import (
	"context"
	"github.com/west2-online/DomTok/app/payment/domain/model"
)

func (uc *paymentUseCase) ProcessPayment(ctx context.Context, orderID int64) (*model.Payment, error) {
	return nil, nil
}
