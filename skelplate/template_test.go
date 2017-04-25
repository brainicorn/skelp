package skelplate_test

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/brainicorn/skelp/skelplate"
)

var tmplTests = []struct {
	in       string             // input
	expected skelplate.Template // expected result
}{
	{`{
  "author": "brainicorn",
  "variables":{
    "beer":"ipa"
  }
}`, skelplate.Template{
		TemplateAuthor:    "brainicorn",
		TemplateVariables: map[string]interface{}{"beer": "ipa"},
	}},
	{`{
  "author": "brainicorn",
  "variables":{
    "beer":"ipa",
	"food":true
  }
}`, skelplate.Template{
		TemplateAuthor:    "brainicorn",
		TemplateVariables: map[string]interface{}{"beer": "ipa", "food": true},
	}},
	{`{
  "author": "brainicorn",
  "variables":{
    "beer":6.5
  }
}`, skelplate.Template{
		TemplateAuthor:    "brainicorn",
		TemplateVariables: map[string]interface{}{"beer": float64(6.5)},
	}},
	{`{
  "author": "brainicorn",
  "variables":{
    "beer":"good",
	"cheese":{
		"required":true,
		"default":"gouda"
	}
  }
}`, skelplate.Template{
		TemplateAuthor: "brainicorn",
		TemplateVariables: map[string]interface{}{
			"beer": "good",
			"cheese": skelplate.Configurable{
				Required: true,
				Default:  "gouda",
			},
		},
	}},
}

func TestParsing(t *testing.T) {
	for _, tt := range tmplTests {
		var actual skelplate.Template
		err := json.Unmarshal([]byte(tt.in), &actual)

		if err != nil {
			t.Errorf("error parsing template: %s\n%s", tt.in, err)
		}
		if !reflect.DeepEqual(tt.expected, actual) {
			t.Errorf("template parse error:\n  expected:\n  %+v\n  actual:\n  %+v", tt.expected, actual)
		}
	}
}
