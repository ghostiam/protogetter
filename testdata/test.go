package testdata

import (
	"fmt"

	"github.com/ghostiam/protogetter/testdata/proto"
)

type Test struct {
	Embedded *proto.Embedded
}

func testInvalid(t *proto.Test) {
	func(...interface{}) {}(t.B, t.D)      // want `avoid direct access to proto field "t\.B" use "t\.GetB\(\)"` `avoid direct access to proto field "t\.D" use "t\.GetD\(\)"`
	func(...interface{}) {}(t.GetB(), t.D) // want `avoid direct access to proto field "t\.D" use "t\.GetD\(\)"`
	func(...interface{}) {}(t.B, t.GetD()) // want `avoid direct access to proto field "t\.B" use "t\.GetB\(\)"`

	_ = t.D                                             // want `avoid direct access to proto field "t\.D" use "t\.GetD\(\)"`
	_ = t.F                                             // want `avoid direct access to proto field "t\.F" use "t\.GetF\(\)"`
	_ = t.I32                                           // want `avoid direct access to proto field "t\.I32" use "t\.GetI32\(\)"`
	_ = t.I64                                           // want `avoid direct access to proto field "t\.I64" use "t\.GetI64\(\)"`
	_ = t.U32                                           // want `avoid direct access to proto field "t\.U32" use "t\.GetU32\(\)"`
	_ = t.U64                                           // want `avoid direct access to proto field "t\.U64" use "t\.GetU64\(\)"`
	_ = t.T                                             // want `avoid direct access to proto field "t\.T" use "t\.GetT\(\)"`
	_ = t.B                                             // want `avoid direct access to proto field "t\.B" use "t\.GetB\(\)"`
	_ = t.S                                             // want `avoid direct access to proto field "t\.S" use "t\.GetS\(\)"`
	_ = t.Embedded                                      // want `avoid direct access to proto field "t\.Embedded" use "t\.GetEmbedded\(\)"`
	_ = t.Embedded.S                                    // want `avoid direct access to proto field "t\.Embedded\.S" use "t\.GetEmbedded\(\)\.GetS\(\)"`
	_ = t.GetEmbedded().S                               // want `avoid direct access to proto field "t\.GetEmbedded\(\)\.S" use "t\.GetEmbedded\(\)\.GetS\(\)"`
	_ = t.Embedded.Embedded                             // want `avoid direct access to proto field "t\.Embedded\.Embedded" use "t\.GetEmbedded\(\)\.GetEmbedded\(\)"`
	_ = t.GetEmbedded().Embedded                        // want `avoid direct access to proto field "t\.GetEmbedded\(\)\.Embedded" use "t\.GetEmbedded\(\)\.GetEmbedded\(\)"`
	_ = t.Embedded.Embedded.S                           // want `avoid direct access to proto field "t\.Embedded\.Embedded\.S" use "t\.GetEmbedded\(\)\.GetEmbedded\(\).GetS\(\)"`
	_ = t.GetEmbedded().GetEmbedded().S                 // want `avoid direct access to proto field "t\.GetEmbedded\(\)\.GetEmbedded\(\)\.S" use "t\.GetEmbedded\(\)\.GetEmbedded\(\)\.GetS\(\)"`
	_ = t.RepeatedEmbeddeds                             // want `avoid direct access to proto field "t\.RepeatedEmbeddeds" use "t\.GetRepeatedEmbeddeds\(\)"`
	_ = t.RepeatedEmbeddeds[0]                          // want `avoid direct access to proto field "t\.RepeatedEmbeddeds" use "t\.GetRepeatedEmbeddeds\(\)"`
	_ = t.RepeatedEmbeddeds[0].S                        // want `avoid direct access to proto field "t\.RepeatedEmbeddeds\[0\]\.S" use "t\.GetRepeatedEmbeddeds\(\)\[0\]\.GetS\(\)"`
	_ = t.GetRepeatedEmbeddeds()[0].S                   // want `avoid direct access to proto field "t\.GetRepeatedEmbeddeds\(\)\[0\]\.S" use "t\.GetRepeatedEmbeddeds\(\)\[0\]\.GetS\(\)"`
	_ = t.RepeatedEmbeddeds[0].Embedded                 // want `avoid direct access to proto field "t\.RepeatedEmbeddeds\[0\]\.Embedded" use "t\.GetRepeatedEmbeddeds\(\)\[0\]\.GetEmbedded\(\)"`
	_ = t.GetRepeatedEmbeddeds()[0].Embedded            // want `avoid direct access to proto field "t\.GetRepeatedEmbeddeds\(\)\[0\]\.Embedded" use "t\.GetRepeatedEmbeddeds\(\)\[0\]\.GetEmbedded\(\)"`
	_ = t.RepeatedEmbeddeds[0].Embedded.S               // want `avoid direct access to proto field "t\.RepeatedEmbeddeds\[0\]\.Embedded\.S" use "t\.GetRepeatedEmbeddeds\(\)\[0\].GetEmbedded\(\).GetS\(\)"`
	_ = t.GetRepeatedEmbeddeds()[0].GetEmbedded().S     // want `avoid direct access to proto field "t\.GetRepeatedEmbeddeds\(\)\[0\].GetEmbedded\(\).S" use "t\.GetRepeatedEmbeddeds\(\)\[0\].GetEmbedded\(\).GetS\(\)"`
	_ = t.RepeatedEmbeddeds[t.I64].Embedded.S           // want `avoid direct access to proto field "t\.RepeatedEmbeddeds\[t.I64\]\.Embedded\.S" use "t\.GetRepeatedEmbeddeds\(\)\[t\.GetI64\(\)\].GetEmbedded\(\).GetS\(\)"`
	_ = t.GetRepeatedEmbeddeds()[t.I64].GetEmbedded().S // want `avoid direct access to proto field "t\.GetRepeatedEmbeddeds\(\)\[t\.I64\]\.GetEmbedded\(\)\.S" use "t\.GetRepeatedEmbeddeds\(\)\[t\.GetI64\(\)\]\.GetEmbedded\(\).GetS\(\)"`

	var many []*proto.Test
	manyIndex := 42

	_ = many[0].T                   // want `avoid direct access to proto field "many\[0\]\.T" use "many\[0\]\.GetT\(\)"`
	_ = many[1].Embedded.S          // want `avoid direct access to proto field "many\[1\]\.Embedded\.S" use "many\[1\]\.GetEmbedded\(\)\.GetS\(\)"`
	_ = many[2].GetEmbedded().S     // want `avoid direct access to proto field "many\[2\]\.GetEmbedded\(\)\.S" use "many\[2\].GetEmbedded\(\)\.GetS\(\)"`
	_ = many[3].Embedded.Embedded.S // want `avoid direct access to proto field "many\[3\]\.Embedded\.Embedded\.S" use "many\[3\].GetEmbedded\(\)\.GetEmbedded\(\)\.GetS\(\)"`
	_ = many[manyIndex].S           // want `avoid direct access to proto field "many\[manyIndex\]\.S" use "many\[manyIndex\]\.GetS\(\)"`

	test := many[0].Embedded.S == "" || t.Embedded.CustomMethod() == nil || t.S == "" || t.Embedded == nil // want `avoid direct access to proto field "many\[0\]\.Embedded\.S" use "many\[0\]\.GetEmbedded\(\).GetS\(\)"` `avoid direct access to proto field "t\.Embedded\.CustomMethod\(\)" use "t\.GetEmbedded\(\)\.CustomMethod\(\)"` `avoid direct access to proto field "t\.S" use "t\.GetS\(\)"` `avoid direct access to proto field "t\.Embedded" use "t\.GetEmbedded\(\)"`
	_ = test

	other := proto.Other{}
	_ = other.MyMethod(nil).S // want `avoid direct access to proto field "other\.MyMethod\(nil\)\.S" use "other\.MyMethod\(nil\)\.GetS\(\)"`

	ems := t.RepeatedEmbeddeds // want `avoid direct access to proto field "t\.RepeatedEmbeddeds" use "t\.GetRepeatedEmbeddeds\(\)"`
	_ = ems[len(ems)-1].S      // want `avoid direct access to proto field "ems\[len\(ems\)-1\]\.S" use "ems\[len\(ems\)-1\]\.GetS\(\)"`

	ch := make(chan string)
	ch <- t.S // want `avoid direct access to proto field "t\.S" use "t\.GetS\(\)"`

	for _, v := range t.RepeatedEmbeddeds { // want `avoid direct access to proto field "t\.RepeatedEmbeddeds" use "t\.GetRepeatedEmbeddeds\(\)"`
		_ = v
	}

	fn := func(...interface{}) bool { return false }
	fn((*proto.Test)(nil).S) // want `avoid direct access to proto field "\(\*proto\.Test\)\(nil\)\.S" use "\(\*proto\.Test\)\(nil\)\.GetS\(\)"`

	var ptrs *[]proto.Test
	_ = (*ptrs)[42].RepeatedEmbeddeds    // want `avoid direct access to proto field "\(\*ptrs\)\[42\]\.RepeatedEmbeddeds" use "\(\*ptrs\)\[42\].GetRepeatedEmbeddeds\(\)"`
	_ = (*ptrs)[t.I64].RepeatedEmbeddeds // want `avoid direct access to proto field "\(\*ptrs\)\[t\.I64\]\.RepeatedEmbeddeds" use "\(\*ptrs\)\[t\.GetI64\(\)\].GetRepeatedEmbeddeds\(\)"`

	var anyType interface{}
	_ = anyType.(*proto.Test).S // want `avoid direct access to proto field "anyType\.\(\*proto\.Test\)\.S" use "anyType\.\(\*proto\.Test\)\.GetS\(\)"`

	t.Embedded.SetS("test")                              // want `avoid direct access to proto field "t\.Embedded\.SetS\(\\"test\\"\)" use "t\.GetEmbedded\(\)\.SetS\(\\"test\\"\)"`
	t.Embedded.SetMap(map[string]string{"test": "test"}) // want `avoid direct access to proto field "t\.Embedded\.SetMap\(map\[string\]string{\\"test\\": \\"test\\"}\)" use "t\.GetEmbedded\(\)\.SetMap\(map\[string\]string{\\"test\\": \\"test\\"}\)"`
}

func testValid(t *proto.Test) {
	func(...interface{}) {}(t.GetB(), t.GetD())
	func(...interface{}) {}(&t.B, &t.D)

	_, t.T = true, true
	_, t.T, _ = true, true, false
	_, _, t.T = true, true, false
	t.T, _ = true, true
	t.D = 2
	t.I32++
	t.I32 += 2

	fmt.Scanf("Test", &t.S, &t.B, &t.T)

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
	t.Embedded.S = "1"
	t.GetEmbedded().S = "1"
	t.Embedded.Embedded = &proto.Embedded{}
	t.GetEmbedded().Embedded = &proto.Embedded{}
	t.Embedded.Embedded.S = "1"
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

	other := proto.Other{}
	other.MyMethod(nil).CustomMethod()
	other.MyMethod(nil).GetS()

	var tt Test
	_ = tt.Embedded.GetS()
	_ = tt.Embedded.GetEmbedded().GetS()

	ems := t.GetRepeatedEmbeddeds()
	_ = ems[len(ems)-1].GetS()

	ch := make(chan string)
	ch <- t.GetS()

	t.Equal(&proto.Test{S: "test", I64: 42})
}
