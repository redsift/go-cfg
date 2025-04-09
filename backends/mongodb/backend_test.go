package mongodb_test

import (
	"context"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/redsift/go-cfg/backends/mongodb"
	"github.com/redsift/go-cfg/dcfg"
	"github.com/stretchr/testify/require"
)

type testdata struct {
	String string
	Int    int
	Bool   bool
}

func connect(t *testing.T) dcfg.Backend {
	t.Helper()
	uri, ok := os.LookupEnv("MONGODB_URI")
	if !ok {
		uri = "mongodb://localhost:27017/?replicaSet=rs0"
	}

	be, err := mongodb.New(uri, "dcfg")
	require.NoError(t, err)
	return be
}

func TestBackend(t *testing.T) {
	var (
		be   = connect(t)
		aKey = dcfg.NewKey[testdata](1, t.Name(), "a", strconv.FormatInt(time.Now().Unix(), 16))
		v    testdata
	)

	// should not exist
	_, err := be.Load(context.TODO(), aKey, &v)
	require.ErrorIs(t, err, dcfg.ErrNotFound)

	// write
	err = be.Store(context.TODO(), aKey, nil, testdata{
		Bool:   true,
		Int:    123,
		String: t.Name(),
	})
	require.NoError(t, err)

	//read
	meta, err := be.Load(context.TODO(), aKey, &v)
	require.NoError(t, err)
	require.Equal(t, true, v.Bool)
	require.Equal(t, 123, v.Int)
	require.Equal(t, t.Name(), v.String)
	require.EqualValues(t, 1, meta.Generation)

	// successful overwrite
	err = be.Store(context.TODO(), aKey, &meta, testdata{
		Bool:   false,
		Int:    234,
		String: "over",
	})
	require.NoError(t, err)
	meta2, err := be.Load(context.TODO(), aKey, &v)
	require.Equal(t, false, v.Bool)
	require.Equal(t, 234, v.Int)
	require.Equal(t, "over", v.String)
	require.EqualValues(t, 2, meta2.Generation)

	// unsuccessful overwrite
	err = be.Store(context.TODO(), aKey, &meta, testdata{
		Bool:   false,
		Int:    345,
		String: "not",
	})
	require.NoError(t, err)
	meta3, err := be.Load(context.TODO(), aKey, &v)
	require.Equal(t, false, v.Bool)
	require.Equal(t, 234, v.Int)
	require.Equal(t, "over", v.String)
	require.EqualValues(t, 2, meta3.Generation)

	// remove
	require.NoError(t, be.Delete(context.TODO(), aKey))

	// ensure not found error
	_, err = be.Load(context.TODO(), aKey, &v)
	require.ErrorIs(t, err, dcfg.ErrNotFound)
}
