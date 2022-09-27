package repository

import (
	"context"
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
)

var (
	sqlInsertSession = squirrel.Insert("user_sessions").Columns("session_id", "user_id", "expired_at")
	sqlSelectSession = squirrel.Select("u.id", "u.email", "u.password", "u.fullname", "u.access_level").From("user_sessions s").LeftJoin("userdata u on s.user_id=u.id")
	sqlDeleteSession = squirrel.Delete("user_sessions")
)

func (repo *repository) FindUserSessionByKey(ctx context.Context, key string) (res *model.Users, err error) {
	expiredUnix := time.Now().Add(constants.CacheSessionDuration).Unix()

	stmt, args, err := sqlSelectSession.Where(squirrel.And{squirrel.Eq{"session_id": key}, squirrel.Lt{"expired_at": expiredUnix}}).ToSql()
	if err != nil {
		repo.logger.Errorf("%s failed to prepare SQL statement: %+v", logTagUserSessionFindByKey, err)
		return
	}

	res = &model.Users{}
	err = repo.mariaDB.QueryRowContext(ctx, stmt, args...).Scan(&res.ID, &res.Email, &res.Password, &res.Fullname, &res.Priv)
	if err != nil {
		repo.logger.Errorf("%s failed to parsed query results: %+v", logTagUserSessionFindByKey, err)
		return
	}

	return
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

	stmt, args, err := sqlDeleteSession.Where(squirrel.Eq{"session_id": sessionKey}).ToSql()
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
