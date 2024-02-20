package handler

import (
	"errors"

	"github.com/labstack/echo/v4"
	"github.com/sawitpro/UserService/generated"
	"github.com/sawitpro/UserService/service"
)

func (s *Server) UpdateProfile(ctx echo.Context) error {
	userID, err := getContextUserID(ctx)
	if err != nil {
		return handleForbiddenAccessJSON(ctx, err)
	}

	request := &generated.UpdateProfileRequest{}
	if err := ctx.Bind(request); err != nil {
		return handleBadRequestJSON(ctx, err)
	}

	if request.FullName == nil && request.PhoneNumber == nil {
		return handleBadRequestJSON(ctx, errors.New("full name or phone number cannot be empty"))
	}

	fullName := ""
	phoneNumber := ""

	if request.FullName != nil {
		fullName = *request.FullName
	}

	if request.PhoneNumber != nil {
		phoneNumber = *request.PhoneNumber
	}

	resp, errSvc := s.Service.UpdateProfile(ctx.Request().Context(), service.UpdateProfileParam{
		FullName:    fullName,
		UserID:      userID,
		PhoneNumber: phoneNumber,
	})

	if errSvc != nil {
		return handleServiceError(ctx, errSvc)
	}

	return handleSuccessJSON(ctx, &generated.UpdateProfileResponse{
		FullName:    resp.FullName,
		PhoneNumber: resp.PhoneNumber,
	})
}
