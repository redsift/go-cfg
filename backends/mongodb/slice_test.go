package mongodb_test

import (
	"context"
	"strconv"
	"testing"
	"time"

	"github.com/redsift/go-cfg/dcfg"
	"github.com/stretchr/testify/require"
)

func TestSlice(t *testing.T) {
	var (
		be   = connect(t)
		sKey = dcfg.NewKey[[]testdata](1, t.Name(), "s", strconv.FormatInt(time.Now().Unix(), 16))
		one  = testdata{
			Bool:   false,
			Int:    123,
			String: "one",
		}
		two = testdata{
			Bool:   true,
			Int:    234,
			String: "two",
		}
	)

	slice, err := dcfg.NewTypedSlice[testdata](be, sKey)
	require.NoError(t, err)

	// ensure value does not exist
	var tmp any
	_, err = be.Load(context.TODO(), sKey, &tmp)
	require.ErrorIs(t, err, dcfg.ErrNotFound)

	// ensure ErrNotFound is mapped to empty value
	_, err = slice.Load(context.TODO())
	require.NoError(t, err)

	require.NoError(t, slice.Append(context.TODO(), one, two))

	// test expected data
	values, err := slice.Load(context.TODO())
	require.NoError(t, err)
	require.Len(t, values, 2)
	require.EqualValues(t, one, values[0])
	require.EqualValues(t, two, values[1])

	// remove one, ensure the other is left
	require.NoError(t, slice.RemoveItems(context.TODO(), one))

	values, err = slice.Load(context.TODO())
	require.NoError(t, err)
	require.Len(t, values, 1)
	require.EqualValues(t, two, values[0])

	// remove of non-existing item should be no-op
	require.NoError(t, slice.RemoveItems(context.TODO(), one))

	values, err = slice.Load(context.TODO())
	require.NoError(t, err)
	require.Len(t, values, 1)
	require.EqualValues(t, two, values[0])

	// re-add one
	require.NoError(t, slice.Append(context.TODO(), one))

	// test expected data
	values, err = slice.Load(context.TODO())
	require.NoError(t, err)
	require.Len(t, values, 2)
	require.EqualValues(t, one, values[1])
	require.EqualValues(t, two, values[0])

	// clean up
	require.NoError(t, be.Delete(context.TODO(), sKey))
}
