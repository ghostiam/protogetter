package testdata

import (
	"github.com/ghostiam/protogetter/testdata/proto"
)

func testInvalidEdition2023(t *proto.TestEdition2023) {
	_ = t.GetValue()
	t.SetValue(t.GetValue())
	t.SetValue(map[string]string{"test": "test"})
	// Issue #16
	t.SetValue(make(map[string]string, 5))
}
