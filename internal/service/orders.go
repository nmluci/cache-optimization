package service

import (
	"context"

	"github.com/nmluci/cache-optimization/pkg/dto"
)

var (
	logTagCheckout = "[OrderService-Checkout]"
)

func (s service) Checkout(ctx context.Context, sessionKey string, payload *dto.PublicCheckout) (err error) {
	return
}
