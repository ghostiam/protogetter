package testdata

import (
	"github.com/ghostiam/protogetter/testdata/proto"
)

// AliasTest is an alias of a message, as a package re-exporting one declares it. Go materialises
// aliases, so an expression of this type holds a types.Alias rather than the types.Named it
// aliases.
type AliasTest = proto.Test

func testAlias_read(t *AliasTest) string {
	return t.S // want `avoid direct access to proto field t\.S, use t\.GetS\(\) instead`
}

func testAlias_readEmbedded(t *AliasTest) string {
	return t.Embedded.S // want `avoid direct access to proto field t\.Embedded\.S, use t\.GetEmbedded\(\)\.GetS\(\) instead`
}

func testAlias_write(t *AliasTest) {
	t.S = "value"
}

func testAlias_getter(t *AliasTest) string {
	return t.GetS()
}
