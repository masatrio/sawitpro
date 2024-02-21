package service

import (
	"context"

	"github.com/sawitpro/UserService/common"
	"github.com/sawitpro/UserService/common/errors"
	"github.com/sawitpro/UserService/repository"
	"github.com/sawitpro/UserService/service"
)

func (s *Service) Register(ctx context.Context, params service.RegisterParam) (*service.RegisterResponse, common.Error) {

	existingUser, err := s.Repository.GetUserByPhone(ctx, params.PhoneNumber)
	if err != nil {
		return nil, errors.NewError(
			err.Error(),
			errors.SystemErrorType)
	}

	if existingUser != nil {
		return nil, errors.NewError(
			errors.NewPhoneAlreadyUsedErrorMessage(params.PhoneNumber),
			errors.BadRequestErrorType)
	}

	hashedPassword, err := s.Hasher.HashPassword(params.Password)
	if err != nil {
		return nil, errors.NewError(
			err.Error(),
			errors.SystemErrorType)
	}

	user, err := s.Repository.InsertUser(ctx, &repository.User{
		FullName:       params.FullName,
		HashedPassword: hashedPassword,
		Phone:          params.PhoneNumber,
	})
	if err != nil {
		return nil, errors.NewError(
			err.Error(),
			errors.SystemErrorType)
	}

	return &service.RegisterResponse{
		UserID: user.ID,
	}, nil
}
