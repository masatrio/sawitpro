package service

import (
	"context"

	"github.com/sawitpro/UserService/common"
	"github.com/sawitpro/UserService/common/errors"
	"github.com/sawitpro/UserService/service"
)

func (s *Service) GetProfile(ctx context.Context, userID int64) (*service.UserInfoResponse, common.Error) {

	user, err := s.Repository.GetUserByID(ctx, userID)

	if err != nil {
		return nil, errors.NewError(
			err.Error(),
			errors.SystemErrorType)
	}

	if user == nil {
		return nil, errors.NewError(
			errors.UserDataNotFoundErrorMessage,
			errors.BadRequestErrorType)
	}

	return &service.UserInfoResponse{
		FullName:    user.FullName,
		PhoneNumber: user.Phone,
	}, nil
}
