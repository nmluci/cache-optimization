package service

import (
	"context"

	"github.com/nmluci/cache-optimization/internal/config"
	"github.com/nmluci/cache-optimization/internal/model"
	"github.com/nmluci/cache-optimization/pkg/dto"
	"github.com/nmluci/cache-optimization/pkg/errs"
	"golang.org/x/crypto/bcrypt"
)

var (
	logTagRegister   = "[AuthService-Register]"
	logTagLogin      = "[AuthService-Login]"
	logTagEditUser   = "[AuthService-EditUser"
	logTagLoginNC    = "[NC-AuthService-Login]"
	logTagEditUserNC = "[NC-AuthService-EditUser"
	logTagDeleteUser = "[AuthService-DeleteUser]"
)

func (s *service) Register(ctx context.Context, payload *dto.PublicUserPayload) (err error) {
	exists, err := s.repository.FindUserByEmail(ctx, payload.Email)
	if err != nil {
		s.logger.Errorf("%s failed to check email duplication: %+v", logTagRegister, err)
		return errs.ErrDuplicated
	}

	if exists != nil {
		s.logger.Errorf("%s user already exists", logTagRegister)
		return errs.ErrDuplicated
	}

	encPassword, err := bcrypt.GenerateFromPassword([]byte(payload.Password), bcrypt.DefaultCost)
	if err != nil {
		s.logger.Errorf("%s failed to hashed new password", logTagRegister)
		return
	}

	usr := &model.Users{
		Email:    payload.Email,
		Fullname: payload.Fullname,
		Password: string(encPassword),
	}

	conf := config.Get()
	if payload.MasterKey != nil && *payload.MasterKey == conf.MasterKey {
		usr.Priv = payload.Priv
	} else {
		usr.Priv = 1
	}

	err = s.repository.InsertNewUser(ctx, usr)
	if err != nil {
		s.logger.Errorf("%s failed to insert new user: %+v", logTagRegister, err)
		return
	}

	return
}

func (s *service) Login(ctx context.Context, payload *dto.PublicUserLoginPayload) (sessionKey string, usr *model.Users, err error) {
	exists, err := s.repository.FindUserByEmail(ctx, payload.Email)
	if err != nil {
		s.logger.Errorf("%s failed to check email duplication: %+v", logTagLogin, err)
		return
	}

	if exists == nil {
		s.logger.Errorf("%s user not registered: %+v", logTagLogin, err)
		err = errs.ErrBadRequest
		return
	}

	sessionKey, err = s.repository.NewUserSession(ctx, exists)
	if err != nil {
		s.logger.Errorf("%s failed generate new session: %+v", logTagLogin, err)
		return "", nil, err
	}

	exists.Password = ""
	return sessionKey, exists, nil
}

func (s *service) ForceLogin(ctx context.Context, payload *dto.PublicUserLoginPayload) (sessionKey string, usr *model.Users, err error) {
	exists, err := s.repository.FindUserByEmail(ctx, payload.Email)
	if err != nil {
		s.logger.Errorf("%s failed to check email duplication: %+v", logTagLoginNC, err)
		return
	}

	if exists == nil {
		s.logger.Errorf("%s user not registered: %+v", logTagLoginNC, err)
		err = errs.ErrBadRequest
		return
	}

	sessionKey, err = s.repository.NewSession(ctx, exists)
	if err != nil {
		s.logger.Errorf("%s failed generate new session: %+v", logTagLoginNC, err)
		return "", nil, err
	}
	exists.Password = ""
	return sessionKey, exists, nil
}

func (s *service) EditUser(ctx context.Context, id uint64, payload *dto.PublicUserPayload) (err error) {
	session, err := s.repository.FindUserBySession(ctx, payload.SessionKey)
	if err != nil {
		s.logger.Errorf("%s failed to fetch userdata from session key: %+v", logTagEditUser, err)
		return
	}

	if session.ID != id && session.Priv != 2 {
		s.logger.Errorf("%s cannot override higher user accounts", logTagEditUser)
		return errs.ErrBadRequest
	}

	user, err := s.repository.FindUserByID(ctx, id)
	if err != nil {
		s.logger.Errorf("%s failed to check email duplication: %+v", logTagEditUser, err)
		return err
	}

	if user == nil {
		s.logger.Errorf("%s user does not exists", logTagRegister)
		err = errs.ErrBadRequest
		return
	}

	encPassword, err := bcrypt.GenerateFromPassword([]byte(payload.Password), bcrypt.DefaultCost)
	if err != nil {
		s.logger.Errorf("%s failed to hashed new password", logTagEditUser)
		return
	}

	usr := &model.Users{
		Email:    payload.Email,
		Fullname: payload.Fullname,
		Password: string(encPassword),
	}

	err = s.repository.UpdateUserByID(ctx, id, usr)
	if err != nil {
		s.logger.Errorf("%s failed to insert new user: %+v", logTagEditUser, err)
		return
	}

	return
}

func (s *service) ForceEditUser(ctx context.Context, id uint64, payload *dto.PublicUserPayload) (err error) {
	session, err := s.repository.FindUserSessionByKey(ctx, payload.SessionKey)
	if err != nil {
		s.logger.Errorf("%s failed to fetch userdata from session key: %+v", logTagEditUserNC, err)
		return
	}

	if session.ID != id && session.Priv != 2 {
		s.logger.Errorf("%s cannot override higher user accounts", logTagEditUserNC)
		return errs.ErrBadRequest
	}

	user, err := s.repository.ForceFindUserByID(ctx, id)
	if err != nil {
		s.logger.Errorf("%s failed to check email duplication: %+v", logTagEditUserNC, err)
		return err
	}

	if user == nil {
		s.logger.Errorf("%s user does not exists", logTagEditUserNC)
		err = errs.ErrBadRequest
		return
	}

	encPassword, err := bcrypt.GenerateFromPassword([]byte(payload.Password), bcrypt.DefaultCost)
	if err != nil {
		s.logger.Errorf("%s failed to hashed new password", logTagEditUserNC)
		return
	}

	usr := &model.Users{
		Email:    payload.Email,
		Fullname: payload.Fullname,
		Password: string(encPassword),
	}

	err = s.repository.UpdateUserByID(ctx, id, usr)
	if err != nil {
		s.logger.Errorf("%s failed to insert new user: %+v", logTagEditUserNC, err)
		return
	}

	return
}

func (s *service) DeleteUser(ctx context.Context, id uint64, sessionKey string) (err error) {
	session, err := s.repository.FindUserBySession(ctx, sessionKey)
	if err != nil {
		s.logger.Errorf("%s failed to fetch userdata from session key: %+v", logTagDeleteUser, err)
		return
	}

	if session.ID != id && session.Priv != 2 {
		s.logger.Errorf("%s cannot override higher user accounts", logTagDeleteUser)
		return errs.ErrBadRequest
	}

	user, err := s.repository.FindUserByID(ctx, id)
	if err != nil {
		s.logger.Errorf("%s failed to check email duplication: %+v", logTagDeleteUser, err)
		return err
	}

	if user == nil {
		s.logger.Errorf("%s user does not exists", logTagDeleteUser)
		return errs.ErrBadRequest
	}

	err = s.repository.DeleteUserByID(ctx, id)
	if err != nil {
		s.logger.Errorf("%s failed to delete user: %+v", logTagDeleteUser, err)
		return
	}

	return
}

func (s *service) ForceDeleteUser(ctx context.Context, id uint64, sessionKey string) (err error) {
	session, err := s.repository.FindUserSessionByKey(ctx, sessionKey)
	if err != nil {
		s.logger.Errorf("%s failed to fetch userdata from session key: %+v", logTagDeleteUser, err)
		return
	}

	if session.ID != id && session.Priv != 2 {
		s.logger.Errorf("%s cannot override higher user accounts", logTagDeleteUser)
		return errs.ErrBadRequest
	}

	user, err := s.repository.FindUserByID(ctx, id)
	if err != nil {
		s.logger.Errorf("%s failed to check email duplication: %+v", logTagDeleteUser, err)
		return err
	}

	if user == nil {
		s.logger.Errorf("%s user does not exists", logTagDeleteUser)
		return errs.ErrBadRequest
	}

	err = s.repository.DeleteUserByID(ctx, id)
	if err != nil {
		s.logger.Errorf("%s failed to delete user: %+v", logTagDeleteUser, err)
		return
	}

	return
}
