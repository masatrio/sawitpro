package service

import (
	"context"
	"time"

	"github.com/sawitpro/UserService/common"
	"github.com/sawitpro/UserService/common/errors"
	"github.com/sawitpro/UserService/helper"
	"github.com/sawitpro/UserService/service"
	"golang.org/x/crypto/bcrypt"
)

func (s *Service) Login(ctx context.Context, params service.LoginParam) (*service.LoginResponse, common.Error) {

	user, err := s.Repository.GetUserByPhone(ctx, params.PhoneNumber)
	if err != nil {
		return nil, errors.NewError(
			err.Error(),
			errors.SystemErrorType)
	}

	if user == nil {
		return nil, errors.NewError(
			errors.WrongPhonePasswordErrorMessage,
			errors.BadRequestErrorType)
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.HashedPassword), []byte(params.Password))
	if err != nil {
		return nil, errors.NewError(
			errors.WrongPhonePasswordErrorMessage,
			errors.BadRequestErrorType)
	}

	err = s.Repository.ExecTransaction(ctx, func(ctx context.Context) error {
		err = s.Repository.InsertUserActivityLog(ctx, user.ID, common.LOGIN_ACTIVITY)
		if err != nil {
			return err
		}

		err = s.Repository.IncrementLoginCount(ctx, user.ID)
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, errors.NewError(
			err.Error(),
			errors.SystemErrorType)
	}

	token, err := helper.CreateToken(user.ID, time.Duration(6)*time.Hour)
	if err != nil {
		return nil, errors.NewError(
			err.Error(),
			errors.SystemErrorType)
	}

	return &service.LoginResponse{
		UserID: user.ID,
		Token:  token,
	}, nil
}
