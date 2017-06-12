package provider

// DataProvider is a function that returns the data to be applied to a template or an error
type DataProvider func(templateRoot string) (interface{}, error)

// OverwriteProvider is a function that returns whether or not the provided file should be overwritten
type OverwriteProvider func(rootDir, relFile string) bool

type BasicAuthProvider func() (string, string)

func DefaultOverwriteProvider(rootDir, relFile string) bool {
	return false
}
