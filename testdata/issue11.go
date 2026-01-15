package testdata

import (
	"github.com/ghostiam/protogetter/testdata/proto"
)

func testIssue11(t *proto.Test) {
	var optBoolVar *bool

	type structWithPtrField struct {
		OptBool *bool
	}

	var foo issue11Foo

	// Invalid

	optBoolVar = t.Embedded.OptBool // want `avoid direct access to proto field t\.Embedded\.OptBool, use t\.GetEmbedded\(\)\.OptBool instead`
	_ = optBoolVar

	_ = structWithPtrField{
		OptBool: t.Embedded.OptBool, // want `avoid direct access to proto field t\.Embedded\.OptBool, use t\.GetEmbedded\(\)\.OptBool instead`
	}

	optionalArgsFunc(t.Embedded.OptBool)                 // want `avoid direct access to proto field t\.Embedded\.OptBool, use t\.GetEmbedded\(\)\.OptBool instead`
	foo.optionalArgsFunc(t.Embedded.OptBool)             // want `avoid direct access to proto field t\.Embedded\.OptBool, use t\.GetEmbedded\(\)\.OptBool instead`
	optionalArgs2Func(t.OptBool, t.Embedded.OptBool)     // want `avoid direct access to proto field t\.Embedded\.OptBool, use t\.GetEmbedded\(\)\.OptBool instead`
	foo.optionalArgs2Func(t.OptBool, t.Embedded.OptBool) // want `avoid direct access to proto field t\.Embedded\.OptBool, use t\.GetEmbedded\(\)\.OptBool instead`
	messageArgsFunc(t.Embedded)                          // want `avoid direct access to proto field t\.Embedded, use t\.GetEmbedded\(\) instead`
	foo.messageArgsFunc(t.Embedded)                      // want `avoid direct access to proto field t\.Embedded, use t\.GetEmbedded\(\) instead`
	nonOptionalArgsFunc(t.T)                             // want `avoid direct access to proto field t\.T, use t\.GetT\(\) instead`
	foo.nonOptionalArgsFunc(t.T)                         // want `avoid direct access to proto field t\.T, use t\.GetT\(\) instead`
	optionalVariadicArgsFunc(
		t.OptBool,
		t.Embedded.OptBool, // want `avoid direct access to proto field t\.Embedded\.OptBool, use t\.GetEmbedded\(\)\.OptBool instead`
	)
	foo.optionalVariadicArgsFunc(
		t.OptBool,
		t.Embedded.OptBool, // want `avoid direct access to proto field t\.Embedded\.OptBool, use t\.GetEmbedded\(\)\.OptBool instead`
	)
	nonOptionalVariadicArgsFunc(
		t.T,                 // want `avoid direct access to proto field t\.T, use t\.GetT\(\) instead`
		*t.Embedded.OptBool, // want `avoid direct access to proto field \*t\.Embedded\.OptBool, use t\.GetEmbedded\(\)\.GetOptBool\(\) instead`
		false,
		true,
	)
	foo.nonOptionalVariadicArgsFunc(
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
	foo.variadicArgsAnyFunc(
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
	foo.variadicArgsInterfaceFunc(
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
	foo.variadicArgsExtInterfaceFunc(
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
	foo.optionalArgsFunc(t.GetEmbedded().OptBool)
	optionalArgs2Func(t.OptBool, t.GetEmbedded().OptBool)
	foo.optionalArgs2Func(t.OptBool, t.GetEmbedded().OptBool)
	messageArgsFunc(t.GetEmbedded())
	foo.messageArgsFunc(t.GetEmbedded())
	nonOptionalArgsFunc(t.GetT())
	foo.nonOptionalArgsFunc(t.GetT())
	optionalVariadicArgsFunc(t.OptBool, t.GetEmbedded().OptBool)
	foo.optionalVariadicArgsFunc(t.OptBool, t.GetEmbedded().OptBool)
	nonOptionalVariadicArgsFunc(t.GetT(), t.GetEmbedded().GetOptBool(), false, true)
	foo.nonOptionalVariadicArgsFunc(t.GetT(), t.GetEmbedded().GetOptBool(), false, true)
	variadicArgsAnyFunc(t.OptBool, t.GetEmbedded().OptBool, t.GetT(), t.GetEmbedded())
	foo.variadicArgsAnyFunc(t.OptBool, t.GetEmbedded().OptBool, t.GetT(), t.GetEmbedded())
	variadicArgsInterfaceFunc(t.OptBool, t.GetEmbedded().OptBool, t.GetT(), t.GetEmbedded())
	foo.variadicArgsInterfaceFunc(t.OptBool, t.GetEmbedded().OptBool, t.GetT(), t.GetEmbedded())
	variadicArgsExtInterfaceFunc(t.OptBool, t.GetEmbedded().OptBool, t.GetT(), t.GetEmbedded())
	foo.variadicArgsExtInterfaceFunc(t.OptBool, t.GetEmbedded().OptBool, t.GetT(), t.GetEmbedded())
}

type ExtInterface interface{}

func optionalArgsFunc(*bool)                       {}
func nonOptionalArgsFunc(bool)                     {}
func optionalArgs2Func(a, b *bool)                 {}
func optionalVariadicArgsFunc(...*bool)            {}
func nonOptionalVariadicArgsFunc(...bool)          {}
func variadicArgsAnyFunc(...any)                   {}
func variadicArgsInterfaceFunc(...interface{})     {}
func messageArgsFunc(*proto.Embedded)              {}
func variadicArgsExtInterfaceFunc(...ExtInterface) {}

type issue11Foo struct{}

func (issue11Foo) optionalArgsFunc(*bool)                       {}
func (issue11Foo) nonOptionalArgsFunc(bool)                     {}
func (issue11Foo) optionalArgs2Func(a, b *bool)                 {}
func (issue11Foo) optionalVariadicArgsFunc(...*bool)            {}
func (issue11Foo) nonOptionalVariadicArgsFunc(...bool)          {}
func (issue11Foo) variadicArgsAnyFunc(...any)                   {}
func (issue11Foo) variadicArgsInterfaceFunc(...interface{})     {}
func (issue11Foo) messageArgsFunc(*proto.Embedded)              {}
func (issue11Foo) variadicArgsExtInterfaceFunc(...ExtInterface) {}
