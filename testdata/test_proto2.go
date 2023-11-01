package testdata

import (
	"github.com/ghostiam/protogetter/testdata/proto"
)

func testInvalidProto2(t *proto.TestProto2) {
	_ = *t.D   // want `avoid direct access to proto field \*t\.D, use t\.GetD\(\) instead`
	_ = *t.F   // want `avoid direct access to proto field \*t\.F, use t\.GetF\(\) instead`
	_ = *t.I32 // want `avoid direct access to proto field \*t\.I32, use t\.GetI32\(\) instead`
	_ = *t.I64 // want `avoid direct access to proto field \*t\.I64, use t\.GetI64\(\) instead`
	_ = *t.U32 // want `avoid direct access to proto field \*t\.U32, use t\.GetU32\(\) instead`
	_ = *t.U64 // want `avoid direct access to proto field \*t\.U64, use t\.GetU64\(\) instead`
	_ = *t.T   // want `avoid direct access to proto field \*t\.T, use t\.GetT\(\) instead`
	_ = t.B    // want `avoid direct access to proto field t\.B, use t\.GetB\(\) instead`
	_ = *t.S   // want `avoid direct access to proto field \*t\.S, use t\.GetS\(\) instead`
}
