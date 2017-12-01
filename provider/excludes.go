package provider

// ExcludesProvider is a function that returns a list of template files that should be excluded from processing
type ExcludesProvider func(templateRoot string) (map[string]bool, error)

func DefaultExcludesProvider(templateRoot string) (map[string]bool, error) {
	return make(map[string]bool), nil
}
