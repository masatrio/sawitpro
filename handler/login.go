package handler

import (
	"github.com/labstack/echo/v4"
	"github.com/sawitpro/UserService/generated"
	"github.com/sawitpro/UserService/service"
)

func (s *Server) Login(ctx echo.Context) error {
	request := new(generated.LoginRequest)
	if err := ctx.Bind(request); err != nil {
		return handleBadRequestJSON(ctx, err)
	}

	resp, errSvc := s.Service.Login(ctx.Request().Context(), service.LoginParam{
		PhoneNumber: request.PhoneNumber,
		Password:    request.Password,
	})

	if errSvc != nil {
		return handleServiceError(ctx, errSvc)
	}

	return handleSuccessJSON(ctx, &generated.LoginResponse{
		UserID: resp.UserID,
		Token:  resp.Token,
	})
}
