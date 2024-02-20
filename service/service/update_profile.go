package service

import (
	"context"

	"github.com/sawitpro/UserService/common"
	"github.com/sawitpro/UserService/common/errors"
	"github.com/sawitpro/UserService/service"
)

func (s *Service) UpdateProfile(ctx context.Context, params service.UpdateProfileParam) (*service.UpdateProfileResponse, common.Error) {

	user, err := s.Repository.GetUserByID(ctx, params.UserID)

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

	if params.FullName == "" {
		params.FullName = user.FullName
	}

	if params.PhoneNumber == "" {
		params.PhoneNumber = user.Phone
	} else {
		existingUser, err := s.Repository.GetUserByPhone(ctx, params.PhoneNumber)
		if err != nil {
			return nil, errors.NewError(
				err.Error(),
				errors.SystemErrorType)
		}

		if existingUser != nil {
			return nil, errors.NewError(
				errors.NewPhoneAlreadyUsedErrorMessage(params.PhoneNumber),
				errors.ConflictErrorType)
		}
	}

	err = s.Repository.UpdateUser(ctx, params.UserID, params.FullName, params.PhoneNumber)
	if err != nil {
		return nil, errors.NewError(
			err.Error(),
			errors.SystemErrorType)
	}

	return &service.UpdateProfileResponse{
		FullName:    params.FullName,
		PhoneNumber: params.PhoneNumber,
	}, nil
}
