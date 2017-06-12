package prompter

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"testing"
)

type fakeSelectingUser struct {
	in         *os.File
	keystrokes []string
}

func newFakeSelectingUser(keystrokes ...string) *fakeSelectingUser {
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

func TestSelectedInput(t *testing.T) {
	user := newFakeSelectingUser("\x0e \x0e \n")
	defer user.done()

	ki := &SelectedInput{
		Prompt: Prompt{
			Question: "hello",
		},
		Options: []string{"you", "world", "me", "dog"},
		IsMulti: true,
	}

	ki.BeforePrompt = user.nextKeystroke
	s, e := ki.Ask()

	if e != nil {
		t.Fatalf("error prompting %s", e)
	}

	fmt.Println(s)
}
