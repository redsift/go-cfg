package mongodb

import (
	"context"
	"strings"

	"github.com/redsift/go-cfg/dcfg"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

// Subscribe implements dcfg.Backend.
func (b *Backend) Subscribe(ctx context.Context, key dcfg.Key) (dcfg.Stream, error) {
	stream := &Stream{ctx: ctx}

	if err := b.withColl(ctx, key, func(coll *mongo.Collection) error {
		in, err := coll.Watch(ctx, mongo.Pipeline{
			{{
				Key: "$match",
				Value: bson.D{
					{Key: "fullDocument.key", Value: strings.Join(key.Elements, "/")},
					{Key: "fullDocument.version", Value: uint(key.Version)},
				},
			}},
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
	var tmp struct {
		OperationType string
		FullDocument  envelope[bson.RawValue]
	}

	// unwrap from event & envelope
	if err := s.in.Decode(&tmp); err != nil {
		return mapError(err)
	}

	return tmp.FullDocument.Value.Unmarshal(target)
}

// Next implements dcfg.Stream.
func (s *Stream) Next(ctx context.Context) bool {
	return s.in.Next(ctx)
}
