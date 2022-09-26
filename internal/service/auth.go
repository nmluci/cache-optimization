package service

import (
	"context"

	"github.com/nmluci/cache-optimization/internal/model"
	"github.com/nmluci/cache-optimization/pkg/dto"
	"golang.org/x/crypto/bcrypt"
)

var (
	logTagRegister   = "[AuthService-Register]"
	logTagLogin      = "[AuthService-Login]"
	logTagEditUser   = "[AuthService-EditUser"
	logTagDeleteUser = "[AuthService-DeleteUser]"
)

func (s *service) Register(ctx context.Context, payload *dto.PublicUserPayload) (err error) {
	exists, err := s.repository.FindUserByEmail(ctx, payload.Email)
	if err != nil {
		s.logger.Errorf("%s failed to check email duplication: %+v", logTagRegister, err)
		return err
	}

	if exists != nil {
		s.logger.Errorf("%s user already exists", logTagRegister)
		return
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

	err = s.repository.InsertNewUser(ctx, usr)
	if err != nil {
		s.logger.Errorf("%s failed to insert new user: %+v", logTagRegister, err)
		return
	}

	return
}

func (s *service) Login(ctx context.Context, payload *dto.PublicUserLoginPayload) (sessionKey string, err error) {
	exists, err := s.repository.FindUserByEmail(ctx, payload.Email)
	if err != nil {
		s.logger.Errorf("%s failed to check email duplication: %+v", logTagLogin, err)
		return
	}

	if exists == nil {
		s.logger.Errorf("%s user not registered: %+v", logTagLogin, err)
		return
	}

	sessionKey, err = s.repository.NewUserSession(ctx, exists)
	if err != nil {
		s.logger.Errorf("%s failed generate new session: %+v", logTagLogin, err)
		return "", err
	}

	return
}

func (s *service) EditUser(ctx context.Context, payload *dto.PublicUserPayload) (err error) {
	exists, err := s.repository.FindUserByEmail(ctx, payload.Email)
	if err != nil {
		s.logger.Errorf("%s failed to check email duplication: %+v", logTagEditUser, err)
		return err
	}

	if exists != nil {
		s.logger.Errorf("%s user already exists", logTagRegister)
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

	err = s.repository.InsertNewUser(ctx, usr)
	if err != nil {
		s.logger.Errorf("%s failed to insert new user: %+v", logTagEditUser, err)
		return
	}

	return
}

func (s *service) DeleteUser(ctx context.Context, id uint64) (err error) {
	exists, err := s.repository.FindUserByID(ctx, id)
	if err != nil {
		s.logger.Errorf("%s failed to check email duplication: %+v", logTagDeleteUser, err)
		return err
	}

	if exists != nil {
		s.logger.Errorf("%s user already exists", logTagDeleteUser)
		return
	}

	err = s.repository.DeleteUserByID(ctx, id)
	if err != nil {
		s.logger.Errorf("%s failed to insert new user: %+v", logTagDeleteUser, err)
		return
	}

	return
}
