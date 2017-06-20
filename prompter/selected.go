package prompter

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/AlecAivazis/survey/core"
	"github.com/AlecAivazis/survey/terminal"
)

type SelectedInput struct {
	Prompt
	IsMulti       bool
	Options       []string
	selectedIndex int
	checked       map[int]bool
	showingHelp   bool
}

var SelectedInputTemplate = `
{{- if .ShowHelp }}{{- color "cyan"}}{{ HelpIcon }} {{ .Help }}{{color "reset"}}{{"\n"}}{{end}}
{{- color "green+hb"}}{{ QuestionIcon }} {{color "reset"}}
{{- color "default+hb"}}{{ .Question }}{{color "reset"}}
{{- if .ShowAnswer}}{{color "cyan"}} {{.Answer}}{{color "reset"}}{{"\n"}}
{{- else }}
  {{- if and .Help (not .ShowHelp)}} {{color "cyan"}}[{{ HelpInputRune }} for help]{{color "reset"}}{{end}}
  {{- "\n"}}
  {{- range $ix, $option := .Options}}
    {{- if eq $ix $.SelectedIndex}}{{color "cyan"}}{{ SelectFocusIcon }}{{color "reset"}}{{else}} {{end}}
    {{- if index $.Checked $ix}}{{color "green"}} {{ MarkedOptionIcon }} {{else}}{{color "default+hb"}} {{ UnmarkedOptionIcon }} {{end}}
    {{- color "reset"}}
    {{- " "}}{{$option}}{{"\n"}}
  {{- end}}
{{- end}}`

type SelectedTemplateData struct {
	InputTemplateData
	Checked       map[int]bool
	SelectedIndex int
	Options       []string
}

//MULTI
// OnChange is called on every keypress.
func (s *SelectedInput) OnChange(line []rune, pos int, key rune) (newLine []rune, newPos int, ok bool) {
	if key == terminal.KeyArrowUp && s.selectedIndex > 0 {
		// decrement the selected index
		s.selectedIndex--
		if !s.IsMulti {
			s.checked = map[int]bool{s.selectedIndex: true}
		}
	} else if key == terminal.KeyArrowDown && s.selectedIndex < len(s.Options)-1 {
		// if the user pressed down and there is room to move
		// increment the selected index
		s.selectedIndex++
		if !s.IsMulti {
			s.checked = map[int]bool{s.selectedIndex: true}
		}
	} else if key == terminal.KeySpace && s.IsMulti {
		if old, ok := s.checked[s.selectedIndex]; !ok {
			// otherwise just invert the current value
			s.checked[s.selectedIndex] = true
		} else {
			// otherwise just invert the current value
			s.checked[s.selectedIndex] = !old
		}
		// only show the help message if we have one to show
	} else if key == core.HelpInputRune && s.Help != "" {
		s.showingHelp = true
	}

	// render the options
	s.Render(
		SelectedInputTemplate,
		SelectedTemplateData{
			InputTemplateData: InputTemplateData{
				Prompt:   s.Prompt,
				ShowHelp: s.showingHelp,
			},
			SelectedIndex: s.selectedIndex,
			Checked:       s.checked,
			Options:       s.Options,
		},
	)

	// if we are not pressing ent
	return line, 0, true
}

func (s *SelectedInput) Ask() (string, error) {
	// if there are no options to render
	if len(s.Options) == 0 {
		// we failed
		return "", errors.New("please provide options to select from")
	}

	// compute the default state
	s.checked = make(map[int]bool)
	defaults := strings.Split(s.Default, ",")

	// if there is a default
	if len(defaults) > 0 {
		for _, dflt := range defaults {
			for i, opt := range s.Options {
				// if the option correponds to the default
				if opt == dflt {
					// we found our initial value
					s.checked[i] = true
					if !s.IsMulti {
						s.selectedIndex = i
					}
					// stop looking
					break
				}
			}
		}
	}

	// hide the cursor
	terminal.CursorHide()
	// show the cursor when we're done
	defer terminal.CursorShow()

	if s.BeforePrompt != nil {
		s.BeforePrompt()
	}

	s.Render(
		SelectedInputTemplate,
		SelectedTemplateData{
			InputTemplateData: InputTemplateData{
				Prompt: s.Prompt,
			},
			SelectedIndex: s.selectedIndex,
			Checked:       s.checked,
			Options:       s.Options,
		},
	)

	rr := terminal.NewRuneReader(os.Stdin)
	rr.SetTermMode()
	defer rr.RestoreTermMode()

	// start waiting for input
	for {
		r, _, _ := rr.ReadRune()
		if r == '\r' || r == '\n' {
			break
		}
		if r == terminal.KeyInterrupt {
			return "", fmt.Errorf("cancelled")
		}
		if r == terminal.KeyEndTransmission {
			return "", fmt.Errorf("cancelled")
		}
		s.OnChange(nil, 0, r)
	}

	answers := []string{}
	for ix, option := range s.Options {
		if val, ok := s.checked[ix]; ok && val {
			answers = append(answers, option)
		}
	}

	return strings.Join(answers, ","), nil
}
