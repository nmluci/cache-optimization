package repository

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/go-redis/redis/v8"
	"github.com/nmluci/cache-optimization/internal/constants"
	"github.com/nmluci/cache-optimization/internal/model"
	"github.com/nmluci/cache-optimization/internal/util/sessionutil"
	"github.com/nmluci/cache-optimization/pkg/errs"
	"golang.org/x/net/context"
)

var (
	logTagAuthFindUserByID    = "[AuthFindUserByID]"
	logTagAuthFindUserByEmail = "[AuthFindUserByEmail]"
	logTagAuthInsertUser      = "[AuthInsertUser]"
	logTagAuthUpdateUser      = "[AuthUpdateUser]"
	logTagAuthDeleteUser      = "[AuthDeleteUser]"
	logTagNewSession          = "[AuthNewSession]"
	logTagFindBySessionKey    = "[AuthFindBySessionKey]"
	logTagInvalidateSession   = "[AuthInvalidateSession]"
)

var (
	sqlInsertUserCred = squirrel.Insert("userdata").Columns("email", "password", "fullname", "created_at", "updated_at", "access_level")
	sqlSelectUserCred = squirrel.Select("id", "email", "password", "fullname", "access_level").From("userdata")
	sqlUpdateUserCred = squirrel.Update("userdata")
)

func (repo *repository) FindUserByEmail(ctx context.Context, email string) (res *model.Users, err error) {
	stmt, args, err := sqlSelectUserCred.Where(squirrel.And{squirrel.Like{"email": fmt.Sprintf("%%%s%%", email)}, squirrel.Eq{"deleted_at": nil}}).ToSql()
	if err != nil {
		repo.logger.Errorf("%s failed to prepare SQL statement: %+v", logTagAuthFindUserByEmail, err)
		return
	}

	res = &model.Users{}
	err = repo.mariaDB.QueryRowContext(ctx, stmt, args...).Scan(&res.ID, &res.Email, &res.Password, &res.Fullname, &res.Priv)
	if err != nil && err != sql.ErrNoRows {
		repo.logger.Errorf("%s failed to parsed query result: %+v", logTagAuthFindUserByEmail, err)
		return
	} else if err == sql.ErrNoRows {
		return nil, nil
	}

	return
}

func (repo *repository) FindUserByID(ctx context.Context, id uint64) (res *model.Users, err error) {
	cacheKey, err := repo.redis.Get(ctx, fmt.Sprintf(constants.CacheUser, id)).Result()
	if err != nil && err != redis.Nil {
		repo.logger.Errorf("%s failed to fetch data from cache: %+v", logTagAuthFindUserByID, err)
		return
	}

	res = &model.Users{}
	if cacheKey != "" {
		if err = json.Unmarshal([]byte(cacheKey), res); err != nil {
			repo.logger.Errorf("%s failed to parse cache data: %+v", logTagAuthFindUserByID, err)
			repo.redis.Del(ctx, fmt.Sprintf(constants.CacheUser, id))
		} else {
			return
		}
	}

	stmt, args, err := sqlSelectUserCred.Where(squirrel.And{squirrel.Eq{"id": id}, squirrel.Eq{"deleted_at": nil}}).ToSql()
	if err != nil {
		repo.logger.Errorf("%s failed to prepare SQL statement: %+v", logTagAuthFindUserByID, err)
		return
	}

	err = repo.mariaDB.QueryRowContext(ctx, stmt, args...).Scan(&res.ID, &res.Email, &res.Password, &res.Fullname, &res.Priv)
	if err != nil && err != sql.ErrNoRows {
		repo.logger.Errorf("%s failed to parsed query result: %+v", logTagAuthFindUserByID, err)
		return
	} else if err == sql.ErrNoRows {
		return nil, nil
	}

	byteData, err := json.Marshal(res)
	if err != nil {
		repo.logger.Errorf("%s failed to encode query result: %+v", logTagAuthFindUserByID, err)
		return nil, err
	}

	err = repo.redis.Set(ctx, fmt.Sprintf(constants.CacheUser, id), string(byteData), constants.CacheSessionDuration).Err()
	if err != nil {
		repo.logger.Errorf("%s failed to insert query result into cache: %+v", logTagAuthFindUserByID, err)
		return res, nil
	}

	return
}

func (repo *repository) ForceFindUserByID(ctx context.Context, id uint64) (res *model.Users, err error) {
	stmt, args, err := sqlSelectUserCred.Where(squirrel.And{squirrel.Eq{"id": id}, squirrel.Eq{"deleted_at": nil}}).ToSql()
	if err != nil {
		repo.logger.Errorf("%s failed to prepare SQL statement: %+v", logTagAuthFindUserByID, err)
		return
	}

	res = &model.Users{}
	err = repo.mariaDB.QueryRowContext(ctx, stmt, args...).Scan(&res.ID, &res.Email, &res.Password, &res.Fullname, &res.Priv)
	if err != nil && err != sql.ErrNoRows {
		repo.logger.Errorf("%s failed to parsed query result: %+v", logTagAuthFindUserByID, err)
		return
	} else if err == sql.ErrNoRows {
		return nil, nil
	}

	return
}

func (repo *repository) InsertNewUser(ctx context.Context, data *model.Users) (err error) {
	tx, err := repo.mariaDB.BeginTx(ctx, &sql.TxOptions{
		Isolation: sql.LevelReadCommitted,
	})
	if err != nil {
		repo.logger.Errorf("%s failed to start new transaction: %+v", logTagAuthInsertUser, err)
		return
	}

	defer tx.Rollback()

	stmt, args, err := sqlInsertUserCred.Values(data.Email, data.Password, data.Fullname, time.Now(), time.Now(), data.Priv).ToSql()
	if err != nil {
		repo.logger.Errorf("%s failed to prepare SQL statement: %+v", logTagAuthInsertUser, err)
		return
	}

	_, err = tx.ExecContext(ctx, stmt, args...)
	if err != nil {
		repo.logger.Errorf("%s failed to parsed query result: %+v", logTagAuthInsertUser, err)
		return
	}

	if err = tx.Commit(); err != nil {
		repo.logger.Errorf("%s failed to committed transactions: %+v", logTagAuthInsertUser, err)
		return
	}

	return
}

func (repo *repository) UpdateUserByID(ctx context.Context, id uint64, data *model.Users) (err error) {
	tx, err := repo.mariaDB.BeginTx(ctx, &sql.TxOptions{
		Isolation: sql.LevelReadCommitted,
	})
	if err != nil {
		repo.logger.Errorf("%s failed to start new transaction: %+v", logTagAuthUpdateUser, err)
		return
	}
	defer tx.Rollback()

	stmt, args, err := sqlUpdateUserCred.SetMap(map[string]interface{}{
		"email":      data.Email,
		"password":   data.Password,
		"fullname":   data.Fullname,
		"updated_at": time.Now(),
	}).Where(squirrel.And{squirrel.Like{"id": id}, squirrel.Eq{"deleted_at": nil}}).ToSql()
	if err != nil {
		repo.logger.Errorf("%s failed to prepare SQL statement: %+v", logTagAuthUpdateUser, err)
		return
	}

	_, err = tx.ExecContext(ctx, stmt, args...)
	if err != nil {
		repo.logger.Errorf("%s failed to parsed query result: %+v", logTagAuthUpdateUser, err)
		return
	}

	if err = tx.Commit(); err != nil {
		repo.logger.Errorf("%s failed to committed transactions: %+v", logTagAuthUpdateUser, err)
		return
	}

	repo.redis.Del(ctx, fmt.Sprintf(constants.CacheUser, id))

	return
}

func (repo *repository) DeleteUserByID(ctx context.Context, id uint64) (err error) {
	tx, err := repo.mariaDB.BeginTx(ctx, &sql.TxOptions{
		Isolation: sql.LevelReadCommitted,
	})
	if err != nil {
		repo.logger.Errorf("%s failed to start new transaction: %+v", logTagAuthDeleteUser, err)
		return
	}
	defer tx.Rollback()

	stmt, args, err := sqlUpdateUserCred.Set("deleted_at", time.Now()).Where(squirrel.And{squirrel.Like{"id": id}, squirrel.Eq{"deleted_at": nil}}).ToSql()
	if err != nil {
		repo.logger.Errorf("%s failed to prepare SQL statement: %+v", logTagAuthDeleteUser, err)
		return
	}

	_, err = tx.ExecContext(ctx, stmt, args...)
	if err != nil {
		repo.logger.Errorf("%s failed to parsed query result: %+v", logTagAuthDeleteUser, err)
		return
	}

	if err = tx.Commit(); err != nil {
		repo.logger.Errorf("%s failed to committed transactions: %+v", logTagAuthDeleteUser, err)
		return
	}

	repo.redis.Del(ctx, fmt.Sprintf(constants.CacheUser, id))
	return
}

func (repo *repository) NewUserSession(ctx context.Context, data *model.Users) (sessionKey string, err error) {
	sKey := ""
	for {
		sKey = sessionutil.GenerateSessionKey()
		_, err = repo.redis.Get(ctx, fmt.Sprintf(constants.CacheSessionUser, sKey)).Result()
		if err != nil && err != redis.Nil {
			repo.logger.Errorf("%s failed to check session duplication: %+v", logTagNewSession, err)
			return "", err
		}

		if err == redis.Nil {
			break
		}
	}

	oldSKey, err := repo.redis.Get(ctx, fmt.Sprintf(constants.CacheSessionIdx, data.ID)).Result()
	if err != nil && err != redis.Nil {
		repo.logger.Errorf("%s failed to fetch old session key: %+v", logTagNewSession, err)
		return
	} else if err != redis.Nil {
		repo.redis.Del(ctx, fmt.Sprintf(constants.CacheSessionIdx, data.ID))
		repo.redis.Del(ctx, fmt.Sprintf(constants.CacheSessionUser, oldSKey))
	}

	encodeData, err := json.Marshal(data)
	if err != nil {
		repo.logger.Errorf("%s failed to encode profile data: %+v", logTagNewSession, err)
		return
	}

	if err = repo.redis.Set(ctx, fmt.Sprintf(constants.CacheSessionUser, sKey), string(encodeData), constants.CacheSessionDuration).Err(); err != nil {
		repo.logger.Errorf("%s failed to save session key: %+v", logTagNewSession, err)
		return
	}

	if err = repo.redis.Set(ctx, fmt.Sprintf(constants.CacheSessionIdx, data.ID), sKey, constants.CacheSessionDuration).Err(); err != nil {
		repo.logger.Errorf("%s failed to save session key: %+v", logTagNewSession, err)
		repo.redis.Del(ctx, fmt.Sprintf(constants.CacheSessionUser, sKey))
		return
	}

	return sKey, nil
}

func (repo *repository) FindUserBySession(ctx context.Context, sessionKey string) (res *model.Users, err error) {
	val, err := repo.redis.Get(ctx, fmt.Sprintf(constants.CacheSessionUser, sessionKey)).Result()
	if err != nil && err != redis.Nil {
		repo.logger.Errorf("%s failed to fetch session id: %+v", logTagFindBySessionKey, err)
		return
	}

	if val == "" || err == redis.Nil {
		repo.logger.Errorf("%s session key not found", logTagFindBySessionKey)
		err = errs.ErrBadRequest
		return
	}

	res = &model.Users{}
	err = json.Unmarshal([]byte(val), res)
	if err != nil {
		repo.logger.Errorf("%s failed to decoded userdata: %+v", logTagFindBySessionKey, err)
		return
	}

	return
}

func (repo *repository) InvalidateUserSession(ctx context.Context, sessionKey string) (err error) {
	val, err := repo.redis.Get(ctx, fmt.Sprintf(constants.CacheSessionUser, sessionKey)).Result()
	if err != nil && err != redis.Nil {
		repo.logger.Errorf("%s failed to fetch session id: %+v", logTagInvalidateSession, err)
		return
	}

	if val == "" || err == redis.Nil {
		return
	}

	err = repo.redis.Del(ctx, fmt.Sprintf(constants.CacheSessionUser, sessionKey)).Err()
	if err != nil {
		repo.logger.Errorf("%s failed to remove user session: %+v", logTagInvalidateSession, err)
		return
	}

	res := &model.Users{}
	err = json.Unmarshal([]byte(val), res)
	if err != nil {
		repo.logger.Errorf("%s failed to decoded userdata: %+v", logTagInvalidateSession, err)
		return
	}

	err = repo.redis.Del(ctx, fmt.Sprintf(constants.CacheSessionIdx, res.ID)).Err()
	if err != nil {
		repo.logger.Errorf("%s failed to remove user session: %+v", logTagInvalidateSession, err)
		return
	}

	return
}
