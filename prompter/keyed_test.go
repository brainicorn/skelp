package prompter

import (
	"io"
	"io/ioutil"
	"os"
	"testing"
)

var stringMinMax = &MinMaxString{Min: float64(2), Max: float64(5)}

type fakeTypingUser struct {
	in         *os.File
	keystrokes []string
}

func newFakeTypingUser(keystrokes []string) *fakeTypingUser {
	in, _ := ioutil.TempFile("", "")
	os.Stdin = in

	return &fakeTypingUser{
		in:         in,
		keystrokes: keystrokes,
	}
}

func (f *fakeTypingUser) nextKeystroke() {
	var keystroke string
	keystroke, f.keystrokes = f.keystrokes[0], f.keystrokes[1:]
	f.in.Truncate(0)
	f.in.Seek(0, os.SEEK_SET)
	io.WriteString(f.in, keystroke+"\n")
	f.in.Seek(0, os.SEEK_SET)
}

func (f *fakeTypingUser) done() {
	f.in.Close()
}

var tmplKeyedTests = []struct {
	ki         *KeyedInput
	keystrokes []string
	expected   string
}{
	{
		&KeyedInput{
			Prompt: Prompt{
				Question: "hello",
				Default:  "world",
				Help:     "just say something",
			},
		},
		[]string{"?", "hi"},
		"hi",
	},
	{
		&KeyedInput{
			Prompt: Prompt{
				Question: "do you like beer?",
				Default:  "n",
			},
			IsConfirm: true,
		},
		[]string{"y"},
		"true",
	},
	{
		&KeyedInput{
			Prompt: Prompt{
				Question: "do you like beer?",
				Default:  "y",
			},
			IsConfirm: true,
		},
		[]string{"n"},
		"false",
	},
	{
		&KeyedInput{
			Prompt: Prompt{
				Question: "do you like beer?",
				Default:  "y",
			},
			IsConfirm: true,
		},
		[]string{"b", "n"},
		"false",
	},
	{
		&KeyedInput{
			Prompt: Prompt{
				Question: "hello",
				Validators: []Validator{
					StringNotBlank,
				},
			},
		},
		[]string{"", "hi"},
		"hi",
	},
	{
		&KeyedInput{
			Prompt: Prompt{
				Question: "how many?",
				Default:  "1",
				Validators: []Validator{
					IsANumber,
				},
			},
		},
		[]string{"two", "2"},
		"2",
	},
	{
		&KeyedInput{
			Prompt: Prompt{
				Question: "hello",
				Validators: []Validator{
					StringNotBlank,
					stringMinMax.CheckMin,
					stringMinMax.CheckMax,
				},
			},
		},
		[]string{"", "hi"},
		"hi",
	},
}

func TestKeyedInput(t *testing.T) {

	for _, tt := range tmplKeyedTests {
		user := newFakeTypingUser(tt.keystrokes)
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
