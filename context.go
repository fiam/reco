package reco

import (
	"context"
	"time"
)

type contextKey string

func (c contextKey) String() string {
	return "reco." + string(c)
}

var (
	debugContextKey   = contextKey("debug")
	startedContextKey = contextKey("started")
)

func SetDebug(ctx context.Context, debug bool) context.Context {
	return context.WithValue(ctx, debugContextKey, debug)
}

func Debug(ctx context.Context) *bool {
	if dbg, ok := ctx.Value(debugContextKey).(bool); ok {
		return &dbg
	}
	return nil
}

func SetStarted(ctx context.Context, started time.Time) context.Context {
	return context.WithValue(ctx, startedContextKey, started)
}

func Started(ctx context.Context) time.Time {
	if started, ok := ctx.Value(startedContextKey).(time.Time); ok {
		return started
	}
	return time.Time{}
}
