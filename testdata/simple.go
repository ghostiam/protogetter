package testdata

import (
	"github.com/ghostiam/protogolint/testdata/proto"
)

func simpleValid(t *proto.Test) {
	t.D = 1.0
	t.F = 1.0
	t.I32 = 1
	t.I64 = 1
	t.U32 = 1
	t.U64 = 1
	t.T = true
	t.B = []byte{1}
	t.S = "1"
	t.Embedded = &proto.Embedded{}
	t.GetEmbedded().S = "1"
	t.GetEmbedded().Embedded = &proto.Embedded{}
	t.GetEmbedded().GetEmbedded().S = "1"
	t.RepeatedEmbeddeds = []*proto.Embedded{{S: "1"}}

	_ = t.GetD()
	_ = t.GetF()
	_ = t.GetI32()
	_ = t.GetI64()
	_ = t.GetU32()
	_ = t.GetU64()
	_ = t.GetT()
	_ = t.GetB()
	_ = t.GetS()
	_ = t.GetEmbedded()
	_ = t.GetEmbedded().GetS()
	_ = t.GetEmbedded().GetEmbedded()
	_ = t.GetEmbedded().GetEmbedded().GetS()
	_ = t.GetRepeatedEmbeddeds()
	_ = t.GetRepeatedEmbeddeds()[0]
	_ = t.GetRepeatedEmbeddeds()[0].GetS()
	_ = t.GetRepeatedEmbeddeds()[0].GetEmbedded()
	_ = t.GetRepeatedEmbeddeds()[0].GetEmbedded().GetS()
}

func simpleInvalid(t *proto.Test) {
	_ = t.D
	_ = t.F
	_ = t.I32
	_ = t.I64
	_ = t.U32
	_ = t.U64
	_ = t.T
	_ = t.B
	_ = t.S
	_ = t.Embedded
	_ = t.Embedded.S
	_ = t.Embedded.Embedded
	_ = t.Embedded.Embedded.S
	_ = t.RepeatedEmbeddeds
	_ = t.RepeatedEmbeddeds[0]
	_ = t.RepeatedEmbeddeds[0].S
	_ = t.RepeatedEmbeddeds[0].Embedded
	_ = t.RepeatedEmbeddeds[0].Embedded.S
}
