package handler

import (
	"context"
	"fmt"

	codegenMiddleware "github.com/deepmap/oapi-codegen/pkg/middleware"
	"github.com/getkin/kin-openapi/openapi3filter"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/sawitpro/UserService/common"
	"github.com/sawitpro/UserService/generated"
	"github.com/sawitpro/UserService/helper"
)

var whitelistPaths = map[string]struct{}{
	"/auth/login":    {},
	"/auth/register": {},
}

func InitMiddleware() ([]echo.MiddlewareFunc, error) {
	spec, err := generated.GetSwagger()
	if err != nil {
		return nil, fmt.Errorf("failed to load OpenAPI spec: %w", err)
	}

	reqValidatorMiddleware := codegenMiddleware.OapiRequestValidatorWithOptions(spec, &codegenMiddleware.Options{
		Options: openapi3filter.Options{
			AuthenticationFunc: func(ctx context.Context, input *openapi3filter.AuthenticationInput) error {
				return nil
			},
		},
	})

	return []echo.MiddlewareFunc{middleware.Recover(), reqValidatorMiddleware, AuthMiddleware}, nil
}

func AuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		if _, ok := whitelistPaths[ctx.Request().URL.Path]; ok {
			return next(ctx)
		}

		token := getBearerToken(ctx)
		claimsToken, err := helper.ValidateToken(token)
		if err != nil {
			return echo.ErrForbidden
		}
		ctx.Set(common.USER_ID_CTX_KEY, claimsToken.UserID)

		return next(ctx)
	}
}

// func CreateAuthMiddlewareFunc() echo.MiddlewareFunc {
// 	return func(next echo.HandlerFunc) echo.HandlerFunc {
// 		return func(ctx echo.Context) error {
// 			if _, ok := whitelistPaths[ctx.Request().URL.Path]; ok {
// 				return next(ctx)
// 			}

// 			token := getBearerToken(ctx)
// 			claimsToken, err := helper.ValidateToken(token)
// 			if err != nil {
// 				return echo.ErrForbidden
// 			}
// 			ctx.Set(common.USER_ID_CTX_KEY, claimsToken.UserID)

// 			return next(ctx)
// 		}
// 	}
// }
