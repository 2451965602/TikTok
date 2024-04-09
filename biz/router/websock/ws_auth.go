package websock

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	"work4/biz/middleware"
)

func _wsAuth() []app.HandlerFunc {
	return append(make([]app.HandlerFunc, 0),
		tokenAuthFunc(),
	)
}

func tokenAuthFunc() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		if !middleware.IsAccessTokenAvailable(ctx, c) {
			c.Abort()
			return
		}
		c.Next(ctx)
	}
}
