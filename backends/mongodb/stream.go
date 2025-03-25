package mongodb

import (
	"context"

	"github.com/redsift/go-cfg/dcfg"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

// Subscribe implements dcfg.Backend.
func (b *Backend) Subscribe(ctx context.Context, type_ dcfg.Type, key dcfg.Key) (dcfg.Stream, error) {
	stream := &Stream{ctx: ctx}

	if err := b.withColl(ctx, type_, key, func(coll *mongo.Collection) error {
		in, err := coll.Watch(ctx, mongo.Pipeline{
			{{Key: "$match", Value: b.filter(key)}},
		})
		if err != nil {
			return err
		}
		stream.in = in
		return nil
	}); err != nil {
		return nil, err
	}

	return stream, nil
}

type Stream struct {
	ctx context.Context
	in  *mongo.ChangeStream
}

// Close implements dcfg.Stream.
func (s *Stream) Close() error {
	return s.in.Close(s.ctx)
}

// Decode implements dcfg.Stream.
func (s *Stream) Decode(target any) error {
	return s.in.Decode(target)
}

// Next implements dcfg.Stream.
func (s *Stream) Next(ctx context.Context) bool {
	return s.in.Next(ctx)
}
