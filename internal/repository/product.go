package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/go-redis/redis/v8"
	"github.com/nmluci/cache-optimization/internal/constants"
	"github.com/nmluci/cache-optimization/internal/model"
)

var (
	logTagProductFindByID = "[ProductFindByID]"
	logTagProductFindAll  = "[ProductFindAll]"
	logTagProductStore    = "[ProductStore]"
	logTagProductEdit     = "[ProductEdit]"
	logTagProductDelete   = "[ProductDelete]"
)

var (
	sqlInsertProduct = squirrel.Insert("products").Columns("name", "category", "description", "unit_price", "qty", "created_at", "updated_at")
	sqlSelectProduct = squirrel.Select("id", "name", "category", "description", "unit_price", "qty").From("products")
	sqlUpdateProduct = squirrel.Update("products")
)

func (repo *repository) FindProducts(ctx context.Context) (res []*model.Product, err error) {
	cacheKey, err := repo.redis.Get(ctx, fmt.Sprintf(constants.CacheProducts, "summary")).Result()
	if err != nil && err != redis.Nil {
		repo.logger.Errorf("%s failed to fetch data from cache: %+v", logTagProductFindAll, err)
		return
	}

	var stmt string
	var args []interface{}
	if cacheKey != "" {
		stmt, args, err = sqlSelectProduct.Where(squirrel.And{squirrel.Expr(fmt.Sprintf("id in %s", cacheKey))}).ToSql()
	} else {
		stmt, args, err = sqlSelectProduct.Where(squirrel.Eq{"deleted_at": nil}).ToSql()
	}

	if err != nil {
		repo.logger.Errorf("%s failed to prepare SQL statement: %+v", logTagProductFindAll, err)
		return
	}

	rows, err := repo.mariaDB.QueryContext(ctx, stmt, args...)
	if err != nil {
		repo.logger.Errorf("%s failed to parsed query results: %+v", logTagProductFindAll, err)
		return
	}

	res = []*model.Product{}
	productID := []string{}
	for rows.Next() {
		temp := &model.Product{}
		err = rows.Scan(&temp.ID, &temp.Name, &temp.Category, &temp.Description, &temp.UnitPrice, &temp.Qty)
		if err != nil {
			repo.logger.Errorf("%s failed to map query result: %+v", logTagProductFindAll, err)
		}

		productID = append(productID, strconv.FormatUint(temp.ID, 10))
		res = append(res, temp)
	}

	err = repo.redis.Set(ctx, fmt.Sprintf(constants.CacheProducts, "summary"), string(strings.Join(productID, ",")), constants.CacheDuration).Err()
	if err != nil {
		repo.logger.Errorf("%s failed to insert query result into cache: %+v", logTagProductFindAll, err)
		return res, nil
	}

	return
}

func (repo *repository) FindProductByID(ctx context.Context, id uint64) (res *model.Product, err error) {
	cacheKey, err := repo.redis.Get(ctx, fmt.Sprintf(constants.CacheProducts, id)).Result()
	if err != nil && err != redis.Nil {
		repo.logger.Errorf("%s failed to fetch data from cache: %+v", logTagProductFindByID, err)
		return
	}

	res = &model.Product{}
	if cacheKey != "" {
		if err = json.Unmarshal([]byte(cacheKey), res); err != nil {
			repo.logger.Errorf("%s failed to parse cache data: %+v", logTagProductFindByID, err)
			repo.redis.Del(ctx, fmt.Sprintf(constants.CacheProducts, id))
		} else {
			return
		}
	}

	stmt, args, err := sqlSelectProduct.Where(squirrel.And{squirrel.Eq{"id": id}, squirrel.Eq{"deleted_at": nil}}).ToSql()
	if err != nil {
		repo.logger.Errorf("%s failed to prepare SQL statement: %+v", logTagProductFindByID, err)
		return
	}

	err = repo.mariaDB.QueryRowContext(ctx, stmt, args...).Scan(&res.ID, &res.Name, &res.Category, &res.Description, &res.UnitPrice, &res.Qty)
	if err != nil {
		repo.logger.Errorf("%s failed to parsed query results: %+v", logTagProductFindByID, err)
		return
	}

	byteData, err := json.Marshal(res)
	if err != nil {
		repo.logger.Errorf("%s failed to encode query result: %+v", logTagProductFindByID, err)
		return nil, err
	}

	err = repo.redis.Set(ctx, fmt.Sprintf(constants.CacheProducts, id), string(byteData), constants.CacheDuration).Err()
	if err != nil {
		repo.logger.Errorf("%s failed to insert query result into cache: %+v", logTagProductFindByID, err)
		return res, nil
	}

	return
}

func (repo *repository) InsertNewProduct(ctx context.Context, data *model.Product) (err error) {
	tx, err := repo.mariaDB.BeginTx(ctx, &sql.TxOptions{
		Isolation: sql.LevelReadCommitted,
	})
	if err != nil {
		repo.logger.Errorf("%s failed to start new transaction: %+v", logTagProductStore, err)
		return
	}

	defer tx.Rollback()
	stmt, args, err := sqlInsertProduct.Values(data.Name, data.Category, data.UnitPrice, data.Qty, time.Now(), time.Now()).ToSql()
	if err != nil {
		repo.logger.Errorf("%s failed to prepare SQL statement: %+v", logTagProductStore, err)
		return
	}

	_, err = tx.ExecContext(ctx, stmt, args...)
	if err != nil {
		repo.logger.Errorf("%s failed to parsed query result: %+v", logTagProductStore, err)
		return
	}

	if err = tx.Commit(); err != nil {
		repo.logger.Errorf("%s failed to committed transactions: %+v", logTagProductStore, err)
		return
	}

	return
}

func (repo *repository) UpdateProduct(ctx context.Context, id uint64, data *model.Product) (err error) {
	tx, err := repo.mariaDB.BeginTx(ctx, &sql.TxOptions{
		Isolation: sql.LevelReadCommitted,
	})
	if err != nil {
		repo.logger.Errorf("%s failed to start new transaction: %+v", logTagProductEdit, err)
		return
	}
	defer tx.Rollback()

	stmt, args, err := sqlUpdateProduct.SetMap(map[string]interface{}{
		"name":        data.Name,
		"category":    data.Category,
		"description": data.Description,
		"unit_price":  data.UnitPrice,
		"qty":         data.Qty,
		"updated_at":  time.Now(),
	}).Where(squirrel.And{squirrel.Like{"id": id}, squirrel.Eq{"deleted_at": nil}}).ToSql()
	if err != nil {
		repo.logger.Errorf("%s failed to prepare SQL statement: %+v", logTagProductEdit, err)
		return
	}

	_, err = tx.ExecContext(ctx, stmt, args...)
	if err != nil {
		repo.logger.Errorf("%s failed to parsed query result: %+v", logTagProductEdit, err)
		return
	}

	if err = tx.Commit(); err != nil {
		repo.logger.Errorf("%s failed to committed transactions: %+v", logTagProductEdit, err)
		return
	}

	repo.redis.Del(ctx, fmt.Sprintf(constants.CacheUser, id))
	return
}

func (repo *repository) DeleteProductByID(ctx context.Context, id uint64) (err error) {
	tx, err := repo.mariaDB.BeginTx(ctx, &sql.TxOptions{
		Isolation: sql.LevelReadCommitted,
	})
	if err != nil {
		repo.logger.Errorf("%s failed to start new transaction: %+v", logTagProductDelete, err)
		return
	}
	defer tx.Rollback()

	stmt, args, err := sqlUpdateProduct.Set("deleted_at", time.Now()).Where(squirrel.And{squirrel.Like{"id": id}, squirrel.Eq{"deleted_at": nil}}).ToSql()
	if err != nil {
		repo.logger.Errorf("%s failed to prepare SQL statement: %+v", logTagProductDelete, err)
		return
	}

	_, err = tx.ExecContext(ctx, stmt, args...)
	if err != nil {
		repo.logger.Errorf("%s failed to parsed query result: %+v", logTagProductDelete, err)
		return
	}

	repo.redis.Del(ctx, fmt.Sprintf(constants.CacheProducts, id))
	if err = tx.Commit(); err != nil {
		repo.logger.Errorf("%s failed to committed transactions: %+v", logTagProductDelete, err)
		return
	}

	return
}
