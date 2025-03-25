package dcfg

import (
	"context"
)

// Stream represents
type Stream interface {
	Close() error
	Next(context.Context) bool
	Decode(any) error
}

func Subscribe[T any](ctx context.Context, b Backend, type_ Type, key Key, fn func(T, error) bool) error {
	stream, err := b.Subscribe(ctx, type_, key)
	if err != nil {
		return err
	}

	go func() {
		defer stream.Close()

		for stream.Next(ctx) {
			var tmp T
			err := stream.Decode(&tmp)
			if !fn(tmp, err) {
				return
			}
		}
	}()

	return nil
}
