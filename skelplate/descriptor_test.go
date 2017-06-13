package skelplate

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"reflect"
	"testing"

	"github.com/AlecAivazis/survey/core"
)

var tmplTests = []struct {
	tmpl     string
	input    []string
	expected map[string]interface{}
}{
	{
		`{
										  "author": "brainicorn",
										"repository": "https://github.com/brainicorn/skelp",
										"created":"2017-06-07T19:15:08+00:00",
										"modified":"2017-06-07T19:15:08+00:00",
										  "variables":[{"name":"beer", "default":"ipa"}]
										}`,
		[]string{"porter"},
		map[string]interface{}{"beer": "porter"},
	},
	{
		`{
										  "author": "brainicorn",
										"repository": "https://github.com/brainicorn/skelp",
										  "variables":[{"name":"beer", "default":"ipa", "prompt":"enter beer style", "min":3,"max":20}]
										}`,
		[]string{"porter"},
		map[string]interface{}{"beer": "porter"},
	},
	{
		`{
										  "author": "brainicorn",
										  "variables":[{"name":"beer","default":"ipa"}
											,{"name":"food", "default":true}
											]
										}`,
		[]string{"ale", "n"},
		map[string]interface{}{"beer": "ale", "food": false},
	},
	{
		`{
										  "author": "brainicorn",
										  "variables":[{"name":"beer","default":"ipa"}
											,{"name":"food", "default":"{{.beer}}"}
											]
										}`,
		[]string{"ale", ""},
		map[string]interface{}{"beer": "ale", "food": "ale"},
	},
	{
		`{
										  "author": "brainicorn",
										  "variables":[{"name":"beer","default":6.5}]
										}`,
		[]string{"7.2"},
		map[string]interface{}{"beer": float64(7.2)},
	},
	{
		`{
										  "author": "brainicorn",
										  "variables":[{"name":"beer","default":6.5,"prompt":"rating","min":1,"max":10,"required":true}]
										}`,
		[]string{"7.2"},
		map[string]interface{}{"beer": float64(7.2)},
	},
	{
		`{
										  "author": "brainicorn",
										  "variables":[{"name":"beer","default":"good"}
											,{"name":"cheese",
												"required":true,
												"default":"gouda"
											}
											]
										}`,
		[]string{"ale", "\n"},
		map[string]interface{}{"beer": "ale", "cheese": "gouda"},
	},
	{
		`{
								  "author": "brainicorn",
								  "variables":[{
									"name":"beer",
									"default":"kolsch",
									"choices":["pale","kolsch","stout"]
									}]
								}`,
		[]string{"\x0e\x0e"},
		map[string]interface{}{"beer": "stout"},
	},
	{
		`{
							  "author": "brainicorn",
							  "variables":[{
								"name":"beer",
								"default":"kolsch",
								"mutlichoice":true,
								"choices":["pale","kolsch","stout"]
								}]
							}`,
		[]string{" \x0e \x0e "},
		map[string]interface{}{"beer": "pale,stout"},
	},
	{
		`{
							  "author": "brainicorn",
							  "variables":[{
								"name":"beer",
								"default":["kolsch"],
								"mutlichoice":true,
								"choices":["pale","kolsch","stout"]
								}]
							}`,
		[]string{" \x0e \x0e "},
		map[string]interface{}{"beer": []interface{}{"pale", "stout"}},
	},
	{
		`{
							  "author": "brainicorn",
							  "variables":[{
							    "name":"beer",
								"required":true,
								"default":["kolsch"],
								"mutlival":true,
								"addprompt":"add another?"
								}]
							}`,
		[]string{"ale", "y", "lager", "n"},
		map[string]interface{}{"beer": []interface{}{"ale", "lager"}},
	},
	{
		`{
							  "author": "brainicorn",
							  "variables":[{
								"name":"rounds",
								"default":[2],
								"mutlichoice":true,
								"choices":["1","2","5","7"]
								}]
							}`,
		[]string{" \x0e \x0e "},
		map[string]interface{}{"rounds": []interface{}{float64(1), float64(5)}},
	},
}

func TestGatherData(t *testing.T) {
	core.DisableColor = true
	core.QuestionIcon = "%"
	in, _ := ioutil.TempFile("", "")
	defer in.Close()

	os.Stdin = in

	var valmap map[string]interface{}
	dp := NewDataProvider(nil)

	for _, tt := range tmplTests {

		var n string
		dp.beforePrompt = func() {
			fmt.Println("before prompt")
			n, tt.input = tt.input[0], tt.input[1:]
			in.Truncate(0)
			in.Seek(0, os.SEEK_SET)
			io.WriteString(in, n+"\n")
			in.Seek(0, os.SEEK_SET)
		}

		var descriptor SkelplateDescriptor
		err := json.Unmarshal([]byte(tt.tmpl), &descriptor)

		if err != nil {
			t.Errorf("error parsing template: %s\n%s", tt.tmpl, err)
		}

		valmap, err = dp.gatherData(descriptor)

		if err != nil {
			t.Errorf("error gathering data: %s\n%s", tt.tmpl, err)
		}

		if !reflect.DeepEqual(tt.expected, valmap) {
			t.Errorf("template parse error:\n  expected:\n  %+v\n  actual:\n  %+v", tt.expected, valmap)
		}

	}
}
