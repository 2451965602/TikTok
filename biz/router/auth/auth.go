package auth

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	"work4/biz/middleware/jwt"
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
		if !jwt.IsAccessTokenAvailable(ctx, c) {
			if !jwt.IsRefreshTokenAvailable(ctx, c) {
				pack.BuildFailResponse(c, errmsg.AuthError)
				c.Abort()
				return
			}
			jwt.GenerateAccessToken(c)
		}

		c.Next(ctx)
	}
}
