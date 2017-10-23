package prompter

import (
	"github.com/AlecAivazis/survey/core"
)

type Prompter interface {
	Ask() (string, error)
	Validate(val string) error
	Error(error) error
}

type Prompt struct {
	core.Renderer
	Question     string
	Default      string
	Help         string
	Validators   []Validator
	BeforePrompt func()
}

func (p *Prompt) Validate(val string) error {
	var err error
	for _, v := range p.Validators {
		err = v(val)

		if err != nil {
			return err
		}
	}

	return err
}

type InputTemplateData struct {
	Prompt
	Answer     string
	ShowAnswer bool
	ShowHelp   bool
}
