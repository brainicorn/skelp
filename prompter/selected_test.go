package prompter

import (
	"io"
	"io/ioutil"
	"os"
	"testing"

	"github.com/AlecAivazis/survey/core"
	"github.com/AlecAivazis/survey/terminal"
)

type fakeSelectingUser struct {
	in         *os.File
	keystrokes []string
}

func newFakeSelectingUser(keystrokes []string) *fakeSelectingUser {
	in, _ := ioutil.TempFile("", "")
	os.Stdin = in

	return &fakeSelectingUser{
		in:         in,
		keystrokes: keystrokes,
	}
}

func (f *fakeSelectingUser) nextKeystroke() {
	var keystroke string

	keystroke, f.keystrokes = f.keystrokes[0], f.keystrokes[1:]

	f.in.Truncate(0)
	f.in.Seek(0, os.SEEK_SET)
	io.WriteString(f.in, keystroke+"\n")
	f.in.Seek(0, os.SEEK_SET)
}

func (f *fakeSelectingUser) done() {
	f.in.Close()
}

var tmplSelectedTests = []struct {
	ki         *SelectedInput
	keystrokes []string
	expected   string
}{
	{
		&SelectedInput{
			Prompt: Prompt{
				Question: "1",
			},
			Options: []string{"one", "two", "three", "four"},
			IsMulti: true,
		},
		[]string{"\x0e \x0e "},
		"two,three",
	},
	{
		&SelectedInput{
			Prompt: Prompt{
				Question: "2",
			},
			Options: []string{"one", "two", "three", "four"},
			IsMulti: true,
		},
		[]string{"\x0e\x0e\x10 \x0e "},
		"two,three",
	},
	{
		&SelectedInput{
			Prompt: Prompt{
				Question: "3",
			},
			Options: []string{"one", "two", "three", "four"},
			IsMulti: true,
		},
		[]string{" \x0e\x0e\x10  \x0e "},
		"one,three",
	},
	{
		&SelectedInput{
			Prompt: Prompt{
				Question: "4",
				Default:  "four",
			},
			Options: []string{"one", "two", "three", "four"},
			IsMulti: true,
		},
		[]string{" \x0e\x0e\x10  \x0e "},
		"one,three,four",
	},
	{
		&SelectedInput{
			Prompt: Prompt{
				Question: "5",
				Help:     "pick one",
			},
			Options: []string{"one", "two", "three", "four"},
			IsMulti: true,
		},
		[]string{"? \x0e\x0e\x10  \x0e "},
		"one,three",
	},
	{
		&SelectedInput{
			Prompt: Prompt{
				Question: "6",
			},
			Options: []string{"one", "two", "three", "four"},
		},
		[]string{"\x0e\x0e"},
		"three",
	},
	{
		&SelectedInput{
			Prompt: Prompt{
				Question: "7",
			},
			Options: []string{"one", "two", "three", "four"},
			IsMulti: false,
		},
		[]string{"\x0e\x0e\x10\x0e"},
		"three",
	},
	{
		&SelectedInput{
			Prompt: Prompt{
				Question: "8",
			},
			Options: []string{"one", "two", "three", "four"},
		},
		[]string{"\x0e\x0e\x10\x0e"},
		"three",
	},
	{
		&SelectedInput{
			Prompt: Prompt{
				Question: "9",
				Default:  "four",
			},
			Options: []string{"one", "two", "three", "four"},
		},
		[]string{"\x0e\x0e\x10\x0e"},
		"four",
	},
	{
		&SelectedInput{
			Prompt: Prompt{
				Question: "10",
				Help:     "pick one",
			},
			Options: []string{"one", "two", "three", "four"},
		},
		[]string{"?\x0e\x0e\x10\x0e"},
		"three",
	},
}

func TestSelectedInput(t *testing.T) {
	core.SelectFocusIcon = "▶"
	for _, tt := range tmplSelectedTests {

		user := newFakeSelectingUser(tt.keystrokes)
		defer user.done()

		tt.ki.BeforePrompt = user.nextKeystroke
		ans, err := tt.ki.Ask()

		if err != nil {
			t.Fatalf("error prompting %s", err)
		}

		if ans != tt.expected {
			t.Fatalf("answers don't match, have (%s) want (%s)", ans, tt.expected)
		}
	}
}

var tmplSelectedErrorTests = []struct {
	ki         *SelectedInput
	keystrokes []string
	expected   string
}{
	{
		&SelectedInput{
			Prompt: Prompt{
				Question: "1",
			},
			Options: []string{"one", "two", "three", "four"},
			IsMulti: true,
		},
		[]string{"\x0e " + string(terminal.KeyInterrupt)},
		"cancelled",
	},
	{
		&SelectedInput{
			Prompt: Prompt{
				Question: "1",
			},
			Options: []string{"one", "two", "three", "four"},
			IsMulti: true,
		},
		[]string{"\x0e " + string(terminal.KeyEndTransmission)},
		"cancelled",
	},
	{
		&SelectedInput{
			Prompt: Prompt{
				Question: "1",
			},
			Options: []string{},
			IsMulti: true,
		},
		[]string{},
		"please provide options to select from",
	},
}

func TestSelectedErrorInput(t *testing.T) {
	core.SelectFocusIcon = "▶"
	for _, tt := range tmplSelectedErrorTests {

		user := newFakeSelectingUser(tt.keystrokes)
		defer user.done()

		tt.ki.BeforePrompt = user.nextKeystroke
		_, err := tt.ki.Ask()

		if err == nil || err.Error() != tt.expected {
			t.Fatalf("wrong error. have (%s) want (%s)", err.Error(), tt.expected)
		}
	}
}
