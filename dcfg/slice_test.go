package dcfg

import (
	"context"
	"reflect"
	"strings"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestSliceSubscribeDiff(t *testing.T) {
	ctrl := gomock.NewController(t)
	be := NewMockBackend(ctrl)
	ms := NewMockSlice(ctrl)
	be.EXPECT().Slice(gomock.Any()).Return(ms)
	s, err := NewTypedSlice[string](be, NewKey[[]string](1, "test", t.Name()), strings.Compare)
	require.NoError(t, err)

	ms.EXPECT().
		Load(gomock.Any(), gomock.Any()).
		DoAndReturn(func(_ context.Context, target any) error {
			reflect.ValueOf(target).Elem().Set(
				reflect.ValueOf([]string{"r", "e", "d", "s", "i", "f", "t"}),
			)
			return nil
		})

	var wg sync.WaitGroup
	wg.Add(2)

	stream := NewMockStream(ctrl)
	stream.EXPECT().Next(gomock.Any()).Return(true)
	stream.EXPECT().Next(gomock.Any()).Return(false)
	stream.EXPECT().Decode(gomock.Any()).Do(func(target any) {
		defer wg.Done()
		reflect.ValueOf(target).Elem().Set(
			reflect.ValueOf([]string{"o", "n", "d", "m", "a", "r", "c"}),
		)
	})
	done := make(chan struct{})
	stream.EXPECT().Close().Do(func() { close(done) })

	be.EXPECT().
		Subscribe(gomock.Any(), gomock.Any()).
		DoAndReturn(func(_ context.Context, key Key) (Stream, error) {
			return stream, nil
		})

	var (
		a, r []string
		e    error
	)
	s.SubscribeDiff(context.TODO(), func(add []string, remove []string, err error) bool {
		defer wg.Done()
		t.Log("add", add)
		a = add
		t.Log("remove", remove)
		r = remove
		t.Log("err", err)
		e = err
		return true
	})
	wg.Wait()
	require.NoError(t, e)
	require.EqualValues(t, []string{"e", "f", "i", "s", "t"}, r, "mismatch in removed")
	require.EqualValues(t, []string{"a", "c", "m", "n", "o"}, a, "mismatch in added")
	<-done
}
