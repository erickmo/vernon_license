package querybus

import (
	"context"
	"fmt"
)

type Handler interface {
	Handle(ctx context.Context, q any) (any, error)
}

type QueryBus struct {
	handlers map[string]Handler
}

func New() *QueryBus {
	return &QueryBus{handlers: make(map[string]Handler)}
}

func (b *QueryBus) Register(qType string, h Handler) {
	b.handlers[qType] = h
}

func Dispatch[R any](ctx context.Context, bus *QueryBus, q any) (R, error) {
	key := fmt.Sprintf("%T", q)
	h, ok := bus.handlers[key]
	if !ok {
		var zero R
		return zero, fmt.Errorf("tidak ada handler untuk %s", key)
	}
	result, err := h.Handle(ctx, q)
	if err != nil {
		var zero R
		return zero, err
	}
	return result.(R), nil
}
