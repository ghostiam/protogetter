package testdata

import (
	"github.com/ghostiam/protogetter/testdata/proto"
)

type issue22Foo struct {
	Bar issue22Bar
}

type issue22Bar struct {
	Test *proto.Test
}

func (f *issue22Foo) GetBar() {}

func testIssue22() {
	foo := &issue22Foo{}
	_ = foo.Bar.Test.Embedded        // want `avoid direct access to proto field foo\.Bar\.Test\.Embedded, use foo\.Bar\.Test\.GetEmbedded\(\) instead`
	_ = foo.Bar.Test.Embedded.S      // want `avoid direct access to proto field foo\.Bar\.Test\.Embedded\.S, use foo\.Bar\.Test\.GetEmbedded\(\)\.GetS\(\) instead`
	_ = foo.Bar.Test.GetEmbedded().S // want `avoid direct access to proto field foo\.Bar\.Test\.GetEmbedded\(\)\.S, use foo\.Bar\.Test\.GetEmbedded\(\)\.GetS\(\) instead`
	_ = foo.Bar.Test.GetEmbedded().GetS()
	_ = foo.Bar.Test.GetEmbedded()

	bar := &issue22Bar{}
	_ = bar.Test.Embedded        // want `avoid direct access to proto field bar\.Test\.Embedded, use bar\.Test\.GetEmbedded\(\) instead`
	_ = bar.Test.Embedded.S      // want `avoid direct access to proto field bar\.Test\.Embedded\.S, use bar\.Test\.GetEmbedded\(\)\.GetS\(\) instead`
	_ = bar.Test.GetEmbedded().S // want `avoid direct access to proto field bar\.Test\.GetEmbedded\(\)\.S, use bar\.Test\.GetEmbedded\(\)\.GetS\(\) instead`
	_ = bar.Test.GetEmbedded().GetS()
	_ = bar.Test.GetEmbedded()
}
