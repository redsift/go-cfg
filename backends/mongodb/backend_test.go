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

	be, err := mongodb.New(uri)
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
	err := be.Load(context.TODO(), aKey, &v)
	require.ErrorIs(t, err, dcfg.ErrNotFound)

	// write
	err = be.Store(context.TODO(), aKey, testdata{
		Bool:   true,
		Int:    123,
		String: t.Name(),
	})
	require.NoError(t, err)

	//read
	err = be.Load(context.TODO(), aKey, &v)
	require.NoError(t, err)
	require.Equal(t, true, v.Bool)
	require.Equal(t, 123, v.Int)
	require.Equal(t, t.Name(), v.String)

	// remove
	require.NoError(t, be.Delete(context.TODO(), aKey))

	// ensure not found error
	err = be.Load(context.TODO(), aKey, &v)
	require.ErrorIs(t, err, dcfg.ErrNotFound)
}
