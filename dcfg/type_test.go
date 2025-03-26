package dcfg_test

import (
	"testing"

	"github.com/redsift/go-cfg/dcfg"
	"github.com/redsift/go-cfg/dcfg/testdata"
	"github.com/stretchr/testify/require"
)

func TestTypeOf(t *testing.T) {
	require.EqualValues(t, "go-cfg/dcfg/testdata/SomeStruct", dcfg.TypeOfValue(testdata.SomeStruct{}))
	require.EqualValues(t, "go-cfg/dcfg/testdata/SomeStruct", dcfg.TypeOfValue(&testdata.SomeStruct{}))

	var tmp testdata.SomeInterface = &testdata.SomeStruct{}
	require.EqualValues(t, "go-cfg/dcfg/testdata/SomeInterface", dcfg.TypeOfValue(tmp))
	require.EqualValues(t, "go-cfg/dcfg/testdata/SomeInterface", dcfg.TypeOfValue[testdata.SomeInterface](&testdata.SomeStruct{}))
	require.EqualValues(t, "go-cfg/dcfg/testdata/SomeInterface", dcfg.TypeOfValue[testdata.SomeInterface](nil))

	require.EqualValues(t, "string_slice", dcfg.TypeOfValue([]string{"test"}))
	require.EqualValues(t, "string_slice", dcfg.TypeOfValue([...]string{"test"}))
	require.EqualValues(t, "map[string]string", dcfg.TypeOf[map[string]string]())
}
