package skelplate

import (
	"encoding/json"
	"strings"
	"testing"
)

var tmplErrTests = []struct {
	tmpl     string
	prefill  map[string]interface{}
	expected string
}{
	{
		`{
			"author": "brainicorn",
			"variables":[{"name":"{{.SomeVar", "default":"ipa"}]
		}`,
		nil,
		"unable to parse variable name template:",
	},
	{
		`{
			"author": "brainicorn",
			"variables":[{"name":"beer", "default":"{{ipa"}]
		}`,
		nil,
		"unable to parse variable default template:",
	},
	{
		`{
			"author": "brainicorn",
			"variables":[{"name":"beer", "default":["{{ipa"]}]
		}`,
		nil,
		"unable to parse variable default template:",
	},
	{
		`{
			"author": "brainicorn",
			"variables":[{"name":"beer", "default":""}]
		}`,
		map[string]interface{}{"beer": "{{ipa"},
		"unable to parse data template:",
	},
	{
		`{
				"author": "brainicorn",
				"variables":[{"name":"beer", "default":""}]
			}`,
		map[string]interface{}{"beer": float64(1)},
		"invalid type for provided data entry",
	},
}

func TestTemplateParseErrors(t *testing.T) {

	for _, tt := range tmplErrTests {
		dp := NewDataProvider(tt.prefill)

		var descriptor SkelplateDescriptor
		err := json.Unmarshal([]byte(tt.tmpl), &descriptor)

		if err != nil {
			t.Fatalf("error parsing descriptor: %s\n%s", tt.tmpl, err)
		}

		_, err = dp.gatherData(descriptor)

		if err == nil {
			t.Fatalf("expected error but was nil: %s", tt.tmpl)
		}

		if !strings.HasPrefix(err.Error(), tt.expected) {
			t.Fatalf("wrong error have (%s) want (%s)", err, tt.expected)
		}

	}
}
