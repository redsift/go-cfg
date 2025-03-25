package dcfg_test

import (
	"testing"

	"github.com/redsift/go-cfg/dcfg"
	"github.com/redsift/go-cfg/dcfg/testdata"
	"github.com/stretchr/testify/require"
)

func TestTypeOf(t *testing.T) {
	require.Equal(t, "go_dcfg_dcfg_testdata_SomeStruct", dcfg.TypeOf(testdata.SomeStruct{}))
	require.Equal(t, "go_dcfg_dcfg_testdata_SomeStruct", dcfg.TypeOf(&testdata.SomeStruct{}))

	var tmp testdata.SomeInterface = &testdata.SomeStruct{}
	require.Equal(t, "go_dcfg_dcfg_testdata_SomeInterface", dcfg.TypeOf(tmp))
	require.Equal(t, "go_dcfg_dcfg_testdata_SomeInterface", dcfg.TypeOf[testdata.SomeInterface](&testdata.SomeStruct{}))
	require.Equal(t, "go_dcfg_dcfg_testdata_SomeInterface", dcfg.TypeOf[testdata.SomeInterface](nil))

	require.Equal(t, "string_slice", dcfg.TypeOf([]string{"test"}))
	require.Equal(t, "string_slice", dcfg.TypeOf([...]string{"test"}))
}
