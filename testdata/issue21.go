package testdata

import (
	"github.com/ghostiam/protogetter/testdata/proto"
)

func testIssue21(t *proto.Test) {
	// Invalid

	var _ bool = *t.OptBool                              // want `avoid direct access to proto field \*t\.OptBool, use t\.GetOptBool\(\) instead`
	var _ bool = *t.GetEmbedded().OptBool                // want `avoid direct access to proto field \*t\.GetEmbedded\(\)\.OptBool, use t\.GetEmbedded\(\).GetOptBool\(\) instead`
	var _, _ bool = *t.OptBool, *t.GetEmbedded().OptBool // want `avoid direct access to proto field \*t\.OptBool, use t\.GetOptBool\(\) instead` `avoid direct access to proto field \*t\.GetEmbedded\(\)\.OptBool, use t\.GetEmbedded\(\).GetOptBool\(\) instead`

	// Valid

	var _ bool = false
	var _ bool
	var _ *bool
	var _ *bool = t.OptBool
	var _ *bool = t.GetEmbedded().OptBool
	var _, _ *bool = t.OptBool, t.GetEmbedded().OptBool
}
