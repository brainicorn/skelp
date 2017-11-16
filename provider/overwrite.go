package provider

// OverwriteProvider is a function that returns whether or not the provided file should be overwritten
type OverwriteProvider func(rootDir, relFile string) bool

func DefaultOverwriteProvider(rootDir, relFile string) bool {
	return false
}

func AlwaysOverwriteProvider(rootDir, relFile string) bool {
	return true
}
