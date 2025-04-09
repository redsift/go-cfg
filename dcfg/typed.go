package dcfg

import "context"

type TypedValue[T any] struct {
	backend Backend
	key     Key
}

func Typed[T any](b Backend, key Key) TypedValue[T] {
	return TypedValue[T]{
		backend: b,
		key:     key,
	}
}

func (tv *TypedValue[T]) Load(ctx context.Context) (meta Meta, out T, err error) {
	meta, err = tv.backend.Load(ctx, tv.key, &out)
	return
}

func (tv *TypedValue[T]) Store(ctx context.Context, meta *Meta, value T) error {
	return tv.backend.Store(ctx, tv.key, meta, value)
}
