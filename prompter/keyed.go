package prompter

import (
	"os"
	"strings"

	"github.com/AlecAivazis/survey/core"
	"github.com/AlecAivazis/survey/terminal"
)

type KeyedInput struct {
	Prompt
	IsConfirm bool
}

var (
	KeyedInputTemplate = `
{{- if .ShowHelp }}{{- color "cyan"}}{{ HelpIcon }} {{ .Help }}{{color "reset"}}{{"\n"}}{{end}}
{{- color "green+hb"}}{{ QuestionIcon }} {{color "reset"}}
{{- color "default+hb"}}{{ .Question }} {{color "reset"}}
{{- if .ShowAnswer}}
  {{- color "cyan"}}{{.Answer}}{{color "reset"}}{{"\n"}}
{{- else }}
  {{- if and .Help (not .ShowHelp)}}{{color "cyan"}}[{{ HelpInputRune }} for help]{{color "reset"}} {{end}}
  {{- color "white"}}
	{{- if .IsConfirm}}
		{{- if .BoolDefault}}(Y/n) {{- else}}(y/N) {{end}}
	{{- else if .Default}}({{.Default}}) {{end}}{{color "reset"}}
{{- end}}`
)

type KeyedTemplateData struct {
	InputTemplateData
	BoolDefault bool
	IsConfirm   bool
}

func trueOrFalseBool(val string) bool {
	tf := false

	if truePattern.Match([]byte(val)) {
		tf = true
	}

	return tf
}

func trueOrFalseString(val string) string {
	tf := "false"

	if truePattern.Match([]byte(val)) {
		tf = "true"
	}

	return tf
}

func (i *KeyedInput) Ask() (string, error) {
	if i.BeforePrompt != nil {
		i.BeforePrompt()
	}

	i.yesNoValidatorIfNeeded()

	// render the template
	err := i.Render(
		KeyedInputTemplate,
		KeyedTemplateData{
			InputTemplateData: InputTemplateData{
				Prompt: i.Prompt,
			},
			BoolDefault: trueOrFalseBool(i.Default),
			IsConfirm:   i.IsConfirm,
		},
	)
	if err != nil {
		return "", err
	}

	rr := terminal.NewRuneReader(os.Stdin)
	rr.SetTermMode()
	defer rr.RestoreTermMode()

	line := []rune{}
	// get the next line
	for {
		line, err = rr.ReadLine(0)
		if err != nil {
			return string(line), err
		}
		// terminal will echo the \n so we need to jump back up one row
		terminal.CursorPreviousLine(1)

		if string(line) == string(core.HelpInputRune) && i.Help != "" {
			if i.BeforePrompt != nil {
				i.BeforePrompt()
			}
			err = i.Render(
				KeyedInputTemplate,
				KeyedTemplateData{
					InputTemplateData: InputTemplateData{
						Prompt:   i.Prompt,
						ShowHelp: true,
					},
					BoolDefault: trueOrFalseBool(i.Default),
					IsConfirm:   i.IsConfirm,
				},
			)

			if err != nil {
				return "", err
			}
			continue
		}
		break
	}

	ans := string(line)

	// if the line is empty
	if len(strings.TrimSpace(ans)) < 1 {
		// use the default value
		ans = i.Default
	}

	// wait for a valid response
	for invalid := i.Validate(ans); invalid != nil; invalid = i.Validate(ans) {
		err = i.Prompt.Error(invalid)
		// if there was a problem
		if err != nil {
			return "", err
		}

		// ask for more input
		ans, err = i.Ask()
		// if there was a problem
		if err != nil {
			return "", err
		}
	}

	i.Render(
		KeyedInputTemplate,
		KeyedTemplateData{
			InputTemplateData: InputTemplateData{
				Prompt:     i.Prompt,
				Answer:     ans,
				ShowAnswer: true,
			},
			BoolDefault: trueOrFalseBool(i.Default),
			IsConfirm:   i.IsConfirm,
		},
	)

	if i.IsConfirm {
		return trueOrFalseString(ans), err
	}
	return ans, err
}

func (i *KeyedInput) yesNoValidatorIfNeeded() {
	if i.IsConfirm {
		if i.Validators == nil {
			i.Validators = []Validator{}
		}

		i.Validators = append(i.Validators, YesOrNo)
	}
}
