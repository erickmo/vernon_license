package commandbus

import (
	"context"
	"fmt"
)

type Handler interface {
	Handle(ctx context.Context, cmd any) error
}

type CommandBus struct {
	handlers map[string]Handler
}

func New() *CommandBus {
	return &CommandBus{handlers: make(map[string]Handler)}
}

func (b *CommandBus) Register(cmdType string, h Handler) {
	b.handlers[cmdType] = h
}

func (b *CommandBus) Dispatch(ctx context.Context, cmd any) error {
	key := fmt.Sprintf("%T", cmd)
	h, ok := b.handlers[key]
	if !ok {
		return fmt.Errorf("tidak ada handler untuk %s", key)
	}
	return h.Handle(ctx, cmd)
}
