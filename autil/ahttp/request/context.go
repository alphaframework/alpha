package request

import (
	"context"

	"github.com/gin-gonic/gin"
)

type key int

const (
	requestIdKey key = iota
)

func WithValue(parent context.Context, key interface{}, val interface{}) context.Context {
	return context.WithValue(parent, key, val)
}

func WithRequestId(parent context.Context, requestId string) context.Context {
	return WithValue(parent, requestIdKey, requestId)
}

func RequestIdFrom(ctx context.Context) (string, bool) {
	namespace, ok := ctx.Value(requestIdKey).(string)
	return namespace, ok
}

func RequestIdValue(ctx context.Context) string {
	requestId, _ := RequestIdFrom(ctx)
	return requestId
}

const stdContextKey = "__std_context"

func ContextToGin(c *gin.Context, ctx context.Context) {
	c.Set(stdContextKey, ctx)
}

func ContextFromGin(c *gin.Context) context.Context {
	if ctx, ok := c.Get(stdContextKey); ok {
		return ctx.(context.Context)
	}

	return context.TODO()
}

func RequestIdToGin(c *gin.Context, requestId string) {
	ctx := ContextFromGin(c)
	ctx = WithRequestId(ctx, requestId)
	ContextToGin(c, ctx)
}

func RequestIdFromGin(c *gin.Context) string {
	return RequestIdValue(ContextFromGin(c))
}
