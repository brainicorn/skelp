package skelplate

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"reflect"
	"testing"

	"github.com/AlecAivazis/survey/core"
	"github.com/xeipuuv/gojsonschema"
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
								"default":["kolsch"],
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
								"default":["kolsch"],
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
									"addPrompt":"add another?"
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
								"choices":[1,2,5,7]
								}]
							}`,
		[]string{" \x0e \x0e "},
		map[string]interface{}{"rounds": []interface{}{float64(1), float64(5)}},
	},
	{
		`{
							  "author": "brainicorn",
							  "variables":[{
								"name":"bar",
								"variables":[{
									"name":"barname",
									"default":"my bar"
									},
									{
									"name":"barslogan",
									"default":"free beer tomorrow"
									}
									]
								}]
							}`,
		[]string{"\n", "\n"},
		map[string]interface{}{"bar": map[string]interface{}{"barname": "my bar", "barslogan": "free beer tomorrow"}},
	},
	{
		`{
							  "author": "brainicorn",
							  "variables":[{
								"name":"database",
								"variables":[{
									"name":"db",
									"required":true,
									"default":"mongo",
									"choices":["mongo","cassandra","dynamo"]
									},
									{
									"name":"namespace",
									"prompt":"Enter a namespace:",
									"default":"",
									"required": true
									},
									{
									"name":"regions",
									"prompt":"Choose your regions:",
									"required": true,
									"default":[""],
									"choices":["east","west","ap"]
									}
								]
							}]
						}`,
		[]string{" \n", "myspace", " \x0e\x0e "},
		map[string]interface{}{"database": map[string]interface{}{"db": "mongo", "namespace": "myspace", "regions": []interface{}{"east", "ap"}}},
	},
}

func TestGatherData(t *testing.T) {
	core.DisableColor = true
	core.QuestionIcon = "%"
	in, _ := ioutil.TempFile("", "")
	defer in.Close()

	os.Stdin = in

	var valmap map[string]interface{}
	dp := NewDataProvider(nil, 0)

	for _, tt := range tmplTests {

		var n string
		dp.beforePrompt = func() {
			n, tt.input = tt.input[0], tt.input[1:]
			in.Truncate(0)
			in.Seek(0, os.SEEK_SET)
			io.WriteString(in, n+"\n")
			in.Seek(0, os.SEEK_SET)
		}

		var descriptor SkelplateDescriptor

		// validate our input
		schemaLoader := gojsonschema.NewStringLoader(GithubComBrainicornSkelpSkelplateSkelplateDescriptor)
		docLoader := gojsonschema.NewBytesLoader([]byte(tt.tmpl))

		schemaValidationResult, verr := gojsonschema.Validate(schemaLoader, docLoader)

		if verr == nil && len(schemaValidationResult.Errors()) > 0 {
			var errBuf bytes.Buffer
			errBuf.WriteString("Error validating skelp descriptor:\n")
			for _, re := range schemaValidationResult.Errors() {
				errBuf.WriteString(fmt.Sprintf("  - %s\n", re))
			}

			t.Error(errBuf.String())
		}

		// input is valid
		err := json.Unmarshal([]byte(tt.tmpl), &descriptor)

		if err != nil {
			t.Errorf("error parsing template: %s\n%s", tt.tmpl, err)
		}

		valmap, err = dp.gatherData(descriptor)

		if err != nil {
			t.Errorf("error gathering data: %s\n%s", tt.tmpl, err)
		}

		// delete default filler data
		delete(valmap, "TemplateAuthor")
		delete(valmap, "TemplateRepo")
		delete(valmap, "TemplateCreated")
		delete(valmap, "TemplateModified")
		delete(valmap, "TemplateDesc")

		if !reflect.DeepEqual(tt.expected, valmap) {
			t.Errorf("template parse error:\n  expected:\n  %+v\n  actual:\n  %+v", tt.expected, valmap)
		}

	}
}
