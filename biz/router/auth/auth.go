package auth

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	"work4/biz/middleware"
	"work4/biz/pack"
	"work4/pkg/errmsg"
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
				pack.BuildFailResponse(c, errmsg.AuthError)
				c.Abort()
				return
			}
			middleware.GenerateAccessToken(c)
		}

		c.Next(ctx)
	}
}
