package handler

import (
	"github.com/labstack/echo/v4"
	"github.com/sawitpro/UserService/generated"
)

func (s *Server) GetProfile(ctx echo.Context) error {
	userID, err := getContextUserID(ctx)
	if err != nil {
		return handleForbiddenAccessJSON(ctx, err)
	}

	resp, errSvc := s.Service.GetProfile(ctx.Request().Context(), userID)
	if errSvc != nil {
		return handleServiceError(ctx, errSvc)
	}

	return handleSuccessJSON(ctx, &generated.GetProfileResponse{
		FullName:    resp.FullName,
		PhoneNumber: resp.PhoneNumber,
	})
}
