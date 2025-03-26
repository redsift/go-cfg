package mongodb

import (
	"context"
	"errors"
	"fmt"
	"reflect"

	"github.com/redsift/go-cfg/dcfg"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

// Map implements dcfg.Backend.
func (b *Backend) Map(key dcfg.Key) dcfg.Map {
	return &Map{backend: b, key: key}
}

// Map implements dcfg.Map.
type Map struct {
	backend *Backend
	key     dcfg.Key
}

// Clear implements dcfg.Map.
func (m *Map) Clear(ctx context.Context) error {
	return m.backend.withColl(ctx, m.key, func(coll *mongo.Collection) error {
		_, err := coll.DeleteOne(ctx, m.backend.filter(m.key))
		return err
	})
}

// DelKey implements dcfg.Map.
func (m *Map) DelKey(ctx context.Context, key string) error {
	op := bson.D{{
		Key: "$unset",
		Value: bson.D{{
			Key: "value." + m.mangleKey(key),
		}},
	}}

	return m.backend.withColl(ctx, m.key, func(coll *mongo.Collection) error {
		_, err := coll.UpdateOne(ctx, m.backend.filter(m.key), op)
		return err
	})
}

// GetKey implements dcfg.Map.
func (m *Map) GetKey(ctx context.Context, key string, target any) error {
	var tmp envelope[bson.Raw]
	if err := m.backend.load(ctx, m.key, &tmp); err != nil {
		return err
	}

	return mapError(bson.Unmarshal(
		tmp.Value.Lookup(m.mangleKey(key)).Value,
		target,
	))
}

// Load implements dcfg.Map.
func (m *Map) Load(ctx context.Context, target any) error {
	tv := reflect.ValueOf(target)
	// fill potential nil pointer
	for tv.Kind() == reflect.Pointer {
		if tv.IsNil() {
			tv.Set(reflect.New(tv.Elem().Type()))
		}
		tv = tv.Elem()
	}

	var (
		errs       []error
		keyType    = tv.Type().Key()
		stringType = reflect.TypeOf("")
	)

	// we only support a map target
	if tv.Kind() != reflect.Map {
		return fmt.Errorf("cannot unmarshal into %T", target)
	}

	// the key must be ~string
	if !stringType.ConvertibleTo(keyType) {
		return fmt.Errorf("cannot unmarshal into map with key type %s", keyType)
	}

	// load data
	if err := m.backend.Load(ctx, m.key, target); err != nil {
		return err
	}

	// demangle keys
	for _, k := range tv.MapKeys() {
		var key = k.Convert(stringType).Interface().(string)

		demangled, changed, err := m.demangleKey(key)
		if err != nil {
			errs = append(errs, err)
			continue
		} else if !changed {
			continue
		}

		v := tv.MapIndex(k)
		tv.SetMapIndex(k, reflect.Value{})

		k := reflect.ValueOf(demangled).Convert(keyType)
		tv.SetMapIndex(k, v)
	}

	return errors.Join(errs...)
}

// SetKey implements dcfg.Map.
func (m *Map) SetKey(ctx context.Context, key string, value any) error {
	op := bson.D{{
		Key: "$set",
		Value: bson.D{{
			Key:   "value." + m.mangleKey(key),
			Value: value,
		}},
	}}

	return m.backend.withColl(ctx, m.key, func(coll *mongo.Collection) error {
		_, err := coll.UpdateOne(ctx, m.backend.filter(m.key), op, options.UpdateOne().SetUpsert(true))
		return err
	})
}

// Update implements dcfg.Map.
func (m *Map) Update(ctx context.Context, values map[string]any) error {
	tmp := make(map[string]any, len(values))
	for key, value := range values {
		mangled := m.mangleKey(key)
		tmp[mangled] = value
	}

	op := bson.D{{
		Key: "$set",
		Value: bson.D{{
			Key:   "value",
			Value: tmp,
		}},
	}}

	return m.backend.withColl(ctx, m.key, func(coll *mongo.Collection) error {
		_, err := coll.UpdateOne(ctx, m.backend.filter(m.key), op, options.UpdateOne().SetUpsert(true))
		return err
	})
}
