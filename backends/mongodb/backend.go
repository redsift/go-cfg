package mongodb

import (
	"context"
	"strings"
	"sync"

	"github.com/redsift/go-cfg/dcfg"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func New(uri string) (dcfg.Backend, error) {
	client, err := mongo.Connect(options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}
	return NewFromClient(client), nil
}

func NewFromClient(client *mongo.Client) dcfg.Backend {
	return &Backend{
		client:      client,
		collections: map[string]*mongo.Collection{},
	}
}

type Backend struct {
	client      *mongo.Client
	lock        sync.Mutex
	collections map[string]*mongo.Collection
}

// Delete implements dcfg.Backend.
func (b *Backend) Delete(ctx context.Context, key dcfg.Key) error {
	return b.withColl(ctx, key, func(coll *mongo.Collection) error {
		_, err := coll.DeleteOne(ctx, b.filter(key))
		return err
	})
}

// Load implements dcfg.Backend.
func (b *Backend) Load(ctx context.Context, key dcfg.Key, target any) error {
	var tmp envelope[bson.RawValue]

	if err := b.load(ctx, key, &tmp); err != nil {
		return err
	}

	return mapError(tmp.Value.Unmarshal(target))
}

func (b *Backend) load(ctx context.Context, key dcfg.Key, target any) error {
	return b.withColl(ctx, key, func(coll *mongo.Collection) error {
		return coll.FindOne(ctx, b.filter(key)).Decode(target)
	})
}

// Store implements dcfg.Backend.
func (b *Backend) Store(ctx context.Context, key dcfg.Key, value any) error {
	tmp := envelope[any]{
		Key:     strings.Join(key.Elements, "/"),
		Version: uint(key.Version),
		Value:   value,
	}
	return b.withColl(ctx, key, func(coll *mongo.Collection) error {
		_, err := coll.ReplaceOne(
			ctx,
			b.filter(key),
			tmp,
			options.Replace().SetUpsert(true),
		)
		return err
	})
}
