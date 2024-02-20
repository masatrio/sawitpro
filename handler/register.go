package handler

import (
	"github.com/labstack/echo/v4"
	"github.com/sawitpro/UserService/generated"
	"github.com/sawitpro/UserService/service"
)

func (s *Server) Register(ctx echo.Context) error {
	request := new(generated.RegisterRequest)
	if err := ctx.Bind(request); err != nil {
		return handleBadRequestJSON(ctx, err)
	}

	resp, errSvc := s.Service.Register(ctx.Request().Context(), service.RegisterParam{
		FullName:    request.FullName,
		PhoneNumber: request.PhoneNumber,
		Password:    request.Password,
	})

	if errSvc != nil {
		return handleServiceError(ctx, errSvc)
	}

	return handleSuccessJSON(ctx, &generated.RegisterResponse{UserID: resp.UserID})
}
