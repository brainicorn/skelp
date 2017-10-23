package skelplate

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"github.com/AlecAivazis/survey/terminal"
)

var tmplErrTests = []struct {
	tmpl     string
	prefill  map[string]interface{}
	expected string
}{
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
	{
		`{
					"author": "brainicorn",
					"variables":[{"name":"beer", "default":"yes", "disabled":"{{ipa"}]
				}`,
		map[string]interface{}{"beer": "yes"},
		"unable to parse variable disabled template:",
	},
	{
		`{
					"author": "brainicorn",
					"variables":[{"name":"beer", "default":"yes", "choices":["{{ipa"]}]
				}`,
		map[string]interface{}{"beer": "yes"},
		"unable to parse choice template:",
	},
	{
		`{
				"author": "brainicorn",
				"variables":[{"name":"beer", "default":"yes", "disabled":"%"}]
			}`,
		map[string]interface{}{"beer": "yes"},
		"unable to convert disabled value to a boolean:",
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

type fakeInterruptingUser struct {
	in         *os.File
	keystrokes []string
}

func newFakeInterruptingUser(keystrokes []string) *fakeInterruptingUser {
	in, _ := ioutil.TempFile("", "")
	os.Stdin = in

	return &fakeInterruptingUser{
		in:         in,
		keystrokes: keystrokes,
	}
}

func (f *fakeInterruptingUser) nextKeystroke() {
	var keystroke string
	keystroke, f.keystrokes = f.keystrokes[0], f.keystrokes[1:]
	f.in.Truncate(0)
	f.in.Seek(0, os.SEEK_SET)
	io.WriteString(f.in, keystroke+"\n")
	f.in.Seek(0, os.SEEK_SET)
}

func (f *fakeInterruptingUser) done() {
	f.in.Close()
}

func TestGatherDataInterrupt(t *testing.T) {
	descJSON := `{
				  "author": "brainicorn",
				  "variables":[{"name":"beer", "default":"ipa"}]
				}`

	dp := NewDataProvider(nil)

	user := newFakeInterruptingUser([]string{string(terminal.KeyInterrupt)})
	defer user.done()

	dp.beforePrompt = user.nextKeystroke

	var descriptor SkelplateDescriptor
	err := json.Unmarshal([]byte(descJSON), &descriptor)

	if err != nil {
		t.Fatalf("error parsing descriptor: %s\n%s", descJSON, err)
	}

	_, err = dp.gatherData(descriptor)

	if err == nil || !strings.HasSuffix(err.Error(), "interrupt") {
		t.Errorf("wrong error: have (%s), want (%s)", err, "interrupt")
	}
}
