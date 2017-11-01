package skelplate

import (
	"testing"
)

func TestBadDotKeyLookup(t *testing.T) {
	datamap := make(map[string]interface{})
	submap := map[string]interface{}{"slug": "slow"}

	datamap["root"] = submap

	_, found := getDotKeyFromMap("root.slug.slow", datamap)

	if found {
		t.Error("should return false for bad key")
	}
}
