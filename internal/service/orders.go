package service

import (
	"context"

	"github.com/nmluci/cache-optimization/internal/indto"
	"github.com/nmluci/cache-optimization/pkg/dto"
	"github.com/nmluci/cache-optimization/pkg/errs"
)

var (
	logTagCheckout = "[OrderService-Checkout]"
)

func (s service) Checkout(ctx context.Context, payload *dto.PublicCheckout) (err error) {
	usr, err := s.repository.FindUserBySession(ctx, payload.SessionKey)
	if err != nil {
		s.logger.Errorf("%s failed to authenticate user's session: %+v", logTagCheckout, err)
		return errs.ErrUnauthorized
	}

	// validate items
	validItems := []*indto.ItemData{}
	for _, itm := range payload.Items {
		prd, err := s.repository.FindProductByID(ctx, uint64(itm.ProductID))
		if err != nil {
			s.logger.Errorf("%s failed to find product for ID: %d: %+v", logTagCheckout, itm.ProductID, err)
			return err
		}

		if prd == nil {
			s.logger.Errorf("%s product does not existed for ID: %d", logTagCheckout, itm.ProductID)
			return errs.ErrBadRequest
		}

		if prd.Qty < uint64(itm.Qty) {
			s.logger.Errorf("%s product exceeded qty in stocks (have: %d, want: %d)", logTagCheckout, prd.Qty, itm.Qty)
			return errs.ErrBadRequest
		}

		validItems = append(validItems, &indto.ItemData{
			Qty:       itm.Qty,
			ProductID: itm.ProductID,
			Price:     itm.Price,
		})
	}

	checkoutPayload := &indto.OrderData{
		Items:      validItems,
		UserdataID: int64(usr.ID),
	}

	if err = s.repository.CheckoutOrders(ctx, checkoutPayload); err != nil {
		s.logger.Errorf("%s an error occured while processing orders: %+v", logTagCheckout, err)
		return
	}

	return
}
