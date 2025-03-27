package mongodb_test

import (
	"context"
	"strconv"
	"testing"
	"time"

	"github.com/redsift/go-cfg/dcfg"
	"github.com/stretchr/testify/require"
)

func TestStream(t *testing.T) {
	var (
		be   = connect(t)
		vKey = dcfg.NewKey[testdata](1, t.Name(), "v", strconv.FormatInt(time.Now().Unix(), 16))
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
		update = make(chan testdata, 5)
	)

	require.NoError(t, dcfg.Subscribe(context.TODO(), be, vKey, func(d testdata, err error) bool {
		t.Log("update", d, err)

		if err != nil {
			panic(err)
		}

		update <- d
		return true
	}))

	requireUpdate := func(expected testdata) {
		select {
		case recv := <-update:
			require.EqualValues(t, expected, recv)
		case <-time.After(time.Second):
			t.Fatalf("missing update event (%v)", expected)
		}
	}

	require.NoError(t, be.Store(context.TODO(), vKey, nil, one))
	requireUpdate(one)

	require.NoError(t, be.Store(context.TODO(), vKey, nil, two))
	requireUpdate(two)

	require.NoError(t, be.Delete(context.TODO(), vKey))
}
