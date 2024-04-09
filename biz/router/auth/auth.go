package auth

import (
	"context"
	"errors"
	"github.com/cloudwego/hertz/pkg/app"
	"work4/biz/middleware"
	"work4/biz/pack"
)

func Auth() []app.HandlerFunc {
	return append(make([]app.HandlerFunc, 0),
		DoubleTokenAuthFunc(),
	)
}

func DoubleTokenAuthFunc() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		if !middleware.IsAccessTokenAvailable(ctx, c) {
			if !middleware.IsRefreshTokenAvailable(ctx, c) {
				pack.SendFailResponse(c, errors.New("TokenIsInavailableError"))
				c.Abort()
				return
			}
			middleware.GenerateAccessToken(ctx, c)
		}

		c.Next(ctx)
	}
}
