package mongodb_test

import (
	"context"
	"strconv"
	"testing"
	"time"

	"github.com/redsift/go-cfg/dcfg"
	"github.com/redsift/go-siftjson"
	"github.com/stretchr/testify/require"
)

func TestMap(t *testing.T) {
	var (
		be      = connect(t)
		mKey    = dcfg.NewKey[map[siftjson.GUID]testdata](1, t.Name(), "m", strconv.FormatInt(time.Now().Unix(), 16))
		oneGUID = siftjson.GUID("sift-one.1")
		one     = testdata{
			Bool:   false,
			Int:    123,
			String: "one",
		}
		twoGUID = siftjson.GUID("sift-two.2")
		two     = testdata{
			Bool:   true,
			Int:    234,
			String: "two",
		}
	)

	m, err := dcfg.NewTypedMap[siftjson.GUID, testdata](be, mKey)
	require.NoError(t, err)

	// ensure value does not exist
	var tmp any
	_, err = be.Load(context.TODO(), mKey, &tmp)
	require.ErrorIs(t, err, dcfg.ErrNotFound)

	// set one item
	require.NoError(t, m.SetKey(context.TODO(), oneGUID, one))
	loaded, err := m.GetKey(context.TODO(), oneGUID)
	require.NoError(t, err)
	require.EqualValues(t, one, loaded)

	// ensure error when loading unknown key
	loaded, err = m.GetKey(context.TODO(), twoGUID)
	require.ErrorIs(t, err, dcfg.ErrNotFound)

	// load whole map
	data, err := m.Load(context.TODO())
	require.NoError(t, err)
	require.Len(t, data, 1)
	require.Contains(t, data, oneGUID)
	require.EqualValues(t, one, data[oneGUID])

	// add two and overwrite
	data[twoGUID] = two
	require.NoError(t, m.Update(context.TODO(), data))

	// ensure update was successful
	loaded, err = m.GetKey(context.TODO(), oneGUID)
	require.NoError(t, err)
	require.EqualValues(t, one, loaded)
	loaded, err = m.GetKey(context.TODO(), twoGUID)
	require.NoError(t, err)
	require.EqualValues(t, two, loaded)

	// load whole map
	data, err = m.Load(context.TODO())
	require.NoError(t, err)
	require.Len(t, data, 2)
	require.Contains(t, data, oneGUID)
	require.EqualValues(t, one, data[oneGUID])
	require.Contains(t, data, twoGUID)
	require.EqualValues(t, two, data[twoGUID])

	// clean up
	require.NoError(t, be.Delete(context.TODO(), mKey))

	_ = m
}
