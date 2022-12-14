package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/nmluci/cache-optimization/internal/constants"
	"github.com/nmluci/cache-optimization/internal/model"
	"github.com/nmluci/cache-optimization/internal/util/sessionutil"
)

var (
	logTagUserSessionNew        = "[SessionNew]"
	logTagUserSessionInvalidate = "[SessionInvalidate]"
	logTagUserSessionFindByKey  = "[SessionFindByKey]"
	// logTagUserSessionInvalidateByID = "[SessionNCInvalidateByID"
)

var (
	sqlInsertSession    = squirrel.Insert("user_sessions").Columns("session_id", "user_id", "expired_at")
	sqlSelectSession    = squirrel.Select("u.id", "u.email", "u.password", "u.fullname", "u.access_level").From("user_sessions s").LeftJoin("userdata u on s.user_id=u.id")
	sqlDeleteSession    = squirrel.Update("user_sessions")
	sqlSelectSessionKey = squirrel.Select("session_id").From("user_sessions")
)

func (repo *repository) FindUserSessionByKey(ctx context.Context, key string) (res *model.Users, err error) {
	expiredUnix := time.Now().Unix()

	stmt, args, err := sqlSelectSession.Where(squirrel.And{squirrel.Eq{"session_id": key}, squirrel.GtOrEq{"expired_at": expiredUnix}, squirrel.Eq{"s.deleted_at": nil}}).ToSql()
	if err != nil {
		repo.logger.Errorf("%s failed to prepare SQL statement: %+v", logTagUserSessionFindByKey, err)
		return
	}

	res = &model.Users{}
	err = repo.mariaDB.QueryRowContext(ctx, stmt, args...).Scan(&res.ID, &res.Email, &res.Password, &res.Fullname, &res.Priv)
	if err != nil && err != sql.ErrNoRows {
		repo.logger.Errorf("%s failed to parsed query [key: %s]: %+v", logTagUserSessionFindByKey, key, err)
		return
	} else if err == sql.ErrNoRows {
		return nil, nil
	}

	return
}

func (repo *repository) findUserSessionByID(ctx context.Context, id uint64) (sKey string, err error) {
	expiredUnix := time.Now().Unix()

	stmt, args, err := sqlSelectSessionKey.Where(squirrel.And{squirrel.Eq{"user_id": id}, squirrel.GtOrEq{"expired_at": expiredUnix}, squirrel.Eq{"deleted_at": nil}}).ToSql()
	if err != nil {
		repo.logger.Errorf("%s failed to prepare SQL statement: %+v", logTagUserSessionFindByKey, err)
		return
	}

	err = repo.mariaDB.QueryRowContext(ctx, stmt, args...).Scan(&sKey)
	if err != nil && err != sql.ErrNoRows {
		repo.logger.Errorf("%s failed to parsed query [id: %d]: %+v", logTagUserSessionFindByKey, id, err)
		return
	}

	return "", nil
}

func (repo *repository) NewSession(ctx context.Context, data *model.Users) (sessionKey string, err error) {
	sKey := ""
	for {
		sKey = sessionutil.GenerateSessionKey()
		val, err := repo.FindUserSessionByKey(ctx, sKey)
		if err != nil {
			repo.logger.Errorf("%s failed to check session duplication: %+v", logTagUserSessionNew, err)
			return "", err
		}

		if val == nil {
			break
		}
	}

	oldKey, err := repo.findUserSessionByID(ctx, data.ID)
	if err != nil {
		repo.logger.Errorf("%s failed to fetch old session key: %+v", logTagNewSession, err)
	} else if oldKey != "" {
		repo.InvalidateSessionKey(ctx, oldKey)
	}

	expiredUnix := time.Now().Add(constants.CacheSessionDuration).Unix()
	stmt, args, err := sqlInsertSession.Values(sKey, data.ID, expiredUnix).ToSql()
	if err != nil {
		repo.logger.Errorf("%s failed to prepare SQL statement: %+v", logTagUserSessionNew, err)
		return
	}

	_, err = repo.mariaDB.ExecContext(ctx, stmt, args...)
	if err != nil {
		repo.logger.Errorf("%s failed to insert new session key: %+v", logTagUserSessionNew, err)
		return
	}

	return sKey, nil
}

func (repo *repository) InvalidateSessionKey(ctx context.Context, sessionKey string) (err error) {
	val, err := repo.FindUserSessionByKey(ctx, sessionKey)
	if err != nil {
		repo.logger.Errorf("%s failed to check for session key: %+v", logTagUserSessionInvalidate, err)
		return
	}

	if val == nil {
		return
	}

	stmt, args, err := sqlDeleteSession.Set("deleted_at", time.Now()).Where(squirrel.Eq{"session_id": sessionKey}).ToSql()
	if err != nil {
		repo.logger.Errorf("%s failed to prepare SQL statement: %+v", logTagUserSessionNew, err)
		return
	}

	_, err = repo.mariaDB.ExecContext(ctx, stmt, args...)
	if err != nil {
		repo.logger.Errorf("%s failed to delete session key: %+v", logTagUserSessionNew, err)
		return
	}

	return
}

// func (repo *repository) invalidateSessionKeyByUserID(ctx context.Context, id uint64) (err error) {
// 	stmt, args, err := sqlDeleteSession.Set("deleted_at", time.Now()).Where(squirrel.Eq{"user_id": id}).ToSql()
// 	if err != nil {
// 		repo.logger.Errorf("%s failed to prepare SQL statement: %+v", logTagUserSessionInvalidateByID, err)
// 		return
// 	}

// 	_, err = repo.mariaDB.ExecContext(ctx, stmt, args...)
// 	if err != nil {
// 		repo.logger.Errorf("%s failed to delete session key: %+v", logTagUserSessionInvalidateByID, err)
// 		return
// 	}

// 	return
// }
