package provider

import "github.com/brainicorn/skelp/prompter"

type BasicAuthProvider func() (string, string)

type DefaultBasicAuthProvider struct {
	BeforePrompt func()
}

func (bap *DefaultBasicAuthProvider) ProvideAuth() (string, string) {
	var u, p string

	userPrompt := &prompter.KeyedInput{
		Prompt: prompter.Prompt{
			BeforePrompt: bap.BeforePrompt,
			Question:     "enter username:",
		},
	}

	passPrompt := &prompter.KeyedInput{
		Prompt: prompter.Prompt{
			BeforePrompt: bap.BeforePrompt,
			Question:     "enter password:",
		},
		IsPassword: true,
	}

	u, _ = userPrompt.Ask()
	p, _ = passPrompt.Ask()

	return u, p
}
