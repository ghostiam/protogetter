package testdata

import (
	"github.com/ghostiam/protogetter/testdata/proto"
)

type issue23Foo struct{}

func (issue23Foo) present(optional *bool) bool {
	return optional != nil
}

func issue23Present(optional *bool) bool {
	return optional != nil
}

func testIssue23(t *proto.Test) {
	var foo issue23Foo
	foo.present(t.OptBool)
	issue23Foo{}.present(t.OptBool)
	issue23Present(t.OptBool)
}
