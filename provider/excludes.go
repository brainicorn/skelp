package provider

// ExcludesProvider is a function that returns a list of template files that should be excluded from processing
type ExcludesProvider func() []string

func DefaultExcludesProvider() []string {
	return []string{}
}
