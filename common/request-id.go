package common

import "context"

type requestIdKey struct{}

func ContextWithRequestId(
	ctx context.Context,
	requestId string,
) context.Context {
	return context.WithValue(ctx, requestIdKey{}, requestId)
}

func RequestIdFromContext(ctx context.Context) string {
	requestId, _ := ctx.Value(requestIdKey{}).(string)

	return requestId
}
