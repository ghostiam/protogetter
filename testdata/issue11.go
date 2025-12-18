package testdata

import (
	"github.com/ghostiam/protogetter/testdata/proto"
)

func testIssue11(t *proto.Test) {
	var optBoolVar *bool

	type structWithPtrField struct {
		OptBool *bool
	}

	// Invalid

	optBoolVar = t.Embedded.OptBool // want `avoid direct access to proto field t\.Embedded\.OptBool, use t\.GetEmbedded\(\)\.OptBool instead`
	_ = optBoolVar

	_ = structWithPtrField{
		OptBool: t.Embedded.OptBool, // want `avoid direct access to proto field t\.Embedded\.OptBool, use t\.GetEmbedded\(\)\.OptBool instead`
	}

	optionalArgsFunc(t.Embedded.OptBool)             // want `avoid direct access to proto field t\.Embedded\.OptBool, use t\.GetEmbedded\(\)\.OptBool instead`
	optionalArgs2Func(t.OptBool, t.Embedded.OptBool) // want `avoid direct access to proto field t\.Embedded\.OptBool, use t\.GetEmbedded\(\)\.OptBool instead`
	messageArgsFunc(t.Embedded)                      // want `avoid direct access to proto field t\.Embedded, use t\.GetEmbedded\(\) instead`
	nonOptionalArgsFunc(t.T)                         // want `avoid direct access to proto field t\.T, use t\.GetT\(\) instead`
	optionalVariadicArgsFunc(
		t.OptBool,
		t.Embedded.OptBool, // want `avoid direct access to proto field t\.Embedded\.OptBool, use t\.GetEmbedded\(\)\.OptBool instead`
	)
	nonOptionalVariadicArgsFunc(
		t.T,                 // want `avoid direct access to proto field t\.T, use t\.GetT\(\) instead`
		*t.Embedded.OptBool, // want `avoid direct access to proto field \*t\.Embedded\.OptBool, use t\.GetEmbedded\(\)\.GetOptBool\(\) instead`
		false,
		true,
	)
	variadicArgsAnyFunc(
		t.OptBool,
		t.Embedded.OptBool, // want `avoid direct access to proto field t\.Embedded\.OptBool, use t\.GetEmbedded\(\)\.OptBool instead`
		t.T,                // want `avoid direct access to proto field t\.T, use t\.GetT\(\) instead`
		t.Embedded,         // want `avoid direct access to proto field t\.Embedded, use t\.GetEmbedded\(\) instead`
	)
	variadicArgsInterfaceFunc(
		t.OptBool,
		t.Embedded.OptBool, // want `avoid direct access to proto field t\.Embedded\.OptBool, use t\.GetEmbedded\(\)\.OptBool instead`
		t.T,                // want `avoid direct access to proto field t\.T, use t\.GetT\(\) instead`
		t.Embedded,         // want `avoid direct access to proto field t\.Embedded, use t\.GetEmbedded\(\) instead`
	)
	variadicArgsExtInterfaceFunc(
		t.OptBool,
		t.Embedded.OptBool, // want `avoid direct access to proto field t\.Embedded\.OptBool, use t\.GetEmbedded\(\)\.OptBool instead`
		t.T,                // want `avoid direct access to proto field t\.T, use t\.GetT\(\) instead`
		t.Embedded,         // want `avoid direct access to proto field t\.Embedded, use t\.GetEmbedded\(\) instead`
	)

	// Valid

	optBoolVar = t.OptBool
	optBoolVar = t.GetEmbedded().OptBool
	_ = optBoolVar

	_ = structWithPtrField{
		OptBool: t.OptBool,
	}
	_ = structWithPtrField{
		OptBool: t.GetEmbedded().OptBool,
	}

	optionalArgsFunc(t.GetEmbedded().OptBool)
	optionalArgs2Func(t.OptBool, t.GetEmbedded().OptBool)
	messageArgsFunc(t.GetEmbedded())
	nonOptionalArgsFunc(t.GetT())
	optionalVariadicArgsFunc(t.OptBool, t.GetEmbedded().OptBool)
	nonOptionalVariadicArgsFunc(t.GetT(), t.GetEmbedded().GetOptBool(), false, true)
	variadicArgsAnyFunc(t.OptBool, t.GetEmbedded().OptBool, t.GetT(), t.GetEmbedded())
	variadicArgsInterfaceFunc(t.OptBool, t.GetEmbedded().OptBool, t.GetT(), t.GetEmbedded())
	variadicArgsExtInterfaceFunc(t.OptBool, t.GetEmbedded().OptBool, t.GetT(), t.GetEmbedded())
}

func optionalArgsFunc(*bool)                   {}
func nonOptionalArgsFunc(bool)                 {}
func optionalArgs2Func(a, b *bool)             {}
func optionalVariadicArgsFunc(...*bool)        {}
func nonOptionalVariadicArgsFunc(...bool)      {}
func variadicArgsAnyFunc(...any)               {}
func variadicArgsInterfaceFunc(...interface{}) {}

type ExtInterface interface{}

func variadicArgsExtInterfaceFunc(...ExtInterface) {}

func messageArgsFunc(*proto.Embedded) {}
