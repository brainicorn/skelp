package provider

// Hooks is a struct for holding the various hook scripts to run
type Hooks struct {
	PreInput []string
	PreGen   []string
	PostGen  []string
}

// HookProvider is a function that returns the hooks to run or an error
type HookProvider func(templateRoot string) (Hooks, error)

func DefaultHookProvider(templateRoot string) (Hooks, error) {
	return Hooks{
		PreInput: []string{},
		PreGen:   []string{},
		PostGen:  []string{},
	}, nil
}
