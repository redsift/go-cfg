package dcfg

import (
	"context"
)

// Stream represents
type Stream interface {
	Close() error
	Next(context.Context) bool
	Decode(any) (Meta, error)
}

func Subscribe[T any](ctx context.Context, b Backend, key Key, fn func(T, Meta, error) bool) error {
	stream, err := b.Subscribe(ctx, key)
	if err != nil {
		return err
	}

	go func() {
		defer stream.Close()

		for stream.Next(ctx) {
			var tmp T
			meta, err := stream.Decode(&tmp)
			if !fn(tmp, meta, err) {
				return
			}
		}
	}()

	return nil
}
