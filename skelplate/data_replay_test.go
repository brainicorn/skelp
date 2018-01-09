package skelplate

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"os"
	"reflect"
	"testing"

	"github.com/AlecAivazis/survey/core"
)

var replayTests = []struct {
	tmpl        string
	input       []string
	templateDir string
	expected    map[string]interface{}
}{
	//	{
	//		`{
	//											  "author": "brainicorn",
	//											"repository": "https://github.com/brainicorn/skelp",
	//											"created":"2017-06-07T19:15:08+00:00",
	//											"modified":"2017-06-07T19:15:08+00:00",
	//											  "variables":[{"name":"beer", "default":"ipa"}]
	//											}`,
	//		[]string{"porter"},
	//		"github.com/brainicorn/skelplates/test1",
	//		map[string]interface{}{"beer": "porter"},
	//	},

	//	{
	//		`{
	//										  "author": "brainicorn",
	//										"repository": "https://github.com/brainicorn/skelp",
	//										  "variables":[{"name":"beer", "default":"ipa"},{"name":"password", "default":"", "prompt":"enter password", "password":true}]
	//										}`,
	//		[]string{"ipa", "ilikecheese"},
	//		"github.com/brainicorn/skelplates/test2",
	//		map[string]interface{}{"beer": "ipa"},
	//	},
	//	{
	//		`{
	//							  "author": "brainicorn",
	//							  "variables":[{
	//								"name":"bar",
	//								"variables":[{
	//									"name":"barname",
	//									"default":"my bar"
	//									},
	//									{
	//									"name":"barslogan",
	//									"default":"free beer tomorrow"
	//									},
	//									{"name":"password", "default":"kitchenpass", "prompt":"enter bar password", "password":true}
	//									]
	//								}]
	//							}`,
	//		[]string{"\n", "\n", "barpass"},
	//		"github.com/brainicorn/skelplates/test3",
	//		map[string]interface{}{"bar": map[string]interface{}{"barname": "my bar", "barslogan": "free beer tomorrow"}},
	//	},
	//	{
	//		`{
	//											  "author": "brainicorn",
	//											  "variables":[{"name":"beer","default":6.5}]
	//											}`,
	//		[]string{"7.2"},
	//		"github.com/brainicorn/skelplates/test4",
	//		map[string]interface{}{"beer": float64(7.2)},
	//	},
	//	{
	//		`{
	//											  "author": "brainicorn",
	//											  "variables":[{"name":"beer","default":"ipa"}
	//												,{"name":"food", "default":true}
	//												]
	//											}`,
	//		[]string{"ale", "n"},
	//		"github.com/brainicorn/skelplates/test5",
	//		map[string]interface{}{"beer": "ale", "food": false},
	//	},
	{
		`{
								  "author": "brainicorn",
								  "variables":[{
									"name":"databases",
									"addPrompt":"Add Another?",
									"variables":[{
										"name":"db",
										"required":true,
										"default":"mongo"
										},
										{
										"name":"namespace",
										"prompt":"Enter a namespace:",
										"default":"",
										"required": true
										},
										{
										"name":"dbpassword",
										"prompt":"Enter a password:",
										"default":"",
										"required": true,
										"password": true
										}
									]
								}]
							}`,
		[]string{"mongo", "myspace", "mypass", "y", "cassandra", "yourspace", "yourpass", "n"},
		"github.com/brainicorn/skelplates/test6",
		map[string]interface{}{"databases": []map[string]interface{}{{"db": "mongo", "namespace": "myspace"}, {"db": "cassandra", "namespace": "yourspace"}}},
	},
	//	{
	//		`{
	//								  "author": "brainicorn",
	//								  "variables":[{
	//									"name":"post",
	//									"required":true,
	//									"default":""
	//								},
	//								{
	//									"name":"tagcollectors",
	//									"addPrompt":"Add Another?",
	//									"variables":[{
	//										"name":"collection",
	//										"required":true,
	//										"variables":[{
	//											"name":"name",
	//											"required":true,
	//											"default":""
	//										   },
	//										   {
	//											"name":"tags",
	//											"required":true,
	//											"addPrompt":"Add Another Tag?",
	//											"default":[""]
	//										   }]
	//										},
	//										{
	//										"name":"author",
	//										"prompt":"Enter your name:",
	//										"default":"",
	//										"required": true
	//									}]
	//								}]
	//							}`,
	//		[]string{"mypost", "collection1", "tag1", "y", "tag2", "n", "me", "y", "collection2", "tag3", "y", "tag4", "n", "you", "n"},
	//		"github.com/brainicorn/skelplates/test7",
	//		map[string]interface{}{
	//			"post": "mypost",
	//			"tagcollectors": []map[string]interface{}{
	//				{"collection": map[string]interface{}{
	//					"name": "collection1",
	//					"tags": []interface{}{"tag1", "tag2"},
	//				},
	//					"author": "me",
	//				},
	//				{"collection": map[string]interface{}{
	//					"name": "collection2",
	//					"tags": []interface{}{"tag3", "tag4"},
	//				},
	//					"author": "you",
	//				},
	//			},
	//		},
	//	},
}

func TestWritingReplayData(t *testing.T) {
	core.DisableColor = true
	core.QuestionIcon = "%"

	in, _ := ioutil.TempFile("", "")
	defer in.Close()

	os.Stdin = in

	tmpDir, _ := ioutil.TempDir("", "skelp-replay-data-test")
	defer os.RemoveAll(tmpDir)

	for _, tt := range replayTests {
		dp := NewDataProvider(nil, 0)
		var n string
		var replayData map[string]interface{}
		var valmap map[string]interface{}
		var expectedmap map[string]interface{}

		dp.beforePrompt = func() {
			n, tt.input = tt.input[0], tt.input[1:]
			in.Truncate(0)
			in.Seek(0, os.SEEK_SET)
			io.WriteString(in, n+"\n")
			in.Seek(0, os.SEEK_SET)
		}

		// validate our input
		descriptor, verr := ValidateDescriptor([]byte(tt.tmpl))

		if verr != nil {
			t.Error(verr)
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

		// This stinks, but json unmarshalling of the replay data uses []interface{} instead of
		// []map[string]interface{} and so we need to json round trip our expected data so we can
		// use DeepEquals.
		jbites, _ := json.Marshal(tt.expected)
		json.Unmarshal(jbites, &expectedmap)

		_, err = dp.ReplayProvider.WriteData(valmap, tmpDir, tt.templateDir)
		replayData, err = dp.ReplayProvider.ReadData(tmpDir, tt.templateDir)

		// delete default filler data
		delete(replayData, "TemplateAuthor")
		delete(replayData, "TemplateRepo")
		delete(replayData, "TemplateCreated")
		delete(replayData, "TemplateModified")
		delete(replayData, "TemplateDesc")

		if err != nil {
			t.Errorf("error gathering data: %s\n%s", tt.tmpl, err)
		}

		if !reflect.DeepEqual(expectedmap, replayData) {
			t.Errorf("template data error:\n  expected:\n  %+v\n  actual replay:\n  %+v", expectedmap, replayData)
		}

	}
}
