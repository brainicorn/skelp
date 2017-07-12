package provider

import "github.com/brainicorn/skelp/prompter"

// DataProvider is a function that returns the data to be applied to a template or an error
type DataProvider func(templateRoot string) (interface{}, error)

// OverwriteProvider is a function that returns whether or not the provided file should be overwritten
type OverwriteProvider func(rootDir, relFile string) bool

type BasicAuthProvider func() (string, string)

func DefaultOverwriteProvider(rootDir, relFile string) bool {
	return false
}

func AlwaysOverwriteProvider(rootDir, relFile string) bool {
	return true
}

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
