package skelplate

import (
	"strings"
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

func TestHooksMissingDescriptor(t *testing.T) {
	dp := NewDataProvider(nil, 0)

	_, err := dp.HookProviderFunc("/some/bad/root/")

	if err == nil || !strings.HasPrefix(err.Error(), "skelp.json not found:") {
		t.Errorf("wrong error: have (%s), want (%s)", err, "skelp.json not found:")
	}
}
