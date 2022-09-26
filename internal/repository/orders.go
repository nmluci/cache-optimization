package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/nmluci/cache-optimization/internal/constants"
	"github.com/nmluci/cache-optimization/internal/indto"
)

var (
	logTagCheckout      = "[OrderCheckout]"
	logTagCreateOrder   = "[Checkout-CreateOrder]"
	logTagCreateDetails = "[Checkout-CreateDetails]"
	logTagUpdateProduct = "[Checkout-UpdateProducts]"
)

var (
	sqlInsertOrder       = squirrel.Insert("orders").Columns("userdata_id", "order_date", "paid_at", "created_at", "updated_at")
	sqlInsertOrderDetail = squirrel.Insert("order_details").Columns("order_id", "product_id", "qty", "price", "created_at", "updated_at")
	// sqlUpdateProductQty  = squirrel.Update("products")
)

func (repo *repository) CheckoutOrders(ctx context.Context, payload *indto.OrderData) (err error) {
	tx, err := repo.mariaDB.BeginTx(ctx, &sql.TxOptions{
		Isolation: sql.LevelReadCommitted,
	})

	if err != nil {
		repo.logger.Errorf("%s failed to start new transaction: %+v", logTagCheckout, err)
		return
	}
	defer tx.Rollback()

	orderID, err := repo.insertOrder(ctx, tx, uint64(payload.UserdataID))
	if err != nil {
		repo.logger.Errorf("%s an error occured while storing new order metadata: %+v", logTagCheckout, err)
		return
	}

	for _, itm := range payload.Items {
		err = repo.insertDetail(ctx, tx, orderID, uint64(itm.ProductID), uint64(itm.Qty), uint64(itm.Price))
		if err != nil {
			repo.logger.Errorf("%s an error occured while storing new order details for productID: %d: %+v", logTagCheckout, itm.ProductID, err)
			return
		}

		err = repo.updateProductQty(ctx, tx, uint64(itm.ProductID), uint64(itm.Qty))
		if err != nil {
			repo.logger.Errorf("%s an error occured while updating product for productID: %d: %+v", logTagCheckout, itm.ProductID, err)
			return
		}
	}

	return
}

func (repo *repository) insertOrder(ctx context.Context, tx *sql.Tx, userID uint64) (orderID uint64, err error) {
	stmt, args, err := sqlInsertOrder.Values(userID, time.Now(), time.Now(), time.Now(), time.Now()).ToSql()
	if err != nil {
		repo.logger.Errorf("%s failed to prepare SQL statement: %+v", logTagCreateOrder, err)
		return
	}

	res, err := tx.ExecContext(ctx, stmt, args...)
	if err != nil {
		repo.logger.Errorf("%s failed to insert new order: %+v", logTagCreateOrder, err)
		return
	}

	lid, err := res.LastInsertId()
	if err != nil {
		repo.logger.Errorf("%s failed to fetch order ID: %+v", logTagCreateOrder, err)
		return
	}

	return uint64(lid), nil
}

func (repo *repository) insertDetail(ctx context.Context, tx *sql.Tx, orderID uint64, productID uint64, qty uint64, price uint64) (err error) {
	// "order_id", "product_id", "qty", "price", "created_at", "updated_at"
	stmt, args, err := sqlInsertOrderDetail.Values(orderID, productID, qty, price, time.Now(), time.Now()).ToSql()
	if err != nil {
		repo.logger.Errorf("%s failed to prepare SQL statement: %+v", logTagCreateDetails, err)
		return
	}

	_, err = tx.ExecContext(ctx, stmt, args...)
	if err != nil {
		repo.logger.Errorf("%s failed to insert order details: %+v", logTagCreateDetails, err)
		return
	}

	return nil
}

func (repo *repository) updateProductQty(ctx context.Context, tx *sql.Tx, productID uint64, qty uint64) (err error) {
	stmt := fmt.Sprintf("update products set qty = (case when qty - %d < 0 then 0 else qty - %d end) where id = %d", qty, qty, productID)

	_, err = tx.ExecContext(ctx, stmt, nil)
	if err != nil {
		repo.logger.Errorf("%s failed to update product listing: %+v", logTagUpdateProduct, err)
		return
	}

	repo.redis.Del(ctx, fmt.Sprintf(constants.CacheProducts, productID))
	return
}
